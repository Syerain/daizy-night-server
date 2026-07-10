package handler

type Request struct {
	requestType RequestType
}

type RequestType int

const (
	LoginRequest RequestType = iota
)
