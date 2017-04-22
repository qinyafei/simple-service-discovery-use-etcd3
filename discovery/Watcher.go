package discovery

import (
	"encoding/json"
	//"flag"
	//"fmt"
	"log"
	//"time"

	"github.com/coreos/etcd/clientv3"
	mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
)

type Master struct {
	Members  map[string]*Servers
	Watcher  *clientv3.Client
	WatchKey string //"servers/"
}

type Servers struct {
	Ip     string
	Port   int
	Active bool
	//name和key必须一样
	Name string
}

func NewMaster(watchKey string, endpoints []string) *Master {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	mas := &Master{
		Members:  make(map[string]*Servers),
		Watcher:  cli,
		WatchKey: watchKey,
	}

	//go mas.WatcherServers()

	return mas
}

func ValueToServer(value []byte) *Servers {
	//log.Println(value)
	server := &Servers{}
	err := json.Unmarshal(value, server)
	if err != nil {
		log.Println(err)
	}

	return server
}

//----------------------------

func (mas *Master) Add(server *Servers) {
	mas.Members[server.Name] = server
}

func (mas *Master) Update(name string, alive bool) {
	if member, ok := mas.Members[name]; ok {
		member.Active = alive
	}
}

func (mas *Master) WatcherServers() {
	log.Println("Master|WatcherServers key:", mas.WatchKey)
	clientList := mas.Watcher.Watch(context.TODO(), mas.WatchKey, clientv3.WithPrefix())
	//for {
	for cli := range clientList {
		//cli := <-clientList
		for _, ev := range cli.Events {
			if ev.Type == mvccpb.PUT {
				server := ValueToServer(ev.Kv.Value)
				if _, ok := mas.Members[server.Name]; ok {
					log.Println("update server:", server.Name)
					mas.Update(server.Name, true)
				} else {
					log.Println("add server:", server.Name)
					mas.Add(server)
				}
			} else if ev.Type == mvccpb.DELETE {
				log.Println("delete server:", string(ev.Kv.Key))
				delete(mas.Members, string(ev.Kv.Key))
			} else {
				log.Println("unkown action")
			}
		}
	}
}

func (mas *Master) GetWatchers() map[string]*Servers {
	return mas.Members
}

func Watch(watcher string) *Master {
	master := NewMaster(watcher, endpoints)
	go master.WatcherServers()

	return master
}
