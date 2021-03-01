package main

import "GoChat/server/store/types"

type sessionLeave struct {
	pkt  *ClientComMessage
	sess *Session
}

type metaReq struct {
	pkt     *ClientComMessage
	sess    *Session
	forUser types.Uid
	state   types.ObjState
}
