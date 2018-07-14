package tcpserver

import (
	"net"
	"io"
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
		time.Sleep(10*time.Second)
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}