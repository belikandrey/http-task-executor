package tasks

import "net/http"

// Executor represents tasks executor.
type Executor interface {
	Execute(value []byte)
}

// ClientProvider represents mechanism for http.Client creation.
type ClientProvider interface {
	Client() *http.Client
}
