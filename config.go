package infrapi

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var Config GlobalConfig

type Redis struct {
	Host string `yaml:"host"`
	Pass string `yaml:"pass"`
	DB   int    `yaml:"db"`
}

type GlobalConfig struct {
	ApiBind string `yaml:"apiBind"`
	Redis   Redis
}

type ProxyConfig struct {
	DomainNames       []string     `json:"domainNames"`
	ListenTo          string       `json:"listenTo"`
	ProxyTo           string       `json:"proxyTo"`
	ProxyBind         string       `json:"proxyBind"`
	ProxyProtocol     bool         `json:"proxyProtocol"`
	RealIP            bool         `json:"realIp"`
	Timeout           int          `json:"timeout"`
	DisconnectMessage string       `json:"disconnectMessage"`
	OnlineStatus      StatusConfig `json:"onlineStatus"`
	OfflineStatus     StatusConfig `json:"offlineStatus"`
}

type StatusConfig struct {
	VersionName    string         `json:"versionName"`
	ProtocolNumber int            `json:"protocolNumber"`
	MaxPlayers     int            `json:"maxPlayers"`
	PlayersOnline  int            `json:"playersOnline"`
	PlayerSamples  []PlayerSample `json:"playerSamples"`
	IconPath       string         `json:"iconPath"`
	MOTD           string         `json:"motd"`
}

type PlayerSample struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

var DefaultConfig = GlobalConfig{
	ApiBind: ":5000",
	Redis: Redis{
		Host: "localhost",
		Pass: "",
		DB:   0,
	},
}

func LoadGlobalConfig() error {
	log.Println("Loading config.yml")
	ymlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return err
	}
	var config = DefaultConfig
	err = yaml.Unmarshal(ymlFile, &config)
	if err != nil {
		return err
	}
	Config = config
	return nil
}
