package discovery

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	//mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
)

var (
	//dialTimeout    = 5 * time.Second
	//requestTimeout = 10 * time.Second
	endpoints = []string{"localhost:2379", "localhost:22379", "localhost:32379"}

//endpoints = []string{"localhost:2379"}
)

type Registry struct {
	Ip   string
	Port int
	Name string
}

type ServerAgent struct {
	Reg   Registry
	Agent *clientv3.Client
}

func NewServerAgent(reg *Registry, endpoints []string) *ServerAgent {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Print(err)
	}

	//serverName := "servers/" + reg.Name
	//reg.Name = serverName

	agent := &ServerAgent{
		Reg:   *reg,
		Agent: cli,
	}
	return agent
}

func (server *ServerAgent) HeartBeat() {
	fmt.Println("ServerAgent|HeartBeat key:", server.Reg.Name)
	for {
		//
		expire, err := server.Agent.Grant(context.TODO(), 15)
		if err != nil {
			log.Println(err)
		}

		key := server.Reg.Name
		value, _ := json.Marshal(server.Reg)
		_, err2 := server.Agent.Put(context.TODO(), key, string(value), clientv3.WithLease(expire.ID))
		//_, err2 := server.Agent.Put(context.TODO(), key, string(value))
		if err2 != nil {
			log.Println(err2)
		}

		time.Sleep(4 * time.Second)
	}

}

func Regist(reg *Registry) *ServerAgent {
	server := NewServerAgent(reg, endpoints)
	go server.HeartBeat()

	fmt.Println("register server:", reg.Name)
	return server
}
