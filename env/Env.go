package env

import "github.com/magiconair/properties"

type AppConfig struct {
	HomeDir      string `properties:"homeDir,default=/home"`
	RootDir      string `properties:"rootDir,default="`
	NodeDir      string `properties:"nodeDir,default=/home/node"`
	ContractsDir string `properties:"contractsDir,default=/root/quorum-maker/contracts"`
	GethLogs     string `properties:"gethLogs,default=/home/node/qdata/gethLogs"`
	PrivacyLogs  string `properties:"privacyLogs,default=/home/node/qdata/constellationLogs"`
}

var appConfig AppConfig

func GetAppConfig(refresh ...bool) AppConfig {

	if (nil != refresh || AppConfig{} == appConfig) {
		filename := "application.conf"
		p := properties.MustLoadFile(filename, properties.UTF8)

		p.Decode(&appConfig)
	}

	return appConfig

}

type SetupConf struct {
	ContextPath         string `properties:"CONTEXT_PATH,default="`
	CurrentIp           string `properties:"CURRENT_IP,default="`
	RpcPort             string `properties:"RPC_PORT,default="`
	WhisperPort         string `properties:"WHISPER_PORT,default="`
	ConstellationPort   string `properties:"CONSTELLATION_PORT,default="`
	RaftPort            string `properties:"RAFT_PORT,default="`
	ThisNodemanagerPort string `properties:"THIS_NODEMANAGER_PORT,default="`
	WsPort              string `properties:"WS_PORT,default="`
	NetworkId           string `properties:"NETWORK_ID,default="`
	RaftId              string `properties:"RAFT_ID,default="`
	Nodename            string `properties:"NODENAME,default="`
	Mode                string `properties:"MODE,default="`
	State               string `properties:"STATE,default="`
	PubKey              string `properties:"PUBKEY,default="`
	ContractAdd         string `properties:"CONTRACT_ADD,default="`
	Registered          string `properties:"REGISTERED,default="`
	Role                string `properties:"ROLE,default="`
	RecipientList       string `properties:"RECIPIENTLIST,default="`
}

var setupConf SetupConf

func GetSetupConf(refresh ...bool) SetupConf {
	if nil != refresh {
		filename := GetAppConfig().HomeDir + "/setup.conf"
		p := properties.MustLoadFile(filename, properties.UTF8)

		p.Decode(&setupConf)
	}
	return setupConf
}
