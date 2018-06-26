package service

import (
	"bytes"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
)

type ConnectionInfo struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	Enode string `json:"enode"`
}

type PendingRequests struct {
	NodeName string `json:"nodeName"`
	Enode    string `json:"enode"`
	Message  string `json:"message"`
	EnodeID  string `json:"enodeid"`
	IP       string `json:"ip"`
}

type NodeInfo struct {
	NodeName       string           `json:"nodeName"`
	NodeCount      int              `json:"nodeCount"`
	TotalNodeCount int              `json:"totalNodeCount"`
	Active         string           `json:"active"`
	ConnectionInfo ConnectionInfo   `json:"connectionInfo"`
	RaftRole       string           `json:"raftRole"`
	RaftID         int              `json:"raftID"`
	BlockNumber    int64            `json:"blockNumber"`
	PendingTxCount int              `json:"pendingTxCount"`
	Genesis        string           `json:"genesis"`
	AdminInfo      client.AdminInfo `json:"adminInfo"`
}

type JoinNetworkRequest struct {
	EnodeID   string `json:"enode-id,omitempty"`
	IPAddress string `json:"ip-address,omitempty"`
	Nodename  string `json:"nodename,omitempty"`
}

type GetGenesisResponse struct {
	ContstellationPort string `json:"contstellation-port"`
	NetID              string `json:"netID"`
	Genesis            string `json:"genesis"`
}

type BlockDetailsResponse struct {
	Number           int64                        `json:"number"`
	Hash             string                       `json:"hash"`
	ParentHash       string                       `json:"parentHash"`
	Nonce            string                       `json:"nonce"`
	Sha3Uncles       string                       `json:"sha3Uncles"`
	LogsBloom        string                       `json:"logsBloom"`
	TransactionsRoot string                       `json:"transactionsRoot"`
	StateRoot        string                       `json:"stateRoot"`
	Miner            string                       `json:"miner"`
	Difficulty       int64                        `json:"difficulty"`
	TotalDifficulty  int64                        `json:"totalDifficulty"`
	ExtraData        string                       `json:"extraData"`
	Size             int64                        `json:"size"`
	GasLimit         int64                        `json:"gasLimit"`
	GasUsed          int64                        `json:"gasUsed"`
	Timestamp        int64                        `json:"timestamp"`
	Transactions     []TransactionDetailsResponse `json:"transactions"`
	Uncles           []string                     `json:"uncles"`
	TimeElapsed      int64                        `json:"TimeElapsed"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      int64  `json:"blockNumber"`
	From             string `json:"from"`
	Gas              int64  `json:"gas"`
	GasPrice         int64  `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            int64  `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex int64  `json:"transactionIndex"`
	Value            int64  `json:"value"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
	TransactionType  string `json:"transactionType"`
	TimeElapsed      int64  `json:"TimeElapsed"`
}

type TransactionReceiptResponse struct {
	BlockHash         string                         `json:"blockHash"`
	BlockNumber       int64                          `json:"blockNumber"`
	ContractAddress   string                         `json:"contractAddress"`
	CumulativeGasUsed int64                          `json:"cumulativeGasUsed"`
	From              string                         `json:"from"`
	Gas               int64                          `json:"gas"`
	GasPrice          int64                          `json:"gasPrice"`
	GasUsed           int64                          `json:"gasUsed"`
	Input             string                         `json:"input"`
	Logs              []Logs                         `json:"logs"`
	LogsBloom         string                         `json:"logsBloom"`
	Nonce             int64                          `json:"nonce"`
	Root              string                         `json:"root"`
	To                string                         `json:"to"`
	TransactionHash   string                         `json:"transactionHash"`
	TransactionIndex  int64                          `json:"transactionIndex"`
	Value             int64                          `json:"value"`
	V                 string                         `json:"v"`
	R                 string                         `json:"r"`
	S                 string                         `json:"s"`
	TransactionType   string                         `json:"transactionType"`
	TimeElapsed       int64                          `json:"TimeElapsed"`
	DecodedInputs     []contractclient.ParamTableRow `json:"decodedInputs"`
}

type Logs struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      int64    `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         int64    `json:"logIndex"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex int64    `json:"transactionIndex"`
}

type JoinNetworkResponse struct {
	EnodeID string `json:"enode-id"`
	Status  string `json:"status"`
}

type ContractJson struct {
	Filename        string `json:"filename"`
	Interface       string `json:"interface"`
	Bytecode        string `json:"bytecode"`
	ContractAddress string `json:"address"`
	Json            string `json:"json"`
}

type CreateNetworkScriptArgs struct {
	Nodename          string `json:"nodename,omitempty"`
	CurrentIP         string `json:"currentIP,omitempty"`
	RPCPort           string `json:"rpcPort,omitempty"`
	WhisperPort       string `json:"whisperPort,omitempty"`
	ConstellationPort string `json:"constellationPort,omitempty"`
	RaftPort          string `json:"raftPort,omitempty"`
	NodeManagerPort   string `json:"nodeManagerPort,omitempty"`
}

type JoinNetworkScriptArgs struct {
	Nodename              string `json:"nodename,omitempty"`
	CurrentIP             string `json:"currentIP,omitempty"`
	RPCPort               string `json:"rpcPort,omitempty"`
	WhisperPort           string `json:"whisperPort,omitempty"`
	ConstellationPort     string `json:"constellationPort,omitempty"`
	RaftPort              string `json:"raftPort,omitempty"`
	NodeManagerPort       string `json:"nodeManagerPort,omitempty"`
	MasterNodeManagerPort string `json:"masterNodeManagerPort,omitempty"`
	MasterIP              string `json:"masterIP,omitempty"`
}

type SuccessResponse struct {
	Status string `json:"statusMessage"`
}

type LatestBlockResponse struct {
	LatestBlockNumber int64 `json:"latestBlockNumber"`
	TimeElapsed       int64 `json:"TimeElapsed"`
}

type NodeList struct {
	NodeName  string `json:"nodeName"`
	Role      string `json:"role,omitempty"`
	PublicKey string `json:"publicKey"`
	IP        string `json:"ip,omitempty"`
	Enode     string `json:"enode,omitempty"`
}

