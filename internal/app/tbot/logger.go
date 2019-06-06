package tbot

type logger interface {
	Info(mes string)
	Error(mes string, err error)
	Panic(mes string, err error)
}
