package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"synechron.com/NodeManagerGo/service"
	"os"
	"synechron.com/NodeManagerGo/contractclient"
	"synechron.com/NodeManagerGo/client"
	"time"
)

var nodeUrl = "http://localhost:22000"
var listenPort = ":8000"
var logPort = ":3000"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

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

	router := mux.NewRouter()
	nodeService := service.NodeServiceImpl{nodeUrl}

	ticker := time.NewTicker(86400 * time.Second)
	go func() {
		for range ticker.C {
			log.Debug("Rotating log for Geth and Constellation.")
			nodeService.LogRotaterGeth()
			nodeService.LogRotaterConst()
		}
	}()

	go func() {
		fs := http.FileServer(http.Dir("/home/node/qdata/gethLogs/"))

		http.Handle("/geth/", http.StripPrefix("/geth", fs))

		fs1 := http.FileServer(http.Dir("/home/node/qdata/constellationLogs/"))

		http.Handle("/constellation/", http.StripPrefix("/constellation", fs1))

		fs2 := http.FileServer(http.Dir("/root/quorum-maker"))

		http.Handle("/contracts/", http.StripPrefix("/contracts", fs2))

		http.ListenAndServe(logPort, nil)
	}()

	go func() {
		time.Sleep(80 * time.Second)
		log.Info("Deploying Network Manager Contract")
		nodeService.NetworkManagerContractDeployer(nodeUrl)
		nodeService.RegisterNodeDetails(nodeUrl)
	}()

	networkMapService := contractclient.NetworkMapContractClient{EthClient: client.EthClient{nodeUrl}}
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
	router.HandleFunc("/latestBlock", nodeService.LatestBlockHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/latency", nodeService.LatencyHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/logs", nodeService.LogsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/txnsearch/{txn_hash}", nodeService.TransactionSearchHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/mailserver", nodeService.MailServerConfigHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/registerNode", networkMapService.RegisterNodeRequestHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/updateNode", networkMapService.UpdateNodeRequestsHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/getNodeDetails/{index}", networkMapService.GetNodeDetailsResponseHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/getNodeList", networkMapService.GetNodeListResponseHandler).Methods("GET", "OPTIONS")

	log.WithFields(log.Fields{"url" : nodeUrl, "port" : listenPort}).Info("Node Manager listening...")
	log.Fatal(http.ListenAndServe(listenPort, router))
}
