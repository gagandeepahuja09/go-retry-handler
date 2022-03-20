package goretryhandler

import "net/http"

type Plugin interface {
	OnRequestStart(*http.Request)
	OnRequestEnd(*http.Request, *http.Response)
	OnError(*http.Request, error)
}
