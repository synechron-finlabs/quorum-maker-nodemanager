package service

import (
	"io/ioutil"
	"log"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"strings"
	"fmt"
	"strconv"
	"github.com/magiconair/properties"
	"bytes"
	"os/exec"
	"regexp"
	"os"
	"time"
	"gopkg.in/gomail.v2"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
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
	BlockHash         string `json:"blockHash"`
	BlockNumber       int64  `json:"blockNumber"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed int64  `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	Gas               int64  `json:"gas"`
	GasPrice          int64  `json:"gasPrice"`
	GasUsed           int64  `json:"gasUsed"`
	Input             string `json:"input"`
	Logs              []Logs `json:"logs"`
	LogsBloom         string `json:"logsBloom"`
	Nonce             int64  `json:"nonce"`
	Root              string `json:"root"`
	To                string `json:"to"`
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  int64  `json:"transactionIndex"`
	Value             int64  `json:"value"`
	V                 string `json:"v"`
	R                 string `json:"r"`
	S                 string `json:"s"`
	TransactionType   string `json:"transactionType"`
	TimeElapsed       int64  `json:"TimeElapsed"`
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
};

type LatencyResponse struct {
	EnodeID string `json:"enode-id"`
	Latency string `json:"latency"`
}

type NodeServiceImpl struct {
	Url string
}

var warning = 0
var mailServerConfig MailServerConfig

func (nsi *NodeServiceImpl) getGenesis(url string) (response GetGenesisResponse) {
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	netId := util.MustGetString("NETWORK_ID", p)
	constl := util.MustGetString("CONSTELLATION_PORT", p)

	b, err := ioutil.ReadFile("/home/node/genesis.json")
	if err != nil {
		log.Fatal(err)
	}
	genesis := string(b)
	genesis = strings.Replace(genesis, "\n", "", -1)

	response = GetGenesisResponse{constl, netId, genesis}
	return response
}

func (nsi *NodeServiceImpl) joinNetwork(enode string, url string) (string) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	raftId := ethClient.RaftAddPeer(enode)
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	collatedInfo := fmt.Sprint(raftId, ":", contractAdd)
	return collatedInfo
}

//@TODO: If this function is repeatedly called from UI, please cache the static informations.
func (nsi *NodeServiceImpl) getCurrentNode(url string) (NodeInfo) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()
	nms := contractclient.NetworkMapContractClient{EthClient: client.EthClient{url}}
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	totalCount := len(nms.GetNodeDetailsList(fromAddress, contractAdd, "", nil))
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

	ipAddr := util.MustGetString("CURRENT_IP", p)
	raftId := util.MustGetString("RAFT_ID", p)
	rpcPort := util.MustGetString("RPC_PORT", p)

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
	responseObj := NodeInfo{nodename, count, totalCount, activeStatus, conn, raftRole, raftIdInt, blockNumberInt, pendingTxCount, genesis, thisAdminInfo}
	return responseObj
}

func (nsi *NodeServiceImpl) getOtherPeer(peerId string, url string) (client.AdminPeers) {
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

func (nsi *NodeServiceImpl) getOtherPeers(url string) ([]client.AdminPeers) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	otherPeersResponse := ethClient.AdminPeers()
	return otherPeersResponse
}

func (nsi *NodeServiceImpl) getPendingTransactions(url string) ([]TransactionDetailsResponse) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	pendingTxResponseClient := ethClient.PendingTransactions()
	pendingTxCount := len(pendingTxResponseClient)
	pendingTxResponse := make([]TransactionDetailsResponse, pendingTxCount)
	for i := 0; i < pendingTxCount; i++ {

		//@TODO: Create a Utility function to convert transaction to readable transaction object and call from here
		pendingTxResponse[i].BlockNumber = util.HexStringtoInt64(pendingTxResponseClient[i].BlockNumber)
		pendingTxResponse[i].Gas = util.HexStringtoInt64(pendingTxResponseClient[i].Gas)
		pendingTxResponse[i].GasPrice = util.HexStringtoInt64(pendingTxResponseClient[i].GasPrice)
		pendingTxResponse[i].TransactionIndex = util.HexStringtoInt64(pendingTxResponseClient[i].TransactionIndex)
		pendingTxResponse[i].Value = util.HexStringtoInt64(pendingTxResponseClient[i].Value)
		pendingTxResponse[i].Nonce = util.HexStringtoInt64(pendingTxResponseClient[i].Nonce)
		pendingTxResponse[i].BlockHash = pendingTxResponseClient[i].BlockHash
		pendingTxResponse[i].From = pendingTxResponseClient[i].From
		pendingTxResponse[i].Hash = pendingTxResponseClient[i].Hash
		pendingTxResponse[i].Input = pendingTxResponseClient[i].Input
		pendingTxResponse[i].To = pendingTxResponseClient[i].To
		pendingTxResponse[i].V = pendingTxResponseClient[i].V
		pendingTxResponse[i].R = pendingTxResponseClient[i].R
		pendingTxResponse[i].S = pendingTxResponseClient[i].S
		if util.HexStringtoInt64(pendingTxResponseClient[i].V) == 37 || util.HexStringtoInt64(pendingTxResponseClient[i].V) == 38 {
			pendingTxResponse[i].TransactionType = "Private or Hash Only"
		} else {
			pendingTxResponse[i].TransactionType = "Public"
		}
	}
	return pendingTxResponse
}

