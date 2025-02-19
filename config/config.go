package config

import (
	"errors"
	"os"
	"path"

	"github.com/creasty/defaults"
)

type JsonRpcConfig struct {
	Url     string            `yaml:"url" json:"url" validate:"omitempty,url"`
	Headers map[string]string `yaml:"headers" json:"headers"`
}

type ScannerConfig struct {
	StartBlock         int           `yaml:"-" json:"_startBlock"`
	EndBlock           int           `yaml:"-" json:"_endBlock"`
	JsonRpc            JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	DisableAutostart   bool          `yaml:"disableAutostart" json:"disableAutostart"`
	BlockRateLimit     int           `yaml:"blockRateLimit" json:"blockRateLimit" default:"200"`
	BlockMaxAgeSeconds int64         `json:"blockMaxAgeSeconds" json:"blockMaxAgeSeconds" default:"600"`
}

type TraceConfig struct {
	JsonRpc JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc"`
	Enabled bool          `yaml:"enabled" json:"enabled"`
}

type RateLimitConfig struct {
	Rate  float64 `yaml:"rate" json:"rate"`
	Burst int     `yaml:"burst" json:"burst" validate:"min=1"`
}

type JsonRpcProxyConfig struct {
	JsonRpc         JsonRpcConfig    `yaml:"jsonRpc" json:"jsonRpc"`
	RateLimitConfig *RateLimitConfig `yaml:"rateLimit" json:"rateLimit"`
}

type LogConfig struct {
	Level       string `yaml:"level" json:"level" default:"info" `
	MaxLogSize  string `yaml:"maxLogSize" json:"maxLogSize" default:"50m" `
	MaxLogFiles int    `yaml:"maxLogFiles" json:"maxLogFiles" default:"10" `
}

