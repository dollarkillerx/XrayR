package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/XrayR-project/XrayR/panel"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use: "XrayR",
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file for XrayR.")
}

// generateDefaultConfig 生成默认配置文件
func generateDefaultConfig(configPath string) error {
	// 从环境变量获取配置，如果没有则使用默认值
	apiHost := getEnv("XRAYR_API_HOST", "https://your-v2board-domain.com")
	apiKey := getEnv("XRAYR_API_KEY", "your-api-key-here")
	nodeID := getEnv("XRAYR_NODE_ID", "1")
	panelType := getEnv("XRAYR_PANEL_TYPE", "V2board")
	nodeType := getEnv("XRAYR_NODE_TYPE", "V2ray")
	listenIP := getEnv("XRAYR_LISTEN_IP", "0.0.0.0")
	sendIP := getEnv("XRAYR_SEND_IP", "0.0.0.0")
	logLevel := getEnv("XRAYR_LOG_LEVEL", "info")

	// 转换nodeID为整数
	nodeIDInt, err := strconv.Atoi(nodeID)
	if err != nil {
		nodeIDInt = 1
	}

	defaultConfig := ""

	vlessConfig := fmt.Sprintf(`Log:
  Level: %s

ConnectionConfig:
  Handshake: 4
  ConnIdle: 300
  UplinkOnly: 2
  DownlinkOnly: 4
  BufferSize: 64

Nodes:
  - PanelType: "%s" # 面板类型，如 V2board、SSpanel 等
    ApiConfig:
      ApiHost: "%s" # 面板 API 地址
      ApiKey: "%s"  # 面板 API 密钥
      NodeID: %d    # 节点 ID
      NodeType: Vless # 节点类型
      Timeout: 30
      EnableVless: true
      VlessFlow: "xtls-rprx-vision"
      SpeedLimit: 0
      DeviceLimit: 0
    ControllerConfig:
      ListenIP: %s
      SendIP: %s
      UpdatePeriodic: 60
      EnableDNS: false
      DNSType: AsIs
      EnableProxyProtocol: false
      EnableREALITY: false
      CertConfig:
        CertMode: none
`, logLevel, panelType, apiHost, apiKey, nodeIDInt, listenIP, sendIP)

	vmessConfig := fmt.Sprintf(`Log:
  Level: %s

ConnectionConfig:
  Handshake: 4
  ConnIdle: 300
  UplinkOnly: 2
  DownlinkOnly: 4
  BufferSize: 64

Nodes:
  - PanelType: "%s"
    ApiConfig:
      ApiHost: "%s"
      ApiKey: "%s"
      NodeID: %d
      NodeType: Vmess
      Timeout: 30
      SpeedLimit: 0
      DeviceLimit: 0
    ControllerConfig:
      ListenIP: %s
      SendIP: %s
      UpdatePeriodic: 60
      EnableDNS: false
      DNSType: AsIs
      EnableProxyProtocol: false
      CertConfig:
        CertMode: none
`, logLevel, panelType, apiHost, apiKey, nodeIDInt, listenIP, sendIP)

	trojanConfig := fmt.Sprintf(`Log:
	Level: %s
  
  ConnectionConfig:
	Handshake: 4
	ConnIdle: 300
	UplinkOnly: 2
	DownlinkOnly: 4
	BufferSize: 64
  
  Nodes:
	- PanelType: "%s"
	  ApiConfig:
		ApiHost: "%s"
		ApiKey: "%s"
		NodeID: %d
		NodeType: Trojan
	  ControllerConfig:
		ListenIP: %s
		SendIP: %s
		UpdatePeriodic: 60
	  CertConfig:
        Enable: true
        CertMode: file
        CertDomain: okr.hacksnews.top
        CertFile: /etc/xrayr-ssl/cert.crt
        KeyFile: /etc/xrayr-ssl/private.key
  `, logLevel, panelType, apiHost, apiKey, nodeIDInt, listenIP, sendIP)

	// TODO: 检查 /etc/xrayr-ssl/cert.crt &  /etc/xrayr-ssl/private.key 是否存在，不存在则生成
	// 注意不存在目录则创建
	if _, err := os.Stat("/etc/xrayr-ssl"); os.IsNotExist(err) {
		if err := os.MkdirAll("/etc/xrayr-ssl", 0755); err != nil {
			log.Warnf("创建证书目录失败: %v", err)
		}
	}
	if _, err := os.Stat("/etc/xrayr-ssl/cert.crt"); os.IsNotExist(err) {
		if err := os.WriteFile("/etc/xrayr-ssl/cert.crt", []byte(cert), 0644); err != nil {
			log.Warnf("写入证书文件失败: %v", err)
		}
	}
	if _, err := os.Stat("/etc/xrayr-ssl/private.key"); os.IsNotExist(err) {
		if err := os.WriteFile("/etc/xrayr-ssl/private.key", []byte(privateKey), 0644); err != nil {
			log.Warnf("写入私钥文件失败: %v", err)
		}
	}

	ssConfig := fmt.Sprintf(`Log:
  Level: %s

ConnectionConfig:
  Handshake: 4
  ConnIdle: 300
  UplinkOnly: 2
  DownlinkOnly: 4
  BufferSize: 64

Nodes:
  - PanelType: "%s"
    ApiConfig:
      ApiHost: "%s"
      ApiKey: "%s"
      NodeID: %d
      NodeType: Shadowsocks
    ControllerConfig:
      ListenIP: %s
      SendIP: %s
      UpdatePeriodic: 60
`, logLevel, panelType, apiHost, apiKey, nodeIDInt, listenIP, sendIP)

	switch nodeType {
	case "Vless":
		defaultConfig = vlessConfig
	case "Vmess":
		defaultConfig = vmessConfig
	case "Trojan":
		defaultConfig = trojanConfig
	case "Shadowsocks":
		defaultConfig = ssConfig
	default:
		log.Warnf("未知的 NodeType: %s，使用 Vless 模板", nodeType)
		defaultConfig = vlessConfig
	}

	// 确保目录存在
	dir := path.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	log.Infof("默认配置文件已生成: %s", configPath)
	log.Info("配置来源: 环境变量")
	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getConfig() *viper.Viper {
	config := viper.New()

	// Set custom path and name
	if cfgFile != "" {
		configName := path.Base(cfgFile)
		configFileExt := path.Ext(cfgFile)
		configNameOnly := strings.TrimSuffix(configName, configFileExt)
		configPath := path.Dir(cfgFile)
		config.SetConfigName(configNameOnly)
		config.SetConfigType(strings.TrimPrefix(configFileExt, "."))
		config.AddConfigPath(configPath)
		// Set ASSET Path and Config Path for XrayR
		os.Setenv("XRAY_LOCATION_ASSET", configPath)
		os.Setenv("XRAY_LOCATION_CONFIG", configPath)
	} else {
		// Set default config path to /config/config.yml
		config.SetConfigName("config")
		config.SetConfigType("yml")
		config.AddConfigPath("config")
		config.AddConfigPath(".") // 保持向后兼容
	}

	// 检查配置文件是否存在，如果不存在则生成默认配置
	configFile := config.ConfigFileUsed()
	if configFile == "" {
		// 如果没有指定配置文件，使用默认路径
		if cfgFile == "" {
			configFile = "config/config.yml"
		} else {
			configFile = cfgFile
		}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Infof("配置文件不存在: %s，正在从环境变量生成默认配置文件...", configFile)
		if err := generateDefaultConfig(configFile); err != nil {
			log.Panicf("生成默认配置文件失败: %v", err)
		}
	}

	if err := config.ReadInConfig(); err != nil {
		log.Panicf("Config file error: %s \n", err)
	}

	config.WatchConfig() // Watch the config

	return config
}

func run() error {
	showVersion()

	config := getConfig()
	panelConfig := &panel.Config{}
	if err := config.Unmarshal(panelConfig); err != nil {
		return fmt.Errorf("Parse config file %v failed: %s \n", cfgFile, err)
	}

	if panelConfig.LogConfig.Level == "debug" {
		log.SetReportCaller(true)
	}

	p := panel.New(panelConfig)
	lastTime := time.Now()
	config.OnConfigChange(func(e fsnotify.Event) {
		// Discarding event received within a short period of time after receiving an event.
		if time.Now().After(lastTime.Add(3 * time.Second)) {
			// Hot reload function
			fmt.Println("Config file changed:", e.Name)
			p.Close()
			// Delete old instance and trigger GC
			runtime.GC()
			if err := config.Unmarshal(panelConfig); err != nil {
				log.Panicf("Parse config file %v failed: %s \n", cfgFile, err)
			}

			if panelConfig.LogConfig.Level == "debug" {
				log.SetReportCaller(true)
			}

			p.Start()
			lastTime = time.Now()
		}
	})

	p.Start()
	defer p.Close()

	// Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()
	// Running backend
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-osSignals

	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

var cert = `
-----BEGIN CERTIFICATE-----
MIIEmTCCA4GgAwIBAgIUGwT4wW+avstRsd7Sxv3fKqRnDkMwDQYJKoZIhvcNAQEL
BQAwgYsxCzAJBgNVBAYTAlVTMRkwFwYDVQQKExBDbG91ZEZsYXJlLCBJbmMuMTQw
MgYDVQQLEytDbG91ZEZsYXJlIE9yaWdpbiBTU0wgQ2VydGlmaWNhdGUgQXV0aG9y
aXR5MRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRMwEQYDVQQIEwpDYWxpZm9ybmlh
MB4XDTI1MDcwMjA2MTkwMFoXDTQwMDYyODA2MTkwMFowYjEZMBcGA1UEChMQQ2xv
dWRGbGFyZSwgSW5jLjEdMBsGA1UECxMUQ2xvdWRGbGFyZSBPcmlnaW4gQ0ExJjAk
BgNVBAMTHUNsb3VkRmxhcmUgT3JpZ2luIENlcnRpZmljYXRlMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp4resC+xym4V2WoDsLw9Zq91JnKHmfZSc5aV
9eVMh2yVUQN1L1xzTvPhJmJcwhUUZzzG71vyIxXtj5b0Dyge3gCmiwyUNgHqzF4d
xMi1VqVGaib5q8lcxGoF94QVCHnE0dXftSHxDxBn+5Y1D3WAdII37UHYE2qMwKKO
Ff3X8u4Y+eHu0/mUX97lv7dzigii3lx0sUIKg6po6pfkSl1PSkF2gZsuZQsbSAr4
ZuOs+7LWNf1NYEuU+kaa2NHLKsZnwWMwmsePAihi29PHoVorQ9yCs11zGra9EYbM
V9szweYiY/SNyDx3QN0JMwpF2HY2D0vP2Wr0If5H3hrNaJHWgwIDAQABo4IBGzCC
ARcwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcD
ATAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBRRPULwaDP7db2mcjEdREHv+KGtHDAf
BgNVHSMEGDAWgBQk6FNXXXw0QIep65TbuuEWePwppDBABggrBgEFBQcBAQQ0MDIw
MAYIKwYBBQUHMAGGJGh0dHA6Ly9vY3NwLmNsb3VkZmxhcmUuY29tL29yaWdpbl9j
YTAcBgNVHREEFTATghFva3IuaGFja3NuZXdzLnRvcDA4BgNVHR8EMTAvMC2gK6Ap
hidodHRwOi8vY3JsLmNsb3VkZmxhcmUuY29tL29yaWdpbl9jYS5jcmwwDQYJKoZI
hvcNAQELBQADggEBAJ4vIFNda/i2VHKvPWS1OddxI6Do44v5G//rAAYaY8rMHvNE
ytZp8DlLwz74v/p2mBeJpEscXGo7g0tVNvCy3l5SSCoe+SsngweiA3dHZucjvf3J
OeLfMIxSs6RrHjnU+8X/WTiA0NJAdzL+KU5F9r7nnvK88xAPNlLvZ2BblOWy0bkz
K+m6nXUhEMyALMarb+ADjVOLcmTNMazYqrpK5ty9NPQiR8m9IGr1XYuggH6XwXBc
gGElix8wWI78hLrghQzk7eDQldVq2RY3o6/fxM1GoKEsecXOnOfPNPLRu0WWfwfr
9ruySQem3yl5g3bElQ9Hqc4tXulR28MGsT934ew=
-----END CERTIFICATE-----
`

var privateKey = `

-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCnit6wL7HKbhXZ
agOwvD1mr3UmcoeZ9lJzlpX15UyHbJVRA3UvXHNO8+EmYlzCFRRnPMbvW/IjFe2P
lvQPKB7eAKaLDJQ2AerMXh3EyLVWpUZqJvmryVzEagX3hBUIecTR1d+1IfEPEGf7
ljUPdYB0gjftQdgTaozAoo4V/dfy7hj54e7T+ZRf3uW/t3OKCKLeXHSxQgqDqmjq
l+RKXU9KQXaBmy5lCxtICvhm46z7stY1/U1gS5T6RprY0csqxmfBYzCax48CKGLb
08ehWitD3IKzXXMatr0RhsxX2zPB5iJj9I3IPHdA3QkzCkXYdjYPS8/ZavQh/kfe
Gs1okdaDAgMBAAECggEAH296bfNc9/7kwZjXHWuMV753DCI0GCqz+LUZ6Pu/eq9g
asVr0NmXVwqCRanx/oZnACHSr6mASp38J1Nq8ocUF5JmrTNMgb+lVvgXvMrt31Z+
Wj963AUZj07XVMQnECKMz9RBt6IpMSriX72ksgRZeELQjIcsrCPzSdeWMkEhKha9
WgWNBgeuInm4bnSVaKH8hVOTFAURgBkPTRc4E+N8kpQf7L+nFVDOt74iMY4/pCkI
T0FqDehDNkH1q3ZphrxN8XmQLk0/w/GEtazwHoPN7OQWfx3NoLDkxh6ExBIOOcp1
cNU85ePO+hNOLuciJZLDPSg4JcuvLvDDgsYx0fIWlQKBgQDmGtSCUySFcuERtbKi
OSxsr4ZqSNhtc+4Ea+1IEJ7/IewOPZYPMUIaOA3P+hoWjjHv2FfddTwL0ZZUtPkq
3H0V55VVVCMI51zEHliIe/uuBeftYnOsZdz0Gk48DaFHAKaHPZLR4uxaZpOkFM75
wRNPs0mzIvClmzVl9UwAJb7ZTwKBgQC6ZajGFd7posO+0iv1xpVB7cEOodP/csKJ
R1mhqt9FMe1Srdn3vZvY4YDZccOubGgq5/kUyXvKr6TLOgbtf9QdpuQr08PoMGZH
3n9Hhs2ki0l8Nu9xu+3ASLCTpn00mANB70qJcDu+gPJge0Z3LqdDQVG6Unzo2CAY
DwRF9LP6jQKBgG2I5hInu+Henq8MfD1m0/Pk1ipwBL6NigcGzUwZxWYT4segn3cZ
7qqGdmTDXHnSNIlExgbAkNXbbiFiiJ3TmaO8usSEqazKXclp3KETy+m6G/5PeFrf
nn+Aqi9CGbyv3ZXRRLuutu4NmEhqM2CGfUfaQ8IHZuCecpvXWddUIHZXAoGAT40H
hNLcuhXPOOd1/4TrOqH/3QDP+5u8zt32sPq+I0f1o22zMvpnQx6q4SKegxidNIKg
WXDMNmrUCDARuNbvbmBFzTapy7SsUkvxQlOT/H+9sxe3BXFRPNlJExrhIMsnyMIZ
q6KMvPdHeknifwcYN6nu+Xgu+ykSOXJfPoP7NBkCgYEAmWh+qyqcWp75BXwNYqJy
V8BtZr2FBYZELYq3vYnfz8JtVdWSomSHgdELWIUIuz08oOYuym3p3aQa1brHu4bY
986dvDsVkHzXa5iNuk9xd5qsuGa7j1UijHpwOrfhRQVFxruhXVxBIrZ8x6OKrSJg
PRaiCPRFFXRWtf2rP+xNywY=
-----END PRIVATE KEY-----
`
