package service

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"bufio"
	"os"
	"strings"
	"time"
	"io"
	"io/ioutil"
	"bytes"
)

type PendingRequests struct {
	PendingJoinRequests 	[]string `json:"pendingJoinRequests,omitempty"`
}

var cCLI = make(chan string, 1)
var cUI = make(chan string, 1)

var peerMap = map[string]string{}

func (nsi *NodeServiceImpl) JoinNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID

	if peerMap[enode] == "" {
		peerMap[enode] = "PENDING"
	}

	if peerMap[enode] == "YES" {
		response := nsi.joinNetwork(enode, nsi.Url)
		json.NewEncoder(w).Encode(response)
	} else if peerMap[enode] == "NO" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Access denied"))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Pending user response"))
	}
}

func (nsi *NodeServiceImpl) GetGenesisHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	foreignIP := request.IPAddress

	if peerMap[enode] == "" {
		peerMap[enode] = "PENDING"
	}

	go func() {
		fmt.Println("Request for Joining this network for Enode",enode,"from IP",foreignIP,"Do you approve ? y/N")

		reader := bufio.NewReader(os.Stdin)
		reply, _ := reader.ReadString('\n')
		reply =  strings.TrimSuffix(reply, "\n")
		if reply == "y" || reply == "Y" {
			peerMap[enode] = "YES"
			response := nsi.getGenesis(nsi.Url)
			json.NewEncoder(w).Encode(response)
		} else {
			peerMap[enode] = "NO"
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
		}
		cCLI <- "CLI response"
	}()

	select {
	case resCLI := <-cCLI:
		fmt.Println(resCLI)
	case resUI := <-cUI:
		if peerMap[enode] == "YES" {
			response := nsi.getGenesis(nsi.Url)
			json.NewEncoder(w).Encode(response)
		} else if peerMap[enode] == "NO" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
		} else {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Pending user response"))
		}
		fmt.Println(resUI)
	case <-time.After(time.Second * 300):
		fmt.Println("Response Timed Out")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Pending user response"))
	}

}

func (nsi *NodeServiceImpl) GetCurrentNodeHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.getCurrentNode(nsi.Url)
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetOtherPeerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["peer_id"] == "all" {
		response := nsi.getOtherPeers(nsi.Url)
		json.NewEncoder(w).Encode(response)
	} else {
		response := nsi.getOtherPeer(params["peer_id"],nsi.Url)
		json.NewEncoder(w).Encode(response)
	}
}

func (nsi *NodeServiceImpl) GetLatestBlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	count := r.FormValue("number")
	reference := r.FormValue("reference")
	response := nsi.getLatestBlockInfo(count,reference,nsi.Url)
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) GetBlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	block, err := strconv.ParseInt(params["block_no"], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	response := nsi.getBlockInfo(block,nsi.Url)
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) GetTransactionInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["txn_hash"] == "pending" {
		response := nsi.getPendingTransactions(nsi.Url)
		json.NewEncoder(w).Encode(response)
	} else {
		response := nsi.getTransactionInfo(params["txn_hash"],nsi.Url)
		json.NewEncoder(w).Encode(response)
	}
}


func (nsi *NodeServiceImpl) GetLatestTransactionInfoHandler(w http.ResponseWriter, r *http.Request) {
	count := r.FormValue("number")
	response := nsi.getLatestTransactionInfo(count,nsi.Url)
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) GetTransactionReceiptHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := nsi.getTransactionReceipt(params["txn_hash"],nsi.Url)
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) PendingJoinRequestsHandler(w http.ResponseWriter, r *http.Request) {
	var pendingRequests PendingRequests
	for key, _ := range peerMap {
		if peerMap[key] == "PENDING" {
			pendingRequests.PendingJoinRequests = append(pendingRequests.PendingJoinRequests, key)
		}
	}
	json.NewEncoder(w).Encode(pendingRequests)
}

func (nsi *NodeServiceImpl) JoinRequestResponseHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkResponse
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	status := request.Status
	response := nsi.joinRequestResponse(enode,status)
	json.NewEncoder(w).Encode(response)
	cUI <- "UI response"
}


func (nsi *NodeServiceImpl) DeployContractHandler(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	var private bool
	var publicKeys []string

	count := r.FormValue("count")
	countInt, err := strconv.Atoi(count)
	if err != nil {
		panic(err)
	}
	fileNames := make([]string, countInt)
	boolVal := r.FormValue("private")
	if boolVal == "true" {
		private = true
	} else {
		private = false
	}

	keys := r.FormValue("privateFor")
	publicKeys = strings.Split(keys, ",")

	for i := 0; i < countInt; i++ {
		keyVal := "file" + strconv.Itoa(i + 1)

		file, header, err := r.FormFile(keyVal)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")

		fileNames[i] = name[0] + ".sol"

		io.Copy(&Buf, file)

		contents := Buf.String()

		fileContent := []byte(contents)
		err = ioutil.WriteFile("./" + name[0] + ".sol" , fileContent, 0775)
		if err != nil {
			panic(err)
		}

		Buf.Reset()
	}

	response := nsi.deployContract(publicKeys, fileNames, private, nsi.Url)
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) CreateNetworkScriptCallHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateNetworkScriptArgs
	_ = json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(request)
	response := nsi.createNetworkScriptCall(request.Nodename,request.CurrentIP,request.RPCPort,request.WhisperPort,request.ConstellationPort,request.RaftPort,request.NodeManagerPort)
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) JoinNetworkScriptCallHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkScriptArgs
	_ = json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(request)
	response := nsi.joinRequestResponseCall(request.Nodename,request.CurrentIP,request.RPCPort,request.WhisperPort,request.ConstellationPort,request.RaftPort,request.NodeManagerPort,request.MasterNodeManagerPort,request.MasterIP)
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) ResetHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.resetCurrentNode()
	json.NewEncoder(w).Encode(response)
}


func (nsi *NodeServiceImpl) RestartHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.restartCurrentNode()
	json.NewEncoder(w).Encode(response)
}