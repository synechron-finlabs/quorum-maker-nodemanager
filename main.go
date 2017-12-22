package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"synechron.com/quorum-manager/service"
	"synechron.com/quorum-manager/client"
	"os"
	"fmt"
)

var nodeUrl = "http://localhost:22000"
var listenPort = ":8000"

func main() {

	if len(os.Args) > 1 {
		nodeUrl = os.Args[1]
	}

	if len(os.Args) > 2 {
		listenPort = ":"  + os.Args[2]
	}

	router := mux.NewRouter()
	nodeService := service.NodeServiceImpl{}
	ethclient:= client.EthClient{nodeUrl}

	router.HandleFunc("/txn/{id}", ethclient.GetTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/txn/pending", ethclient.GetPendingTransactionsHandler).Methods("GET")
	router.HandleFunc("/block/{id}", ethclient.GetBlockInfoHandler).Methods("GET")
	router.HandleFunc("/genesis", nodeService.GetGenesisHandler).Methods("GET")
	//router.HandleFunc("/sendGenesis", nodeService.GetGenesisHandler).Methods("GET")
	router.HandleFunc("/peer/{id}", ethclient.GetOtherPeerHandler).Methods("GET")
	router.HandleFunc("/peer", nodeService.JoinNetworkHandler).Methods("POST")
	//router.HandleFunc("/joinNetwork", nodeService.JoinNetworkHandler).Methods("POST")
	router.HandleFunc("/peer", ethclient.GetCurrentNodeHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(listenPort, router))
}
