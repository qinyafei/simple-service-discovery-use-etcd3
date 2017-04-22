package main

import (
	"fmt"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3/MyTest/discovery"
)

func main() {
	//server name
	regname := os.Args[1]

	reg := &discovery.Registry{
		Ip:   "192.168.8.251",
		Port: 8090,
		Name: regname,
	}

	agentwatch := discovery.Watch(regname)
	time.Sleep(1 * time.Second)

	agent := discovery.Regist(reg)
	fmt.Println("register information:", agent.Reg.Ip, agent.Reg.Port, agent.Reg.Name)

	for {
		time.Sleep(5 * time.Second)
		servers := make(map[string]*discovery.Servers)
		servers = agentwatch.GetWatchers()
		fmt.Println(servers[regname])
	}
}
