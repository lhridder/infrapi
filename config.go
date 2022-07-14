package infrapi

import (
	"encoding/json"
	"os"
)

var Config GlobalConfig

type GlobalConfig struct {
	ApiBind   string `json:"apiBind"`
	RedisHost string `json:"redisHost"`
	RedisDB   int    `json:"redisDB"`
	RedisPass string `json:"redisPass"`
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
	ApiBind:   ":5000",
	RedisHost: "localhost",
	RedisDB:   0,
	RedisPass: "",
}

func LoadGlobalConfig() error {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	var config = DefaultConfig
	jsonParser := json.NewDecoder(jsonFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return err
	}
	Config = config
	_ = jsonFile.Close()
	return nil
}
