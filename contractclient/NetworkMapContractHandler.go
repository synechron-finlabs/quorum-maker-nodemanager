package contractclient

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
)

type NodeDetailsSelf struct {
	Name      string `json:"nodeName,omitempty"`
	Role      string `json:"role,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	Enode     string `json:"enode,omitempty"`
	IP        string `json:"ip,omitempty"`
	ID        string `json:"id,omitempty"`
	Self      string `json:"self,omitempty"`
}

func (nms *NetworkMapContractClient) UpdateNodeRequestsHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	var request NodeDetails
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.Enode
	role := request.Role
	nodeName := request.Name
	publickey := request.PublicKey
	ip := request.IP
	id := request.ID
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)

	response := nms.UpdateNode(nodeName, role, publickey, enode, ip, id, coinbase, contractAdd, "", nil)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) RegisterNodeRequestHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	var request NodeDetails
	_ = json.NewDecoder(r.Body).Decode(&request)

	enode := request.Enode
	role := request.Role
	nodeName := request.Name
	publickey := request.PublicKey
	ip := request.IP
	id := request.ID
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)

	response := nms.RegisterNode(nodeName, role, publickey, enode, ip, id, coinbase, contractAdd, "", nil)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) GetNodeDetailsResponseHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	params := mux.Vars(r)
	index, err := strconv.ParseInt(params["index"], 10, 64)
	i := int(index)
	if err != nil {
		fmt.Println(err)
	}
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)

	response := nms.GetNodeDetails(i, coinbase, contractAdd, "", nil)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) GetNodeListResponseHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	response := nms.GetNodeDetailsList(coinbase, contractAdd, "", nil)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) GetNodeListSelfResponseHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	contractAdd := util.MustGetString("CONTRACT_ADD", p)
	nodename := util.MustGetString("NODENAME", p)
	nodeList := nms.GetNodeDetailsList(coinbase, contractAdd, "", nil)
	response := make([]NodeDetailsSelf, len(nodeList))
	for i := 0; i < len(nodeList); i++ {
		response[i].ID = nodeList[i].ID
		response[i].IP = nodeList[i].IP
		response[i].PublicKey = nodeList[i].PublicKey
		response[i].Enode = nodeList[i].Enode
		response[i].Role = nodeList[i].Role
		response[i].Name = nodeList[i].Name
		if nodeList[i].Name == nodename {
			response[i].Self = "true"
		} else {
			response[i].Self = "false"
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}
