package datastorage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type savableQuote struct {
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
	Score     int       `json:"score"`
}

type userQuotes struct {
	UserID  string                  `json:"userID"`
	GuildID string                  `json:"guildID"`
	Quotes  map[string]savableQuote `json:"quotes"` //hash map / hashset
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
	storeMutex   sync.Mutex
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
		storeMutex:   sync.Mutex{},
	}, nil
}

//helper functions
func feedMutexes(dir *os.File) (map[string]sync.RWMutex, error) {

	infos, err := dir.Readdir(0) //read all infos in the directory at once
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0, len(infos))

	for _, info := range infos {
		if info.IsDir() {
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
	qs.storeMutex.Lock()
	defer qs.storeMutex.Unlock() // unlock mutex at the end of the method
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

func (qs *fileQuoteStore) aquireRead(guildID string) *sync.RWMutex {
	if mutex, ok := qs.guildMutexes[guildID]; ok {
		mutex.RLock()
		return &mutex
	}
	return nil
}

func (qs *fileQuoteStore) getUserQuotes(userID string, guildID string) (*userQuotes, error) {
	var res userQuotes
	filepath := qs.rootPath + string(os.PathSeparator) + guildID + string(os.PathSeparator) + userID + fileExtension
	data, err := ioutil.ReadFile(filepath)
	if err != nil { //return default userQuotes
		res.GuildID = guildID
		res.UserID = userID
		res.Quotes = map[string]savableQuote{}
	} else {
		err = json.Unmarshal(data, &res)
		if err != nil {
			return nil, errors.New("malformed file " + filepath)
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
	savable := savableQuote{
		Timestamp: quote.Timestamp,
		Content:   quote.Content,
	}
	userQuotes.Quotes[quote.QuoteId] = savable

	//Save modified quotes
	return qs.saveUserQuotes(userQuotes)
}

func (qs *fileQuoteStore) Delete(quoteID string, userID string, guildID string) error {
	mutex := qs.aquireWrite(guildID)
	defer mutex.Unlock()

	userQuotes, err := qs.getUserQuotes(userID, guildID)
	if err != nil {
		return err
	}

	delete(userQuotes.Quotes, quoteID)

	return qs.saveUserQuotes(userQuotes)
}

func (qs *fileQuoteStore) GetQuotesFromUser(userID string, guildID string) ([]Quote, error) {
	mutex := qs.aquireRead(guildID)
	defer mutex.RUnlock()

	userQuotes, err := qs.getUserQuotes(userID, guildID)
	if err != nil {
		return nil, err
	}

	res := make([]Quote, len(userQuotes.Quotes))
	index := 0
	for id, quote := range userQuotes.Quotes {
		res[index] = Quote{
			QuoteId:   id,
			UserID:    userQuotes.UserID,
			GuildID:   userQuotes.GuildID,
			Timestamp: quote.Timestamp,
			Score:     quote.Score,
			Content:   quote.Content,
		}
		index++
	}

	return res, nil
}
