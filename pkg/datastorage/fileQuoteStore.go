package datastorage

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type saveableQuote struct {
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

func (quote *saveableQuote) hash() string {
	data, err := json.Marshal(quote) //not efficient
	if err != nil {
		log.Fatal("error while computing hash for", quote, err)
	}
	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		log.Fatal("error while computing hash for", quote, err)
	}
	return string(h.Sum(nil))
}

type userQuotes struct {
	UserID  string                   `json:"userID"`
	GuildID string                   `json:"guildID"`
	Quotes  map[string]saveableQuote `json:"quotes"` //hash map / hashset
}

const dirPerm = 0777
const filePerm = 0666
const fileExtension = ".json"

//root/
//   +-{guildID}/
//	 |  +{userID}.json

type fileQuoteStore struct {
	rootPath     string
	guildMutexes map[string]sync.RWMutex
	storMutex    sync.Mutex
}

//constructor
func newfileQuoteStore(uri string) (*fileQuoteStore, error) {

	//Open directory, if can't try to create it
	rootdir, err := os.Open(uri)
	if err != nil {
		e2 := os.Mkdir(uri, dirPerm) //+rwx
		if e2 != nil {
			return nil, e2
		}
		rootdir, e2 = os.Open(uri)
		if e2 != nil {
			return nil, e2
		}
	}

	//check if rootdir is a directory
	stat, err := rootdir.Stat()
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New(uri + " is not a directory")
	}

	//feed mutexes (one per guild directory)

	mutexes, err := feedMutexes(rootdir)
	if err != nil {
		return nil, err
	}

	return &fileQuoteStore{
		rootPath:     uri,
		guildMutexes: mutexes,
		storMutex:    sync.Mutex{},
	}, nil
}

//helper functions
func feedMutexes(dir *os.File) (map[string]sync.RWMutex, error) {

	infos, err := dir.Readdir(0) //read all infos in the directory in one time
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0, len(infos))

	for _, info := range infos {
		if !info.IsDir() {
			dirs = append(dirs, info.Name())
		}
	}

	mutexes := make(map[string]sync.RWMutex, len(dirs)) // size of rootdir directories
	for _, guildID := range dirs {
		mutexes[guildID] = sync.RWMutex{}
	}

	return mutexes, nil
}

func (qs *fileQuoteStore) newGuild(guildID string) {
	qs.storMutex.Lock()
	defer qs.storMutex.Unlock() // unlock mutex at the end of the method
	if _, ok := qs.guildMutexes[guildID]; !ok {
		os.Mkdir(qs.rootPath+string(os.PathSeparator)+guildID, dirPerm)
		qs.guildMutexes[guildID] = sync.RWMutex{}
	}
}

func (qs *fileQuoteStore) aquireWrite(guildID string) *sync.RWMutex {
	for { //active wait until mutex is created and aquire
		if mutex, ok := qs.guildMutexes[guildID]; ok {
			mutex.Lock()
			return &mutex
		}
		qs.newGuild(guildID)
	}
}

func (qs *fileQuoteStore) getUserQuotes(userID string, guildID string) (*userQuotes, error) {
	var res userQuotes
	fullpath := qs.rootPath + string(os.PathSeparator) + guildID + string(os.PathSeparator) + userID + fileExtension
	data, err := ioutil.ReadFile(fullpath)
	if err != nil { //return default userQuotes
		res.GuildID = guildID
		res.UserID = userID
		res.Quotes = map[string]saveableQuote{}
	} else {
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil, errors.New("malformed file " + fullpath)
		}
	}
	return &res, nil
}

func (qs *fileQuoteStore) saveUserQuotes(quotes *userQuotes) error {
	data, err := json.MarshalIndent(quotes, "", "	")
	if err != nil {
		return err
	}
	fullpath := qs.rootPath + string(os.PathSeparator) + quotes.GuildID + string(os.PathSeparator) + quotes.UserID + fileExtension
	return ioutil.WriteFile(fullpath, data, filePerm)
}

//quoteStoreMethods

func (qs *fileQuoteStore) Save(quote *Quote) error {
	//first lock guild data
	mutex := qs.aquireWrite(quote.GuildID)
	defer mutex.Unlock()

	//load quotes
	userQuotes, err := qs.getUserQuotes(quote.UserID, quote.GuildID)
	if err != nil {
		return err
	}

	//update quotes
	saveable := saveableQuote{
		Timestamp: quote.Timestamp,
		Content:   quote.Content,
	}
	userQuotes.Quotes[quote.QuoteId] = saveable

	//Save modified quotes
	return qs.saveUserQuotes(userQuotes)
}

func (qs *fileQuoteStore) Forget(quote *Quote) error {
	mutex := qs.aquireWrite(quote.GuildID)
	defer mutex.Unlock()

	userQuotes, err := qs.getUserQuotes(quote.UserID, quote.GuildID)
	if err != nil {
		return err
	}

	delete(userQuotes.Quotes, quote.QuoteId)

	return qs.saveUserQuotes(userQuotes)
}
