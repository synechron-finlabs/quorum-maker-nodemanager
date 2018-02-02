package client

import (
	"fmt"
	"github.com/ybbus/jsonrpc"
	"log"
)

type AdminInfo struct {
	ID   		string 		`json:"id,omitempty"`
	Name 		string 		`json:"name,omitempty"`
	Enode 		string 		`json:"enode,omitempty"`
	IP 		string 		`json:"ip,omitempty"`
	Ports 		Ports 		`json:"ports,omitempty"`
	ListenAddr 	string 		`json:"listenAddr,omitempty"`
	Protocols 	Protocols	`json:"protocols,omitempty"`
}

type Ports struct {
	Discovery int `json:"discovery"`
	Listener  int `json:"listener,omitempty"`
}

type AdminPeers struct {
	ID      	string   	`json:"id,omitempty"`
	Name    	string   	`json:"name,omitempty"`
	Caps    	[]string 	`json:"caps,omitempty"`
	Network 	Network 	`json:"network,omitempty"`
	Protocols 	Protocols 	`json:"protocols,omitempty"`
}

type Protocols struct {
	Eth Eth `json:"eth,omitempty"`
}

type Eth struct {
	Network    int	  `json:"network,omitempty"`
 	Version    int    `json:"version,omitempty"`
	Difficulty int    `json:"difficulty,omitempty"`
	Genesis    string `json:"genesis,omitempty"`
	Head       string `json:"head,omitempty"`
}

type Network struct {
	LocalAddress  string `json:"localAddress,omitempty"`
	RemoteAddress string `json:"remoteAddress,omitempty"`
}

type BlockDetailsResponse struct {
	Number           string                       `json:"number,omitempty"`
	Hash             string                       `json:"hash,omitempty"`
	ParentHash       string                       `json:"parentHash,omitempty"`
	Nonce            string                       `json:"nonce,omitempty"`
	Sha3Uncles       string                       `json:"sha3Uncles,omitempty"`
	LogsBloom        string                       `json:"logsBloom,omitempty"`
	TransactionsRoot string                       `json:"transactionsRoot,omitempty"`
	StateRoot        string                       `json:"stateRoot,omitempty"`
	Miner            string                       `json:"miner,omitempty"`
	Difficulty       string                       `json:"difficulty,omitempty"`
	TotalDifficulty  string                       `json:"totalDifficulty,omitempty"`
	ExtraData        string                       `json:"extraData,omitempty"`
	Size             string                       `json:"size,omitempty"`
	GasLimit         string                       `json:"gasLimit,omitempty"`
	GasUsed          string                       `json:"gasUsed,omitempty"`
	Timestamp        string                       `json:"timestamp,omitempty"`
	Transactions     []TransactionDetailsResponse `json:"transactions,omitempty"`
	Uncles           []string                     `json:"uncles,omitempty"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash,omitempty"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from,omitempty"`
	Gas              string `json:"gas,omitempty"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash,omitempty"`
	Input            string `json:"input,omitempty"`
	Nonce            string `json:"nonce"`
	To               string `json:"to,omitempty"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value,omitempty"`
	V                string `json:"v,omitempty"`
	R                string `json:"r,omitempty"`
	S                string `json:"s,omitempty"`
}

type TransactionReceiptResponse struct {
	BlockHash         	string          `json:"blockHash"`
	BlockNumber       	string          `json:"blockNumber"`
	ContractAddress   	string 		`json:"contractAddress"`
	CumulativeGasUsed 	string          `json:"cumulativeGasUsed"`
	From              	string          `json:"from"`
	GasUsed           	string          `json:"gasUsed"`
	Logs              	[]Logs 		`json:"logs"`
	LogsBloom        	string 		`json:"logsBloom"`
	Root             	string 		`json:"root"`
	To               	string 		`json:"to"`
	TransactionHash  	string 		`json:"transactionHash"`
	TransactionIndex 	string 		`json:"transactionIndex"`
}

type Logs struct {
	Address         	string        `json:"address"`
	BlockHash       	string        `json:"blockHash"`
	BlockNumber   		string        `json:"blockNumber"`
	Data 			string        `json:"data"`
	LogIndex          	string        `json:"logIndex"`
	Topics           	[]string      `json:"topics"`
	TransactionHash         string        `json:"transactionHash"`
	TransactionIndex  	string        `json:"transactionIndex"`
}

type EthClient struct {
	Url string
}

func (ec *EthClient) GetTransactionByHash(txNo string) (TransactionDetailsResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_getTransactionByHash", txNo)

	if err != nil {
		fmt.Println(err)
	}
	txResponse := TransactionDetailsResponse{}
	err = response.GetObject(&txResponse)
	if err != nil {
		fmt.Println(err)
	}
	return txResponse
}