func (nsi *NodeServiceImpl) getBlockInfo(blockno int64, url string) (BlockDetailsResponse) {
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
	for i := 0; i < txnNo; i++ {

		//@TODO: call the utility function to convert to readeable object
		txGetClient := ethClient.GetTransactionReceipt(blockResponseClient.Transactions[i].Hash)
		txResponse[i].BlockNumber = util.HexStringtoInt64(blockResponseClient.Transactions[i].BlockNumber)
		txResponse[i].Gas = util.HexStringtoInt64(blockResponseClient.Transactions[i].Gas)
		txResponse[i].GasPrice = util.HexStringtoInt64(blockResponseClient.Transactions[i].GasPrice)
		txResponse[i].TransactionIndex = util.HexStringtoInt64(blockResponseClient.Transactions[i].TransactionIndex)
		txResponse[i].Value = util.HexStringtoInt64(blockResponseClient.Transactions[i].Value)
		txResponse[i].Nonce = util.HexStringtoInt64(blockResponseClient.Transactions[i].Nonce)
		txResponse[i].BlockHash = blockResponseClient.Transactions[i].BlockHash
		txResponse[i].From = blockResponseClient.Transactions[i].From
		txResponse[i].Hash = blockResponseClient.Transactions[i].Hash
		txResponse[i].Input = blockResponseClient.Transactions[i].Input
		txResponse[i].To = blockResponseClient.Transactions[i].To
		txResponse[i].V = blockResponseClient.Transactions[i].V
		txResponse[i].R = blockResponseClient.Transactions[i].R
		txResponse[i].S = blockResponseClient.Transactions[i].S
		if util.HexStringtoInt64(txResponse[i].V) == 37 || util.HexStringtoInt64(txResponse[i].V) == 38 {
			if len(txGetClient.Logs) == 0 {
				txResponse[i].TransactionType = "Hash only"
			} else {
				txResponse[i].TransactionType = "Private"
			}
		} else {
			txResponse[i].TransactionType = "Public"
		}
	}
	blockResponse.Transactions = txResponse
	return blockResponse
}

func (nsi *NodeServiceImpl) getLatestBlockInfo(count string, reference string, url string) ([]BlockDetailsResponse) {
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
		for i := 0; i < txnNo; i++ {
			//@TODO: Call the utility function to convert to readable transaction object
			txGetClient := ethClient.GetTransactionReceipt(blockResponseClient.Transactions[i].Hash)
			txResponse[i].BlockNumber = util.HexStringtoInt64(blockResponseClient.Transactions[i].BlockNumber)
			txResponse[i].Gas = util.HexStringtoInt64(blockResponseClient.Transactions[i].Gas)
			txResponse[i].GasPrice = util.HexStringtoInt64(blockResponseClient.Transactions[i].GasPrice)
			txResponse[i].TransactionIndex = util.HexStringtoInt64(blockResponseClient.Transactions[i].TransactionIndex)
			txResponse[i].Value = util.HexStringtoInt64(blockResponseClient.Transactions[i].Value)
			txResponse[i].Nonce = util.HexStringtoInt64(blockResponseClient.Transactions[i].Nonce)
			txResponse[i].BlockHash = blockResponseClient.Transactions[i].BlockHash
			txResponse[i].From = blockResponseClient.Transactions[i].From
			txResponse[i].Hash = blockResponseClient.Transactions[i].Hash
			txResponse[i].Input = blockResponseClient.Transactions[i].Input
			txResponse[i].To = blockResponseClient.Transactions[i].To
			txResponse[i].V = blockResponseClient.Transactions[i].V
			txResponse[i].R = blockResponseClient.Transactions[i].R
			txResponse[i].S = blockResponseClient.Transactions[i].S
			if util.HexStringtoInt64(txResponse[i].V) == 37 || util.HexStringtoInt64(txResponse[i].V) == 38 {
				if len(txGetClient.Logs) == 0 {
					txResponse[i].TransactionType = "Hash Only"
				} else {
					txResponse[i].TransactionType = "Private"
				}
			} else {
				txResponse[i].TransactionType = "Public"
			}
		}
		blockResponse[blockNumber-i].Transactions = txResponse
	}
	return blockResponse
}

