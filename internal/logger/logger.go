package logger

import (
	"log"
)

type HTTPReqInfo struct {
	Time         string
	Method       string
	ResponseCode string
	IPadress     string
	RequestID    any
}

func (i *HTTPReqInfo) Info() {
	log.Println(i)
}
