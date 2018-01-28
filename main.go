package main

import (
	"synechron.com/NodeManagerGo/service"
	"os"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var nodeUrl = "http://localhost:22000"
var listenPort = ":8000"

func main() {

	if len(os.Args) > 1 {
		nodeUrl = os.Args[1]
	}

	if len(os.Args) > 2 {
		listenPort = ":" + os.Args[2]
	}

	router := mux.NewRouter()
	nodeService := service.NodeServiceImpl{nodeUrl}

	router.HandleFunc("/txn/{txn_hash}", nodeService.GetTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/block/{block_no}", nodeService.GetBlockInfoHandler).Methods("GET")
	router.HandleFunc("/genesis", nodeService.GetGenesisHandler).Methods("GET")
	router.HandleFunc("/peer/{peer_id}", nodeService.GetOtherPeerHandler).Methods("GET")
	router.HandleFunc("/peer", nodeService.JoinNetworkHandler).Methods("POST")
	router.HandleFunc("/peer", nodeService.GetCurrentNodeHandler).Methods("GET")
	router.HandleFunc("/txnrcpt/{txn_hash}", nodeService.GetTransactionReceiptHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(listenPort, router))

}


