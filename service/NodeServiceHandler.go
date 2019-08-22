package service

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"strings"
	"io"
	"io/ioutil"
	"bytes"
	"time"
	"github.com/magiconair/properties"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	log "github.com/sirupsen/logrus"
	"os"
)

type contractJSON struct {
	Abi       []interface{} `json:"abi"`
	Interface []interface{} `json:"interface"`
	Bytecode  string        `json:"bytecode"`
}

type genesisJSON struct {
	Config configField `json:"config"`
}

type configField struct {
	ChainId int `json:"chainId"`
}

type accountPassword struct {
	Password string `json:"password"`
}

type connectedIP struct {
	IP          string `json:"ip"`
	Whitelisted bool   `json:"whitelisted"`
	Count       int    `json:"count"`
}

type IPList struct {
	WhiteList     []string      `json:"whiteList"`
	ConnectedList []connectedIP `json:"connectedList"`
}

var pendCount = 0
var whiteList []string
var allowedIPs = map[string]bool{}
var nameMap = map[string]string{}
var peerMap = map[string]string{}
var channelMap = make(map[string](chan string))

func (nsi *NodeServiceImpl) IPWhitelister() {
	go func() {
		if _, err := os.Stat("/root/quorum-maker/contracts/.whiteList"); os.IsNotExist(err) {
			util.CreateFile("/root/quorum-maker/contracts/.whiteList")
		}
		whitelistedIPs, _ := util.File2lines("/root/quorum-maker/contracts/.whiteList")
		whiteList = append(whiteList, whitelistedIPs...)
		for _, ip := range whitelistedIPs {
			allowedIPs[ip] = true
		}
		log.Info("Adding whitelisted IPs")
	}()
}

