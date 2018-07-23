package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
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

func handleReq(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var args *RequestArgs
	err := json.Unmarshal(body, &args)
	if err != nil {
		log.Println("Unmarshal: ", err)
	} else {
		switch args.Op {
		case OPPUT:
			data := []byte("SET " + args.Key + " " + args.Value)
			ioutil.WriteFile("buf.bin", data, 0644)
			cmd := exec.Command("sh", "-c", "redis-cli < buf.bin")
			c, err := cmd.Output()
			if err != nil {
				log.Println("Redis SET: ", err)
			}
			fmt.Fprint(w, string(c))
		case OPGET:
			data := []byte("GET " + args.Key)
			ioutil.WriteFile("buf.bin", data, 0644)
			cmd := exec.Command("sh", "-c", "redis-cli < buf.bin")
			c, err := cmd.Output()
			if err != nil {
				log.Println("Redis SET: ", err)
			}
			fmt.Fprint(w, string(c))
		default:
			log.Println("Invalid Request")
		}
	}
}

func main() {
	http.HandleFunc("/", handleReq)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
