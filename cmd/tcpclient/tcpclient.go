package tcpclient

import (
	"net"
	"fmt"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
	"sync"
)

func ConnectServer(s *utils.ServerInfo) {
	//data := "this is client...\n"
	addr := fmt.Sprintf("%s:%d",s.Host,s.Port)
	conn, err := net.DialTimeout(s.Proto, addr,time.Second*2)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	utils.SetTimeOut(conn, 2)
	var wg  sync.WaitGroup
	wg.Add(1)
	go sender(conn, &wg)
	wg.Wait()

}

func sender(conn net.Conn, wg *sync.WaitGroup){
	defer wg.Done()
	for i:=0;i< 300;i++{
		data := fmt.Sprintf("{\"ID\":%d, \"Name\":\"user-%d\"}", i,i)
		n , err := conn.Write([]byte(data))
		if err != nil{
			log.Println(err.Error())
			return
		}
		fmt.Printf("write %d data\n",n)
		//
		//content, err := bufio.NewReader(conn).ReadString('\n')
		//if err != nil{
		//	log.Fatal(err)
		//}
		//fmt.Println(content)
	}
}