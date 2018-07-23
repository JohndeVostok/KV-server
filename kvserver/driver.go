package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"sync"
)

const (
	OPERR = 0
	OPPUT = 1
	OPGET = 2
	OPADD = 3
	OPREM = 4
)

const (
	SYNCED   = 0
	UNSYNCED = 1
)

const (
	BASE = 307
	MODS = 100000007
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

type WorkerEntry struct {
	client *rpc.Client
	name   string
	status int
}

type DriverServer struct {
	mu      sync.Mutex
	logs    []LogEntry
	workers []WorkerEntry
}

var driver *DriverServer

func hash(str string) int {
	t := 0
	p := 1
	for _, ch := range []rune(str) {
		t = (t*p + int(ch)) % MODS
		p = p * BASE % MODS
	}
	return t
}

func (driver *DriverServer) StartServer() {
	driver.logs = make([]LogEntry, 0)
	driver.workers = make([]WorkerEntry, 0)
	log.Println("Start server.")
}

func (driver *DriverServer) AddWorker(args *RequestArgs) {
	for _, worker := range driver.workers {
		if args.Key == worker.name {
			log.Println("Add worker: Duplicate worker name: ", args.Key)
			return
		}
	}
	client, err := rpc.DialHTTP("tcp", args.Value)
	if err != nil {
		log.Println("Add worker: tcp error: ", err)
		return
	}
	worker := WorkerEntry{client, args.Key, 0}
	driver.workers = append(driver.workers, worker)
	err = worker.client.Call("WorkerServer.StartServer", &RequestArgs{0, "", ""}, &ResponseArgs{0, ""})
	if err != nil {
		log.Println("Add worker: Rpc call ", err)
	} else {
		log.Println("Add worker:", args.Key, args.Value)
	}
}

func (driver *DriverServer) RemoveServer(args *RequestArgs) {
	idx := -1
	for i, worker := range driver.workers {
		if worker.name == args.Key {
			idx = i
			break
		}
	}
	if idx == -1 {
		log.Println("Remove worker: Invalid worker name: ", args.Key)
		return
	}
	driver.workers[idx].client.Close()
	driver.workers = append(driver.workers[:idx], driver.workers[idx+1:]...)
}

func (driver *DriverServer) Sync() {
	//TODO
}

func (driver *DriverServer) Put(args *RequestArgs) *ResponseArgs {
	reply := new(ResponseArgs)
	idx := hash(args.Key) % len(driver.workers)
	err := driver.workers[idx].client.Call("WorkerServer.Put", &args, &reply)
	if err != nil {
		log.Println("Put: Rpc call ", err)
		return &ResponseArgs{0, ""}
	}
	log.Println("Put: ", args.Key)
	return reply
}

func (driver *DriverServer) Get(args *RequestArgs) *ResponseArgs {
	reply := new(ResponseArgs)
	idx := hash(args.Key) % len(driver.workers)
	err := driver.workers[idx].client.Call("WorkerServer.Get", &args, &reply)
	if err != nil {
		log.Println("Get: Rpc call ", err)
		return &ResponseArgs{0, ""}
	}
	log.Println("Get: ", args.Key)
	return reply
}

func handleReq(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var args *RequestArgs
	err := json.Unmarshal(body, &args)
	if err != nil {
		log.Println("Unmarshal: ", err)
	} else {
		switch args.Op {
		case OPPUT:
			reply := driver.Put(args)
			c, _ := json.Marshal(reply)
			fmt.Fprint(w, string(c))
		case OPGET:
			reply := driver.Get(args)
			c, _ := json.Marshal(reply)
			fmt.Fprint(w, string(c))
		case OPADD:
			driver.AddWorker(args)
		default:
			log.Println("Invalid Request")
		}
	}
}

func main() {
	driver = new(DriverServer)
	driver.StartServer()
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
