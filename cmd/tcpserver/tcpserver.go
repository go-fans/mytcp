package tcpserver

import (
	"net"
	"fmt"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
)

//TcpServerCreate xxx
func TcpServerCreate(s *utils.ServerInfo){
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	addr := fmt.Sprintf("%s:%d",s.Host,s.Port)
	fmt.Println("Listening....")
	l, err := net.Listen(s.Proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		fmt.Println("Accept....")
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("come a new connection from %s\n", conn.RemoteAddr().String())
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleNewFunc(conn)
	}
}


func handleNewFunc(c net.Conn){
	defer c.Close()
	var buff = make([]byte,1024)
	for{
		n, err := c.Read(buff)
		if err != nil{
			log.Println(err)
			return
		}
		fmt.Println(n ,string(buff[:n]))
		//n ,err = c.Write(buff)
		//if err != nil{
		//	log.Println(err)
		//	return
		//}
		//fmt.Printf("write %d data\n", n)
	}
}


func handleFunc(c net.Conn){
	defer c.Close()
	time.Sleep(10*time.Second)
	var buff = make([]byte,1024)
	n, err := c.Read(buff)
	//content, err := bufio.NewReader(c).ReadString('\n')
	if err != nil{
		log.Println(err)
		return
	}
	fmt.Println(string(buff))

	n ,err = c.Write(buff)
	if err != nil{
		log.Println(err)
		return
	}
	fmt.Printf("write %d data\n", n)
}

