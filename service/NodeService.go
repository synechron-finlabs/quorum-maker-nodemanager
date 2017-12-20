package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type JoinNetworkRequest struct {
	EnodeID    string `json:"enode-id,omitempty"`
	EthAccount string `json:"eth-account,omitempty"`
}

type JoinNetworkResponse struct {
	RaftID             string `json:"raftID,omitempty"`
	ContstellationPort string `json: "contstellation-port, omitempty"`
	Genesis            string `json: "genesis, omitempty"`
}


type NodeService interface {
	JoinNetwork(request JoinNetworkRequest) JoinNetworkResponse
}

type NodeServiceImpl struct {
}

func (nsi NodeServiceImpl) JoinNetwork(request JoinNetworkRequest) (response JoinNetworkResponse) {

	b, err := ioutil.ReadFile("genesis.json")

	if err != nil {
		log.Fatal(err)
	}

	genesis := string(b)

	response = JoinNetworkResponse{"619", "8080", genesis}

	return
}

func (nsi NodeServiceImpl) JoinNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkRequest
	_ = json.NewDecoder(r.Body).Decode(&request)

	response := nsi.JoinNetwork(request)

	json.NewEncoder(w).Encode(response)
}


