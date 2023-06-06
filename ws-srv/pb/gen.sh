# grpc generate
protoc --go_out=./ --go-grpc_out=./ --go-grpc_opt=require_unimplemented_servers=false ws.proto
ls ws.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'