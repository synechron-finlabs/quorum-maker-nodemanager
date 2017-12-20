package client

import (
	"fmt"
	"github.com/ybbus/jsonrpc"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)
type BlockDetailsResponse struct {
	Hash             string `json:"hash,omitempty"`
	Nonce            string `json:"nonce,omitempty"`
	BlockHash        string `json:"blockHash,omitempty"`
	BlockNumber      string `json:"blockNumber,omitempty"`
	TransactionIndex string `json:"transactionIndex,omitempty"`
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
	Value            string `json:"value,omitempty"`
	Gas              string `json:"gas,omitempty"`
	GasPrice         string `json:"gasPrice,omitempty"`
	Input            string `json:"input,omitempty"`
}

type TransactionDetailsResponse struct {
	BlockHash        string      `json:"blockHash,omitempty"`
	BlockNumber      string `json:"blockNumber,omitempty"`
	From             string      `json:"from,omitempty"`
	Gas              string      `json:"gas,omitempty"`
	GasPrice         string      `json:"gasPrice,omitempty"`
	Hash             string      `json:"hash,omitempty"`
	Input            string      `json:"input,omitempty"`
	Nonce            string      `json:"nonce,omitempty"`
	To               string      `json:"to,omitempty"`
	TransactionIndex string `json:"transactionIndex,omitempty"`
	Value            string      `json:"value,omitempty"`
	V                string      `json:"v,omitempty"`
	R                string      `json:"r,omitempty"`
	S                string      `json:"s,omitempty"`
}

type EthClient struct {
	Url string
}

//func (ec *EthClient) GetCoinbaseAddress() (string, error) {
//	rpcClient := jsonrpc.NewRPCClient(ec.Url)
//
//	response, _ :=rpcClient.Call("eth_coinbase")
//
//	return response.GetString()
//}
//
//func (ec *EthClient) AddRaftPeer(enodeId string) (string, error){
//	rpcClient := jsonrpc.NewRPCClient(ec.Url)
//
//	response, err :=rpcClient.Call("admin_addPeer", enodeId)
//
//	fmt.Print(response)
//
//	if err != nil {
//		fmt.Println(err)
//	}
//	return response.GetString()
//}

func (ec *EthClient) GetTransactionInfo(txno string) (TransactionDetailsResponse){
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err1 :=rpcClient.Call("eth_getTransactionByHash", txno)
	if err1 != nil {
		fmt.Println(err1)
	}
	txresponse := TransactionDetailsResponse{}
	err := response.GetObject(&txresponse)
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

func (ec *EthClient) GetBlockInfo(blockno int64) (BlockDetailsResponse){
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	var blocknohex	string = strconv.FormatInt(blockno, 16)
	bnohex := fmt.Sprint("0x",blocknohex)
	response, err1 :=rpcClient.Call("eth_getBlockByNumber", bnohex,true)
	if err1 != nil {
		fmt.Println(err1)
	}
	blockresponse := BlockDetailsResponse{}
	err := response.GetObject(&blockresponse)
	if err != nil {
		fmt.Println(err)
	}
	return blockresponse
}

func (ec *EthClient) GetBlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	block,err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	response := ec.GetBlockInfo(block)
	fmt.Print(response)
	json.NewEncoder(w).Encode(response)
}

func (ec *EthClient) GetPendingTransactions() ([]TransactionDetailsResponse){
	rpcClient := jsonrpc.NewRPCClient(ec.Url)
	response, err1 :=rpcClient.Call("eth_pendingTransactions")
	if err1 != nil {
		fmt.Println(err1)
	}
	pendingtxresponse := []TransactionDetailsResponse{}
	err := response.GetObject(&pendingtxresponse)
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