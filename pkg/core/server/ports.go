package server

type Repository interface {
}

type Service interface {
	ListenProtocol() error
	StopServer() error
}
