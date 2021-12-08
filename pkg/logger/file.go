package logger

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileLogger struct {
	loggingEnabled bool
	Logger         io.ReadWriteCloser
	Name           string
}

func (logger *FileLogger) StartLogging() {
	if logger.Name == "" {
		logger.Name = time.Now().Format("log/2006-01-02_15-04-05_nes.log")
	}

	logger.ensureLogDir(logger.Name)
	fo, err := os.Create(logger.Name)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		return
	}
	defer plz.Close(fo)

	f, err := os.OpenFile(logger.Name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
		return
	}
	logger.Logger = f
	log.SetOutput(f)
	logger.loggingEnabled = true
}

func (logger *FileLogger) StopLogging() {
	logger.loggingEnabled = false
	log.SetOutput(os.Stdout)
	plz.Close(logger.Logger)
	logger.Logger = nil
}

func (logger *FileLogger) ensureLogDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			log.Fatal(merr)
		}
	}
}

func (logger *FileLogger) LoggingEnabled() bool {
	return logger.loggingEnabled
}

func (logger *FileLogger) LogLine(logLine string) {
	if logger.Logger == nil {
		return
	}
	plz.Just(fmt.Fprintln(logger.Logger, logLine))
}

