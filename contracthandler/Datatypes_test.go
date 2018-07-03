package contracthandler

import "testing"

func TestIsSuported(t *testing.T) {

	for _, sig := range getDatatypesTestData() {

		if !IsSuported(sig) {
			t.Error(sig + " Not suported")
		}
	}

}


func getDatatypesTestData() []string {

	data := []string {
		"Func(uint256)",
		"createConfig(bytes32[],uint256,bytes32[3],uint256)",
		"createTrade(uint256,string,bool,uint256,uint256)",
		"getConfig()",
		"updateConfig(address,bytes32[],uint256,bytes32,uint256)",
		"updateConfig(address[],bytes32[],address[10],bytes32,uint256)",
	}

	return data
}