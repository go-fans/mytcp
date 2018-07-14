package tcpserver2

import (
	"net"
	"fmt"
	"log"
	"mytcp/cmd/pkg/utils"
	"time"
	"mytcp/cmd/pkg/process"
)

//TODO: replace fmt.Printx with log.Printx or glog @songjiang

func init(){
	log.SetFlags(log.Lshortfile| log.LstdFlags)
}

const(
	//TODO: read from env or config file @songjiang
	MaxQueueSize = 50000
	MaxWorkers = 50000
)

type Job struct{
	conn net.Conn
}

var jobQueue = make(chan Job, MaxQueueSize)
var quit = make(chan bool)

type Worker struct{
	ID int
	WorkPool chan chan Job
	JobChannel chan Job
	Quit chan bool
}

func NewWorker(workerPool chan chan Job, id int ) Worker {
	return Worker{
		ID:id,
		WorkPool:workerPool,
		JobChannel:make(chan Job),
		Quit: make(chan bool),
	}
}

func (w Worker) Start(){
	go func(){
		for {
			w.WorkPool <- w.JobChannel
			fmt.Printf("w.WorkPool <- w.JobChannel id: %d\n", w.ID)
			select{
			case job := <- w.JobChannel:
				//job.Do
				handleNewFunc(job.conn)
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w Worker)Stop(){
	go func(){
		w.Quit <- true
	}()
}

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
}

func NewDispatcher() *Dispatcher {
	pool := make(chan chan Job, MaxWorkers)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, i)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

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
	dispatcher := NewDispatcher()
	dispatcher.Run()

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
			fmt.Println("job queue full....................................................")
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
		tmpBuff = process.UnPacket(append(tmpBuff,buff[:n]...), readerChan)
	}
}


func handleFunc(c net.Conn){
	defer c.Close()
	time.Sleep(10*time.Second)
	var buff = make([]byte,1024)
	n, err := c.Read(buff)
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

