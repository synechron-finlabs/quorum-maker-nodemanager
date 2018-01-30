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
)

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
		fmt.Println("Request for genesis.JSON for Enode",enode,"from IP",foreignIP,"Do you approve ? y/N")

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
	case <-time.After(time.Second * 100):
		fmt.Println("Response Timed Out")
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
	}

}

func (nsi *NodeServiceImpl) GetCurrentNodeHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.getCurrentNode(nsi.Url)
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetOtherPeerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := nsi.getOtherPeer(params["peer_id"],nsi.Url)
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

func (nsi *NodeServiceImpl) GetTransactionReceiptHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := nsi.getTransactionReceipt(params["txn_hash"],nsi.Url)
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) PendingJoinRequestsHandler(w http.ResponseWriter, r *http.Request) {
	pendingEnodes := []string{}
	for key, _ := range peerMap {
		if peerMap[key] == "PENDING" {
			pendingEnodes = append(pendingEnodes, key)
		}
	}
	json.NewEncoder(w).Encode(pendingEnodes)
}

func (nsi *NodeServiceImpl) JoinRequestResponseHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkResponse
	var reply string
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	status := request.Status
	response := nsi.joinRequestResponse(enode,status)
	if response == true {
		reply = fmt.Sprintf("Successfully updated status of %s to %s",enode, status)
	} else {
		reply = fmt.Sprintf("Failed to update status of %s to %s ",enode, status)
	}
	json.NewEncoder(w).Encode(reply)
	cUI <- "UI response"
}