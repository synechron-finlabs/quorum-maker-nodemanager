package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"synechron.com/NodeManagerGo/service"
	"os"
	"os/exec"
	"fmt"
	"time"
)

var nodeUrl = "http://localhost:22000"
var listenPort = ":8000"
var logPort = ":3000"

func main() {

	if len(os.Args) > 1 {
		nodeUrl = os.Args[1]
	}

	if len(os.Args) > 2 {
		listenPort = ":" + os.Args[2]
	}

	if len(os.Args) > 3 {
		logPort = ":" + os.Args[3]
	}
	ticker := time.NewTicker(86400 * time.Second)
	go func() {
		for range ticker.C {
			logRotater()
		}
	}()

	go func() {
		fs := http.FileServer(http.Dir("/home/node/qdata/logs/"))

		http.Handle("/logs/", http.StripPrefix("/logs", fs))

		fs1 := http.FileServer(http.Dir("/root/quorum-maker"))

		http.Handle("/contracts/", http.StripPrefix("/contracts", fs1))

		http.ListenAndServe(logPort, nil)
	}()

	router := mux.NewRouter()
	nodeService := service.NodeServiceImpl{nodeUrl}
	router.HandleFunc("/txn/{txn_hash}", nodeService.GetTransactionInfoHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/txn", nodeService.GetLatestTransactionInfoHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/block/{block_no}", nodeService.GetBlockInfoHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/block", nodeService.GetLatestBlockInfoHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/genesis", nodeService.GetGenesisHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/peer/{peer_id}", nodeService.GetOtherPeerHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/peer", nodeService.JoinNetworkHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/peer", nodeService.GetCurrentNodeHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/txnrcpt/{txn_hash}", nodeService.GetTransactionReceiptHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/pendingJoinRequests", nodeService.PendingJoinRequestsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/joinRequestResponse", nodeService.JoinRequestResponseHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/createNetwork", nodeService.CreateNetworkScriptCallHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/joinNetwork", nodeService.JoinNetworkScriptCallHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/deployContract", nodeService.DeployContractHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/reset", nodeService.ResetHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/restart", nodeService.RestartHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/nodeList", nodeService.NodeListHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/latestBlock", nodeService.LatestBlockHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/latency", nodeService.LatencyHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/logs", nodeService.LogsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/txnsearch/{txn_hash}", nodeService.TransactionSearchHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/mailserver", nodeService.MailServerConfigHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/pubkeys", nodeService.PublicKeysHandler).Methods("GET", "OPTIONS")

	log.Fatal(http.ListenAndServe(listenPort, router))
}

func logRotater() {
	command := "cat $(ls | grep log | grep -v _) > Geth_$(date| sed -e 's/ /_/g')"

	command1 := "echo -en '' > $(ls | grep log | grep -v _)"

	command2 := "cat $(ls | grep log | grep _) > Constellation_$(date| sed -e 's/ /_/g')"

	command3 := "echo -en '' > $(ls | grep log | grep _)"

	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = "/home/node/qdata/logs"
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd1 := exec.Command("bash", "-c", command1)
	cmd1.Dir = "/home/node/qdata/logs"
	err1 := cmd1.Run()
	if err1 != nil {
		fmt.Println(err)
	}

	cmd2 := exec.Command("bash", "-c", command2)
	cmd2.Dir = "/home/node/qdata/logs"
	err2 := cmd2.Run()
	if err2 != nil {
		fmt.Println(err)
	}

	cmd3 := exec.Command("bash", "-c", command3)
	cmd3.Dir = "/home/node/qdata/logs"
	err3 := cmd3.Run()
	if err3 != nil {
		fmt.Println(err)
	}
}