func (nsi *NodeServiceImpl) UpdateWhitelistHandler(w http.ResponseWriter, r *http.Request) {
	var ipList []string
	_ = json.NewDecoder(r.Body).Decode(&ipList)
	whiteList = whiteList[:0]
	for k := range allowedIPs {
		delete(allowedIPs, k)
	}
	for _, ip := range ipList {
		allowedIPs[ip] = true
	}
	whiteList = append(whiteList, ipList ...)
	response := nsi.updateWhitelist(ipList)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetWhitelistedIPsHandler(w http.ResponseWriter, r *http.Request) {
	var ipList IPList
	var connectedIPList []connectedIP
	var whiteListedIPs []string
	activeIPs := nsi.getNodeIPs(nsi.Url)
	for _, ip := range activeIPs {
		var connected connectedIP
		connected.IP = ip.IP
		if allowedIPs[ip.IP] {
			connected.Whitelisted = true
		}
		connected.Count = ip.Count
		connectedIPList = append(connectedIPList, connected)
	}
	for _, ip := range whiteList {
		whiteListedIPs = append(whiteListedIPs, ip)
	}
	ipList.ConnectedList = connectedIPList
	ipList.WhiteList = whiteListedIPs
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(ipList)
}

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
	if allowedIPs[foreignIP] {
		peerMap[enode] = "YES"
		exists := util.PropertyExists("RECIPIENTLIST", "/home/setup.conf")
		if exists != "" {
			go func() {
				b, err := ioutil.ReadFile("/root/quorum-maker/JoinRequestTemplate.txt")

				if err != nil {
					log.Println(err)
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
			}()
		}
	}
	log.Info(fmt.Sprint("Join request received from node: ", nodename, " with IP: ", foreignIP, " and enode: ", enode))
	if peerMap[enode] == "YES" {
		response := nsi.getGenesis(nsi.Url)
		json.NewEncoder(w).Encode(response)
	} else if peerMap[enode] == "NO" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Access denied"))
	} else {
		exists := util.PropertyExists("RECIPIENTLIST", "/home/setup.conf")
		if exists != "" {
			go func() {
				b, err := ioutil.ReadFile("/root/quorum-maker/JoinRequestTemplate.txt")

				if err != nil {
					log.Println(err)
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
			}()
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
		fmt.Println(err)
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
			fmt.Println(err)
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")

		fileNames[i] = name[0] + ".sol"

		io.Copy(&Buf, file)

		contents := Buf.String()

		fileContent := []byte(contents)
		err = ioutil.WriteFile("./"+name[0]+".sol", fileContent, 0775)
		if err != nil {
			fmt.Println(err)
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

//func (nsi *NodeServiceImpl) LogsHandler(w http.ResponseWriter, r *http.Request) {
//	response := nsi.logs()
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	json.NewEncoder(w).Encode(response)
//}

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

func (nsi *NodeServiceImpl) GetChartDataHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.GetChartData(nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetContractListHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.ContractList()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetContractCountHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.ContractCount()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) ContractDetailsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var Buf bytes.Buffer
	contractAddress := r.FormValue("address")
	contractName := r.FormValue("name")
	description := r.FormValue("description")
	file, header, err := r.FormFile("abi")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	io.Copy(&Buf, file)
	content := Buf.String()

	var jsonContent contractJSON

	json.Unmarshal([]byte(content), &jsonContent)
	abiContent, _ := json.Marshal(jsonContent.Abi)
	abiString := make([]string, len(abiContent))
	for i := 0; i < len(abiContent); i++ {
		abiString[i] = string(abiContent[i])
	}
	abiData := fmt.Sprint(strings.Join(abiString, ""))

	interfaceContent, _ := json.Marshal(jsonContent.Interface)
	interfaceString := make([]string, len(interfaceContent))
	for i := 0; i < len(interfaceContent); i++ {
		interfaceString[i] = string(interfaceContent[i])
	}
	interfaceData := fmt.Sprint(strings.Join(interfaceString, ""))
	bytecodeData := jsonContent.Bytecode
	var data string
	if len(abiData) != 4 {
		data = abiData
	} else if len(interfaceData) != 4 {
		data = interfaceData
	} else {
		data = content
		data = strings.Replace(data, "\n", "", -1)
	}

	jsonString := util.ComposeJSON(data, bytecodeData, contractAddress)

	path := "./contracts"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	}
	path = "./contracts/" + contractAddress + "_" + name[0]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	}

	filePath := path + "/" + name[0] + ".json"
	jsByte := []byte(jsonString)
	err = ioutil.WriteFile(filePath, jsByte, 0775)
	if err != nil {
		fmt.Println(err)
	}

	Buf.Reset()
	response := nsi.updateContractDetails(contractAddress, contractName, data, description)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) AttachedNodeDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var successResponse SuccessResponse
	var Buf bytes.Buffer
	gethLogsDirectory := r.FormValue("gethPath")
	constellationLogsDirectory := r.FormValue("constellationPath")

	file, _, err := r.FormFile("genesis")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	io.Copy(&Buf, file)
	content := Buf.String()

	filePath := "/home/node/genesis.json"
	jsByte := []byte(content)
	err = ioutil.WriteFile(filePath, jsByte, 0775)
	if err != nil {
		fmt.Println(err)
	}

	var jsonContent genesisJSON
	json.Unmarshal([]byte(content), &jsonContent)
	chainIdAppend := fmt.Sprint("NETWORK_ID=", jsonContent.Config.ChainId, "\n")
	util.AppendStringToFile("/home/setup.conf", chainIdAppend)
	util.InsertStringToFile("/home/start.sh", "	   -v "+gethLogsDirectory+":/home/node/qdata/gethLogs \\\n", 13)
	util.InsertStringToFile("/home/start.sh", "	   -v "+constellationLogsDirectory+":/home/node/qdata/constellationLogs \\\n", 13)

	Buf.Reset()
	fmt.Println("Updates have been saved. Please press Ctrl+C to exit from this container and run start.sh to apply changes")
	state := currentState()
	if state == "NI" {
		util.DeleteProperty("STATE=NI", "/home/setup.conf")
		stateInitialized := fmt.Sprint("STATE=I\n")
		util.AppendStringToFile("/home/setup.conf", stateInitialized)
	}
	successResponse.Status = "Updates have been saved. Please press Ctrl+C from CLI to exit from this container and run start.sh to apply changes"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(successResponse)
}

func (nsi *NodeServiceImpl) InitializationHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.returnCurrentInitializationState()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var request accountPassword
	_ = json.NewDecoder(r.Body).Decode(&request)
	password := request.Password
	response := nsi.createAccount(password, nsi.Url)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-Name, Cache-Control")
	json.NewEncoder(w).Encode(response)
}

func (nsi *NodeServiceImpl) GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	response := nsi.getAccounts(nsi.Url)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
