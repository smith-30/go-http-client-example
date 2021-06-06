package restclient

type StatusError struct {
	Message string
}

func (e *StatusError) Error() string {
	return e.Message
}
