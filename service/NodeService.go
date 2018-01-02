package service

import (
	"io/ioutil"
	"log"
	"synechron.com/quorum-manager/client"
	"strings"
	"github.com/magiconair/properties"
	"fmt"
	"strconv"
)

type ConnectionInfo struct {
	IP 	string 	`json:"ip,omitempty"`
	Port 	int 	`json:"port,omitempty"`
	Enode 	string 	`json:"enode,omitempty"`
}

type NodeInfo struct {
	ConnectionInfo  ConnectionInfo		`json:"connectionInfo,omitempty"`
	RaftRole 	string 				`json:"raftRole,omitempty"`
	RaftID 		int      			`json:"raftID,omitempty"`
	BlockNumber 	int64				`json:"blockNumber"`
	PendingTxCount 	int 				`json:"pendingTxCount"`
	Genesis 	string				`json:"genesis,omitempty"`
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
	blocknumber = strings.TrimSuffix(blocknumber, "\n")
	blocknumber = strings.TrimPrefix(blocknumber, "0x")
	blocknumberInt, err := strconv.ParseInt(blocknumber, 16, 64)
	if err != nil {
		fmt.Println(err)
	}

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


func (nsi *NodeServiceImpl) GetPendingTransactions(url string) ([]client.TransactionDetailsResponse) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	pendingtxresponse := Ethclient.PendingTransactions()
	return pendingtxresponse
}


func (nsi *NodeServiceImpl) GetBlockInfo(blockno int64, url string) (client.BlockDetailsResponse) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	blocknohex  := strconv.FormatInt(blockno, 16)
	bnohex := fmt.Sprint("0x", blocknohex)
	blockresponse := Ethclient.GetBlockByNumber(bnohex)
	return blockresponse
}


func (nsi *NodeServiceImpl) GetTransactionInfo(txno string, url string) (client.TransactionDetailsResponse) {
	var nodeUrl = url
	Ethclient := client.EthClient{nodeUrl}
	txresponse := Ethclient.GetTransactionByHash(txno)
	return txresponse
}