type MailServerConfig struct {
	Host          string `json:"smtpServerHost"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	RecipientList string `json:"recipientList"`
}

type LatencyResponse struct {
	EnodeID string `json:"enode-id"`
	Latency string `json:"latency"`
}

type NodeServiceImpl struct {
	Url string
}

type ChartInfo struct {
	TimeStamp        int `json:"timeStamp"`
	BlockCount       int `json:"blockCount"`
	TransactionCount int `json:"transactionCount"`
}

type ContractTableRow struct {
	ContractAdd  string `json:"contractAddress"`
	ContractName string `json:"contractName"`
	ABIContent   string `json:"abi"`
	Sender       string `json:"sender"`
	ContractType string `json:"contractType"`
	Description  string `json:"description"`
}

type ContractCounter struct {
	TotalContracts int `json:"totalContracts"`
	ABIavailable   int `json:"abis"`
}

var txnMap = map[string]TransactionReceiptResponse{}
var abiMap = map[string]string{}
var contractCrawlerMutex = 0

var contDescriptionMap = map[string]string{}
var contTypeMap = map[string]string{}
var contSenderMap = map[string]string{}
var contNameMap = map[string]string{}
var chartSize = 10
var warning = 0
var lastCrawledBlock = 0
var mailServerConfig MailServerConfig

func (nsi *NodeServiceImpl) getGenesis(url string) (response GetGenesisResponse) {
	var netId, constl string
	existsA := util.PropertyExists("NETWORK_ID", "/home/setup.conf")
	existsB := util.PropertyExists("CONSTELLATION_PORT", "/home/setup.conf")
	if existsA != "" && existsB != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		netId = util.MustGetString("NETWORK_ID", p)
		constl = util.MustGetString("CONSTELLATION_PORT", p)
	}
	b, err := ioutil.ReadFile("/home/node/genesis.json")
	if err != nil {
		log.Fatal(err)
	}
	genesis := string(b)
	genesis = strings.Replace(genesis, "\n", "", -1)

	response = GetGenesisResponse{constl, netId, genesis}
	return response
}

func (nsi *NodeServiceImpl) joinNetwork(enode string, url string) string {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	raftId := ethClient.RaftAddPeer(enode)
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}
	collatedInfo := fmt.Sprint(raftId, ":", contractAdd)
	return collatedInfo
}

//@TODO: If this function is repeatedly called from UI, please cache the static informations.
func (nsi *NodeServiceImpl) getCurrentNode(url string) NodeInfo {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()
	var contractAdd string
	var p *properties.Properties
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p = properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, contractAdd, "", nil}}

	totalCount := len(nms.GetNodeDetailsList())
	var activeStatus string
	active := ethClient.NetListening()
	if active == true {
		activeStatus = "active"
	} else {
		activeStatus = "inactive"
	}
	otherPeersResponse := ethClient.AdminPeers()
	count := len(otherPeersResponse)
	count = count + 1

	//@TODO: Cant the regex be simply start*.sh ?
	r, _ := regexp.Compile("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]")

	//@TODO: Use of absolute path (starting with "/") is highly error prone
	files, err := ioutil.ReadDir("/home/node")
	if err != nil {
		log.Fatal(err)
	}
	var nodename string
	for _, f := range files {
		match, _ := regexp.MatchString("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]", f.Name())
		if match {
			nodename = r.FindString(f.Name())
		}
	}

	//@TODO: use grouping in regex to find the name of the node. Refer to whats app message.
	nodename = strings.TrimSuffix(nodename, ".sh")
	nodename = strings.TrimPrefix(nodename, "start_")

	var ipAddr, raftId, rpcPort, nodeName string
	existsA := util.PropertyExists("CURRENT_IP", "/home/setup.conf")
	existsB := util.PropertyExists("RAFT_ID", "/home/setup.conf")
	existsC := util.PropertyExists("RPC_PORT", "/home/setup.conf")
	existsD := util.PropertyExists("NODENAME", "/home/setup.conf")
	if existsA != "" && existsB != "" && existsC != "" && existsD != "" {
		ipAddr = util.MustGetString("CURRENT_IP", p)
		raftId = util.MustGetString("RAFT_ID", p)
		rpcPort = util.MustGetString("RPC_PORT", p)
		nodeName = util.MustGetString("NODENAME", p)
	}
	raftIdInt, err := strconv.Atoi(raftId)
	if err != nil {
		log.Fatal(err)
	}

	rpcPortInt, err := strconv.Atoi(rpcPort)
	if err != nil {
		log.Fatal(err)
	}

	thisAdminInfo := ethClient.AdminNodeInfo()
	enode := thisAdminInfo.Enode

	pendingTxResponse := ethClient.PendingTransactions()
	pendingTxCount := len(pendingTxResponse)

	blockNumber := ethClient.BlockNumber()
	blockNumberInt := util.HexStringtoInt64(blockNumber)

	raftRole := ethClient.RaftRole()

	raftRole = strings.TrimSuffix(raftRole, "\n")

	b, err := ioutil.ReadFile("/home/node/genesis.json")

	if err != nil {
		log.Fatal(err)
	}

	genesis := string(b)
	genesis = strings.Replace(genesis, "\n", "", -1)
	conn := ConnectionInfo{ipAddr, rpcPortInt, enode}
	responseObj := NodeInfo{nodeName, count, totalCount, activeStatus, conn, raftRole, raftIdInt, blockNumberInt, pendingTxCount, genesis, thisAdminInfo}
	return responseObj
}

func (nsi *NodeServiceImpl) getOtherPeer(peerId string, url string) client.AdminPeers {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	otherPeersResponse := ethClient.AdminPeers()
	for _, item := range otherPeersResponse {
		if item.ID == peerId {
			peerResponse := item
			return peerResponse
		}
	}
	return client.AdminPeers{}
}

func (nsi *NodeServiceImpl) getOtherPeers(url string) []client.AdminPeers {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	otherPeersResponse := ethClient.AdminPeers()
	return otherPeersResponse
}

func (nsi *NodeServiceImpl) getPendingTransactions(url string) []TransactionDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	clientPendingTxResponses := ethClient.PendingTransactions()

	pendingTxResponse := make([]TransactionDetailsResponse, len(clientPendingTxResponses))

	for i, clientPendingTxResponse := range clientPendingTxResponses {

		pendingTxResponse[i] = ConvertToReadable(clientPendingTxResponse, true, true)
	}

	return pendingTxResponse
}

func (nsi *NodeServiceImpl) getBlockInfo(blockno int64, url string) BlockDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	blockNoHex := strconv.FormatInt(blockno, 16)
	bNoHex := fmt.Sprint("0x", blockNoHex)
	var blockResponse BlockDetailsResponse
	blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	blockResponse.TimeElapsed = elapsedTime

	//@TODO: Create a utility function to convert block object to readable object.
	blockResponse.Number = util.HexStringtoInt64(blockResponseClient.Number)
	blockResponse.Difficulty = util.HexStringtoInt64(blockResponseClient.Difficulty)
	blockResponse.TotalDifficulty = util.HexStringtoInt64(blockResponseClient.TotalDifficulty)
	blockResponse.Size = util.HexStringtoInt64(blockResponseClient.Size)
	blockResponse.GasLimit = util.HexStringtoInt64(blockResponseClient.GasLimit)
	blockResponse.GasUsed = util.HexStringtoInt64(blockResponseClient.GasUsed)
	blockResponse.Timestamp = util.HexStringtoInt64(blockResponseClient.Timestamp)
	blockResponse.Hash = blockResponseClient.Hash
	blockResponse.ParentHash = blockResponseClient.ParentHash
	blockResponse.Nonce = blockResponseClient.Nonce
	blockResponse.Sha3Uncles = blockResponseClient.Sha3Uncles
	blockResponse.LogsBloom = blockResponseClient.LogsBloom
	blockResponse.TransactionsRoot = blockResponseClient.TransactionsRoot
	blockResponse.StateRoot = blockResponseClient.StateRoot
	blockResponse.Miner = blockResponseClient.Miner
	blockResponse.ExtraData = blockResponseClient.ExtraData
	blockResponse.Uncles = blockResponseClient.Uncles
	txnNo := len(blockResponseClient.Transactions)
	txResponse := make([]TransactionDetailsResponse, txnNo)
	for i, clientTransactions := range blockResponseClient.Transactions {

		txGetClient := ethClient.GetTransactionByHash(clientTransactions.Hash)
		private := ethClient.GetQuorumPayload(txGetClient.Input)
		txResponse[i] = ConvertToReadable(clientTransactions, false, (private == "0x"))

	}
	blockResponse.Transactions = txResponse
	return blockResponse
}

func (nsi *NodeServiceImpl) getLatestBlockInfo(count string, reference string, url string) []BlockDetailsResponse {
	countVal := util.HexStringtoInt64(count)
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	var blockNumber int64
	if reference == "" {
		blockNumber = util.HexStringtoInt64(ethClient.BlockNumber())
	} else {
		var err error
		blockNumber, err = strconv.ParseInt(reference, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		blockNumber = blockNumber - 1
	}
	start := blockNumber - countVal + 1
	blockResponse := make([]BlockDetailsResponse, countVal)

	//@TODO: call the utility function to convert to readable block object
	for i := start; i <= blockNumber; i++ {
		blockNoHex := strconv.FormatInt(i, 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		blockResponse[blockNumber-i].Number = util.HexStringtoInt64(blockResponseClient.Number)
		blockResponse[blockNumber-i].Hash = blockResponseClient.Hash
		currentTime := time.Now().Unix()
		creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
		creationTimeUnix := creationTime / 1000000000
		elapsedTime := currentTime - creationTimeUnix
		blockResponse[blockNumber-i].TimeElapsed = elapsedTime
		txnNo := len(blockResponseClient.Transactions)
		txResponse := make([]TransactionDetailsResponse, txnNo)

		for i, clientTransactions := range blockResponseClient.Transactions {

			txGetClient := ethClient.GetTransactionByHash(clientTransactions.Hash)
			private := ethClient.GetQuorumPayload(txGetClient.Input)
			txResponse[i] = ConvertToReadable(clientTransactions, false, (private == "0x"))

		}

		blockResponse[blockNumber-i].Transactions = txResponse
	}
	return blockResponse
}

func (nsi *NodeServiceImpl) getLatestTransactionInfo(count string, url string) []BlockDetailsResponse {
	countVal := util.HexStringtoInt64(count)
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	blockNumber := util.HexStringtoInt64(ethClient.BlockNumber())
	start := blockNumber - countVal + 1
	blockResponse := make([]BlockDetailsResponse, countVal)
	for i := start; i <= blockNumber; i++ {
		blockNoHex := strconv.FormatInt(i, 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		currentTime := time.Now().Unix()
		creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
		creationTimeUnix := creationTime / 1000000000
		elapsedTime := currentTime - creationTimeUnix
		blockResponse[blockNumber-i].TimeElapsed = elapsedTime
		blockResponse[blockNumber-i].Number = util.HexStringtoInt64(blockResponseClient.Number)
		txnNo := len(blockResponseClient.Transactions)
		txResponse := make([]TransactionDetailsResponse, txnNo)

		for i, clientTransactions := range blockResponseClient.Transactions {

			txGetClient := ethClient.GetTransactionByHash(clientTransactions.Hash)
			private := ethClient.GetQuorumPayload(txGetClient.Input)
			txResponse[i] = ConvertToReadable(clientTransactions, false, (private == "0x"))

		}

		blockResponse[blockNumber-i].Transactions = txResponse
	}
	return blockResponse
}

func (nsi *NodeServiceImpl) getTransactionInfo(txno string, url string) TransactionDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	var txResponse TransactionDetailsResponse
	txResponseClient := ethClient.GetTransactionByHash(txno)

	private := ethClient.GetQuorumPayload(txResponseClient.Input)
	txResponse = ConvertToReadable(txResponseClient, false, (private == "0x"))

	blockResponseClient := ethClient.GetBlockByNumber(txResponseClient.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txResponse.TimeElapsed = elapsedTime
	return txResponse
}

func (nsi *NodeServiceImpl) getTransactionReceipt(txno string, url string) TransactionReceiptResponse {
	if txnMap[txno].TransactionHash == "" {
		txResponse := populateTransactionObject(txno, url)
		decodeTransactionObject(&txResponse, url)
		return txResponse
	} else {
		txnDetails := txnMap[txno]
		calculateTimeElapsed(&txnDetails, url)
		txnMap[txno] = txnDetails
		return txnDetails
	}
}

func populateTransactionObject(txno string, url string) TransactionReceiptResponse {
	ethClient := client.EthClient{url}
	getTransaction := ethClient.GetTransactionByHash(txno)
	var txResponse TransactionReceiptResponse
	getTransactionReceipt := ethClient.GetTransactionReceipt(txno)
	txResponse.BlockNumber = util.HexStringtoInt64(getTransactionReceipt.BlockNumber)
	txResponse.CumulativeGasUsed = util.HexStringtoInt64(getTransactionReceipt.CumulativeGasUsed)
	txResponse.GasUsed = util.HexStringtoInt64(getTransactionReceipt.GasUsed)
	txResponse.TransactionIndex = util.HexStringtoInt64(getTransactionReceipt.TransactionIndex)
	txResponse.BlockHash = getTransactionReceipt.BlockHash
	txResponse.From = getTransactionReceipt.From
	txResponse.ContractAddress = getTransactionReceipt.ContractAddress
	txResponse.LogsBloom = getTransactionReceipt.LogsBloom
	txResponse.Root = getTransactionReceipt.Root
	txResponse.To = getTransactionReceipt.To
	txResponse.TransactionHash = getTransactionReceipt.TransactionHash
	txResponse.Gas = util.HexStringtoInt64(getTransaction.Gas)
	txResponse.GasPrice = util.HexStringtoInt64(getTransaction.GasPrice)
	txResponse.Input = getTransaction.Input
	txResponse.Nonce = util.HexStringtoInt64(getTransaction.Nonce)
	txResponse.Value = util.HexStringtoInt64(getTransaction.Value)
	txResponse.V = getTransaction.V
	txResponse.R = getTransaction.R
	txResponse.S = getTransaction.S
	eventNo := len(getTransactionReceipt.Logs)
	txResponseBuffer := make([]Logs, eventNo)
	for i := 0; i < eventNo; i++ {
		txResponseBuffer[i].BlockNumber = util.HexStringtoInt64(getTransactionReceipt.Logs[i].BlockNumber)
		txResponseBuffer[i].LogIndex = util.HexStringtoInt64(getTransactionReceipt.Logs[i].LogIndex)
		txResponseBuffer[i].TransactionIndex = util.HexStringtoInt64(getTransactionReceipt.Logs[i].TransactionIndex)
		txResponseBuffer[i].Address = getTransactionReceipt.Logs[i].Address
		txResponseBuffer[i].BlockHash = getTransactionReceipt.Logs[i].BlockHash
		txResponseBuffer[i].Data = getTransactionReceipt.Logs[i].Data
		txResponseBuffer[i].TransactionHash = getTransactionReceipt.Logs[i].TransactionHash
		txResponseBuffer[i].Topics = getTransactionReceipt.Logs[i].Topics
	}
	txResponse.Logs = txResponseBuffer
	blockResponseClient := ethClient.GetBlockByNumber(getTransactionReceipt.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txResponse.TimeElapsed = elapsedTime
	return txResponse
}

func decodeTransactionObject(txnDetails *TransactionReceiptResponse, url string) {
	dataDirect, err := ioutil.ReadFile("./SimpleStorage.json")

	if err != nil {
		fmt.Println(err)
	}

	abiString := string(dataDirect)
	abiString = strings.Replace(abiString, "\n", "", -1)
	abiMap["0xdf73d861f5c41756d28fa6c57bea5d04903ea2e9"] = abiString
	var quorumPayload string
	var decoded bool

	ethClient := client.EthClient{url}

	if util.HexStringtoInt64(txnDetails.V) == 37 || util.HexStringtoInt64(txnDetails.V) == 38 {
		quorumPayload = ethClient.GetQuorumPayload(txnDetails.Input)
		if quorumPayload == "0x" {
			txnDetails.TransactionType = "Hash Only"
		} else {
			txnDetails.TransactionType = "Private"
		}
	} else {
		txnDetails.TransactionType = "Public"

	}
	if txnDetails.ContractAddress == "" {
		if txnDetails.TransactionType == "Private" && abiMap[txnDetails.To] != "" && abiMap[txnDetails.To] != "missing" {
			txnDetails.Input = quorumPayload
			txnDetails.DecodedInputs = contractclient.ABIParser(txnDetails.To, abiMap[txnDetails.To], quorumPayload)
			decoded = true
		} else if txnDetails.TransactionType == "Public" && abiMap[txnDetails.To] != "" && abiMap[txnDetails.To] != "missing" {
			txnDetails.DecodedInputs = contractclient.ABIParser(txnDetails.To, abiMap[txnDetails.To], txnDetails.Input)
			decoded = true
		} else if txnDetails.TransactionType == "Hash Only" {
			decodeFail := make([]contractclient.ParamTableRow, 1)
			decodeFail[0].Key = "NA"
			decodeFail[0].Value = "Hash Only Transaction"
			txnDetails.DecodedInputs = decodeFail
			decoded = true
		}
	}
	if abiMap[txnDetails.To] == "" {
		decodeFail := make([]contractclient.ParamTableRow, 1)
		decodeFail[0].Key = "NA"
		decodeFail[0].Value = "Decode in Progress"
		txnDetails.DecodedInputs = decodeFail
	} else if abiMap[txnDetails.To] == "missing" {
		decodeFail := make([]contractclient.ParamTableRow, 1)
		decodeFail[0].Key = "NA"
		decodeFail[0].Value = "ABI is Missing"
		txnDetails.DecodedInputs = decodeFail
	}
	if decoded {
		txnMap[txnDetails.TransactionHash] = *txnDetails
	}
}

func calculateTimeElapsed(txnDetails *TransactionReceiptResponse, url string) {
	ethClient := client.EthClient{url}
	getTransactionReceipt := ethClient.GetTransactionReceipt(txnDetails.TransactionHash)
	blockResponseClient := ethClient.GetBlockByNumber(getTransactionReceipt.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txnDetails.TimeElapsed = elapsedTime
}

func (nsi *NodeServiceImpl) joinRequestResponse(enode string, status string) SuccessResponse {
	var successResponse SuccessResponse
	peerMap[enode] = status
	var enodeString []string
	var ipString []string

	//@TODO: Use regex grouping to extract parts
	enodeVal := strings.TrimPrefix(enode, "enode://")
	enodeString = strings.Split(enodeVal, "@")
	ipString = strings.Split(enodeString[1], ":")
	ip := ipString[0]
	successResponse.Status = fmt.Sprintf("Successfully updated status of %s node with IP: %s to %s", nameMap[enode], ip, status)
	return successResponse
}

func (nsi *NodeServiceImpl) deployContract(pubKeys []string, fileName []string, private bool, url string) []ContractJson {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()

	//@TODO: Dont use absolute paths
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}
	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, contractAdd, "", nil}}
	if private == true && pubKeys[0] == "" {
		enode := ethClient.AdminNodeInfo().ID
		peerNo := len(nms.GetNodeDetailsList())
		publicKeys := make([]string, peerNo-1)
		for i := 0; i < peerNo; i++ {
			if enode != nms.GetNodeDetails(i).Enode {
				publicKeys[i-1] = nms.GetNodeDetails(i).PublicKey
			}
		}
		//pubKeys = []string{"R1fOFUfzBbSVaXEYecrlo9rENW0dam0kmaA2pasGM14=", "Er5J8G+jXQA9O2eu7YdhkraYM+j+O5ArnMSZ24PpLQY="}
		pubKeys = publicKeys
	}
	var solc string
	if solc == "" { //@TODO huh ???
		solc = "solc"
	}
	fileNo := len(fileName)
	contractJsonArr := make([]ContractJson, fileNo)
	for i := 0; i < fileNo; i++ {
		var binOut bytes.Buffer
		var errorstring bytes.Buffer
		cmd := exec.Command(solc, "-o", ".", "--overwrite", "--bin", fileName[i])
		cmd = exec.Command(solc, "--bin", fileName[i])
		cmd.Stdout = &binOut
		cmd.Stderr = &errorstring
		err := cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + errorstring.String())
			contractJsonArr[i].Filename = strings.Replace(fileName[i], ".sol", "", -1)
			contractJsonArr[i].Json = "Compilation Failed: JSON could not be created"
			contractJsonArr[i].Bytecode = "Compilation Failed: " + errorstring.String()
		}
		var abiOut bytes.Buffer
		cmd = exec.Command(solc, "-o", ".", "--overwrite", "--abi", fileName[i])
		cmd = exec.Command(solc, "--abi", fileName[i])
		cmd.Stdout = &abiOut
		cmd.Stderr = &errorstring
		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + errorstring.String())
			contractJsonArr[i].Interface = "Compilation Failed: " + errorstring.String()
			contractJsonArr[i].ContractAddress = "0x"
			continue
		}

		byteCode := binOut.String()

		re := regexp.MustCompile("Binary?")
		index := re.FindStringIndex(byteCode)
		var start int
		start = index[1] + 3
		byteCode = byteCode[start:]
		byteCode = "0x" + byteCode

		abiStr := abiOut.String()

		re = regexp.MustCompile("ABI?")
		index = re.FindStringIndex(abiStr)
		start = index[1] + 2
		abiStr = abiStr[start:]

		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		byteCode = reg.ReplaceAllString(byteCode, "")

		contractAddress := ethClient.DeployContracts(byteCode, pubKeys, private)

		path := "./contracts"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0775)
		}
		path = "./contracts/" + contractAddress + "_" + strings.Replace(fileName[i], ".sol", "", -1)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0775)
		}
		contractJsonArr[i].Interface = strings.Replace(strings.Replace(abiStr, "\n", "", -1), "\\", "", -1)
		contractJsonArr[i].Bytecode = byteCode
		contractJsonArr[i].ContractAddress = contractAddress

		js := util.ComposeJSON(abiStr, byteCode, contractAddress)

		contractJsonArr[i].Filename = strings.Replace(fileName[i], ".sol", "", -1)
		contractJsonArr[i].Json = strings.Replace(strings.Replace(js, "\n", "", -1), "\\", "", -1)
		filePath := path + "/" + strings.Replace(fileName[i], ".sol", "", -1) + ".json"
		jsByte := []byte(js)
		err = ioutil.WriteFile(filePath, jsByte, 0775)
		if err != nil {
			panic(err)
		}
		abi := []byte(abiStr)
		err = ioutil.WriteFile(path+"/ABI", abi, 0775)
		if err != nil {
			panic(err)
		}
		bin := []byte(byteCode)
		err = ioutil.WriteFile(path+"/BIN", bin, 0775)
		if err != nil {
			panic(err)
		}
		contNameMap[contractAddress] = contractJsonArr[i].Filename
		abiMap[contractAddress] = abiStr
	}
	return contractJsonArr
}

func (nsi *NodeServiceImpl) createNetworkScriptCall(nodename string, currentIP string, rpcPort string, whisperPort string, constellationPort string, raftPort string, nodeManagerPort string) SuccessResponse {
	var successResponse SuccessResponse
	cmd := exec.Command("./setup.sh", "1", nodename)
	cmd.Dir = "./Setup"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	var setupConf string
	setupConf = "CURRENT_IP=" + currentIP + "\n" + "RPC_PORT=" + rpcPort + "\n" + "WHISPER_PORT=" + whisperPort + "\n" + "CONSTELLATION_PORT=" + constellationPort + "\n" + "RAFT_PORT=" + raftPort + "\n" + "NODEMANAGER_PORT=" + nodeManagerPort + "\n"
	setupConfByte := []byte(setupConf)
	err = ioutil.WriteFile("./Setup/"+nodename+"/setup.conf", setupConfByte, 0775)
	if err != nil {
		panic(err)
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) joinRequestResponseCall(nodename string, currentIP string, rpcPort string, whisperPort string, constellationPort string, raftPort string, nodeManagerPort string, masterNodeManagerPort string, masterIP string) SuccessResponse {
	var successResponse SuccessResponse
	cmd := exec.Command("./setup.sh", "2", nodename, masterIP, masterNodeManagerPort, currentIP, rpcPort, whisperPort, constellationPort, raftPort, nodeManagerPort)
	cmd.Dir = "./Setup"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	var setupConf string
	setupConf = "CURRENT_IP=" + currentIP + "\n" + "RPC_PORT=" + rpcPort + "\n" + "WHISPER_PORT=" + whisperPort + "\n" + "CONSTELLATION_PORT=" + constellationPort + "\n" + "RAFT_PORT=" + raftPort + "\n" + "THIS_NODEMANAGER_PORT=" + nodeManagerPort + "\n" + "MASTER_IP=" + masterIP + "\n" + "NODEMANAGER_PORT=" + masterNodeManagerPort + "\n"
	setupConfByte := []byte(setupConf)
	err = ioutil.WriteFile("./Setup/"+nodename+"/setup.conf", setupConfByte, 0775)
	if err != nil {
		panic(err)
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) resetCurrentNode() SuccessResponse {
	var successResponse SuccessResponse
	cmd := exec.Command("./reset_chain.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		successResponse.Status = "failure"
		return successResponse
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) restartCurrentNode() SuccessResponse {
	var successResponse SuccessResponse
	r, _ := regexp.Compile("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]")
	files, err := ioutil.ReadDir("/home/node")
	if err != nil {
		log.Fatal(err)
	}
	var filename string
	for _, f := range files {
		match, _ := regexp.MatchString("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]", f.Name())
		if match {
			filename = r.FindString(f.Name())
		}
	}
	filepath := fmt.Sprint("./", filename)
	cmd := exec.Command(filepath)
	cmd.Dir = "/home/node"
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
		successResponse.Status = "failure"
		return successResponse
	}
	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) latestBlockDetails(url string) LatestBlockResponse {
	var latestBlockResponse LatestBlockResponse
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	currentTime := time.Now().Unix()
	blockNumber := ethClient.BlockNumber()
	blockResponseClient := ethClient.GetBlockByNumber(blockNumber)
	blockNumberInt := util.HexStringtoInt64(blockNumber)
	creationTime := blockResponseClient.Timestamp
	creationTimeInt := util.HexStringtoInt64(creationTime)
	creationTimeUnix := creationTimeInt / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	latestBlockResponse.LatestBlockNumber = blockNumberInt
	latestBlockResponse.TimeElapsed = elapsedTime
	return latestBlockResponse
}

//func (nsi *NodeServiceImpl) latency(url string) ([]LatencyResponse) {
//	var nodeUrl = url
//	ethClient := client.EthClient{nodeUrl}
//	otherPeersResponse := ethClient.AdminPeers()
//	peerCount := len(otherPeersResponse)
//	latencyResponse := make([]LatencyResponse, peerCount+1)
//	for i := 0; i < peerCount+1; i++ {
//		var latOut bytes.Buffer
//		var ip string
//		if i == peerCount {
//			ip = "localhost"
//			thisAdminInfo := ethClient.AdminNodeInfo()
//			latencyResponse[i].EnodeID = thisAdminInfo.ID
//		} else {
//			ip = otherPeersResponse[i].Network.LocalAddress
//			ipString := strings.Split(ip, ":")
//			ip = ipString[0]
//			latencyResponse[i].EnodeID = otherPeersResponse[i].ID
//		}
//		command := fmt.Sprint("ping -c 4 ", ip, " |  awk -F'/' '{ print $5 }' | tail -1")
//		cmd := exec.Command("bash", "-c", command)
//		cmd.Stdout = &latOut
//		err := cmd.Run()
//		if err != nil {
//			fmt.Println(err)
//		}
//		latencyString := strings.TrimSuffix(latOut.String(), "\n")
//		latency, err := strconv.ParseFloat(latencyString, 10)
//		latency = latency * 1000
//		latencyStr := strconv.FormatFloat(latency, 'f', 0, 64)
//		latencyResponse[i].Latency = latencyStr
//	}
//	return latencyResponse
//}

func (nsi *NodeServiceImpl) latency(url string) []LatencyResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, contractAdd, "", nil}}

	peerNo := len(nms.GetNodeDetailsList())

	latencyResponse := make([]LatencyResponse, peerNo)
	for i := 0; i < peerNo; i++ {
		var latOut bytes.Buffer
		ip := nms.GetNodeDetails(i).IP
		latencyResponse[i].EnodeID = nms.GetNodeDetails(i).Enode
		command := fmt.Sprint("ping -c 4 ", ip, " |  awk -F'/' '{ print $5 }' | tail -1")
		cmd := exec.Command("bash", "-c", command)
		cmd.Stdout = &latOut
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		latencyString := strings.TrimSuffix(latOut.String(), "\n")
		latency, err := strconv.ParseFloat(latencyString, 10)
		latency = latency * 1000
		latencyStr := strconv.FormatFloat(latency, 'f', 0, 64)
		latencyResponse[i].Latency = latencyStr
	}
	return latencyResponse
}

func (nsi *NodeServiceImpl) transactionSearchDetails(txno string, url string) BlockDetailsResponse {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	txGetClient := ethClient.GetTransactionReceipt(txno)
	blockNumber := util.HexStringtoInt64(txGetClient.BlockNumber)
	blockDetailsResponse := nsi.getBlockInfo(blockNumber, url)
	return blockDetailsResponse
}

func (nsi *NodeServiceImpl) emailServerConfig(host string, port string, username string, password string, recipientList string, url string) SuccessResponse {
	var successResponse SuccessResponse

	mailServerConfig.Host = host
	mailServerConfig.Port = port
	mailServerConfig.Username = username
	mailServerConfig.Password = password
	mailServerConfig.RecipientList = recipientList

	registered := fmt.Sprint("RECIPIENTLIST=", recipientList, "\n")
	util.AppendStringToFile("/home/setup.conf", registered)

	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			//fmt.Println("Healthcheck done at: ", t)
			if warning > 0 {
				//fmt.Println("Ticker stopped")
				ticker.Stop()
			}
			nsi.healthCheck(url)

		}
	}()

	successResponse.Status = "success"
	return successResponse
}

func (nsi *NodeServiceImpl) healthCheck(url string) {
	ethClient := client.EthClient{url}
	blockNumber := ethClient.BlockNumber()
	if blockNumber == "" {
		if warning > 0 {
			exists := util.PropertyExists("RECIPIENTLIST", "/home/setup.conf")
			if exists != "" {
				p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
				recipientList := util.MustGetString("RECIPIENTLIST", p)
				recipients := strings.Split(recipientList, ",")

				b, err := ioutil.ReadFile("/root/quorum-maker/NodeUnavailableTemplate.txt")

				if err != nil {
					log.Fatal(err)
				}

				mailCont := string(b)
				mailCont = strings.Replace(mailCont, "\n", "", -1)
				for i := 0; i < len(recipients); i++ {
					nsi.sendMail(mailServerConfig.Host, mailServerConfig.Port, mailServerConfig.Username, mailServerConfig.Password, "Node is not responding", mailCont, recipients[i])
				}
			}
		}
		warning++
	}
}

func (nsi *NodeServiceImpl) sendMail(host string, port string, username string, password string, subject string, mailContent string, to string) {
	//fmt.Println("Called with warning =", warning)
	portNo, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", to)
	//m.SetAddressHeader("Cc", host, host)
	m.SetHeader("Subject", subject)
	m.SetBody("text", mailContent)
	//m.Attach("")

	d := gomail.NewDialer(host, int(portNo), username, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}

//func (nsi *NodeServiceImpl) logs() SuccessResponse {
//	var successResponse SuccessResponse
//	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
//	ipAddr := util.MustGetString("CURRENT_IP", p)
//	logPort := util.MustGetString("THIS_NODEMANAGER_PORT", p)
//	successResponse.Status = fmt.Sprint(ipAddr, ":", logPort)
//	return successResponse
//}

//@TODO: Implement logrotate command to do this.
func (nsi *NodeServiceImpl) LogRotaterGeth() {
	command := "cat $(ls | grep log | grep -v _) > Geth_$(date| sed -e 's/ /_/g')"

	command1 := "echo -en '' > $(ls | grep log | grep -v _)"

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = "/home/node/qdata/gethLogs"
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd1 := exec.Command("bash", "-c", command1)
	cmd1.Dir = "/home/node/qdata/gethLogs"
	err1 := cmd1.Run()
	if err1 != nil {
		fmt.Println(err)
	}

}

func (nsi *NodeServiceImpl) LogRotaterConst() {

	command := "cat $(ls | grep log | grep _) > Constellation_$(date| sed -e 's/ /_/g')"

	command1 := "echo -en '' > $(ls | grep log | grep _)"

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = "/home/node/qdata/constellationLogs"
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd1 := exec.Command("bash", "-c", command1)
	cmd1.Dir = "/home/node/qdata/constellationLogs"
	err1 := cmd1.Run()
	if err1 != nil {
		fmt.Println(err)
	}

}

func (nsi *NodeServiceImpl) RegisterNodeDetails(url string) {
	var nodeUrl = url
	var registeredVal string
	exists := util.PropertyExists("REGISTERED", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		registeredVal = util.MustGetString("REGISTERED", p)
	}
	if registeredVal != "TRUE" {
		ethClient := client.EthClient{nodeUrl}

		enode := ethClient.AdminNodeInfo().ID
		fromAddress := ethClient.Coinbase()
		var ipAddr, nodename, pubKey, role, id, contractAdd string
		existsA := util.PropertyExists("CURRENT_IP", "/home/setup.conf")
		existsB := util.PropertyExists("NODENAME", "/home/setup.conf")
		existsC := util.PropertyExists("PUBKEY", "/home/setup.conf")
		existsD := util.PropertyExists("ROLE", "/home/setup.conf")
		existsE := util.PropertyExists("RAFT_ID", "/home/setup.conf")
		existsF := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
		if existsA != "" && existsB != "" && existsC != "" && existsD != "" && existsE != "" && existsF != "" {
			p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
			ipAddr = util.MustGetString("CURRENT_IP", p)
			nodename = util.MustGetString("NODENAME", p)
			pubKey = util.MustGetString("PUBKEY", p)
			role = util.MustGetString("ROLE", p)
			id = util.MustGetString("RAFT_ID", p)
			contractAdd = util.MustGetString("CONTRACT_ADD", p)
		}
		//fmt.Println(ipAddr, nodename, pubKey, role, enode, fromAddress, contractAdd)
		registered := fmt.Sprint("REGISTERED=TRUE", "\n")
		util.AppendStringToFile("/home/setup.conf", registered)
		util.DeleteProperty("REGISTERED=", "/home/setup.conf")
		util.DeleteProperty("ROLE=Unassigned", "/home/setup.conf")
		nms := contractclient.NetworkMapContractClient{client.EthClient{url}, contracthandler.ContractParam{fromAddress, contractAdd, "", nil}}
		nms.RegisterNode(nodename, role, pubKey, enode, ipAddr, id)
	}
}

func (nsi *NodeServiceImpl) NetworkManagerContractDeployer(url string) {
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}
	if contractAdd == "" {
		log.Info("Deploying Network Manager Contract")
		filename := []string{"NetworkManagerContract.sol"}
		deployedContract := nsi.deployContract(nil, filename, false, url)
		contAdd := deployedContract[0].ContractAddress
		contAddAppend := fmt.Sprint("CONTRACT_ADD=", contAdd, "\n")
		util.AppendStringToFile("/home/setup.conf", contAddAppend)
		util.DeleteProperty("CONTRACT_ADD=", "/home/setup.conf")
	}
}

func ConvertToReadable(p client.TransactionDetailsResponse, pending bool, hash bool) TransactionDetailsResponse {
	var readableTransactionDetailsResponse TransactionDetailsResponse

	readableTransactionDetailsResponse.BlockNumber = util.HexStringtoInt64(p.BlockNumber)
	readableTransactionDetailsResponse.Gas = util.HexStringtoInt64(p.Gas)
	readableTransactionDetailsResponse.GasPrice = util.HexStringtoInt64(p.GasPrice)
	readableTransactionDetailsResponse.TransactionIndex = util.HexStringtoInt64(p.TransactionIndex)
	readableTransactionDetailsResponse.Value = util.HexStringtoInt64(p.Value)
	readableTransactionDetailsResponse.Nonce = util.HexStringtoInt64(p.Nonce)
	readableTransactionDetailsResponse.BlockHash = p.BlockHash
	readableTransactionDetailsResponse.From = p.From
	readableTransactionDetailsResponse.Hash = p.Hash
	readableTransactionDetailsResponse.Input = p.Input
	readableTransactionDetailsResponse.To = p.To
	readableTransactionDetailsResponse.V = p.V
	readableTransactionDetailsResponse.R = p.R
	readableTransactionDetailsResponse.S = p.S
	if util.HexStringtoInt64(p.V) == 37 || util.HexStringtoInt64(p.V) == 38 {
		if pending {
			readableTransactionDetailsResponse.TransactionType = "Private or Hash Only"
		} else if hash {
			readableTransactionDetailsResponse.TransactionType = "Hash Only"
		} else {
			readableTransactionDetailsResponse.TransactionType = "Private"
		}

	} else {
		readableTransactionDetailsResponse.TransactionType = "Public"
	}

	return readableTransactionDetailsResponse
}

func (nsi *NodeServiceImpl) CheckGethStatus(url string) bool {
	ethClient := client.EthClient{url}
	var coinbase string
	for coinbase == "" {
		time.Sleep(1 * time.Second)
		coinbase = ethClient.Coinbase()
	}
	return true
}

func (nsi *NodeServiceImpl) GetChartData(url string) []ChartInfo {
	ethClient := client.EthClient{url}
	chartResponse := make([]ChartInfo, chartSize)
	currentTimeRaw := time.Now().Unix()
	currentTime := currentTimeRaw - (currentTimeRaw % 60)
	startTime := currentTime
	currentBlockNumber := util.HexStringtoInt64(ethClient.BlockNumber())
	bucketTime := currentTime - 60
	stopTime := currentTime - int64(60*chartSize)
	i := 0
	lastBlockNoHex := strconv.FormatInt(currentBlockNumber, 16)
	lastBNoHex := fmt.Sprint("0x", lastBlockNoHex)
	blockResponseClient := ethClient.GetBlockByNumber(lastBNoHex)
	lastCreationTimeRaw := util.HexStringtoInt64(blockResponseClient.Timestamp)
	lastCreationTime := lastCreationTimeRaw / 1000000000
	lastCreationTimeSec := lastCreationTime - (lastCreationTime % 60)
	if lastCreationTimeSec > stopTime {
		for currentTime > stopTime {
			blockCount := 0
			txnCount := 0
			for currentTime > bucketTime {
				blockNoHex := strconv.FormatInt(currentBlockNumber, 16)
				bNoHex := fmt.Sprint("0x", blockNoHex)
				blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
				creationTimeRaw := util.HexStringtoInt64(blockResponseClient.Timestamp)
				creationTimeRaw = creationTimeRaw / 1000000000
				currentTime = creationTimeRaw - (creationTimeRaw % 60)
				if currentTime > bucketTime {
					currentBlockNumber = currentBlockNumber - 1
					txnCount = txnCount + len(blockResponseClient.Transactions)
					blockCount++
				}
			}
			chartResponse[i].BlockCount = blockCount
			chartResponse[i].TransactionCount = txnCount
			chartResponse[i].TimeStamp = (int(bucketTime) + 60) * 1000
			bucketTime = bucketTime - 60
			i++
		}
	}
	for i := 0; i < chartSize; i++ {

		if chartResponse[i].TimeStamp == 0 {
			chartResponse[i].TimeStamp = int(startTime) * 1000
		}
		startTime = startTime - 60
	}
	for i := 0; i < len(chartResponse)/2; i++ {
		j := len(chartResponse) - i - 1
		chartResponse[i], chartResponse[j] = chartResponse[j], chartResponse[i]
	}

	return chartResponse
}

func (nsi *NodeServiceImpl) ContractCrawler(url string) {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			if contractCrawlerMutex == 0 {
				getContracts(url)
			}
		}
	}()
}

func getContracts(url string) {
	contractCrawlerMutex = 1
	ethClient := client.EthClient{url}
	blockNumber := int(util.HexStringtoInt64(ethClient.BlockNumber()))
	for i := lastCrawledBlock + 1; i <= blockNumber; i++ {
		blockNoHex := strconv.FormatInt(int64(i), 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		for _, clientTransactions := range blockResponseClient.Transactions {
			txGetClient := ethClient.GetTransactionReceipt(clientTransactions.Hash)
			if txGetClient.ContractAddress != "" {
				if abiMap[txGetClient.ContractAddress] == "" {
					abiMap[txGetClient.ContractAddress] = "missing"
				}

				if util.HexStringtoInt64(clientTransactions.V) == 37 || util.HexStringtoInt64(clientTransactions.V) == 38 {
					private := ethClient.GetQuorumPayload(clientTransactions.Input)
					if private == "0x" {
						contTypeMap[txGetClient.ContractAddress] = "Hash Only"
					} else {
						contTypeMap[txGetClient.ContractAddress] = "Private"
					}
				} else {
					contTypeMap[txGetClient.ContractAddress] = "Public"

				}
				contSenderMap[txGetClient.ContractAddress] = clientTransactions.From
			}
		}
	}

	lastCrawledBlock = blockNumber
	contractCrawlerMutex = 0
}

func (nsi *NodeServiceImpl) ContractList() []ContractTableRow {
	contractList := make([]ContractTableRow, len(abiMap))
	i := 0
	for key := range abiMap {
		contractList[i].ContractAdd = key
		contractList[i].ABIContent = abiMap[key]
		if abiMap[key] == "missing" {
			contractList[i].ABIContent = ""
		}
		contractList[i].ContractName = contNameMap[key]
		contractList[i].ContractType = contTypeMap[key]
		contractList[i].Sender = contSenderMap[key]
		contractList[i].Description = contDescriptionMap[key]
		i++
	}

	return contractList
}

func (nsi *NodeServiceImpl) ContractCount() ContractCounter {
	availableABIs := 0
	totalContracts := 0
	for key := range abiMap {
		if abiMap[key] != "missing" {
			availableABIs++
		}
		totalContracts++
	}
	var contractCount ContractCounter
	contractCount.TotalContracts = totalContracts
	contractCount.ABIavailable = availableABIs
	return contractCount
}

func (nsi *NodeServiceImpl) updateContractDetails(contractAddress string, contractName string, abi string, description string) SuccessResponse {
	var successResponse SuccessResponse
	contNameMap[contractAddress] = contractName
	abiMap[contractAddress] = abi
	contDescriptionMap[contractAddress] = description
	successResponse.Status = "Successfully updated contract details"
	return successResponse
}
