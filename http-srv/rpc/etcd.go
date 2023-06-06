package rpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func EtcdRegist(serviceName, addr string, etcdEndpoints []string, ttl int64) (net.Listener, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: etcdEndpoints,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	em, err := endpoints.NewManager(client, serviceName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	lease, err := client.Grant(context.TODO(), ttl)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = em.AddEndpoint(context.TODO(), serviceName+"/"+addr, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(lease.ID))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	alive, err := client.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	go func() {
		for _ = range alive {
			log.Println("etcd client keep alive")
		}
	}()
	
	//关闭信号处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <- ch
		err = EtcdUnRegist(client, serviceName, addr)
		if err != nil {
			panic(err)
		}
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
	// 获取tcp连接
	lis, err := net.Listen("tcp", addr)
	return lis, err
}

func EtcdUnRegist(etcdClient *clientv3.Client, serviceName, addr string) error {
	log.Printf("etcdUnRegister %s\n", serviceName)
	if etcdClient != nil {
		em, err := endpoints.NewManager(etcdClient, serviceName)
		if err != nil {
			log.Println(err)
			return err
		}
		err = em.DeleteEndpoint(context.TODO(), serviceName+"/"+addr)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	return nil
}