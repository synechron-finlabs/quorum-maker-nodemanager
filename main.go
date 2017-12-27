package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"synechron.com/quorum-manager/service"
	"synechron.com/quorum-manager/client"
	"os"
	"bytes"
	"os/exec"
	"strings"
	"fmt"
)

var ipaddr string
var rpcport string
var nodeUrl string
var listenPort = ":8000"

func init(){
	var out bytes.Buffer
	cmd := exec.Command("./get_rpc.sh")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	rpcport = out.String()
	rpcport = strings.TrimSuffix(rpcport, "\n")

	var outip bytes.Buffer
	cmd = exec.Command("./get_ipaddr.sh")
	cmd.Stdout = &outip
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	ipaddr = outip.String()
	ipaddr = strings.TrimSuffix(ipaddr, "\n")
}

func main() {
    //We can use the 4 commented lines below if we want to remove dependency on fmt
	//s := []string{"http:",ipaddr}
	//var halfUrl = strings.Join(s, "//")
	//s = []string{halfUrl,rpcport}
	//nodeUrl = strings.Join(s, ":")
	halfUrl := fmt.Sprint("http://", ipaddr)
	nodeUrl = fmt.Sprint(halfUrl, ":",rpcport)

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
