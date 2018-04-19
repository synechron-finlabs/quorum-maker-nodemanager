package service

import (
	"io/ioutil"
	"log"
	"synechron.com/NodeManagerGo/client"
	"synechron.com/NodeManagerGo/util"
	"strings"
	"fmt"
	"strconv"
	"github.com/magiconair/properties"
	"bytes"
	"os/exec"
	"regexp"
	"os"
)

type ConnectionInfo struct {
	IP 	string 	`json:"ip,omitempty"`
	Port 	int 	`json:"port,omitempty"`
	Enode 	string 	`json:"enode,omitempty"`
}

type NodeInfo struct {
	NodeName	string		        `json:"nodeName,omitempty"`
	NodeCount	int			`json:"nodeCount,omitempty"`
	ConnectionInfo  ConnectionInfo		`json:"connectionInfo,omitempty"`
	RaftRole 	string 			`json:"raftRole,omitempty"`
	RaftID 		int      		`json:"raftID,omitempty"`
	BlockNumber 	int64			`json:"blockNumber"`
	PendingTxCount 	int 			`json:"pendingTxCount"`
	Genesis 	string			`json:"genesis,omitempty"`
	AdminInfo	client.AdminInfo	`json:"adminInfo,omitempty"`
}

type JoinNetworkRequest struct {
	EnodeID    string `json:"enode-id,omitempty"`
	IPAddress  string `json:"ip-address,omitempty"`
}

type GetGenesisResponse struct {
	ContstellationPort string `json:"contstellation-port, omitempty"`
	NetID              string `json:"netID,omitempty"`
	Genesis            string `json:"genesis, omitempty"`
}

type BlockDetailsResponse struct {
	Number           int64                                  `json:"number,omitempty"`
	Hash             string                                 `json:"hash,omitempty"`
	ParentHash       string                                 `json:"parentHash,omitempty"`
	Nonce            string                                 `json:"nonce,omitempty"`
	Sha3Uncles       string                                 `json:"sha3Uncles,omitempty"`
	LogsBloom        string                                 `json:"logsBloom,omitempty"`
	TransactionsRoot string                                 `json:"transactionsRoot,omitempty"`
	StateRoot        string                                 `json:"stateRoot,omitempty"`
	Miner            string                                 `json:"miner,omitempty"`
	Difficulty       int64                                  `json:"difficulty,omitempty"`
	TotalDifficulty  int64                                  `json:"totalDifficulty,omitempty"`
	ExtraData        string                                 `json:"extraData,omitempty"`
	Size             int64                                  `json:"size,omitempty"`
	GasLimit         int64                                  `json:"gasLimit,omitempty"`
	GasUsed          int64                                  `json:"gasUsed,omitempty"`
	Timestamp        int64                                  `json:"timestamp,omitempty"`
	Transactions     []TransactionDetailsResponse           `json:"transactions,omitempty"`
	Uncles           []string                               `json:"uncles,omitempty"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash,omitempty"`
	BlockNumber      int64  `json:"blockNumber"`
	From             string `json:"from,omitempty"`
	Gas              int64  `json:"gas,omitempty"`
	GasPrice         int64  `json:"gasPrice"`
	Hash             string `json:"hash,omitempty"`
	Input            string `json:"input,omitempty"`
	Nonce            int64  `json:"nonce"`
	To               string `json:"to,omitempty"`
	TransactionIndex int64  `json:"transactionIndex"`
	Value            int64  `json:"value,omitempty"`
	V                string `json:"v,omitempty"`
	R                string `json:"r,omitempty"`
	S                string `json:"s,omitempty"`
	TransactionType	 string	`json:"transactionType,omitempty"`
}

type TransactionReceiptResponse struct {
	BlockHash         	string      	        `json:"blockHash"`
	BlockNumber       	int64      		`json:"blockNumber"`
	ContractAddress   	string 			`json:"contractAddress"`
	CumulativeGasUsed 	int64      		`json:"cumulativeGasUsed"`
	From              	string      	        `json:"from"`
	GasUsed           	int64      		`json:"gasUsed"`
	Logs              	[]Logs 			`json:"logs"`
	LogsBloom        	string 			`json:"logsBloom"`
	Root             	string 			`json:"root"`
	To               	string 			`json:"to"`
	TransactionHash  	string 			`json:"transactionHash"`
	TransactionIndex 	int64 			`json:"transactionIndex"`
	TransactionType	 	string			`json:"transactionType,omitempty"`
}

type Logs struct {
	Address         	string    	`json:"address"`
	BlockHash       	string  	`json:"blockHash"`
	BlockNumber   		int64    	`json:"blockNumber"`
	Data 			string    	`json:"data"`
	LogIndex          	int64      	`json:"logIndex"`
	Topics          	[]string   	`json:"topics"`
	TransactionHash   	string 		`json:"transactionHash"`
	TransactionIndex  	int64      	`json:"transactionIndex"`
}

