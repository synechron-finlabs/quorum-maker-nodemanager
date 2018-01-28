package service

import (
	"synechron.com/NodeManagerGo/client"
	"fmt"
	"strconv"
)

const registerNodeFunSig = "0xffe0abeb"

const getNodeDetailsFunSig = "0x5afbba47"

type NetworkMapService struct {
	client.EthClient
	client.ContractParam
}
type NodeDetails struct {
	Enode string
	Name  string
	Role  string
}

type GetNodeDetailsParam int

func (nms *NetworkMapService) RegisterNode(nd NodeDetails) string {

	return nms.SendTransaction(nms.ContractParam, nd)

}

func (nms *NetworkMapService) GetNodeDetails(i int) NodeDetails {
	nd := NodeDetails{}

	nms.EthCall(nms.ContractParam, GetNodeDetailsParam(i), &nd)

	return nd
}

func (nd NodeDetails) Encode() string {

	return registerNodeFunSig + client.EncodeAndPad(nd.Enode) + client.EncodeAndPad(nd.Name) + client.EncodeAndPad(nd.Role)
}

func (nd *NodeDetails) Decode(r string) {

	b := []byte(r)

	nd.Enode = client.Decode(client.Field(b, 0))

	nd.Name = client.Decode(client.Field(b, 1))

	nd.Role = client.Decode(client.Field(b, 2))

}

func (p GetNodeDetailsParam) Encode() string {

	i, _ := strconv.Atoi(fmt.Sprintf("%v", p))

	return getNodeDetailsFunSig + client.EncodeAndPad(i)
}
