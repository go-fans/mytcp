package tcpserver

import (
	"net"
	"fmt"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
	"encoding/binary"
	"bytes"
)

const(
	MaxQueueSize = 5
)

type Job struct{
	conn net.Conn
}

var jobQueue = make(chan Job, MaxQueueSize)

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
	go func(){
		for{
			select{
			case job := <- jobQueue:
				go handleNewFunc(job.conn)	//do job
			}
		}
	}()

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
		//go handleNewFunc(conn)
		if len(jobQueue) == MaxQueueSize{
			fmt.Println("job queue full....")
		}
		jobQueue <- Job{
			conn:conn,
		}
	}
}



func reader(readerCh chan []byte){
	for {
		select{
		case msg :=<- readerCh:
			fmt.Println(string(msg))
		}
	}
}

//bytesToInt get msg length
func bytesToInt(b []byte)int{
	byteBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(byteBuffer, binary.BigEndian, &x)
	return int(x)
}

func unPacket(buffer []byte, readerCh chan []byte)[]byte{
	length := len(buffer)
	var i int
	for i = 0;i< length;i++{
		if length < i + 4{
			break
		}
		msgLength := bytesToInt(buffer[i:i+4])

		if length < i + 4 + msgLength{
			break
		}
		data := buffer[i+4:i+4+msgLength]
		readerCh <- data
		i += msgLength + 4 - 1
	}
	if i == length{
		return make([]byte,0)
	}
	return buffer[i:]
}

func handleNewFunc(c net.Conn){
	var(
	 buff = make([]byte,256)
	 tmpBuff = make([]byte,0)
	)
	defer c.Close()
	readerChan := make(chan []byte,20)
	go reader(readerChan)

	for{
		n, err := c.Read(buff)	//?
		if err != nil{
			log.Println(err)
			return
		}
		tmpBuff = unPacket(append(tmpBuff,buff[:n]...), readerChan)
		//fmt.Println(n ,string(buff[:n]))
		//n ,err = c.Write(buff[:n])
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

