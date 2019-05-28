package tbot

import (
	"io"
	"log"
	"os"
)

type flog struct {
	*log.Logger
	f io.Writer
}

func newFileLogger() logger {
	f, err := os.OpenFile("tbot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return &flog{log.New(f, "", 0), f}

}

func (l *flog) logInfo(mes string) {
	l.Println("Info: " + mes)
}

func (l *flog) logError(mes string, err error) {
	l.Println("Error: "+mes, err)
}

func (l *flog) logPanic(mes string, err error) {
	l.Panicln("Panic: "+mes, err)
}