type JoinNetworkResponse struct {
	EnodeID string `json:"enode-id,omitempty"`
	Status 	string `json:"status,omitempty"`
}

type ContractJson struct {
	Interface       string `json:"interface"`
	Bytecode        string `json:"bytecode"`
	ContractAddress string `json:"address"`
}

type CreateNetworkScriptArgs struct {
    Nodename            string                                 `json:"nodename,omitempty"`
    CurrentIP           string                                 `json:"currentIP,omitempty"`
    RPCPort             string                                 `json:"rpcPort,omitempty"`
    WhisperPort         string                                 `json:"whisperPort,omitempty"`
    ConstellationPort   string                                 `json:"constellationPort,omitempty"`
    RaftPort            string                                 `json:"raftPort,omitempty"`
    NodeManagerPort     string                                 `json:"nodeManagerPort,omitempty"`
}

type JoinNetworkScriptArgs struct {
    Nodename                string                                 `json:"nodename,omitempty"`
    CurrentIP               string                                 `json:"currentIP,omitempty"`
    RPCPort                 string                                 `json:"rpcPort,omitempty"`
    WhisperPort             string                                 `json:"whisperPort,omitempty"`
    ConstellationPort       string                                 `json:"constellationPort,omitempty"`
    RaftPort                string                                 `json:"raftPort,omitempty"`
    NodeManagerPort         string                                 `json:"nodeManagerPort,omitempty"`
    MasterNodeManagerPort   string                                 `json:"masterNodeManagerPort,omitempty"`
    MasterIP                string                                 `json:"masterIP,omitempty"`
}


type SuccessResponse struct {
	Status 	string `json:"status,omitempty"`
}

type NodeServiceImpl struct {
	Url string
}


func (nsi *NodeServiceImpl) getGenesis(url string) (response GetGenesisResponse) {
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	netId := util.MustGetString("NETWORK_ID", p)
	constl := util.MustGetString("CONSTELLATION_PORT", p)
	
	b, err := ioutil.ReadFile("/home/node/genesis.json")
	if err != nil {
		log.Fatal(err)
	}
	genesis := string(b)
	genesis = strings.Replace(genesis, "\n","",-1)

	response = GetGenesisResponse{constl, netId, genesis}
	return response
}


func (nsi *NodeServiceImpl) joinNetwork(enode string, url string) (int) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	raftId := ethClient.RaftAddPeer(enode)
	return raftId
}


func (nsi *NodeServiceImpl) getCurrentNode (url string) (NodeInfo) {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}
	otherPeersResponse := ethClient.AdminPeers()
	count := len(otherPeersResponse)
	count = count + 1
	r, _ := regexp.Compile("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]")
	files, err := ioutil.ReadDir("/home/node")
	if err != nil {
		log.Fatal(err)
	}
	var nodename string
	for _, f := range files {
		match, _ := regexp.MatchString("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]", f.Name())
		if(match) {
			nodename = r.FindString(f.Name())
		}
	}
	nodename = strings.TrimSuffix(nodename, ".sh")
	nodename = strings.TrimPrefix(nodename, "start_")
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
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
	blockNumberInt :=  util.HexStringtoInt64(blockNumber)

	raftRole := ethClient.RaftRole()

	raftRole = strings.TrimSuffix(raftRole, "\n")

	b, err := ioutil.ReadFile("/home/node/genesis.json")

	if err != nil {
		log.Fatal(err)
	}

	genesis := string(b)
	genesis = strings.Replace(genesis, "\n","",-1)
	conn := ConnectionInfo{ipAddr,rpcPortInt,enode}
	responseObj := NodeInfo{nodename,count,conn,raftRole,raftIdInt,blockNumberInt,pendingTxCount,genesis,thisAdminInfo}
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
	blockNoHex  := strconv.FormatInt(blockno, 16)
	bNoHex := fmt.Sprint("0x", blockNoHex)
	var blockResponse BlockDetailsResponse
	blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
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


