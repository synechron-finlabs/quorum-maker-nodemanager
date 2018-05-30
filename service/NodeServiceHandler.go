package service

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	//"bufio"
	//"os"
	"strings"
	"io"
	"io/ioutil"
	"bytes"
	"time"
	"github.com/magiconair/properties"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"log"
)

var pendCount = 0
var nameMap = map[string]string{}
var peerMap = map[string]string{}
var channelMap = make(map[string](chan string))

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
	//recipients := strings.Split(mailServerConfig.RecipientList, ",")

	b, err := ioutil.ReadFile("/root/quorum-maker/JoinRequestTemplate.txt")

	if err != nil {
		log.Fatal(err)
	}

	mailCont := string(b)
	mailCont = strings.Replace(mailCont, "\n", "", -1)

	p := properties.MustLoadFile("/home/setup.conf", properties.UTF8)
	recipientList := util.MustGetString("RECIPIENTLIST", p)
	recipients := strings.Split(recipientList, ",")
	for i := 0; i < len(recipients); i++ {
		message := fmt.Sprintf(mailCont, nodename, enode, foreignIP)
		nsi.sendMail(mailServerConfig.Host, mailServerConfig.Port, mailServerConfig.Username, mailServerConfig.Password, "Incoming Join Request", message, recipients[i])
	}
	var cUIresp = make(chan string, 1)
	channelMap[enode] = cUIresp
	nameMap[enode] = nodename
	if peerMap[enode] == "" {
		peerMap[enode] = "PENDING"
	}

	//go func() {
	//	fmt.Println("Request for Joining this network from", nodename, "with Enode", enode, "from IP", foreignIP, "Do you approve ? y/N")
	//
	//	reader := bufio.NewReader(os.Stdin)
	//	reply, _ := reader.ReadString('\n')
	//	stop := timer.Stop()
	//	if stop {
	//		//fmt.Println("Timer stopped: Ticket ", ticket)
	//	}
	//	reader.Reset(os.Stdin)
	//	reply = strings.TrimSuffix(reply, "\n")
	//	if reply == "y" || reply == "Y" {
	//		peerMap[enode] = "YES"
	//		response := nsi.getGenesis(nsi.Url)
	//		json.NewEncoder(w).Encode(response)
	//	} else {
	//		peerMap[enode] = "NO"
	//		w.WriteHeader(http.StatusForbidden)
	//		w.Write([]byte("Access denied"))
	//	}
	//	cCLI <- fmt.Sprintf("CLI resp: Ticket %d", ticket)
	//}()

	select {
	//case <-cCLI:
	//fmt.Println(resCLI)
	case uiResp := <-cUIresp:
		if uiResp == "YES" {
			response := nsi.getGenesis(nsi.Url)
			json.NewEncoder(w).Encode(response)
		} else if uiResp == "NO" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access denied"))
		} else {
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Pending user response"))
		}
		//fmt.Println(resUI)
	case <-time.After(time.Second * 300):
		//fmt.Println("Response Timed Out")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Pending user response"))
		//fmt.Println(resTimer)
	}
}

func (nsi *NodeServiceImpl) JoinRequestResponseHandler(w http.ResponseWriter, r *http.Request) {
	var request JoinNetworkResponse
	_ = json.NewDecoder(r.Body).Decode(&request)
	enode := request.EnodeID
	status := request.Status
	response := nsi.joinRequestResponse(enode, status)
	channelMap[enode] <- status
	delete(channelMap, enode);
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
	response := nsi.emailServerConfig(request.Host, request.Port, request.Username, request.Password, request.RecipientList, nsi.Url)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	response := "Options Handled"
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}
