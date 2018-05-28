package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/service"
	"os"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"time"
)

var nodeUrl = "http://localhost:22000"
var listenPort = ":8000"

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
		time.Sleep(80 * time.Second)
		log.Info("Deploying Network Manager Contract")
		nodeService.NetworkManagerContractDeployer(nodeUrl)
		nodeService.RegisterNodeDetails(nodeUrl)
	}()

	networkMapService := contractclient.NetworkMapContractClient{EthClient: client.EthClient{nodeUrl}}
	router.HandleFunc("/txn/{txn_hash}", nodeService.GetTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/txn", nodeService.GetLatestTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/block/{block_no}", nodeService.GetBlockInfoHandler).Methods("GET")
	router.HandleFunc("/block", nodeService.GetLatestBlockInfoHandler).Methods("GET")
	router.HandleFunc("/genesis", nodeService.GetGenesisHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/peer/{peer_id}", nodeService.GetOtherPeerHandler).Methods("GET")
	router.HandleFunc("/peer", nodeService.JoinNetworkHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/peer", nodeService.GetCurrentNodeHandler).Methods("GET")
	router.HandleFunc("/txnrcpt/{txn_hash}", nodeService.GetTransactionReceiptHandler).Methods("GET")
	router.HandleFunc("/pendingJoinRequests", nodeService.PendingJoinRequestsHandler).Methods("GET")
	router.HandleFunc("/joinRequestResponse", nodeService.JoinRequestResponseHandler).Methods("POST")
	router.HandleFunc("/joinRequestResponse", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/createNetwork", nodeService.CreateNetworkScriptCallHandler).Methods("POST")
	router.HandleFunc("/createNetwork", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/joinNetwork", nodeService.JoinNetworkScriptCallHandler).Methods("POST")
	router.HandleFunc("/joinNetwork", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/deployContract", nodeService.DeployContractHandler).Methods("POST")
	router.HandleFunc("/reset", nodeService.ResetHandler).Methods("GET")
	router.HandleFunc("/restart", nodeService.RestartHandler).Methods("GET")
	router.HandleFunc("/latestBlock", nodeService.LatestBlockHandler).Methods("GET")
	router.HandleFunc("/latency", nodeService.LatencyHandler).Methods("GET")
	router.HandleFunc("/logs", nodeService.LogsHandler).Methods("GET")
	router.HandleFunc("/txnsearch/{txn_hash}", nodeService.TransactionSearchHandler).Methods("GET")
	router.HandleFunc("/mailserver", nodeService.MailServerConfigHandler).Methods("POST")
	router.HandleFunc("/mailserver", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/registerNode", networkMapService.RegisterNodeRequestHandler).Methods("POST")
	router.HandleFunc("/updateNode", networkMapService.UpdateNodeRequestsHandler).Methods("POST")
	router.HandleFunc("/getNodeDetails/{index}", networkMapService.GetNodeDetailsResponseHandler).Methods("GET")
	router.HandleFunc("/getNodeList", networkMapService.GetNodeListSelfResponseHandler).Methods("GET")

	router.PathPrefix("/contracts").Handler(http.StripPrefix("/contracts", http.FileServer(http.Dir("/root/quorum-maker/contracts"))))
	router.PathPrefix("/geth").Handler(http.StripPrefix("/geth", http.FileServer(http.Dir("/home/node/qdata/gethLogs"))))
	router.PathPrefix("/constellation").Handler(http.StripPrefix("/constellation", http.FileServer(http.Dir("/home/node/qdata/constellationLogs"))))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("NodeManagerUI"))))

	log.WithFields(log.Fields{"url": nodeUrl, "port": listenPort}).Info("Node Manager listening...")
	log.Fatal(http.ListenAndServe(listenPort, router))
}
