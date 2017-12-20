package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"synechron.com/quorum-manager/service"
	"synechron.com/quorum-manager/client"
)

func main() {
	router := mux.NewRouter()
	nodeService := service.NodeServiceImpl{}
	ethclient:=client.EthClient{}
	ethclient.Url= "http://10.34.14.243:22000"
	router.HandleFunc("/joinNetwork", nodeService.JoinNetworkHandler).Methods("POST")
	router.HandleFunc("/txn/{id}", ethclient.GetTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/txn/pending", ethclient.GetPendingTransactionsHandler).Methods("GET")
	router.HandleFunc("/block/{id}", ethclient.GetBlockInfoHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
