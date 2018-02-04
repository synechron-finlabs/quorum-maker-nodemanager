package contractclient

import (
	"synechron.com/NodeManagerGo/client"
	"synechron.com/NodeManagerGo/contracthandler"
)

const registerNodeFunSig = "0x39817edb"
const updateNodeFunSig = "0x5eec2813"
const getNodeDetailsFunSig = "0x5afbba47"

type NodeDetails struct {
	Enode string `json:"enodeid,omitempty"`
	Name  string `json:"name,omitempty"`
	Role  string `json:"role,omitempty"`
}
type NetworkMapContractClient struct {
	client.EthClient
	contractParam contracthandler.ContractParam
}

type GetNodeDetailsParam int

func (nmc *NetworkMapContractClient) RegisterNode(enode string, role string, name string) string {

	nd := NodeDetails{enode, role, name}
	return nmc.SendTransaction(nmc.contractParam, RegisterUpdateNodeFuncHandler{nd, registerNodeFunSig})

}

func (nmc *NetworkMapContractClient) GetNodeDetails(i int) (NodeDetails) {

	encoderDecoder := GetNodeDetailsFuncHandler{index: i, funcSig: getNodeDetailsFunSig}
	nmc.EthCall(nmc.contractParam, encoderDecoder, &encoderDecoder)

	return encoderDecoder.result
}
func (nmc *NetworkMapContractClient) UpdateNode(enode string, role string, name string) string {

	nd := NodeDetails{enode, role, name}
	return nmc.SendTransaction(nmc.contractParam, RegisterUpdateNodeFuncHandler{nd, updateNodeFunSig})
}

type RegisterUpdateNodeFuncHandler struct {
	nd      NodeDetails
	funcSig string
}

func (h RegisterUpdateNodeFuncHandler) Encode() string {

	sig := "string,string,string"

	param := []interface{}{h.nd.Enode, h.nd.Name, h.nd.Role}

	return h.funcSig + contracthandler.FunctionProcessor{sig, param, ""}.GetData()
}

type GetNodeDetailsFuncHandler struct {
	index   int
	funcSig string
	result  NodeDetails
}

func (g *GetNodeDetailsFuncHandler) Decode(r string) {

	sig := "string,string,string"

	resultArray := contracthandler.FunctionProcessor{sig, nil, r}.GetResults()

	g.result = NodeDetails{resultArray[0].(string), resultArray[1].(string), resultArray[2].(string)}
}

func (g GetNodeDetailsFuncHandler) Encode() string {

	sig := "uint256"

	param := []interface{}{g.index}

	return g.funcSig + contracthandler.FunctionProcessor{sig, param, ""}.GetData()
}
