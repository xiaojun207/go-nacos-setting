package main

import (
	"github.com/xiaojun207/go-nacos-setting/nacos"
	"log"
	"time"
)

func main() {
	nacosSetting := nacos.NacosSetting{
		AppId:         "nacos-demo",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8080,
		OnConfigLoad:  OnConfigLoad,
		Username:      "nacos",
		Password:      "nacos",
	}

	nacosSetting = *nacos.Init(nacosSetting)

	go func() {
		//for {
		time.Sleep(time.Second)
		serviceAddress, err := nacosSetting.GetServiceAddress("nacos-demo")
		if err != nil {
			log.Println("err:", err)
		}
		log.Println("instance:", serviceAddress)
		//}
	}()

	select {}
}

func OnConfigLoad(conf nacos.NacosConfig) {
	log.Println("-----OnConfigLoad----------------------------------------------------------------------------------")
	log.Println(conf.Content)
	//log.Println("conf.Content.length:", len(conf.Content))
	//printlnStrMap(conf.YAML)
}
