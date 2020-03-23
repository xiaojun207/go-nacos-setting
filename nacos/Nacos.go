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

	ShowLog bool // 【选填】，默认：false，因为nacos go sdk设置了log输出到日志文件，不会显示到控制台。当ShowLog=true，日志会显示到控制台
}

type NacosConfig struct {
	DataId     string `param:"dataId"`
	Group      string `param:"group"`
	Content    string `param:"content"`
	Namespace  string `param:"namespace"`
	ConfigType string `param:"configType"`
	Properties map[string]string
	JSON       map[string]interface{}
	YAML       map[string]interface{}
	OnChange   func(namespace, group, dataId, data string)
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

func Init(nacosSetting NacosSetting) {

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

			nacosConfig := NacosConfig{
				ConfigType: nacosSetting.ConfigType,
				Content:    data,
			}
			// JSON、YAML、Properties
			if nacosSetting.ConfigType == "Properties" {
				nacosConfig.Properties = Properties(data)
			} else if nacosSetting.ConfigType == "YAML" {
				nacosConfig.YAML = Yaml(data)
			} else if nacosSetting.ConfigType == "JSON" {
				nacosConfig.JSON = make(map[string]interface{})
				utils.JsonToMap(data, nacosConfig.JSON)
			} else {
				nacosConfig.Properties = Properties(data)
			}

			nacosSetting.OnConfigLoad(nacosConfig)
		},
	})
	if nacosSetting.ShowLog {
		log.SetOutput(os.Stdout)
	}
}

func (e *NacosConfig) GetValue(key string) string {
	return e.Properties[key]
}

func (e *NacosConfig) GetString(key, defalueValue string) string {
	value := e.GetValue(key)
	if value == "" {
		value = defalueValue
	}
	return value
}

func (e *NacosConfig) GetFloat64(key string, defalueValue float64) float64 {
	str := e.GetValue(key)
	return utils.StrToFloat64Def(str, defalueValue)
}

func (e *NacosConfig) GetBool(key string, defalueValue bool) bool {
	str := e.GetValue(key)
	return utils.StrToBoolDef(str, defalueValue)
}

func (e *NacosConfig) GetInt(key string, defalueValue int) int {
	str := e.GetValue(key)
	return utils.StrToIntDef(str, defalueValue)
}

func (e *NacosConfig) GetInt64(key string, defalueValue int64) int64 {
	str := e.GetValue(key)
	return utils.StrToInt64Def(str, defalueValue)
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

func GetYamlValue(conf map[string]interface{}, key, defalueValue string) string {
	value := conf[key]
	if value == "" {
		value = defalueValue
	}
	return value.(string)
}
