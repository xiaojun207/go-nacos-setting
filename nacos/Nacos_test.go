package nacos

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	nacosSetting := NacosSetting{
		AppId:         "nacos-demo",
		NacosServerIp: "127.0.0.1",
		ClientPort:    8080,
		ShowLog:       true,
		ConfigType:    "YAML",
		OnConfigLoad:  OnConfigLoad,
	}

	Init(nacosSetting)

	select {}
}

func OnConfigLoad(conf map[string]interface{}) {
	log.Println("-----OnConfigLoad----------------------------------------------------------------------------------")
	log.Println(len(conf))
	printlnStrMap(conf)
}

func printlnStrMap(conf map[string]interface{}) {
	root := ""
	for key, value := range conf {
		if isMap(value) {
			printMap(root, key, value.(map[interface{}]interface{}))
		} else if isArr(value) {
			printArr(root, key, value.([]interface{}))
		} else {
			fmt.Println(root+key+":", value)
		}
	}
}

func printMap(root, rootkey string, conf map[interface{}]interface{}) {
	fmt.Println(root + rootkey + ":")
	root = root + "  "
	for key, value := range conf {
		if isMap(value) {
			printMap(root, key.(string), value.(map[interface{}]interface{}))
		} else if isArr(value) {
			printArr(root, key.(string), value.([]interface{}))
		} else {
			fmt.Println(root+key.(string)+":", value)
		}
	}
}

func printArr(root, rootkey string, conf []interface{}) {
	fmt.Println(root + rootkey + ":")
	root = root + "  "
	var key interface{} = "-"
	for _, value := range conf {
		if isMap(value) {
			printMap(root, key.(string), value.(map[interface{}]interface{}))
		} else if isArr(value) {
			printArr(root, key.(string), value.([]interface{}))
		} else {
			fmt.Println(root+key.(string)+"", value)
		}
	}
}

func isMap(value interface{}) bool {
	var t map[interface{}]interface{}
	return reflect.TypeOf(value) == reflect.TypeOf(t)
}

func isArr(value interface{}) bool {
	var t []interface{}
	return reflect.TypeOf(value) == reflect.TypeOf(t)
}
