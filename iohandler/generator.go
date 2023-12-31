package iohandler

type Generator interface {
	Open(name string) (err error)
	Close() error
	Write(folders []*EffectItems) (err error)
}
