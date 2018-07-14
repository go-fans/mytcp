package tcpclient

import (
	"net"
	"fmt"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
	"sync"
	"bytes"
	"encoding/binary"
)

func ConnectServer(s *utils.ServerInfo) {
	//data := "this is client...\n"
	addr := fmt.Sprintf("%s:%d",s.Host,s.Port)
	var wg  sync.WaitGroup

	for i :=0;i < 10;i++{
		wg.Add(1)
		conn, err := net.DialTimeout(s.Proto, addr,time.Second*2)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer conn.Close()
		utils.SetTimeOut(conn, 2)
		go sender(conn, &wg)
	}
	wg.Wait()
}

func intToByte(n int)[]byte{
	x := int32(n)
	byteBuffer := bytes.NewBuffer([]byte{})
	binary.Write(byteBuffer, binary.BigEndian,x)
	return byteBuffer.Bytes()
}

func Packet(msg []byte)[]byte{
	s := make([]byte,0)
	return append(append(s, intToByte(len(msg))...), msg...)
}

func sender(conn net.Conn, wg *sync.WaitGroup){
	defer wg.Done()
	for i:=0;i< 300;i++{
		data := fmt.Sprintf("{\"ID\":%d, \"Name\":\"user-%d\"}", i,i)

		n , err := conn.Write(Packet([]byte(data)))	// ?
		if err != nil{
			log.Println(err.Error())
			continue
		}
		fmt.Printf("index %d, write %d data\n",i,n)
		////
		//buf := make([]byte, n)
		//n , err = conn.Read(buf)
		////content, err := bufio.NewReader(conn).ReadString('\n')
		//if err != nil{
		//	log.Println(err.Error())
		//	continue
		//}
		//fmt.Println(n,string(buf))
	}
}