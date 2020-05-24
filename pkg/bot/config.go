package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Config config for the bot
type Config struct {
	OwnersID []string `json:"ownersid"`
	DataPath string   `json:"datapath"`
}

type fileNotFound error

func loadConfig(path string) (*Config, error) {
	conf, err := extract(path)
	if err != nil {
		if e, ok := err.(fileNotFound); ok {
			log.Println(e, "recover by creating default config to", path)
			conf, err = saveDefaultConfig(path)
			if err != nil {
				if conf != nil {
					log.Println("error while saving the default config the bot will continue, please change the path to the config file", err)
					return conf, nil
				}
				return nil, err
			}
		}
		return nil, err
	}
	return conf, nil
}

func extract(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fileNotFound(err)
	}
	var conf Config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func saveDefaultConfig(path string) (*Config, error) {
	defaultConf := &Config{
		OwnersID: []string{},
		DataPath: "data",
	}

	//try to save
	data, err := json.MarshalIndent(defaultConf, "", "	")
	if err != nil {
		return defaultConf, err
	}

	err = ioutil.WriteFile(path, data, 0655)

	return defaultConf, err
}
