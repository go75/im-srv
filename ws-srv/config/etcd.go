package config

type Etcd struct {
	ServiceName string
	Endpoints []string
	TTL int64
}