type RegistryConfig struct {
	JsonRpc              JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc" default:"{\"url\": \"https://polygon-rpc.com\"}"`
	IPFS                 IPFSConfig    `yaml:"ipfs" json:"ipfs"`
	ContractAddress      string        `yaml:"contractAddress" json:"contractAddress" validate:"eth_addr"`
	ContainerRegistry    string        `yaml:"containerRegistry" json:"containerRegistry" validate:"hostname|hostname_port" default:"disco.forta.network" `
	Username             string        `yaml:"username" json:"username"`
	Password             string        `yaml:"password" json:"password"`
	Disable              bool          `yaml:"disable" json:"disable"` // for testing situations
	CheckIntervalSeconds int           `yaml:"checkIntervalSeconds" json:"checkIntervalSeconds" default:"15"`
}

type IPFSConfig struct {
	GatewayURL string `yaml:"gatewayUrl" json:"gatewayUrl" validate:"url" default:"https://ipfs.forta.network" `
	APIURL     string `yaml:"apiUrl" json:"apiUrl" validate:"url" default:"https://ipfs.forta.network" `
	Username   string `yaml:"username" json:"username"`
	Password   string `yaml:"password" json:"password"`
}

type BatchConfig struct {
	SkipEmpty       bool `yaml:"skipEmpty" json:"skipEmpty"`
	IntervalSeconds *int `yaml:"intervalSeconds" json:"intervalSeconds" default:"15" `
	MaxAlerts       *int `yaml:"maxAlerts" json:"maxAlerts" default:"1000" `
}

type TestAlertsConfig struct {
	Disable    bool   `yaml:"disable" json:"disable"`
	WebhookURL string `yaml:"webhookUrl" json:"webhookUrl" validate:"omitempty,url"`
}

type PublisherConfig struct {
	SkipPublish bool             `yaml:"skipPublish" json:"skipPublish" default:"false"`
	APIURL      string           `yaml:"apiUrl" json:"apiUrl" default:"https://alerts.forta.network" validate:"url"`
	IPFS        IPFSConfig       `yaml:"ipfs" json:"ipfs" validate:"required_unless=SkipPublish true"`
	Batch       BatchConfig      `yaml:"batch" json:"batch"`
	TestAlerts  TestAlertsConfig `yaml:"testAlerts" json:"testAlerts"`
}

type ResourcesConfig struct {
	DisableAgentLimits bool    `yaml:"disableAgentLimits" json:"disableAgentLimits" default:"false" `
	AgentMaxMemoryMiB  int     `yaml:"agentMaxMemoryMib" json:"agentMaxMemoryMib" validate:"omitempty,min=100"`
	AgentMaxCPUs       float64 `yaml:"agentMaxCpus" json:"agentMaxCpus" validate:"omitempty,gt=0"`
}

type ENSConfig struct {
	DefaultContract bool          `yaml:"defaultContract" json:"defaultContract" default:"false" `
	ContractAddress string        `yaml:"contractAddress" json:"contractAddress" validate:"omitempty,eth_addr" default:"0x08f42fcc52a9C2F391bF507C4E8688D0b53e1bd7"`
	JsonRpc         JsonRpcConfig `yaml:"jsonRpc" json:"jsonRpc" default:"{\"url\": \"https://polygon-rpc.com\"}" `
	Override        bool          `yaml:"override" json:"override" default:"false"`
}

type TelemetryConfig struct {
	URL     string `yaml:"url" json:"url" default:"https://alerts.forta.network/telemetry" validate:"url"`
	Disable bool   `yaml:"disable" json:"disable"`
}

type AutoUpdateConfig struct {
	Disable              bool `yaml:"disable" json:"disable"`
	CheckIntervalSeconds int  `yaml:"checkIntervalSeconds" json:"checkIntervalSeconds" default:"60"`
}

type AgentLogsConfig struct {
	URL     string `yaml:"url" json:"url" default:"https://alerts.forta.network/logs/agents" validate:"url"`
	Disable bool   `yaml:"disable" json:"disable"`
}

type ContainerRegistryConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

type PrivateModeConfig struct {
	Enable            bool                     `yaml:"enable" json:"enable"`
	AgentImages       []string                 `yaml:"agentImages" json:"agentImages" validate:"required_if=Enable true"`
	WebhookURL        string                   `yaml:"webhookUrl" json:"webhookUrl" validate:"required_if=Enable true"`
	ContainerRegistry *ContainerRegistryConfig `yaml:"containerRegistry" json:"containerRegistry"`
}

type Config struct {
	// runtime values

	Development                    bool           `yaml:"-" json:"_development"`
	FortaDir                       string         `yaml:"-" json:"_fortaDir"`
	KeyDirPath                     string         `yaml:"-" json:"_keyDirPath"`
	Passphrase                     string         `yaml:"-" json:"_passphrase"`
	ExposeNats                     bool           `yaml:"-" json:"_exposeNats"`
	LocalAgentsPath                string         `yaml:"-" json:"_localAgentsPath"`
	LocalAgents                    []*AgentConfig `yaml:"-" json:"_localAgents"`
	AgentRegistryContractAddress   string         `yaml:"-" json:"_agentRegistryContractAddress"`
	ScannerVersionContractAddress  string         `yaml:"-" json:"_scannerVersionContractAddress"`
	ScannerRegistryContractAddress string         `yaml:"-" json:"_scannerRegistryContractAddress"`

	// yaml config values

	ChainID int `yaml:"chainId" json:"chainId" default:"1" `

	Scan  ScannerConfig `yaml:"scan" json:"scan"`
	Trace TraceConfig   `yaml:"trace" json:"trace"`

	Registry          RegistryConfig     `yaml:"registry" json:"registry"`
	Publish           PublisherConfig    `yaml:"publish" json:"publish"`
	JsonRpcProxy      JsonRpcProxyConfig `yaml:"jsonRpcProxy" json:"jsonRpcProxy"`
	Log               LogConfig          `yaml:"log" json:"log"`
	ResourcesConfig   ResourcesConfig    `yaml:"resources" json:"resources"`
	ENSConfig         ENSConfig          `yaml:"ens" json:"ens"`
	TelemetryConfig   TelemetryConfig    `yaml:"telemetry" json:"telemetry"`
	AutoUpdate        AutoUpdateConfig   `yaml:"autoUpdate" json:"autoUpdate"`
	AgentLogsConfig   AgentLogsConfig    `yaml:"agentLogs" json:"agentLogs"`
	PrivateModeConfig PrivateModeConfig  `yaml:"privateMode" json:"privateMode"`
}

func (cfg *Config) ConfigFilePath() string {
	return path.Join(cfg.FortaDir, DefaultConfigFileName)
}

//GetConfigForContainer is how a container gets the forta configuration (file or env var)
func GetConfigForContainer() (Config, error) {
	var cfg Config
	if _, err := os.Stat(DefaultContainerConfigPath); os.IsNotExist(err) {
		return cfg, errors.New("config file not found")
	}
	cfg, err := getConfigFromFile(DefaultContainerConfigPath)
	if err != nil {
		return Config{}, err
	}
	applyContextDefaults(&cfg)
	return cfg, nil
}

// apply defaults that apply in certain contexts
func applyContextDefaults(cfg *Config) {
	if cfg.ChainID == 1 {
		cfg.Trace.Enabled = true
	}
	if cfg.ENSConfig.DefaultContract {
		cfg.ENSConfig.ContractAddress = ""
	}
	cfg.FortaDir = DefaultContainerFortaDirPath
	cfg.KeyDirPath = path.Join(cfg.FortaDir, DefaultKeysDirName)
}

func getConfigFromFile(filename string) (Config, error) {
	var cfg Config
	if err := readFile(filename, &cfg); err != nil {
		return Config{}, err
	}
	if err := defaults.Set(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