func (nsi *NodeServiceImpl) getLatestBlockInfo(count string,reference string, url string) ([]BlockDetailsResponse) {
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
	for i :=start ; i <= blockNumber; i++ {
		blockNoHex  := strconv.FormatInt(i, 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		blockResponse[blockNumber - i].Number = util.HexStringtoInt64(blockResponseClient.Number)
		blockResponse[blockNumber - i].Hash = blockResponseClient.Hash
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
		blockResponse[blockNumber - i].Transactions = txResponse
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
	for i :=start ; i <= blockNumber; i++ {
		blockNoHex  := strconv.FormatInt(i, 16)
		bNoHex := fmt.Sprint("0x", blockNoHex)
		blockResponseClient := ethClient.GetBlockByNumber(bNoHex)
		blockResponse[blockNumber - i].Number = util.HexStringtoInt64(blockResponseClient.Number)
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
		blockResponse[blockNumber - i].Transactions = txResponse
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
	return txResponse
}


func (nsi *NodeServiceImpl) joinRequestResponse(enode string, status string) (SuccessResponse) {
	var successResponse SuccessResponse
	peerMap[enode] = status
	successResponse.Status = fmt.Sprintf("Successfully updated status of %s to %s",enode, status)
	return successResponse
}


func (nsi *NodeServiceImpl) deployContract(pubKeys []string, fileName []string, private bool, url string) []ContractJson {
	var nodeUrl = url
	ethClient := client.EthClient{nodeUrl}

	if private == true && pubKeys[0] == "" {
		pubKeys = []string{"F/vdZBFpbIzi7vyRCRba0jvvEpYHGeZdBbKYwIiy1SE=","IgD5KZV+kZBhxfIReFR24JDQZVtz3UBeA+llw8vLpT4=","E/vdZBFpbIzi7vyRCRba0jvvEpYHGeZdBbKYwIiy1SF="}
	}

	var solc string
	if solc == "" {
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
		byteCode = byteCode[start:len(byteCode)]
		byteCode = "0x" + byteCode

		abiStr := abiOut.String()

		re = regexp.MustCompile("ABI?")
		index = re.FindStringIndex(abiStr)
		start = index[1] + 2
		abiStr = abiStr[start:len(abiStr)]

		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		byteCode = reg.ReplaceAllString(byteCode, "")

		contractAddress := ethClient.DeployContracts(byteCode, pubKeys, private)

		path := "./" + contractAddress
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0775)
		}

		contractJsonArr[i].Interface = abiStr
		contractJsonArr[i].Bytecode = byteCode
		contractJsonArr[i].ContractAddress = contractAddress

		js := util.ComposeJSON(abiStr, byteCode, contractAddress)

		filePath := path + "/" + strings.Replace(fileName[i], ".sol","",-1) + ".json"
		jsByte := []byte(js)
		err = ioutil.WriteFile(filePath, jsByte, 0775)
		if err != nil {
			panic(err)
		}
		abi := []byte(abiStr)
		err = ioutil.WriteFile(path + "/ABI", abi, 0775)
		if err != nil {
			panic(err)
		}
		bin := []byte(byteCode)
		err = ioutil.WriteFile(path + "/BIN", bin, 0775)
		if err != nil {
			panic(err)
		}
	}
	return contractJsonArr
}


func (nsi *NodeServiceImpl) createNetworkScriptCall(nodename string, currentIP string, rpcPort string, whisperPort string, constellationPort string, raftPort string, nodeManagerPort string) (SuccessResponse) {
	var successResponse SuccessResponse
	//cmd := exec.Command("./setup.sh","1", nodename)
	//cmd.Dir = "./Setup"
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var setupConf string
	//setupConf = "CURRENT_IP=" + currentIP + "\n" + "RPC_PORT=" + rpcPort + "\n" + "WHISPER_PORT=" + whisperPort + "\n" + "CONSTELLATION_PORT=" + constellationPort + "\n" + "RAFT_PORT=" + raftPort + "\n" + "NODEMANAGER_PORT=" + nodeManagerPort + "\n"
	//setupConfByte := []byte(setupConf)
	//err = ioutil.WriteFile("./Setup/" + nodename + "/setup.conf", setupConfByte, 0775)
	//if err != nil {
	//	panic(err)
	//}
	successResponse.Status = "success"
	return successResponse
}


func (nsi *NodeServiceImpl) joinRequestResponseCall(nodename string, currentIP string, rpcPort string, whisperPort string, constellationPort string, raftPort string, nodeManagerPort string, masterNodeManagerPort string, masterIP string) (SuccessResponse) {
	var successResponse SuccessResponse
	//cmd := exec.Command("./setup.sh","2", nodename, masterIP, masterNodeManagerPort, currentIP, rpcPort, whisperPort, constellationPort, raftPort, nodeManagerPort)
	//cmd.Dir = "./Setup"
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var setupConf string
	//setupConf = "CURRENT_IP=" + currentIP + "\n" + "RPC_PORT=" + rpcPort + "\n" + "WHISPER_PORT=" + whisperPort + "\n" + "CONSTELLATION_PORT=" + constellationPort + "\n" + "RAFT_PORT=" + raftPort + "\n" + "THIS_NODEMANAGER_PORT=" + nodeManagerPort + "\n" + "MASTER_IP=" + masterIP + "\n"+ "NODEMANAGER_PORT=" + masterNodeManagerPort + "\n"
	//setupConfByte := []byte(setupConf)
	//err = ioutil.WriteFile("./Setup/" + nodename + "/setup.conf", setupConfByte, 0775)
	//if err != nil {
	//	panic(err)
	//}
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
		if(match) {
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