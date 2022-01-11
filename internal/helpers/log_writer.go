package helpers

import (
	"errors"
	"github.com/mdanialr/webhook/internal/config"
	"log"
	"os"
)

var (
	NzLogInf *log.Logger
	NzLogErr *log.Logger
)

// InitNzLog init and setup log file to write log about this app
// internal log.
func InitNzLog() error {
	fl, err := os.OpenFile(config.Conf.LogDir+"app-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		errMsg := "failed to open|create log file: " + err.Error()
		return errors.New(errMsg)
	}

	NzLogInf = log.New(fl, "[INFO] ", log.Ldate|log.Ltime)
	NzLogErr = log.New(fl, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
