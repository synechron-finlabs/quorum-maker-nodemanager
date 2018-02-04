package contractclient

import (
	"testing"
	"synechron.com/NodeManagerGo/contracthandler"
	"synechron.com/NodeManagerGo/client"
	"fmt"
)

func getContractParam() contracthandler.ContractParam {
	return contracthandler.ContractParam{
		"0x2c049a350bc1284a662de7296d79c8c486867bdc",
		"0xe8160467b4f498e4cec391c9dbee74c7bd506acf",
		"",
		[]string{
			"GzNM4wJ+eJdJU+PwNAfABTo99zR6U50SbFcO8jdxPGk=",
			"gxkUoQw9hhvTWq2fk5UJZWTpYUl4SMxkrfUAJrPjBg8=",
		},
	}
}

func TestRegisterNode(t *testing.T) {

	cp := getContractParam()

	ec := client.EthClient{"http://localhost:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	txRec := nmc.RegisterNode(
		"10d93926bcd78f37dbbb2c95d65e7bc4c723a75e66fe60317525f0607c7111d3d9088c8a33c944c3e2e24b3281115a8f688312cee61573119e47555e8cd31e30",
		"JPM",
		"Custodian")

	if txRec == "" {
		t.Error("Error Registering Node")
	}
}

func TestGetNodeDetails(t *testing.T)  {
	cp := getContractParam()

	ec := client.EthClient{"http://localhost:22000"}

	nmc := NetworkMapContractClient{ec, cp}

	nd := nmc.GetNodeDetails(1)

	fmt.Println(nd)
}
