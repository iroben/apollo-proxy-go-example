package main

import (
	"apollo-proxy-go-example/config"
	"github.com/philchia/agollo/v4"
	"log"
)

func main() {
	namespaces := []string{
		"application",
		"apollo-proxy-public",
	}

	configInstance := config.NewConfig(agollo.Conf{
		AppID:          "apollo-proxy",
		Cluster:        "default",
		NameSpaceNames: namespaces,
	}, ConfigUpdate)
	configInstance.LoadConfig()
	configInstance.Json()

	// 本地观察配置更新，触发 ConfigUpdate
	// time.Sleep(time.Minute * 5)
}
func ConfigUpdate() {
	log.Println("配置更新了，服务类的实例要重新创建")
}
