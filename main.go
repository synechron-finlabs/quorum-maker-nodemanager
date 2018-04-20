package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"synechron.com/NodeManagerGo/service"
	"os"
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
	nodeService := service.NodeServiceImpl{nodeUrl}
	router.HandleFunc("/txn/{txn_hash}", nodeService.GetTransactionInfoHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/txn", nodeService.GetLatestTransactionInfoHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/block/{block_no}", nodeService.GetBlockInfoHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/block", nodeService.GetLatestBlockInfoHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/genesis", nodeService.GetGenesisHandler).Methods("POST","OPTIONS")
	router.HandleFunc("/peer/{peer_id}", nodeService.GetOtherPeerHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/peer", nodeService.JoinNetworkHandler).Methods("POST","OPTIONS")
	router.HandleFunc("/peer", nodeService.GetCurrentNodeHandler).Methods("GET","OPTIONS")
        router.HandleFunc("/txnrcpt/{txn_hash}", nodeService.GetTransactionReceiptHandler).Methods("GET","OPTIONS")
        router.HandleFunc("/pendingJoinRequests", nodeService.PendingJoinRequestsHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/joinRequestResponse", nodeService.JoinRequestResponseHandler).Methods("POST","OPTIONS")
	router.HandleFunc("/createNetwork", nodeService.CreateNetworkScriptCallHandler).Methods("POST","OPTIONS")
	router.HandleFunc("/joinNetwork", nodeService.JoinNetworkScriptCallHandler).Methods("POST","OPTIONS")
	router.HandleFunc("/deployContract", nodeService.DeployContractHandler).Methods("POST","OPTIONS")
	router.HandleFunc("/reset", nodeService.ResetHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/restart", nodeService.RestartHandler).Methods("GET","OPTIONS")
	log.Fatal(http.ListenAndServe(listenPort, router))
}