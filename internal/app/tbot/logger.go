package tbot

type logger interface {
	logInfo(mes string)
	logError(mes string, err error)
	logPanic(mes string, err error)
}
