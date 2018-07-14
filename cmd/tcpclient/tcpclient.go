package tcpclient

import (
	"net"
	"fmt"
	"bufio"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
)

func ConnectServer(s *utils.ServerInfo) {
	addr := fmt.Sprintf("%s:%d",s.Host,s.Port)
	conn, err := net.DialTimeout(s.Proto, addr,time.Second*2)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	utils.SetTimeOut(conn, 2)
	fmt.Fprintf(conn, "this is client...\n")

	content, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(content)
}