func (ec *EthClient) GetBlockByNumber(blockNo string) (BlockDetailsResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_getBlockByNumber", blockNo, true)
	if err != nil {
		fmt.Println(err)
	}
	blockResponse := BlockDetailsResponse{}
	err = response.GetObject(&blockResponse)
	if err != nil {
		fmt.Println(err)
	}
	return blockResponse
}

func (ec *EthClient) PendingTransactions() ([]TransactionDetailsResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_pendingTransactions")
	if err != nil {
		fmt.Println(err)
	}
	pendingTxResponse := []TransactionDetailsResponse{}
	err = response.GetObject(&pendingTxResponse)
	if err != nil {
		fmt.Println(err)
	}
	return pendingTxResponse
}

func (ec *EthClient) AdminPeers() ([]AdminPeers) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("admin_peers")
	if err != nil {
		fmt.Println(err)
	}
	otherPeersResponse := []AdminPeers{}
	err = response.GetObject(&otherPeersResponse)
	if err != nil {
		fmt.Println(err)
	}
	return otherPeersResponse
}

func (ec *EthClient) AdminNodeInfo () (AdminInfo) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("admin_nodeInfo")
	if err != nil {
		fmt.Println(err)
	}
	thisAdminInfo := AdminInfo{}
	err = response.GetObject(&thisAdminInfo)
	return thisAdminInfo
}

func (ec *EthClient) BlockNumber() (string) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_blockNumber")
	if err != nil {
		fmt.Println(err)
	}
	var blockNumber string;
	err = response.GetObject(&blockNumber)
	if err != nil {
		fmt.Println(err)
	}
	return blockNumber
}

func (ec *EthClient) RaftRole() (string) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("raft_role")
	if err != nil {
		fmt.Println(err)
	}
	var raftRole string;
	err = response.GetObject(&raftRole)
	if err != nil {
		fmt.Println(err)
	}
	return raftRole
}

func (ec *EthClient) RaftAddPeer(request string) (int) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("raft_addPeer",request)
	var raftId int
	err = response.GetObject(&raftId)
	if err != nil {
		log.Fatal(err)
	}
	return raftId
}

func (ec *EthClient) GetTransactionReceipt(txNo string) (TransactionReceiptResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_getTransactionReceipt", txNo)

	if err != nil {
		fmt.Println(err)
	}
	txResponse := TransactionReceiptResponse{}
	err = response.GetObject(&txResponse)
	if err != nil {
		fmt.Println(err)
	}
	return txResponse
}


type Payload struct {
	From       string   `json:"from"`
	To         string   `json:"to"`
	Data       string   `json:"data"`
	PrivateFor []string `json:"privateFor"`
}

type CallPayload struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

func (ec *EthClient) SendTransaction(param ContractParam, rh RequestHandler) string {

	rpcClient := jsonrpc.NewRPCClient(ec.Url)

	response, err := rpcClient.Call("personal_unlockAccount", param.From, param.Passwd, nil)

	if err != nil {
		fmt.Println(err)
	}

	p := Payload{param.From,
		param.To,
		rh.Encode(),
		param.Parties}

	response, err = rpcClient.Call("eth_sendTransaction", p)

	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s", response.Result)

}

func (ec *EthClient) EthCall(param ContractParam, req RequestHandler, res ResponseHandler) interface{} {

	rpcClient := jsonrpc.NewRPCClient(ec.Url)

	p := CallPayload{param.To, req.Encode()}

	response, err := rpcClient.Call("eth_call", p, "latest")

	if err != nil {
		fmt.Println(err)
	}

	res.Decode(fmt.Sprintf("%v", response.Result))

	return res

}

type SendTransaction struct {
	From        string      `json:"from"`
	PrivateFor	[]string    `json:"privateFor"`
	Gas         string      `json:"gas"`
	Data        string      `json:"data"`
}

func (ec *EthClient) DeployContracts(byteCode string, addr []string) string {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_coinbase")
	if err != nil {
		fmt.Println("can't get coinbase" + err.Error())
	}
	var ownerAddress string
	err = response.GetObject(&ownerAddress)

	if err != nil {
		fmt.Println(err)
	}
	var privateFor []string
	for i:=0; i<len(addr);i++{
		privateFor = append(privateFor,addr[i])
	}

	response, err = rpcClient.Call("personal_unlockAccount", ownerAddress, "",nil)
	if err != nil {
		fmt.Println( err)
	}

	sendTransaction := SendTransaction{
		ownerAddress,
		privateFor,
		"100000000",
		byteCode}

	response, err = rpcClient.Call("eth_sendTransaction", sendTransaction)

	if err != nil {
		fmt.Println("transaction failed" + err.Error())
	}

	var txHash string
	err = response.GetObject(&txHash)

	if err != nil {
		fmt.Println(err)
	}

	contractAdd := ec.GetTransactionReceipt(txHash).ContractAddress
	return contractAdd
}