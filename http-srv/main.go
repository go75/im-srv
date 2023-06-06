package main

import (
	"im-srv/http-srv/initialize"
)

func main() {
	initialize.Config()
	initialize.DB()
	initialize.RPC()
}