package tcpclient

import (
	"net"
	"fmt"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
	"sync"
	"mytcp/cmd/pkg/process"
)


func init(){
	log.SetFlags(log.Llongfile|log.LstdFlags)
}

func ConnectServer(s *utils.ServerInfo) {
	addr := fmt.Sprintf("%s:%d",s.Host,s.Port)
	var wg  sync.WaitGroup
	//TODO xxx
	for i :=0;i < 1000;i++{
		wg.Add(1)
		conn, err := net.DialTimeout(s.Proto, addr,time.Second*5)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer conn.Close()
		//utils.SetTimeOut(conn, 5)
		go sender(conn, &wg)
	}
	wg.Wait()
}

func sender(conn net.Conn, wg *sync.WaitGroup){
	defer wg.Done()
	for i:=0;i< 1000;i++{
		data := fmt.Sprintf("{\"ID\":%d, \"Name\":\"user-%d\"}", i,i)
		n , err := conn.Write(process.Packet([]byte(data)))
		if err != nil{
			log.Println(err.Error())
			continue
		}
		fmt.Printf("index %d, write %d data\n",i,n)
	}
}