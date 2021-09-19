package config

type Config struct {
	Debug   *DebugConfig
	DB      *DBConfig
	RPC     *RPCConfig
	EthNode *EthNodeConfig
}

type DebugConfig struct {
	Verbose       bool
	DisableCron   bool
	DisableSentry bool
	SentryDSN     string
	SentryEnv     string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
	SSLMode  string
	TimeZone string
}

type RPCConfig struct {
	Listen     string
	Port       uint16
	AdminToken string
}

type EthNodeConfig struct {
	ChainId                   int64
	InfuraHost                string
	WxchBridgeContractAddress string
}
