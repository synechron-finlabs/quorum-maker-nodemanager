package client

import (
	"fmt"
	"github.com/ybbus/jsonrpc"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)

type AdminPeers struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Caps    []string `json:"caps,omitempty"`
	Network Network `json:"network,omitempty"`
	Protocols Protocols `json:"protocols,omitempty"`
}

type Protocols struct {
	Eth Eth `json:"eth,omitempty"`
}

type Eth struct {
	Version    int    `json:"version,omitempty"`
	Difficulty int    `json:"difficulty,omitempty"`
	Head       string `json:"head,omitempty"`
}

type Network struct {
	LocalAddress  string `json:"localAddress,omitempty"`
	RemoteAddress string `json:"remoteAddress,omitempty"`
}

type BlockDetailsResponse struct {
	Number           string                       `json:"number"`
	Hash             string                       `json:"hash"`
	ParentHash       string                       `json:"parentHash"`
	Nonce            string                       `json:"nonce"`
	Sha3Uncles       string                       `json:"sha3Uncles"`
	LogsBloom        string                       `json:"logsBloom"`
	TransactionsRoot string                       `json:"transactionsRoot"`
	StateRoot        string                       `json:"stateRoot"`
	Miner            string                       `json:"miner"`
	Difficulty       string                       `json:"difficulty"`
	TotalDifficulty  string                       `json:"totalDifficulty"`
	ExtraData        string                       `json:"extraData"`
	Size             string                       `json:"size"`
	GasLimit         string                       `json:"gasLimit"`
	GasUsed          string                       `json:"gasUsed"`
	Timestamp        string                       `json:"timestamp"`
	Transactions     []TransactionDetailsResponse `json:"transactions"`
	Uncles           []string                     `json:"uncles"`
}

type TransactionDetailsResponse struct {
	BlockHash        string `json:"blockHash,omitempty"`
	BlockNumber      string `json:"blockNumber,omitempty"`
	From             string `json:"from,omitempty"`
	Gas              string `json:"gas,omitempty"`
	GasPrice         string `json:"gasPrice,omitempty"`
	Hash             string `json:"hash,omitempty"`
	Input            string `json:"input,omitempty"`
	Nonce            string `json:"nonce,omitempty"`
	To               string `json:"to,omitempty"`
	TransactionIndex string `json:"transactionIndex,omitempty"`
	Value            string `json:"value,omitempty"`
	V                string `json:"v,omitempty"`
	R                string `json:"r,omitempty"`
	S                string `json:"s,omitempty"`
}

type EthClient struct {
	Url string
}

func (ec *EthClient) GetTransactionInfo(txno string) (TransactionDetailsResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_getTransactionByHash", txno)

	if err != nil {
		fmt.Println(err)
	}

	txresponse := TransactionDetailsResponse{}

	err = response.GetObject(&txresponse)
	if err != nil {
		fmt.Println(err)
	}
	return txresponse
}

func (ec *EthClient) GetTransactionInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := ec.GetTransactionInfo(params["id"])
	fmt.Print(response)
	json.NewEncoder(w).Encode(response)
}

func (ec *EthClient) GetBlockInfo(blockno int64) (BlockDetailsResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	blocknohex  := strconv.FormatInt(blockno, 16)
	bnohex := fmt.Sprint("0x", blocknohex)

	response, err := rpcClient.Call("eth_getBlockByNumber", bnohex, true)
	if err != nil {
		fmt.Println(err)
	}

	blockresponse := BlockDetailsResponse{}
	err = response.GetObject(&blockresponse)
	if err != nil {
		fmt.Println(err)
	}
	return blockresponse
}

func (ec *EthClient) GetBlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	block, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	response := ec.GetBlockInfo(block)
	fmt.Print(response)
	json.NewEncoder(w).Encode(response)
}

func (ec *EthClient) GetPendingTransactions() ([]TransactionDetailsResponse) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("eth_pendingTransactions")
	if err != nil {
		fmt.Println(err)
	}
	pendingtxresponse := []TransactionDetailsResponse{}
	err = response.GetObject(&pendingtxresponse)
	if err != nil {
		fmt.Println(err)
	}
	return pendingtxresponse
}

func (ec *EthClient) GetPendingTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	response := ec.GetPendingTransactions()
	fmt.Print(response)
	json.NewEncoder(w).Encode(response)
}

func (ec *EthClient) GetOtherPeer(peerid string) (AdminPeers) {
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err := rpcClient.Call("admin_peers")
	if err != nil {
		fmt.Println(err)
	}
	otherpeersresponse := []AdminPeers{}
	err = response.GetObject(&otherpeersresponse)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range otherpeersresponse {
		if item.ID == peerid {
			peerresponse := item
			return peerresponse
		}
	}
	return AdminPeers{}
}

func (ec *EthClient) GetOtherPeerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := ec.GetOtherPeer(params["id"])
	fmt.Print(response)
	json.NewEncoder(w).Encode(response)
}