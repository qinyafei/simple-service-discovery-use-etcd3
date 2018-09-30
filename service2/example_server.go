package main

import (
	"fmt"
	//"os"
	"time"

	"github.com/coreos/etcd/clientv3/MyTest/discovery"
)

func main() {
	//server name
	regname := "service2" //os.Args[1]

	reg := &discovery.Registry{
		Ip:   "192.168.14.40",
		Port: 8091,
		Name: regname,
	}

	agent := discovery.Regist(reg)
	fmt.Println("register information:", agent.Reg.Ip, agent.Reg.Port, agent.Reg.Name)

	for {
		time.Sleep(5 * time.Second)
	}
}
