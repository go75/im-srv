package initialize

import (
	"context"
	"im-srv/http-srv/global"
	"im-srv/http-srv/pb"
	"im-srv/http-srv/rpc"
	"log"

	"google.golang.org/grpc"
)

func RPC() {
	etcdConfig := global.Config.Etcd
	serverConfig := global.Config.Server
	lis, err := rpc.EtcdRegist(etcdConfig.ServiceName, serverConfig.Addr, etcdConfig.Endpoints, etcdConfig.TTL)
	if err != nil {
		panic(err)
	}
	// 创建grpc的server
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryInterceptor))
	// 注册grpc的服务
	pb.RegisterHTTPServer(server, rpc.HTTPServer{})
	log.Println("service start addr: ", serverConfig.Addr)
	// 开启grpc服务
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

// 自定义拦截器
func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("call %s\n", info.FullMethod)
	resp, err = handler(ctx, req)
	return resp, err
}
