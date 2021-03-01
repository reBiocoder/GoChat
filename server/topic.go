package main

// Session Update: user agent change或者background session becoming normal
//if sess is nil then user agent change
type sessionUpdate struct {
	sess      *Session
	userAgent string
}
