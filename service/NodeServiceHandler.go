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
	"io"
	"io/ioutil"
	"bytes"
	"time"
)

var context = 0
var ticketCheck = 0

var ticketMap = map[string]int{}
var nameMap = map[string]string{}
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
	nodename := request.Nodename

	message := fmt.Sprint("Request for joining network has come in from node ", nodename, " with enode ", enode, " and ip-address ", foreignIP)
	nsi.sendMail(mailServerConfig.Host, mailServerConfig.Port, mailServerConfig.Username, mailServerConfig.Password, "Incoming Join Request", message)
	var cUIresp = make(chan string, 1)
	var cTimer = make(chan string, 1)
	var cCLI = make(chan string, 1)
	nameMap[enode] = nodename
	if peerMap[enode] == "" {
		peerMap[enode] = "PENDING"
	}

	context++
	ticket := context
	fmt.Println(context)

	ticketMap[enode] = ticket

	fmt.Println("context:", context)
	fmt.Println("ticket:", ticket)
	timer := time.NewTimer(300 * time.Second)
	go func() {
		<-timer.C
		fmt.Println("Timer expired Ticket ", ticket)
		cTimer <- fmt.Sprintf("Timer resp: Ticket %d", ticket)
	}()

	go func() {
		for ticket != ticketCheck {
		}
		stop := timer.Stop()
		if stop {
			fmt.Println("Timer stopped: Ticket ", ticketCheck)
		}
		cUIresp <- fmt.Sprintf("UI resp: Ticket %d", ticket)
	}()

	go func() {
		fmt.Println("Request for Joining this network from", nodename, "with Enode", enode, "from IP", foreignIP, "Do you approve ? y/N")

		reader := bufio.NewReader(os.Stdin)
		reply, _ := reader.ReadString('\n')
		stop := timer.Stop()
		if stop {
			fmt.Println("Timer stopped: Ticket ", ticket)
		}
		reader.Reset(os.Stdin)
		reply = strings.TrimSuffix(reply, "\n")
		if reply == "y" || reply == "Y" {
			peerMap[enode] = "YES"
			response := nsi.getGenesis(nsi.Url)
			json.NewEncoder(w).Encode(response)
		} else {
			peerMap[enode] = "NO"
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
		}
		cCLI <- fmt.Sprintf("CLI resp: Ticket %d", ticket)
	}()

	select {
	case resCLI := <-cCLI:
		fmt.Println(resCLI)
	case resUI := <-cUIresp:
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
	case resTimer := <-cTimer:
		fmt.Println("Response Timed Out")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Pending user response"))
		fmt.Println(resTimer)
	}
}

func (nsi *NodeServiceImpl) JoinRequestResponseHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkResponse
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	status := request.Status
	response := nsi.joinRequestResponse(enode, status)
	ticketCheck = ticketMap[enode]
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetCurrentNodeHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.getCurrentNode(nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetOtherPeerHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["peer_id"] == "all" {
		response := nsi.getOtherPeers(nsi.Url)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	} else {
		response := nsi.getOtherPeer(params["peer_id"], nsi.Url)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	}
}

