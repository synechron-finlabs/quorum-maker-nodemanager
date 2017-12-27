package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"bytes"
	"os/exec"
	"strings"
)

type JoinNetworkRequest struct {
	EnodeID    string `json:"enode-id,omitempty"`
	EthAccount string `json:"eth-account,omitempty"`
}

type GetGenesisResponse struct {
	ContstellationPort string `json: "contstellation-port, omitempty"`
	NetID              string `json: "netID,omitempty"`
	Genesis            string `json: "genesis, omitempty"`
}

type NodeService interface {
	GetGenesis() GetGenesisResponse
	JoinNetwork(request JoinNetworkRequest) string
}

type NodeServiceImpl struct {
}

func (nsi *NodeServiceImpl) GetGenesis() (response GetGenesisResponse) {
	cmd := exec.Command("./get_const.sh")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	constl := out.String()
	constl = strings.TrimSuffix(constl, "\n")
	var outnet bytes.Buffer
	cmd = exec.Command("./get_netid.sh")
	cmd.Stdout = &outnet
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	netid := outnet.String()
	netid = strings.TrimSuffix(netid, "\n")
	b, err := ioutil.ReadFile("/home/node/genesis.json")
	if err != nil {
		log.Fatal(err)
	}
	genesis := string(b)
	genesis = strings.Replace(genesis, "\n","",-1)
	response = GetGenesisResponse{constl, netid, genesis}
	return
}

func (nsi *NodeServiceImpl) GetGenesisHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.GetGenesis()
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) JoinNetwork(request string) (response string) {
	cmd := exec.Command("./add_peer.sh",request)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	response = out.String()
	response = strings.TrimSuffix(response, "\n")
	return
}

func (nsi *NodeServiceImpl) JoinNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	response := nsi.JoinNetwork(enode)
	json.NewEncoder(w).Encode(response)
}