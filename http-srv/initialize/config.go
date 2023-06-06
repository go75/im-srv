package initialize

import (
	"encoding/json"
	"im-srv/http-srv/global"
	"os"
)

func Config() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, global.Config)
	if err != nil {
		panic(err)
	}
}