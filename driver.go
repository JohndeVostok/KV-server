package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
}

type KVServer struct {
	mu   sync.mutex
	data map[string]string
}

func (kv *KVServer) Put(args RequestArgs) {

}

func (kv *KVServer) Get(args RequestArgs) {

}

func handleReq(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var args Args
	err := json.Unmarshal(body, &args)
	if err != nil {
		fmt.Println("Unmarshal: ", err)
	} else {
		if args.Op == OPPUT {
			kv.Put(args)
		}
		if args.Op == OPGET {
			kv.Get(args)
		}
		ret, _ := json.Marshal(args)
		fmt.Fprint(w, string(ret))
	}
}

func StartServer() *KVServer {
	kv := new(KVServer)
	kv.data = make(map[string]string)
	return kv
}

func main() {
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
