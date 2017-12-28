package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/ybbus/jsonrpc"
	"strings"
	"github.com/magiconair/properties"
	"regexp"
	"fmt"
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
	JoinNetwork (request JoinNetworkRequest) string
}

type NodeServiceImpl struct {
	Url string
}

func (nsi *NodeServiceImpl) GetGenesis() (response GetGenesisResponse) {

	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	var filename string
	constl := p.MustGetString("CONSTELLATION_PORT")

	//Alternate regex that can be used is (start_)(\w)*.sh
	r, _ := regexp.Compile("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]")
	files, err := ioutil.ReadDir("/home/node")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		match, _ := regexp.MatchString("[s][t][a][r][t][_][A-Za-z0-9]*[.][s][h]", f.Name())
		if(match) {
			filename = r.FindString(f.Name())
		}
	}

	filepath := fmt.Sprint("/home/node/", filename)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(content), "\n")

	netidline := lines[3]

	lines = strings.Split(string(netidline), "=")

	netid := lines[1]

	constl = strings.TrimSuffix(constl, "\n")

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

func (nsi *NodeServiceImpl) JoinNetwork(request string) (int) {
	rpcClient := jsonrpc.NewRPCClient(nsi.Url)
	response, err := rpcClient.Call("raft_addPeer",request)
	fmt.Println(response)
	var raftid int
	err = response.GetObject(&raftid)
	if err != nil {
		log.Fatal(err)
	}
	return raftid
}

func (nsi *NodeServiceImpl) JoinNetworkHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	response := nsi.JoinNetwork(enode)
	json.NewEncoder(w).Encode(response)
}