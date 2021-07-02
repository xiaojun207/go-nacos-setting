package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/xiaojun207/go-base-utils/utils"
	yaml "gopkg.in/yaml.v2"
	"log"
	"strings"
)

var (
	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient
)

type NacosSetting struct {
	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient

	AppId           string // 【必填】，例如：bj-yun-nacos-demo
	NacosServerIp   string // 【选填】，默认: 127.0.0。1
	NacosServerPort uint64 // 【选填】，默认: 8848
	ClientIp        string // 【选填】，默认：获取本机IP，可以自己设定
	ClientPort      uint64 // 【选填】，默认：80

	ServiceName string // 【选填】，默认：{AppId}
	ClusterName string // 【选填】，默认：default

	ConfigDataId string                 // 【选填】，默认：{AppId}
	ConfigGroup  string                 // 【选填】，默认：DEFAULT_GROUP
	ConfigType   string                 // 【选填】，默认：Properties，支持：JSON、YAML、Properties，所有的配置均以map[string]interface{}回调
	OnConfigLoad func(conf NacosConfig) // 【选填】，配置更新回调

	ShowLog  bool   // 【选填】，默认：false，因为nacos go sdk设置了log输出到日志文件，不会显示到控制台。当ShowLog=true，日志会显示到控制台
	AESKey   string // 【选填】，默认："", 当不为空时，会检测配置的值，如果是AESEncrypt()，包括起来的，则尝试解密
	DESKey   string // 【选填】，默认："", 当不为空时，会检测配置的值，如果是DESEncrypt()，包括起来的，则尝试解密
	Username string
	Password string

	Metadata map[string]string
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
		nacosSetting.ClusterName = "DEFAULT"
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

	if nacosSetting.Metadata == nil {
		nacosSetting.Metadata = map[string]string{"secure": "false", "preserved.register.source": "GO_NACOS_SETTING"}
	}

	if nacosSetting.OnConfigLoad == nil {
		nacosSetting.OnConfigLoad = func(conf NacosConfig) {
			log.Println("There is no OnConfigLoad function!")
		}
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

func Init(nacosSetting NacosSetting) *NacosSetting {

	nacosSetting = setDefaultSetting(nacosSetting)

	// 可以没有，采用默认值
	clientConfig := constant.ClientConfig{
		TimeoutMs:      10 * 1000,
		ListenInterval: 5 * 1000,
		BeatInterval:   5 * 1000,
		LogDir:         ".nacos/logs",
		CacheDir:       ".nacos/cache",
		SecretKey:      "",
		Username:       nacosSetting.Username,
		Password:       nacosSetting.Password,
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
		log.Println("CreateNamingClient.err", err)
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	if err != nil {
		log.Println("CreateConfigClient.err", err)
	}

	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          nacosSetting.ClientIp,
		Port:        nacosSetting.ClientPort,
		ServiceName: nacosSetting.ServiceName,
		Weight:      10,
		ClusterName: nacosSetting.ClusterName,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    nacosSetting.Metadata,
	})
	log.Println("RegisterInstance:", success)
	if err != nil {
		log.Println("RegisterInstance.err:", err)
	}

	configParam := vo.ConfigParam{
		DataId: nacosSetting.ConfigDataId,
		Group:  nacosSetting.ConfigGroup,
		OnChange: func(namespace, group, dataId, data string) {

			nacosConfig := NacosConfig{
				ConfigType: nacosSetting.ConfigType,
				Namespace:  namespace,
				Group:      group,
				DataId:     dataId,
				DESKey:     nacosSetting.DESKey,
				AESKey:     nacosSetting.AESKey,
			}
			nacosConfig.LoadData(data)

			nacosSetting.OnConfigLoad(nacosConfig)
		},
	}
	configClient.ListenConfig(configParam)
	log.Println("configClient.ListenConfig")

	context, err := configClient.GetConfig(configParam)
	if err != nil {
		log.Println("GetConfig.err", err)
	} else {
		configParam.OnChange("", configParam.Group, configParam.DataId, context)
		log.Println("GetConfig:", context)
	}

	nacosSetting.NamingClient = namingClient
	nacosSetting.ConfigClient = configClient

	NamingClient = namingClient
	ConfigClient = configClient
	return &nacosSetting
}

// Dep
func GetInstance(serviceName, clusterName string) (*model.Instance, error) {
	instance, err := NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		Clusters:    []string{clusterName},
	})
	return instance, err
}

func (e *NacosSetting) GetInstance(serviceName string) (*model.Instance, error) {
	return e.NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		Clusters:    []string{e.ClusterName},
	})
}

func (e *NacosSetting) GetServiceAddress(appId string) (address string, err error) {
	instance, err := e.GetInstance(appId)
	if err != nil {
		return
	}
	address = instance.Ip + ":" + utils.Uint64ToStr(instance.Port)
	return
}

/**
  Properties文本转map，支持注释
*/
func Properties(data string) map[string]string {
	var resultMap = make(map[string]string)
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, "=") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
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

func GetYamlValue(conf map[string]interface{}, key, defalueValue string) string {
	value := conf[key]
	if value == "" {
		value = defalueValue
	}
	return value.(string)
}
