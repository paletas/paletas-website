package server

type Server interface {
	Start(listenPort string) error
}
