package contractclient

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"net/http"
	"strconv"
)

type NodeDetailsSelf struct {
	Name      string `json:"nodeName,omitempty"`
	Role      string `json:"role,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	Enode     string `json:"enode,omitempty"`
	IP        string `json:"ip,omitempty"`
	ID        string `json:"id,omitempty"`
	Self      string `json:"self,omitempty"`
	Active    string `json:"active,omitempty"`
}

type ActiveNodes struct {
	NodeCount      int `json:"nodeCount,omitempty"`
	TotalNodeCount int `json:"totalNodeCount,omitempty"`
}

type UpdateNode struct {
	NodeName string `json:"nodeName,omitempty"`
	Role     string `json:"role,omitempty"`
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
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}

	nms.SetContractParam(cp)

	response := nms.UpdateNode(nodeName, role, publickey, enode, ip, id)
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
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}
	nms.SetContractParam(cp)

	response := nms.RegisterNode(nodeName, role, publickey, enode, ip, id)

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
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}
	nms.SetContractParam(cp)

	response := nms.GetNodeDetails(i)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) GetNodeListResponseHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}
	nms.SetContractParam(cp)

	response := nms.GetNodeDetailsList()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) GetNodeListSelfResponseHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	enode := nms.EthClient.AdminNodeInfo().ID
	//var contractAdd, nodename string
	var contractAdd string
	//existsA := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	//existsB := util.PropertyExists("NODENAME", "/home/setup.conf")
	//if existsA != "" && existsB != "" {
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
		//nodename = util.MustGetString("NODENAME", p)
	}

	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}
	nms.SetContractParam(cp)
	adminPeers := nms.EthClient.AdminPeers()

	var peerEnodes = map[string]bool{}
	for i := 0; i < len(adminPeers); i++ {
		peerEnodes[adminPeers[i].ID] = true
	}
	nodeList := nms.GetNodeDetailsList()
	response := make([]NodeDetailsSelf, len(nodeList))
	for i := 0; i < len(nodeList); i++ {
		response[i].ID = nodeList[i].ID
		response[i].IP = nodeList[i].IP
		response[i].PublicKey = nodeList[i].PublicKey
		response[i].Enode = nodeList[i].Enode
		response[i].Role = nodeList[i].Role
		response[i].Name = nodeList[i].Name
		//if nodeList[i].Name == nodename {
		if nodeList[i].Enode == enode {
			response[i].Self = "true"
			response[i].Active = "true"
		} else {
			response[i].Self = "false"
			if peerEnodes[nodeList[i].Enode] {
				response[i].Active = "true"
			} else {
				response[i].Active = "false"
			}
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) ActiveNodesHandler(w http.ResponseWriter, r *http.Request) {
	coinbase := nms.EthClient.Coinbase()
	var contractAdd string
	exists := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	if exists != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
	}

	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}
	nms.SetContractParam(cp)
	adminPeers := nms.EthClient.AdminPeers()
	activeNodes := len(adminPeers) + 1
	contractResponse := nms.GetNodeDetailsList()
	totalNodes := len(contractResponse)
	response := ActiveNodes{activeNodes, totalNodes}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	response := "Options Handled"
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nms *NetworkMapContractClient) UpdateNodeHandler(w http.ResponseWriter, r *http.Request) {
	var request UpdateNode
	_ = json.NewDecoder(r.Body).Decode(&request)
	nodeName := request.NodeName
	role := request.Role
	coinbase := nms.EthClient.Coinbase()
	enode := nms.EthClient.AdminNodeInfo().ID
	//var contractAdd, publickey, ip, id, oldName string
	var contractAdd, publickey, ip, id string
	existsA := util.PropertyExists("CONTRACT_ADD", "/home/setup.conf")
	existsB := util.PropertyExists("PUBKEY", "/home/setup.conf")
	existsC := util.PropertyExists("CURRENT_IP", "/home/setup.conf")
	existsD := util.PropertyExists("RAFT_ID", "/home/setup.conf")
	//existsE := util.PropertyExists("NODENAME", "/home/setup.conf")
	//if existsA != "" && existsB != "" && existsC != "" && existsD != "" && existsE != "" {
	if existsA != "" && existsB != "" && existsC != "" && existsD != "" {
		p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
		contractAdd = util.MustGetString("CONTRACT_ADD", p)
		publickey = util.MustGetString("PUBKEY", p)
		ip = util.MustGetString("CURRENT_IP", p)
		id = util.MustGetString("RAFT_ID", p)
		//oldName = util.MustGetString("NODENAME", p)
	}
	cp := contracthandler.ContractParam{coinbase, contractAdd, "", nil}
	nms.SetContractParam(cp)
	response := nms.UpdateNode(nodeName, role, publickey, enode, ip, id)
	//registered := fmt.Sprint("NODENAME=", nodeName, "\n")
	//util.AppendStringToFile("/home/setup.conf", registered)
	//oldProperty := fmt.Sprint("NODENAME=", oldName)
	//util.DeleteProperty(oldProperty, "/home/setup.conf")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}
