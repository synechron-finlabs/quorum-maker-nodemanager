package service

import (
	"io/ioutil"
	"log"
	"synechron.com/quorum-manager/client"
	"strings"
	"github.com/magiconair/properties"
	"fmt"
	"strconv"
	"synechron.com/quorum-manager/util"
)

type ConnectionInfo struct {
	IP 	string 	`json:"ip,omitempty"`
	Port 	int 	`json:"port,omitempty"`
	Enode 	string 	`json:"enode,omitempty"`
}

type NodeInfo struct {
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
	EthAccount string `json:"eth-account,omitempty"`
}

type GetGenesisResponse struct {
	ContstellationPort string `json: "contstellation-port, omitempty"`
	NetID              string `json: "netID,omitempty"`
	Genesis            string `json: "genesis, omitempty"`
}

type BlockDetailsResponse struct {
	Number           int64                       `json:"number,omitempty"`
	Hash             string                       `json:"hash,omitempty"`
	ParentHash       string                       `json:"parentHash,omitempty"`
	Nonce            string                       `json:"nonce,omitempty"`
	Sha3Uncles       string                       `json:"sha3Uncles,omitempty"`
	LogsBloom        string                       `json:"logsBloom,omitempty"`
	TransactionsRoot string                       `json:"transactionsRoot,omitempty"`
	StateRoot        string                       `json:"stateRoot,omitempty"`
	Miner            string                       `json:"miner,omitempty"`
	Difficulty       int64                       `json:"difficulty,omitempty"`
	TotalDifficulty  int64                       `json:"totalDifficulty,omitempty"`
	ExtraData        string                       `json:"extraData,omitempty"`
	Size             int64                       `json:"size,omitempty"`
	GasLimit         int64                       `json:"gasLimit,omitempty"`
	GasUsed          int64                       `json:"gasUsed,omitempty"`
	Timestamp        int64                       `json:"timestamp,omitempty"`
	Transactions     []client.TransactionDetailsResponse `json:"transactions,omitempty"`
	Uncles           []string                     `json:"uncles,omitempty"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash,omitempty"`
	BlockNumber      int64 `json:"blockNumber"`
	From             string `json:"from,omitempty"`
	Gas              int64 `json:"gas,omitempty"`
	GasPrice         int64 `json:"gasPrice"`
	Hash             string `json:"hash,omitempty"`
	Input            string `json:"input,omitempty"`
	Nonce            int64 `json:"nonce"`
	To               string `json:"to,omitempty"`
	TransactionIndex int64 `json:"transactionIndex"`
	Value            int64 `json:"value,omitempty"`
	V                string `json:"v,omitempty"`
	R                string `json:"r,omitempty"`
	S                string `json:"s,omitempty"`
}

type NodeServiceImpl struct {
	Url string
}


func (nsi *NodeServiceImpl) GetGenesis(url string) (response GetGenesisResponse) {
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	constl := p.MustGetString("CONSTELLATION_PORT")
	constl = strings.TrimSuffix(constl, "\n")
	netid := p.MustGetString("NETWORK_ID")
	netid = strings.TrimSuffix(netid, "\n")

	b, err := ioutil.ReadFile("/home/node/genesis.json")
	if err != nil {
		log.Fatal(err)
	}
	genesis := string(b)
	genesis = strings.Replace(genesis, "\n","",-1)

	response = GetGenesisResponse{constl, netid, genesis}
	return response
}


func (nsi *NodeServiceImpl) JoinNetwork(request string, url string) (int) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	raftid := Ethclient.RaftAddPeer(request)
	return raftid
}


func (nsi *NodeServiceImpl) GetCurrentNode (url string) (NodeInfo) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}

	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	ipaddr := p.MustGetString("CURRENT_IP")
	raftid := p.MustGetString("RAFT_ID")
	rpcport := p.MustGetString("RPC_PORT")

	ipaddr = strings.TrimSuffix(ipaddr, "\n")
	raftid = strings.TrimSuffix(raftid, "\n")
	rpcport = strings.TrimSuffix(rpcport, "\n")

	raftidInt, err := strconv.Atoi(raftid)
	if err != nil {
		log.Fatal(err)
	}

	rpcportInt, err := strconv.Atoi(rpcport)
	if err != nil {
		log.Fatal(err)
	}

	thisadmininfo := Ethclient.AdminNodeInfo()
	enode := thisadmininfo.Enode

	pendingtxresponse := Ethclient.PendingTransactions()
	pendingtxcount := len(pendingtxresponse)

	blocknumber := Ethclient.BlockNumber()
	blocknumberInt :=  util.HexStringtoInt64(blocknumber)

	raftrole := Ethclient.RaftRole()

	raftrole = strings.TrimSuffix(raftrole, "\n")

	b, err := ioutil.ReadFile("/home/node/genesis.json")

	if err != nil {
		log.Fatal(err)
	}

	genesis := string(b)
	genesis = strings.Replace(genesis, "\n","",-1)
	conn := ConnectionInfo{ipaddr,rpcportInt,enode}
	responseobj := NodeInfo{conn,raftrole,raftidInt,blocknumberInt,pendingtxcount,genesis,thisadmininfo}
	return responseobj
}


