package contractclient

import (
	"synechron.com/NodeManagerGo/client"
	"synechron.com/NodeManagerGo/contracthandler"
)

const registerNodeFunSig = "0x3072b1b2"
const updateNodeFunSig = "0xaeffe3b7"
const getNodeDetailsFunSig = "0x0d0d583b"

type NodeDetails struct {
	Name      string `json:"nodeName,omitempty"`
	Role      string `json:"role,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	Enode     string `json:"enode,omitempty"`
	IP        string `json:"ip,omitempty"`
	ID        string `json:"id,omitempty"`
}

type NetworkMapContractClient struct {
	client.EthClient
	contractParam contracthandler.ContractParam
}

type GetNodeDetailsParam int

func (nmc *NetworkMapContractClient) RegisterNode(name string, role string, publicKey string, enode string, ip string, id string, from string, to string, passwd string, pvtFor []string) string {
	nmc.contractParam.From = from
	nmc.contractParam.To = to
	nmc.contractParam.Passwd = passwd
	nmc.contractParam.Parties = pvtFor
	nd := NodeDetails{name, role, publicKey, enode, ip, id}
	return nmc.SendTransaction(nmc.contractParam, RegisterUpdateNodeFuncHandler{nd, registerNodeFunSig})

}

func (nmc *NetworkMapContractClient) GetNodeDetails(i int, from string, to string, passwd string, pvtFor []string) (NodeDetails) {
	nmc.contractParam.From = from
	nmc.contractParam.To = to
	nmc.contractParam.Passwd = passwd
	nmc.contractParam.Parties = pvtFor
	encoderDecoder := GetNodeDetailsFuncHandler{index: i, funcSig: getNodeDetailsFunSig}
	nmc.EthCall(nmc.contractParam, encoderDecoder, &encoderDecoder)

	return encoderDecoder.result
}

func (nmc *NetworkMapContractClient) GetNodeDetailsList(from string, to string, passwd string, pvtFor []string) ([]NodeDetails) {
	nmc.contractParam.From = from
	nmc.contractParam.To = to
	nmc.contractParam.Passwd = passwd
	nmc.contractParam.Parties = pvtFor
	var list []NodeDetails

	for i := 0; true; i++ {
		encoderDecoder := GetNodeDetailsFuncHandler{index: i, funcSig: getNodeDetailsFunSig}
		nmc.EthCall(nmc.contractParam, encoderDecoder, &encoderDecoder)

		if encoderDecoder.result.Enode != "" && len(encoderDecoder.result.Enode) > 0 {
			list = append(list, encoderDecoder.result)
		} else {
			return list
		}
	}

	return list
}

func (nmc *NetworkMapContractClient) UpdateNode(name string, role string, publicKey string, enode string, ip string, id string, from string, to string, passwd string, pvtFor []string) string {
	nmc.contractParam.From = from
	nmc.contractParam.To = to
	nmc.contractParam.Passwd = passwd
	nmc.contractParam.Parties = pvtFor
	nd := NodeDetails{name, role, publicKey, enode, ip, id}
	return nmc.SendTransaction(nmc.contractParam, RegisterUpdateNodeFuncHandler{nd, updateNodeFunSig})
}

type RegisterUpdateNodeFuncHandler struct {
	nd      NodeDetails
	funcSig string
}

func (h RegisterUpdateNodeFuncHandler) Encode() string {

	sig := "string,string,string,string,string,string"

	param := []interface{}{h.nd.Name, h.nd.Role, h.nd.PublicKey, h.nd.Enode, h.nd.IP, h.nd.ID,}

	data := h.funcSig + contracthandler.FunctionProcessor{sig, param, ""}.GetData()

	return data
}

type GetNodeDetailsFuncHandler struct {
	index   int
	funcSig string
	result  NodeDetails
}

func (g *GetNodeDetailsFuncHandler) Decode(r string) {
	var nd NodeDetails

	if r == "" || len(r) < 1 {
		g.result = nd
		return
	}

	sig := "string,string,string,string,string,string,uint256"

	resultArray := contracthandler.FunctionProcessor{sig, nil, r}.GetResults()

	g.result = NodeDetails{resultArray[0].(string), resultArray[1].(string), resultArray[2].(string), resultArray[4].(string), resultArray[3].(string), resultArray[5].(string)}
}

func (g GetNodeDetailsFuncHandler) Encode() string {

	sig := "uint256"

	param := []interface{}{g.index}

	return g.funcSig + contracthandler.FunctionProcessor{sig, param, ""}.GetData()
}

type DeployContractHandler struct {
	binary string
}

func (d DeployContractHandler) Encode() string {

	return d.binary
}
