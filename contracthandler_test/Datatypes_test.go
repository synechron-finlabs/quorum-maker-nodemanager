package contracthandler_test

import (
	"testing"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
)

func TestIsSuported(t *testing.T) {

	for _, sig := range getDatatypesTestData() {

		if !contracthandler.IsSupported(sig) {
			t.Error(sig + " Not suported")
		}
	}

}


func getDatatypesTestData() []string {

	data := []string {
		"Func(uint256)",
		"uint256",
		"uint256,address",
		"uint256,address,",
		"",
		"uint256,address,bytes",
		"createConfig(bytes32[],uint256,bytes32[3],uint256)",
		"createConfig(bytes32[],uint256,bytes32[3],uint256)",
		"createTrade(uint256,string,bool,uint256,uint256)",
		"getConfig()",
		"updateConfig(address,bytes32[],uint256,bytes32,uint256)",
		"updateConfig(address[],bytes32[],address[10],bytes32,uint256)",
	}

	return data
}