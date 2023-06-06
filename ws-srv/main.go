package main

import "im-srv/ws-srv/initialize"

func main() {
	initialize.Config()
	initialize.DB()
	initialize.RPC()
}