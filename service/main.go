package service

type Interface interface {
	Status() status.StatusInterface
}

type ClientSet struct {
	
}