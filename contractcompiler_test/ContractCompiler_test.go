package contractcompiler

import (
	"encoding/json"
	"fmt"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contractcompiler"
	"testing"
)

func TestCompile(t *testing.T) {

	filename := "NetworkManagerContract.sol"

	ethClient := client.EthClient{"http://localhost:22000"}

	contractcompiler.Compile(filename, []string{}, ethClient, false, map[string]string{}, map[string]string{}, map[string]string{})
}

type Person struct {
	Name    interface{} `json:"name"`         // field tag for Name
	Aadhar  int         `json:"aadhar"`       // field tag for Aadhar
	Street  string      `json:"street"`       // field tag for Street
	HouseNo int         `json:"house_number"` // field tag for HouseNO
}

func TestJson(t *testing.T) {

	var p Person

	p.Name = map[string]string{"Name": "Bob", "Food": "Pickle"}
	p.Aadhar = 1234123412341234
	p.Street = "XYZ"
	p.HouseNo = 10

	fmt.Println(p)

	// returns []byte which is p in JSON form.
	jsonStr, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(jsonStr))
}

func TestComplexJson(t *testing.T) {
	s := `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"nodeName","type":"string"},{"indexed":false,"internalType":"string","name":"role","type":"string"},{"indexed":false,"internalType":"string","name":"publickey","type":"string"},{"indexed":false,"internalType":"string","name":"enode","type":"string"},{"indexed":false,"internalType":"string","name":"ip","type":"string"},{"indexed":false,"internalType":"string","name":"id","type":"string"}],"name":"print","type":"event"},{"constant":true,"inputs":[{"internalType":"uint16","name":"_index","type":"uint16"}],"name":"getNodeDetails","outputs":[{"internalType":"string","name":"n","type":"string"},{"internalType":"string","name":"r","type":"string"},{"internalType":"string","name":"p","type":"string"},{"internalType":"string","name":"ip","type":"string"},{"internalType":"string","name":"e","type":"string"},{"internalType":"string","name":"id","type":"string"},{"internalType":"uint256","name":"i","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"uint16","name":"i","type":"uint16"}],"name":"getNodeList","outputs":[{"internalType":"string","name":"n","type":"string"},{"internalType":"string","name":"r","type":"string"},{"internalType":"string","name":"p","type":"string"},{"internalType":"string","name":"ip","type":"string"},{"internalType":"string","name":"e","type":"string"},{"internalType":"string","name":"id","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getNodesCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"string","name":"n","type":"string"},{"internalType":"string","name":"r","type":"string"},{"internalType":"string","name":"p","type":"string"},{"internalType":"string","name":"e","type":"string"},{"internalType":"string","name":"ip","type":"string"},{"internalType":"string","name":"id","type":"string"}],"name":"registerNode","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"string","name":"n","type":"string"},{"internalType":"string","name":"r","type":"string"},{"internalType":"string","name":"p","type":"string"},{"internalType":"string","name":"e","type":"string"},{"internalType":"string","name":"ip","type":"string"},{"internalType":"string","name":"id","type":"string"}],"name":"updateNode","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`

	var i interface{}

	json.Unmarshal([]byte(s), &i)

	fmt.Println(i)
}

func TestArray(t *testing.T) {
	s := []string{}

	inner(s)

	for a := range s {
		fmt.Println(a)
	}
}

func inner(t []string){
	t = append(t,"D")

	for _, a := range t {
		fmt.Printf("Inner %s\n", a)
	}
}
