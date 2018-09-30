package main

import (
	"fmt"
	//"os"
	"time"

	"github.com/coreos/etcd/clientv3/MyTest/discovery"
)

func main() {
	//server name
	regname := "service" //os.Args[1]

	agentwatch := discovery.Watch(regname)
	for {
		time.Sleep(3 * time.Second)
		servers := agentwatch.GetWatchers()
		fmt.Println(servers)
	}
}
