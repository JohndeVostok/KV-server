package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type RequestArgs struct {
	Op    int
	Key   string
	Value string
}

type ResponseArgs struct {
	Status int
	Value  string
}

type LogEntry struct {
	Key   string
	Value string
	Id    int
}

type WorkerServer struct {
	mu   sync.Mutex
	logs []LogEntry
	data map[string]string
}

var worker *WorkerServer

func (worker *WorkerServer) StartServer(args *RequestArgs, reply *ResponseArgs) error {
	worker.logs = make([]LogEntry, 0)
	worker.data = make(map[string]string)
	log.Println("StartServer")
	return nil
}

func (worker *WorkerServer) Put(args *RequestArgs, reply *ResponseArgs) error {
	worker.mu.Lock()
	worker.data[args.Key] = args.Value
	worker.mu.Unlock()
	reply.Status = 1
	reply.Value = ""
	log.Println("Put: ", args.Key)
	return nil
}

func (worker *WorkerServer) Get(args *RequestArgs, reply *ResponseArgs) error {
	worker.mu.Lock()
	v, ok := worker.data[args.Key]
	worker.mu.Unlock()
	if ok {
		reply.Status = 1
		reply.Value = v
	} else {
		reply.Status = 0
		reply.Value = ""
	}
	log.Println("Get: ", args.Key)
	return nil
}

func main() {
	worker = new(WorkerServer)
	rpc.Register(worker)
	rpc.HandleHTTP()
	var port string
	l := len(os.Args)
	port = os.Args[l-1]
	port = ":" + port
	log.Println(port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("listen: ", err)
	}
	go http.Serve(listener, nil)
	for {
		time.Sleep(time.Second)
	}
}
