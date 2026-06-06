package web

type Pairable[P any] interface {
	*P
	Name() string
	Offline() bool
}
