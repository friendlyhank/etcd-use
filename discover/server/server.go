package server

import "net"

//Server -
type Server struct{
	ln net.Listener
}

func NewServer()(*Server,error){
	return &Server{},nil
}
