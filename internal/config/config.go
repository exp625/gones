package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	LastROMFile = "last_rom_file"
)

var config map[string]string

func init() {
	config = make(map[string]string)
	filePath, err := FilePath()
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0777); err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("could not close file: ", err.Error())
		}
	}()

	// Check if file is empty
	stat, _ := file.Stat()
	if stat.Size() == 0 {
		// File is empty, create empty config
		config = make(map[string]string, 0)
	} else {
		// File is not empty, try to parse config
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			log.Fatal(err)
		}
	}
}

func FilePath() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "gones", "gones.json"), nil
}

func Set(key string, value string) error {
	filePath, err := FilePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0777); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("failed to close file: ", err.Error())
		}
	}()

	config[key] = value
	return json.NewEncoder(file).Encode(&config)
}

func Get(key string) (string, bool) {
	value, ok := config[key]
	return value, ok
}
