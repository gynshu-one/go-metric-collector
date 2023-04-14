package main

import (
	ag "github.com/gynshu-one/go-metric-collector/internal/controller/http/agent"
	"github.com/gynshu-one/go-metric-collector/internal/domain/service"
	"sync"
)

var (
	agent   ag.Handler
	storage service.MemStorage
)

func main() {
	agent = ag.NewAgent(service.NewMemService(&sync.Map{}))
	agent.Start()
	//
	//var data string
	//
	//data = fmt.Sprintf("%s:%s:%d", "GetSetZip170", "counter", 2287779143)
	//h := hmac.New(sha256.New, []byte(""))
	//h.Write([]byte(data))
	//fmt.Printf("%s", hex.EncodeToString(h.Sum(nil)))

}