func (nsi *NodeServiceImpl) GetOtherPeer(peerid string, url string) (client.AdminPeers) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	otherpeersresponse := Ethclient.AdminPeers()
	for _, item := range otherpeersresponse {
		if item.ID == peerid {
			peerresponse := item
			return peerresponse
		}
	}
	return client.AdminPeers{}
}


func (nsi *NodeServiceImpl) GetPendingTransactions(url string) ([]TransactionDetailsResponse) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	pendingtxresponseclient := Ethclient.PendingTransactions()
	pendingtxcount := len(pendingtxresponseclient)
	pendingtxresponse := make([]TransactionDetailsResponse, pendingtxcount)
	for i := 0; i < pendingtxcount; i++ {
		pendingtxresponse[i].BlockNumber = util.HexStringtoInt64(pendingtxresponseclient[i].BlockNumber)
		pendingtxresponse[i].Gas = util.HexStringtoInt64(pendingtxresponseclient[i].Gas)
		pendingtxresponse[i].GasPrice = util.HexStringtoInt64(pendingtxresponseclient[i].GasPrice)
		pendingtxresponse[i].TransactionIndex = util.HexStringtoInt64(pendingtxresponseclient[i].TransactionIndex)
		pendingtxresponse[i].Value = util.HexStringtoInt64(pendingtxresponseclient[i].Value)
		pendingtxresponse[i].Nonce = util.HexStringtoInt64(pendingtxresponseclient[i].Nonce)
		pendingtxresponse[i].BlockHash = pendingtxresponseclient[i].BlockHash
		pendingtxresponse[i].From = pendingtxresponseclient[i].From
		pendingtxresponse[i].Hash = pendingtxresponseclient[i].Hash
		pendingtxresponse[i].Input = pendingtxresponseclient[i].Input
		pendingtxresponse[i].To = pendingtxresponseclient[i].To
		pendingtxresponse[i].V = pendingtxresponseclient[i].V
		pendingtxresponse[i].R = pendingtxresponseclient[i].R
		pendingtxresponse[i].S = pendingtxresponseclient[i].S
	}
	return pendingtxresponse
}


func (nsi *NodeServiceImpl) GetBlockInfo(blockno int64, url string) (BlockDetailsResponse) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	blocknohex  := strconv.FormatInt(blockno, 16)
	bnohex := fmt.Sprint("0x", blocknohex)
	var blockresponse BlockDetailsResponse
	blockresponseclient := Ethclient.GetBlockByNumber(bnohex)
	blockresponse.Number = util.HexStringtoInt64(blockresponseclient.Number)
	blockresponse.Difficulty = util.HexStringtoInt64(blockresponseclient.Difficulty)
	blockresponse.TotalDifficulty = util.HexStringtoInt64(blockresponseclient.TotalDifficulty)
	blockresponse.Size = util.HexStringtoInt64(blockresponseclient.Size)
	blockresponse.GasLimit = util.HexStringtoInt64(blockresponseclient.GasLimit)
	blockresponse.GasUsed = util.HexStringtoInt64(blockresponseclient.GasUsed)
	blockresponse.Timestamp = util.HexStringtoInt64(blockresponseclient.Timestamp)
	blockresponse.Hash = blockresponseclient.Hash
	blockresponse.ParentHash = blockresponseclient.ParentHash
	blockresponse.Nonce = blockresponseclient.Nonce
	blockresponse.Sha3Uncles = blockresponseclient.Sha3Uncles
	blockresponse.LogsBloom = blockresponseclient.LogsBloom
	blockresponse.TransactionsRoot = blockresponseclient.TransactionsRoot
	blockresponse.StateRoot = blockresponseclient.StateRoot
	blockresponse.Miner = blockresponseclient.Miner
	blockresponse.ExtraData = blockresponseclient.ExtraData
	blockresponse.Transactions = blockresponseclient.Transactions
	blockresponse.Uncles = blockresponseclient.Uncles
	return blockresponse
}


func (nsi *NodeServiceImpl) GetTransactionInfo(txno string, url string) (TransactionDetailsResponse) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	var txresponse TransactionDetailsResponse
	txresponseclient := Ethclient.GetTransactionByHash(txno)
	txresponse.BlockNumber = util.HexStringtoInt64(txresponseclient.BlockNumber)
	txresponse.Gas = util.HexStringtoInt64(txresponseclient.Gas)
	txresponse.GasPrice = util.HexStringtoInt64(txresponseclient.GasPrice)
	txresponse.TransactionIndex = util.HexStringtoInt64(txresponseclient.TransactionIndex)
	txresponse.Value = util.HexStringtoInt64(txresponseclient.Value)
	txresponse.Nonce = util.HexStringtoInt64(txresponseclient.Nonce)
	txresponse.BlockHash = txresponseclient.BlockHash
	txresponse.From = txresponseclient.From
	txresponse.Hash = txresponseclient.Hash
	txresponse.Input = txresponseclient.Input
	txresponse.To = txresponseclient.To
	txresponse.V = txresponseclient.V
	txresponse.R = txresponseclient.R
	txresponse.S = txresponseclient.S
	return txresponse
}