func (nsi *NodeServiceImpl) getLatestTransactionInfo(count string, url string) ([]BlockDetailsResponse) {
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
		for i := 0; i < txnNo; i++ {
			txGetClient := ethClient.GetTransactionReceipt(blockResponseClient.Transactions[i].Hash)
			txResponse[i].BlockNumber = util.HexStringtoInt64(blockResponseClient.Transactions[i].BlockNumber)
			txResponse[i].Gas = util.HexStringtoInt64(blockResponseClient.Transactions[i].Gas)
			txResponse[i].GasPrice = util.HexStringtoInt64(blockResponseClient.Transactions[i].GasPrice)
			txResponse[i].TransactionIndex = util.HexStringtoInt64(blockResponseClient.Transactions[i].TransactionIndex)
			txResponse[i].Value = util.HexStringtoInt64(blockResponseClient.Transactions[i].Value)
			txResponse[i].Nonce = util.HexStringtoInt64(blockResponseClient.Transactions[i].Nonce)
			txResponse[i].BlockHash = blockResponseClient.Transactions[i].BlockHash
			txResponse[i].From = blockResponseClient.Transactions[i].From
			txResponse[i].Hash = blockResponseClient.Transactions[i].Hash
			txResponse[i].Input = blockResponseClient.Transactions[i].Input
			txResponse[i].To = blockResponseClient.Transactions[i].To
			txResponse[i].V = blockResponseClient.Transactions[i].V
			txResponse[i].R = blockResponseClient.Transactions[i].R
			txResponse[i].S = blockResponseClient.Transactions[i].S
			if util.HexStringtoInt64(txResponse[i].V) == 37 || util.HexStringtoInt64(txResponse[i].V) == 38 {
				if len(txGetClient.Logs) == 0 {
					txResponse[i].TransactionType = "Hash Only"
				} else {
					txResponse[i].TransactionType = "Private"
				}
			} else {
				txResponse[i].TransactionType = "Public"
			}
		}
		blockResponse[blockNumber-i].Transactions = txResponse
	}
	return blockResponse
}

func (nsi *NodeServiceImpl) getTransactionInfo(txno string, url string) (TransactionDetailsResponse) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	txGetClient := ethClient.GetTransactionReceipt(txno)
	var txResponse TransactionDetailsResponse
	txResponseClient := ethClient.GetTransactionByHash(txno)
	txResponse.BlockNumber = util.HexStringtoInt64(txResponseClient.BlockNumber)
	txResponse.Gas = util.HexStringtoInt64(txResponseClient.Gas)
	txResponse.GasPrice = util.HexStringtoInt64(txResponseClient.GasPrice)
	txResponse.TransactionIndex = util.HexStringtoInt64(txResponseClient.TransactionIndex)
	txResponse.Value = util.HexStringtoInt64(txResponseClient.Value)
	txResponse.Nonce = util.HexStringtoInt64(txResponseClient.Nonce)
	txResponse.BlockHash = txResponseClient.BlockHash
	txResponse.From = txResponseClient.From
	txResponse.Hash = txResponseClient.Hash
	txResponse.Input = txResponseClient.Input
	txResponse.To = txResponseClient.To
	txResponse.V = txResponseClient.V
	txResponse.R = txResponseClient.R
	txResponse.S = txResponseClient.S
	if util.HexStringtoInt64(txResponseClient.V) == 37 || util.HexStringtoInt64(txResponseClient.V) == 38 {
		if len(txGetClient.Logs) == 0 {
			txResponse.TransactionType = "Hash Only"
		} else {
			txResponse.TransactionType = "Private"
		}
	} else {
		txResponse.TransactionType = "Public"
	}
	blockResponseClient := ethClient.GetBlockByNumber(txResponseClient.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txResponse.TimeElapsed = elapsedTime
	return txResponse
}

