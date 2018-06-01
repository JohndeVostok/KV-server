package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	OPERR = 0
	OPPUT = 1
	OPGET = 2
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

type KVServer struct {
	mu   sync.Mutex
	data map[string]string
}

var kv *KVServer

func (kv *KVServer) Put(args RequestArgs) *ResponseArgs {
	kv.mu.Lock()
	kv.data[args.Key] = args.Value
	kv.mu.Unlock()
	var resp *ResponseArgs = &ResponseArgs{1, args.Value}
	return resp
}

func (kv *KVServer) Get(args RequestArgs) *ResponseArgs {
	kv.mu.Lock()
	v, ok := kv.data[args.Key]
	kv.mu.Unlock()
	var resp *ResponseArgs
	if ok {
		resp = &ResponseArgs{1, v}
	} else {
		resp = &ResponseArgs{0, ""}
	}
	return resp
}

func handleReq(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var args RequestArgs
	err := json.Unmarshal(body, &args)
	if err != nil {
		log.Println("Unmarshal: ", err)
	} else {
		var resp *ResponseArgs
		switch args.Op {
		case OPPUT:
			resp = kv.Put(args)
		case OPGET:
			resp = kv.Get(args)
		}
		content, _ := json.Marshal(resp)
		fmt.Fprint(w, string(content))
	}
}

func StartServer() *KVServer {
	kv := new(KVServer)
	kv.data = make(map[string]string)
	return kv
}

func main() {
	kv = StartServer()
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
