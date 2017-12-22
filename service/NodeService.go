package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"bytes"
	"os/exec"
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
	GetGenesis  () GetGenesisResponse
	JoinNetwork (request JoinNetworkRequest) string
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

	var outnet bytes.Buffer
	cmd = exec.Command("./get_netid.sh")
	cmd.Stdout = &outnet
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	netid := outnet.String()

	b, err := ioutil.ReadFile("/home/node/genesis.json")

	if err != nil {
		log.Fatal(err)
	}

	genesis := string(b)
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
	return
}

func (nsi *NodeServiceImpl) JoinNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	response := nsi.JoinNetwork(enode)
	json.NewEncoder(w).Encode(response)
}