func (nsi *NodeServiceImpl) getTransactionReceipt(txno string, url string) (TransactionReceiptResponse) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	txGetClient := ethClient.GetTransactionByHash(txno)
	var txResponse TransactionReceiptResponse
	txResponseClient := ethClient.GetTransactionReceipt(txno)
	txResponse.BlockNumber = util.HexStringtoInt64(txResponseClient.BlockNumber)
	txResponse.CumulativeGasUsed = util.HexStringtoInt64(txResponseClient.CumulativeGasUsed)
	txResponse.GasUsed = util.HexStringtoInt64(txResponseClient.GasUsed)
	txResponse.TransactionIndex = util.HexStringtoInt64(txResponseClient.TransactionIndex)
	txResponse.BlockHash = txResponseClient.BlockHash
	txResponse.From = txResponseClient.From
	txResponse.ContractAddress = txResponseClient.ContractAddress
	txResponse.LogsBloom = txResponseClient.LogsBloom
	txResponse.Root = txResponseClient.Root
	txResponse.To = txResponseClient.To
	txResponse.TransactionHash = txResponseClient.TransactionHash
	txResponse.Gas = util.HexStringtoInt64(txGetClient.Gas)
	txResponse.GasPrice = util.HexStringtoInt64(txGetClient.GasPrice)
	txResponse.Input = txGetClient.Input
	txResponse.Nonce = util.HexStringtoInt64(txGetClient.Nonce)
	txResponse.Value = util.HexStringtoInt64(txGetClient.Value)
	txResponse.V = txGetClient.V
	txResponse.R = txGetClient.R
	txResponse.S = txGetClient.S
	eventNo := len(txResponseClient.Logs)
	if util.HexStringtoInt64(txGetClient.V) == 37 || util.HexStringtoInt64(txGetClient.V) == 38 {
		if eventNo == 0 {
			txResponse.TransactionType = "Hash Only"
		} else {
			txResponse.TransactionType = "Private"
		}
	} else {
		txResponse.TransactionType = "Public"
	}
	txResponseBuffer := make([]Logs, eventNo)
	for i := 0; i < eventNo; i++ {
		txResponseBuffer[i].BlockNumber = util.HexStringtoInt64(txResponseClient.Logs[i].BlockNumber)
		txResponseBuffer[i].LogIndex = util.HexStringtoInt64(txResponseClient.Logs[i].LogIndex)
		txResponseBuffer[i].TransactionIndex = util.HexStringtoInt64(txResponseClient.Logs[i].TransactionIndex)
		txResponseBuffer[i].Address = txResponseClient.Logs[i].Address
		txResponseBuffer[i].BlockHash = txResponseClient.Logs[i].BlockHash
		txResponseBuffer[i].Data = txResponseClient.Logs[i].Data
		txResponseBuffer[i].TransactionHash = txResponseClient.Logs[i].TransactionHash
		txResponseBuffer[i].Topics = txResponseClient.Logs[i].Topics
	}
	txResponse.Logs = txResponseBuffer
	blockResponseClient := ethClient.GetBlockByNumber(txResponseClient.BlockNumber)
	currentTime := time.Now().Unix()
	creationTime := util.HexStringtoInt64(blockResponseClient.Timestamp)
	creationTimeUnix := creationTime / 1000000000
	elapsedTime := currentTime - creationTimeUnix
	txResponse.TimeElapsed = elapsedTime
	return txResponse
}

