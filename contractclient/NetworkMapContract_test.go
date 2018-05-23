package contractclient

import (
	"testing"
	"synechron.com/NodeManagerGo/contracthandler"
	"synechron.com/NodeManagerGo/client"
	"fmt"
	"synechron.com/NodeManagerGo/util"
	"time"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"os"
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

	ec := client.EthClient{"http://localhost:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	txRec := nmc.RegisterNode(
		"RBC1",
		"Bank",
		"0x4308aa060c5e193191ea96650a3c7b44ef1e9090",
		"0x4308aa060c5e193191ea96650a3c7b44ef1e9091",
		)

	if txRec == "" {
		t.Error("Error Registering Node")
	}
}

//func TestUpdateNode(t *testing.T) {
//
//	cp := getContractParam()
//
//	ec := client.EthClient{"http://localhost:22000"}
//
//	nmc := NetworkMapContractClient{ec, cp}
//
//	txRec := nmc.UpdateNode(
//		"c5f4b39a1c40c5affc99ec6f7be64e7c20d78c96ac55f20ee1156ce87175732a9c2b518aa6897f1590ea78b911be0c6a524d8496a420107651251048332bb04e",
//		"BB&T",
//		"Bank")
//
//	if txRec == "" {
//		t.Error("Error Updating Node")
//	}
//}

func TestGetNodeDetails(t *testing.T) {

	defer util.TotalTime(time.Now().Nanosecond())
	cp := getContractParam()

	ec := client.EthClient{"http://localhost:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	nd := nmc.GetNodeDetails(0)

	fmt.Println(nd.Name)
}

func TestGetNodeDetailsList(t *testing.T) {

	defer util.TotalTime(time.Now().Nanosecond())
	cp := getContractParam()

	ec := client.EthClient{"http://localhost:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	for _, nd := range nmc.GetNodeDetailsList() {
		fmt.Println(nd.Name)
	}

}

func getContractParam() contracthandler.ContractParam {
	return contracthandler.ContractParam{
		"0x044802df9659bbfa78615dc2af3d4740464cb714",
		"0x32072a0bfac753633914661def1c0ae31839dc28",
		"",
		nil,
	}
}
