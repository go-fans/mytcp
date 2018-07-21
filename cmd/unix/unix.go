package unix

import (
	"net"
	"fmt"
	"os"
)

const (
	DefaultBuffSize = 1024
	DefaultSockFileName = "/tmp/tmp.sock"
)

//UnixServerInfo xxx
type UnixServerInfo struct{
	SockFileName string
	SockBuffSize int
}

var us = UnixServerInfo{
	SockFileName:DefaultSockFileName,
	SockBuffSize:DefaultBuffSize,
}

//New xxx
func New(filename string,size int)(*UnixServerInfo){
	if len(filename) >0 {
		us.SockFileName = filename
	}
	if size >0 {
		us.SockBuffSize = size
	}
	return &us
}

func (us *UnixServerInfo) ServerCreate() error {
	os.Remove(us.SockFileName)
	addr, err := net.ResolveUnixAddr("unix",us.SockFileName )
	if err != nil{
		return err
	}
	listener , err := net.ListenUnix("unix", addr)
	if err != nil{
		return err
	}
	defer listener.Close()
	fmt.Println("listening on ", listener.Addr().String())

	for{
		c , err := listener.Accept()
		if err != nil{
			panic(err)
		}
		go us.handle(c)
	}
	return nil
}

func (us *UnixServerInfo) handle(c net.Conn)  {
	defer c.Close()
	buf := make([]byte,us.SockBuffSize )
	n , err := c.Read(buf)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("From [%s] recv [%s]\n", c.RemoteAddr().String(), string(buf[:n]))
}

//ClientCreate xxx
func (us *UnixServerInfo) ClientCreate() error {
	addr, err := net.ResolveUnixAddr("unix",us.SockFileName )
	if err != nil{
		return err
	}
	conn , err := net.DialUnix("unix", nil, addr)
	if err != nil{
		return err
	}
	defer conn.Close()
	_ , err = conn.Write([]byte("hello"))
	if err != nil{
		return err
	}
	return nil
}