func (nsi *NodeServiceImpl) joinRequestResponse(enode string, status string) (SuccessResponse) {
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
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	nms := contractclient.NetworkMapContractClient{EthClient: client.EthClient{url}}
	if private == true && pubKeys[0] == "" {
		enode := ethClient.AdminNodeInfo().ID
		peerNo := len(nms.GetNodeDetailsList(fromAddress, contractAdd, "", nil))
		publicKeys := make([]string, peerNo-1)
		for i := 0; i < peerNo; i++ {
			if enode != nms.GetNodeDetails(i, fromAddress, contractAdd, "", nil).Enode {
				publicKeys[i-1] = nms.GetNodeDetails(i, fromAddress, contractAdd, "", nil).PublicKey
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

		path := "./" + contractAddress + "_" + strings.Replace(fileName[i], ".sol", "", -1)
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
	}
	return contractJsonArr
}

func (nsi *NodeServiceImpl) createNetworkScriptCall(nodename string, currentIP string, rpcPort string, whisperPort string, constellationPort string, raftPort string, nodeManagerPort string) (SuccessResponse) {
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

func (nsi *NodeServiceImpl) joinRequestResponseCall(nodename string, currentIP string, rpcPort string, whisperPort string, constellationPort string, raftPort string, nodeManagerPort string, masterNodeManagerPort string, masterIP string) (SuccessResponse) {
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

func (nsi *NodeServiceImpl) resetCurrentNode() (SuccessResponse) {
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

func (nsi *NodeServiceImpl) restartCurrentNode() (SuccessResponse) {
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

func (nsi *NodeServiceImpl) latestBlockDetails(url string) (LatestBlockResponse) {
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

func (nsi *NodeServiceImpl) latency(url string) ([]LatencyResponse) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	fromAddress := ethClient.Coinbase()
	nms := contractclient.NetworkMapContractClient{EthClient: client.EthClient{url}}
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	peerNo := len(nms.GetNodeDetailsList(fromAddress, contractAdd, "", nil))

	latencyResponse := make([]LatencyResponse, peerNo)
	for i := 0; i < peerNo; i++ {
		var latOut bytes.Buffer
		ip := nms.GetNodeDetails(i, fromAddress, contractAdd, "", nil).IP
		latencyResponse[i].EnodeID = nms.GetNodeDetails(i, fromAddress, contractAdd, "", nil).Enode
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

func (nsi *NodeServiceImpl) transactionSearchDetails(txno string, url string) (BlockDetailsResponse) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	txGetClient := ethClient.GetTransactionReceipt(txno)
	blockNumber := util.HexStringtoInt64(txGetClient.BlockNumber)
	blockDetailsResponse := nsi.getBlockInfo(blockNumber, url)
	return blockDetailsResponse
}

func (nsi *NodeServiceImpl) emailServerConfig(host string, port string, username string, password string, recipientList string, url string) (SuccessResponse) {
	var successResponse SuccessResponse

	mailServerConfig.Host = host
	mailServerConfig.Port = port
	mailServerConfig.Username = username
	mailServerConfig.Password = password
	mailServerConfig.RecipientList = recipientList

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
			recipients := strings.Split(mailServerConfig.RecipientList, ",")
			for i := 0; i < len(recipients); i++ {
				nsi.sendMail(mailServerConfig.Host, mailServerConfig.Port, mailServerConfig.Username, mailServerConfig.Password, "Node is not responding", "Unfortunately this node has stopped responding", recipients[i])
			}
		}
		warning ++
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

func (nsi *NodeServiceImpl) logs() (SuccessResponse) {
	var successResponse SuccessResponse
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	ipAddr := util.MustGetString("CURRENT_IP", p)
	logPort := util.MustGetString("NODEMANAGER_PORT", p)
	successResponse.Status = fmt.Sprint(ipAddr, ":", logPort)
	return successResponse
}

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
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	registeredVal := util.MustGetString("REGISTERED", p)
	if registeredVal != "TRUE" {
		ethClient := client.EthClient{nodeUrl}
		nms := contractclient.NetworkMapContractClient{EthClient: client.EthClient{url}}
		enode := ethClient.AdminNodeInfo().ID
		fromAddress := ethClient.Coinbase()
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		ipAddr := util.MustGetString("CURRENT_IP", p)
		nodename := util.MustGetString("NODENAME", p)
		pubKey := util.MustGetString("PUBKEY", p)
		role := util.MustGetString("ROLE", p)
		id := util.MustGetString("RAFT_ID", p)
		contractAdd := util.MustGetString("CONTRACT_ADD", p)
		//fmt.Println(ipAddr, nodename, pubKey, role, enode, fromAddress, contractAdd)
		registered := fmt.Sprint("REGISTERED=TRUE", "\n")
		util.AppendStringToFile("/home/setup.conf", registered)
		nms.RegisterNode(nodename, role, pubKey, enode, ipAddr, id, fromAddress, contractAdd, "", nil)
	}
}

func (nsi *NodeServiceImpl) NetworkManagerContractDeployer(url string) {
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	if contractAdd == "" {
		filename := []string{"NetworkManagerContract.sol"}
		deployedContract := nsi.deployContract(nil, filename, false, url)
		contAdd := deployedContract[0].ContractAddress
		contAddAppend := fmt.Sprint("CONTRACT_ADD=", contAdd, "\n")
		util.AppendStringToFile("/home/setup.conf", contAddAppend)
	}
}
