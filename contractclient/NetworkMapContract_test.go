package contractclient

import (
	"fmt"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// needed to retrieve requests that arrived at httpServer for further investigation
var requestChan = make(chan *RequestData, 1)

// the request datastructure that can be retrieved for test assertions
type RequestData struct {
	request *http.Request
	body    string
}

var responseBody = ""

var httpServer *httptest.Server

// start the testhttp server and stop it when tests are finished
func TestMain(m *testing.M) {
	httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		// put request and body to channel for the client to investigate them
		//requestChan <- &RequestData{r, string(data)}

		fmt.Println(string(data))

		fmt.Fprintf(w, responseBody)
	}))
	defer httpServer.Close()

	os.Exit(m.Run())
}

func TestRegisterNode(t *testing.T) {

	cp := getContractParam()

	ec := client.EthClient{"http://10.50.0.2:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	txRec := nmc.RegisterNode(
		"RBC1",
		"Bank",
		"0x4308aa060c5e193191ea96650a3c7b44ef1e9090",
		"0x4308aa060c5e193191ea96650a3c7b44ef1e9091",
		"10.0.2.15",
		"1",
	)

	if txRec == "" {
		t.Error("Error Registering Node")
	}
}

func TestUpdateNode(t *testing.T) {

	cp := getContractParam()

	ec := client.EthClient{"http://10.50.0.2:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	txRec := nmc.UpdateNode(
		"RBC1",
		"Custodian",
		"0x4308aa060c5e193191ea96650a3c7b44ef1e9090",
		"0x4308aa060c5e193191ea96650a3c7b44ef1e9091",
		"10.0.2.15",
		"1",
	)

	if txRec == "" {
		t.Error("Error Updating Node")
	}
}

func TestGetNodeDetails(t *testing.T) {

	defer util.TotalTime(time.Now().Nanosecond())
	cp := getContractParam()

	ec := client.EthClient{"http://10.50.0.2:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	nd := nmc.GetNodeDetails(1)

	fmt.Println(nd.Name)
}

func TestGetNodeDetailsList(t *testing.T) {

	defer util.TotalTime(time.Now().Nanosecond())
	cp := getContractParam()

	ec := client.EthClient{"http://10.50.0.2:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	for _, nd := range nmc.GetNodeDetailsList() {
		fmt.Println(nd.Name)
	}

}

func getContractParam() contracthandler.ContractParam {
	return contracthandler.ContractParam{
		"0xfaf947731d6256b9f1da3e3bede1750881e21092",
		"0x56776f1d261e75253f110ac9091367956bce2f32",
		"",
		nil,
	}
}
