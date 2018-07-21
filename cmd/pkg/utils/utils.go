package utils

import (
	"fmt"
	"net"
	"time"
)


//ServerInfo xxx
type ServerInfo struct{
	Host string
	Port int
	Proto string
}



func NewServer(host,proto string, port int)(*ServerInfo,error){
	if len(host) == 0|| len(proto) == 0||port <=0{
		return nil, fmt.Errorf("host,proto or port error")
	}
	return &ServerInfo{
		Host :host,
		Port : port,
		Proto :proto,
	}, nil
}

func SetTimeOut(conn net.Conn, timeout int){
	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}