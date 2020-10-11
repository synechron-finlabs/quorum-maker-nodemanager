package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/env"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/service"
)

var nodeUrl = "http://localhost:22000"
var listenPort = ":8000"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	env.GetAppConfig(true)
	env.GetSetupConf(true)
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
	router.HandleFunc(env.GetSetupConf().ContextPath+"/txn/{txn_hash}", nodeService.GetTransactionInfoHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/txn", nodeService.GetLatestTransactionInfoHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/block/{block_no}", nodeService.GetBlockInfoHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/block", nodeService.GetLatestBlockInfoHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/genesis", nodeService.GetGenesisHandler).Methods("POST", "OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/peer/{peer_id}", nodeService.GetOtherPeerHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/peer", nodeService.JoinNetworkHandler).Methods("POST", "OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/peer", nodeService.GetCurrentNodeHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/txnrcpt/{txn_hash}", nodeService.GetTransactionReceiptHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/pendingJoinRequests", nodeService.PendingJoinRequestsHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/joinRequestResponse", nodeService.JoinRequestResponseHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/joinRequestResponse", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/createNetwork", nodeService.CreateNetworkScriptCallHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/createNetwork", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/joinNetwork", nodeService.JoinNetworkScriptCallHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/joinNetwork", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/deployContract", nodeService.DeployContractHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/reset", nodeService.ResetHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/restart", nodeService.RestartHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/latestBlock", nodeService.LatestBlockHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/latency", nodeService.LatencyHandler).Methods("GET")
	//router.HandleFunc("/logs", nodeService.LogsHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/txnsearch/{txn_hash}", nodeService.TransactionSearchHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/mailserver", nodeService.MailServerConfigHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/mailserver", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/registerNode", networkMapService.RegisterNodeRequestHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/updateNode", networkMapService.UpdateNodeHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/updateNode", networkMapService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/getNodeDetails/{index}", networkMapService.GetNodeDetailsResponseHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/getNodeList", networkMapService.GetNodeListSelfResponseHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/activeNodes", networkMapService.ActiveNodesHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/chartData", nodeService.GetChartDataHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/contractList", nodeService.GetContractListHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/contractCount", nodeService.GetContractCountHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/updateContractDetails", nodeService.ContractDetailsUpdateHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/attachedNodeDetails", nodeService.AttachedNodeDetailsHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/initialized", nodeService.InitializationHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/createAccount", nodeService.CreateAccountHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/createAccount", nodeService.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/getAccounts", nodeService.GetAccountsHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/getWhitelist", nodeService.GetWhitelistedIPsHandler).Methods("GET")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/updateWhitelist", nodeService.UpdateWhitelistHandler).Methods("POST")
	router.HandleFunc(env.GetSetupConf().ContextPath+"/updateWhitelist", nodeService.OptionsHandler).Methods("OPTIONS")

	redirectPaths := []string{"/", "/index.html", "/index.htm" , "/qm/dashboard"}
	
	for _, path := range redirectPaths {
		router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("http://%s/qm/", r.Host)
			http.Redirect(w, r, fmt.Sprintf("http://%s/qm/", r.Host), http.StatusPermanentRedirect)
		})
	}



	router.PathPrefix(env.GetSetupConf().ContextPath + "/contracts").Handler(http.StripPrefix(env.GetSetupConf().ContextPath+"/contracts", http.FileServer(http.Dir(env.GetAppConfig().ContractsDir))))
	router.PathPrefix(env.GetSetupConf().ContextPath + "/geth").Handler(http.StripPrefix(env.GetSetupConf().ContextPath+"/geth", http.FileServer(http.Dir(env.GetAppConfig().GethLogs))))
	router.PathPrefix(env.GetSetupConf().ContextPath + "/constellation").Handler(http.StripPrefix(env.GetSetupConf().ContextPath+"/constellation", http.FileServer(http.Dir(env.GetAppConfig().PrivacyLogs))))
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

	if r.URL.Path == "" && env.GetSetupConf().ContextPath != "" || "/"+r.URL.Path == env.GetSetupConf().ContextPath+"/dashboard" {
		http.Redirect(w, r, env.GetSetupConf().ContextPath, 301)
		return
	}
	_, err := os.Open(mf.name + "/" + r.URL.Path)
	if err != nil {
		http.Error(w, "File Not Found", 404)
	}

	mf.handler.ServeHTTP(w, r)
}
