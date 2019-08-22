package main

import (
	"fmt"
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/service"
	"net/http"
	"os"
	"os/signal"
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
		nodeService.CheckGethStatus(nodeUrl)
		//log.Info("Deploying Network Manager Contract")
		nodeService.NetworkManagerContractDeployer(nodeUrl)
		nodeService.RegisterNodeDetails(nodeUrl)
		nodeService.ContractCrawler(nodeUrl)
		nodeService.ABICrawler(nodeUrl)
		nodeService.IPWhitelister()
	}()

	networkMapService := contractclient.NetworkMapContractClient{EthClient: client.EthClient{nodeUrl}}
	router.HandleFunc("/qm/txn/{txn_hash}", nodeService.GetTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/qm/txn", nodeService.GetLatestTransactionInfoHandler).Methods("GET")
	router.HandleFunc("/qm/block/{block_no}", nodeService.GetBlockInfoHandler).Methods("GET")
	router.HandleFunc("/qm/block", nodeService.GetLatestBlockInfoHandler).Methods("GET")
	router.HandleFunc("/qm/genesis", nodeService.GetGenesisHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/qm/peer/{peer_id}", nodeService.GetOtherPeerHandler).Methods("GET")
	router.HandleFunc("/qm/peer", nodeService.JoinNetworkHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/qm/peer", nodeService.GetCurrentNodeHandler).Methods("GET")
	router.HandleFunc("/qm/txnrcpt/{txn_hash}", nodeService.GetTransactionReceiptHandler).Methods("GET")
	router.HandleFunc("/qm/pendingJoinRequests", nodeService.PendingJoinRequestsHandler).Methods("GET")
	router.HandleFunc("/qm/joinRequestResponse", nodeService.JoinRequestResponseHandler).Methods("POST")
	router.HandleFunc("/qm/joinRequestResponse", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/qm/createNetwork", nodeService.CreateNetworkScriptCallHandler).Methods("POST")
	router.HandleFunc("/qm/createNetwork", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/qm/joinNetwork", nodeService.JoinNetworkScriptCallHandler).Methods("POST")
	router.HandleFunc("/qm/joinNetwork", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/qm/deployContract", nodeService.DeployContractHandler).Methods("POST")
	router.HandleFunc("/qm/reset", nodeService.ResetHandler).Methods("GET")
	router.HandleFunc("/qm/restart", nodeService.RestartHandler).Methods("GET")
	router.HandleFunc("/qm/latestBlock", nodeService.LatestBlockHandler).Methods("GET")
	router.HandleFunc("/qm/latency", nodeService.LatencyHandler).Methods("GET")
	//router.HandleFunc("/logs", nodeService.LogsHandler).Methods("GET")
	router.HandleFunc("/qm/txnsearch/{txn_hash}", nodeService.TransactionSearchHandler).Methods("GET")
	router.HandleFunc("/qm/mailserver", nodeService.MailServerConfigHandler).Methods("POST")
	router.HandleFunc("/qm/mailserver", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/qm/registerNode", networkMapService.RegisterNodeRequestHandler).Methods("POST")
	router.HandleFunc("/qm/updateNode", networkMapService.UpdateNodeHandler).Methods("POST")
	router.HandleFunc("/qm/updateNode", networkMapService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/qm/getNodeDetails/{index}", networkMapService.GetNodeDetailsResponseHandler).Methods("GET")
	router.HandleFunc("/qm/getNodeList", networkMapService.GetNodeListSelfResponseHandler).Methods("GET")
	router.HandleFunc("/qm/activeNodes", networkMapService.ActiveNodesHandler).Methods("GET")
	router.HandleFunc("/qm/chartData", nodeService.GetChartDataHandler).Methods("GET")
	router.HandleFunc("/qm/contractList", nodeService.GetContractListHandler).Methods("GET")
	router.HandleFunc("/qm/contractCount", nodeService.GetContractCountHandler).Methods("GET")
	router.HandleFunc("/qm/updateContractDetails", nodeService.ContractDetailsUpdateHandler).Methods("POST")
	router.HandleFunc("/qm/attachedNodeDetails", nodeService.AttachedNodeDetailsHandler).Methods("POST")
	router.HandleFunc("/qm/initialized", nodeService.InitializationHandler).Methods("GET")
	router.HandleFunc("/qm/createAccount", nodeService.CreateAccountHandler).Methods("POST")
	router.HandleFunc("/qm/createAccount", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/qm/getAccounts", nodeService.GetAccountsHandler).Methods("GET")
	router.HandleFunc("/qm/getWhitelist", nodeService.GetWhitelistedIPsHandler).Methods("GET")
	router.HandleFunc("/qm/updateWhitelist", nodeService.UpdateWhitelistHandler).Methods("POST")
	router.HandleFunc("/qm/updateWhitelist", nodeService.OptionsHandler).Methods("OPTIONS")

	router.PathPrefix("/qm/contracts").Handler(http.StripPrefix("/contracts", http.FileServer(http.Dir("/root/quorum-maker/contracts"))))
	router.PathPrefix("/qm/geth").Handler(http.StripPrefix("/geth", http.FileServer(http.Dir("/home/node/qdata/gethLogs"))))
	router.PathPrefix("/qm/constellation").Handler(http.StripPrefix("/constellation", http.FileServer(http.Dir("/home/node/qdata/constellationLogs"))))
	router.PathPrefix("/").Handler(http.StripPrefix("/", NewFileServer("NodeManagerUI")))

	log.Info(fmt.Sprintf("Node Manager listening on %s...", listenPort))

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0" + listenPort,

		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
		//IdleTimeout:  time.Second * 60,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("Node Manager Shutting down")
	os.Exit(0)
}

type MyFileServer struct {
	name    string
	handler http.Handler
}

func NewFileServer(file string) *MyFileServer {

	return &MyFileServer{file, http.FileServer(http.Dir(file))}

}
func (mf *MyFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	_, err := os.Open(mf.name + "/" + r.URL.Path)
	if err != nil {
		r.URL.Path = "/"
	}

	mf.handler.ServeHTTP(w, r)
}
