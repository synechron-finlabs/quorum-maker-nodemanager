package contractclient

import (
	"fmt"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractclient"
	"io/ioutil"
	"strings"
	"testing"
)

type test_data struct {
	contAdd    string
	abiString  string
	txnPayload string
	funcSig    string
}

func TestABIParser(t *testing.T) {
	for _, table := range getInputTestData() {
		decodedInputs, funcSig := contractclient.ABIParser(table.contAdd, table.abiString, table.txnPayload)
		if funcSig != table.funcSig {
			t.Errorf("Decoding process failed, got: %s, wanted: %s.", funcSig, table.funcSig)
		}
		fmt.Println("decoded inputs:", decodedInputs)
		fmt.Println("function signature:", funcSig)
	}
}

func getInputTestData() []test_data {
	abiString := readAbiFromDisk("../testdata/testabi")
	abiString1 := readAbiFromDisk("../testdata/testabi1")
	tables := []test_data{
		{"0x692a70d2e424a56d2c6c27aa97d1a86395877b3a", abiString, "0xbb83a3142d170000000000000000000000000000000000000000000000000000000000002b18000000000000000000000000000000000000000000000000000000000000291b000000000000000000000000000000000000000000000000000000000000", "createConfig(bytes24[3])"},
		{"0x692a70d2e424a56d2c6c27aa97d1a86395877b3a", abiString1, "0xaf6cb86a000000000000000000000000000000000000000000000000000000000000004500000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002d00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000007617273656e616c00000000000000000000000000000000000000000000000000", "createTrade(uint256,string,bool,uint256,uint256)"}, {"0x692a70d2e424a56d2c6c27aa97d1a86395877b3a", abiString, "0x204ea05300000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000452d170000000000000000000000000000000000000000000000000000000000002b18000000000000000000000000000000000000000000000000000000000000291b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002d000000000000000000000000000000000000000000000000000000000000000217000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000", ""}, {"0x692a70d2e424a56d2c6c27aa97d1a86395877b3a", abiString1, "0x204ea05300000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000452d170000000000000000000000000000000000000000000000000000000000002b18000000000000000000000000000000000000000000000000000000000000291b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002d00000000000000000000000000000000000000000000000000000000000000041700000000000000000000000000000000000000000000000000000000000000202d00000000000000000000000000000000000000000000000000000000000017000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000", "createConfig(bytes32[],uint256,bytes32[3],uint256)"}, {"0x692a70d2e424a56d2c6c27aa97d1a86395877b3a", abiString1, "0xf5810e2c000000000000000000000000692a70d2e424a56d2c6c27aa97d1a86395877b3a00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000452d17000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002d00000000000000000000000000000000000000000000000000000000000000041700000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000017000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000", "updateConfig(address,bytes32[],uint256,bytes32,uint256)"},
	}
	return tables
}

func readAbiFromDisk(filepath string) string {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}
	abiString := string(data)
	abiString = strings.Replace(abiString, "\n", "", -1)
	return abiString
}
