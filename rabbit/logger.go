package rabbit

import (
	"io/ioutil"
	"log"
	"os"
)

func (a *Rabbit) setupLogger() {
	if a.config.Debug {
		a.logger = log.New(os.Stdout, "", 0)
	} else {
		a.logger = log.New(ioutil.Discard, "", 0)
	}
}

func (a *Rabbit) Logger() *log.Logger {
	return a.logger
}
