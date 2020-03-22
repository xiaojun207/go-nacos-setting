package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/xiaojun207/go-base-utils/utils"
	yaml "gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

type NacosSetting struct {
	AppId           string
	NacosServerIp   string // 默认: 127.0.0。1
	NacosServerPort uint64 // 默认: 8848
	ClientIp        string // 默认获取本机IP
	ClientPort      uint64 // 默认80

	ServiceName string
	ClusterName string

	ConfigDataId string
	ConfigGroup  string
	ConfigType   string // JSON、YAML、Properties
	ShowLog      bool
}

func setDefaultSetting(nacosSetting NacosSetting) NacosSetting {
	if nacosSetting.AppId == "" {
		fmt.Println("AppId is empty, use default 'nacos-demo' ")
		nacosSetting.AppId = "nacos-demo"
	}

	if nacosSetting.NacosServerIp == "" {
		nacosSetting.NacosServerIp = "127.0.0.1"
	}

	if nacosSetting.NacosServerPort == 0 {
		nacosSetting.NacosServerPort = 8848
	}

	if nacosSetting.ClientPort == 0 {
		nacosSetting.ClientPort = 80 // 默认80
	}

	if nacosSetting.ServiceName == "" {
		nacosSetting.ServiceName = nacosSetting.AppId
	}

	if nacosSetting.ClusterName == "" {
		nacosSetting.ClusterName = "default"
	}

	if nacosSetting.ConfigDataId == "" {
		nacosSetting.ConfigDataId = nacosSetting.AppId
	}

	if nacosSetting.ConfigGroup == "" {
		nacosSetting.ConfigGroup = "DEFAULT_GROUP"
	}

	if nacosSetting.ConfigType == "" {
		nacosSetting.ConfigType = "Properties"
	}

	if nacosSetting.ClientIp == "" {
		ip, err := utils.ExternalIP()
		if err != nil {
			fmt.Println(err)
		}
		nacosSetting.ClientIp = ip.String()
	}

	return nacosSetting
}

func Init(nacosSetting NacosSetting, OnConfigLoad func(conf map[string]interface{})) {

	nacosSetting = setDefaultSetting(nacosSetting)

	// 可以没有，采用默认值
	clientConfig := constant.ClientConfig{
		TimeoutMs:      10 * 1000,
		ListenInterval: 30 * 1000,
		BeatInterval:   5 * 1000,
		LogDir:         "nacos/logs",
		CacheDir:       "nacos/cache",
		SecretKey:      "",
	}

	// 至少一个
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacosSetting.NacosServerIp,
			ContextPath: "/nacos",
			Port:        nacosSetting.NacosServerPort,
		},
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		fmt.Println(err)
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	if err != nil {
		fmt.Println(err)
	}

	success, _ := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          nacosSetting.ClientIp,
		Port:        nacosSetting.ClientPort,
		ServiceName: nacosSetting.ServiceName,
		Weight:      10,
		ClusterName: nacosSetting.ClusterName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	log.Println("namingClient:", success)

	_, err = configClient.GetConfig(vo.ConfigParam{
		DataId: nacosSetting.ConfigDataId,
		Group:  nacosSetting.ConfigGroup})

	if err != nil {
		log.Println("configClient:", err)
	}
	//fmt.Println("configClient:", content)

	configClient.ListenConfig(vo.ConfigParam{
		DataId: nacosSetting.ConfigDataId,
		Group:  nacosSetting.ConfigGroup,
		OnChange: func(namespace, group, dataId, data string) {
			// JSON、YAML、Properties
			conf := make(map[string]interface{})
			if nacosSetting.ConfigType == "Properties" {
				conf = Properties(data)
			} else if nacosSetting.ConfigType == "Yaml" {
				conf = Yaml(data)
			} else if nacosSetting.ConfigType == "JSON" {
				utils.JsonToMap(data, conf)
			} else {
				conf = Properties(data)
			}
			OnConfigLoad(conf)
		},
	})
	if nacosSetting.ShowLog {
		log.SetOutput(os.Stdout)
	}
}

/**
  Properties文本转map，支持注释
*/
func Properties(data string) map[string]interface{} {
	var resultMap = make(map[string]interface{})
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		kv := strings.Split(line, "=")
		key := strings.TrimSpace(kv[0])
		value := kv[1]
		idx := strings.Index(value, "#")
		if idx >= 0 {
			value = strings.Split(value, "#")[0]
		}
		value = strings.TrimSpace(value)

		resultMap[key] = value
	}
	return resultMap
}

func Yaml(data string) map[string]interface{} {
	resultMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return resultMap
}
