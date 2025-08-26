package tasks

import "net/http"

type Executor interface {
	Execute(value []byte)
}

type ClientProvider interface {
	Client() *http.Client
}