func (nsi *NodeServiceImpl) GetLatestBlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	count := r.FormValue("number")
	reference := r.FormValue("reference")
	response := nsi.getLatestBlockInfo(count, reference, nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetBlockInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	block, err := strconv.ParseInt(params["block_no"], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	response := nsi.getBlockInfo(block, nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetTransactionInfoHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["txn_hash"] == "pending" {
		response := nsi.getPendingTransactions(nsi.Url)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	} else {
		response := nsi.getTransactionInfo(params["txn_hash"], nsi.Url)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	}
}

func (nsi *NodeServiceImpl) GetLatestTransactionInfoHandler(w http.ResponseWriter, r *http.Request) {
	count := r.FormValue("number")
	response := nsi.getLatestTransactionInfo(count, nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetTransactionReceiptHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := nsi.getTransactionReceipt(params["txn_hash"], nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) PendingJoinRequestsHandler(w http.ResponseWriter, r *http.Request) {
	size := 0
	for key := range peerMap {
		if peerMap[key] == "PENDING" {
			size++
		}
	}
	pendingRequests := make([]PendingRequests, size)
	i := 1
	for key := range peerMap {
		if peerMap[key] == "PENDING" {
			var enodeString []string
			var ipString []string
			pendingRequests[i-1].NodeName = nameMap[key]
			pendingRequests[i-1].Enode = key
			enode := strings.TrimPrefix(key, "enode://")
			enodeString = strings.Split(enode, "@")
			enodeVal := enodeString[0]
			ipString = strings.Split(enodeString[1], ":")
			ip := ipString[0]
			pendingRequests[i-1].Message = fmt.Sprintf("Request for joining network from node %s", nameMap[key])
			pendingRequests[i-1].EnodeID = fmt.Sprintf("%s", enodeVal)
			pendingRequests[i-1].IP = fmt.Sprintf("%s", ip)
			i++
		}
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(pendingRequests)
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
		keyVal := "file" + strconv.Itoa(i+1)

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
		err = ioutil.WriteFile("./"+name[0]+".sol", fileContent, 0775)
		if err != nil {
			panic(err)
		}

		Buf.Reset()
	}

	response := nsi.deployContract(publicKeys, fileNames, private, nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) CreateNetworkScriptCallHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateNetworkScriptArgs
	_ = json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(request)
	response := nsi.createNetworkScriptCall(request.Nodename, request.CurrentIP, request.RPCPort, request.WhisperPort, request.ConstellationPort, request.RaftPort, request.NodeManagerPort)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) JoinNetworkScriptCallHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkScriptArgs
	_ = json.NewDecoder(r.Body).Decode(&request)
	fmt.Println(request)
	response := nsi.joinRequestResponseCall(request.Nodename, request.CurrentIP, request.RPCPort, request.WhisperPort, request.ConstellationPort, request.RaftPort, request.NodeManagerPort, request.MasterNodeManagerPort, request.MasterIP)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) ResetHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.resetCurrentNode()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) RestartHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.restartCurrentNode()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) NodeListHandler(w http.ResponseWriter, r *http.Request) {
	nodeList := make([]NodeList, 3)
	nodeList[0].IP = "10.34.15.90"
	nodeList[0].Enode = "5089ccd054af243df3cc7eabed10192c147cd4af611f523bd5447ee39bdd8670768f522232537456de8f15ec64cde8ccc758956e3875f605fe4979f329820b26"
	nodeList[0].NodeName = "Markit"
	nodeList[0].PublicKey = "0OKTpgAGfXN9x8N4K6bIiupI4Lhb/MZaqcD/f4FTMHI="
	nodeList[0].Role = "Custodian"
	nodeList[1].IP = "10.34.15.90"
	nodeList[1].Enode = "e27a60a7f1abca56819574a99acad6b8d40a1866b961ee5ebb8a3f5b5c1b74f7a2abbf5ac3f0bd7f1409cccc0df0f1bdb60889d7ddc4150d45021b940f3440c6"
	nodeList[1].NodeName = "Synechron"
	nodeList[1].PublicKey = "R1fOFUfzBbSVaXEYecrlo9rENW0dam0kmaA2pasGM14="
	nodeList[1].Role = "Custodian"
	nodeList[2].IP = "10.34.15.90"
	nodeList[2].Enode = "964d6583a58ddb04f9773d0a013860ff7d36fd18e4178c6ee0a70c2ee08bb89e047cf734810f743fda139c9a5ec1f2bb6ab97e9a054f55dc4e99610dc1aeb57d"
	nodeList[2].NodeName = "JPMC"
	nodeList[2].PublicKey = "Er5J8G+jXQA9O2eu7YdhkraYM+j+O5ArnMSZ24PpLQY="
	nodeList[2].Role = "Custodian"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(nodeList)
}

func (nsi *NodeServiceImpl) LatestBlockHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.latestBlockDetails(nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) LatencyHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.latency(nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) LogsHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.logs()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) TransactionSearchHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	response := nsi.transactionSearchDetails(params["txn_hash"], nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) MailServerConfigHandler(w http.ResponseWriter, r *http.Request) {
	var request MailServerConfig
	_ = json.NewDecoder(r.Body).Decode(&request)
	response := nsi.emailServerConfig(request.Host, request.Port, request.Username, request.Password, nsi.Url)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) PublicKeysHandler(w http.ResponseWriter, r *http.Request) {
	nodeList := make([]NodeList, 3)
	nodeList[0].NodeName = "Markit"
	nodeList[0].PublicKey = "0OKTpgAGfXN9x8N4K6bIiupI4Lhb/MZaqcD/f4FTMHI="
	nodeList[1].NodeName = "Synechron"
	nodeList[1].PublicKey = "R1fOFUfzBbSVaXEYecrlo9rENW0dam0kmaA2pasGM14="
	nodeList[2].NodeName = "JPMC"
	nodeList[2].PublicKey = "Er5J8G+jXQA9O2eu7YdhkraYM+j+O5ArnMSZ24PpLQY="
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(nodeList)
}
