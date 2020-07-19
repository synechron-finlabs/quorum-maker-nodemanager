package contractcompiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/client"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/env"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const SOLC string = "solc"

type ContractJson struct {
	Filename        string `json:"filename"`
	Interface       string `json:"interface"`
	Bytecode        string `json:"bytecode"`
	ContractAddress string `json:"address"`
	Json            string `json:"json"`
}

type JsonFormatted struct {
	Abi      interface{} `json:"interface"`
	Bytecode string      `json:"bytecode"`
	Address  string      `json:"address"`
}

func Compile(fileName string, publicKeys []string, ethClient client.EthClient, private bool, contNameMap map[string]string, contTimeMap map[string]string, abiMap map[string]string) []ContractJson {

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	results := make([]ContractJson, 0)

	cmd := exec.Command(SOLC, "--combined-json", "abi,bin", fileName)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()

	if nil != err {
		error := ContractJson{}

		error.Filename = strings.Replace(fileName, ".sol", "", -1)
		error.Json = "Compilation Failed: JSON could not be created"
		error.Bytecode = "Compilation Failed: " + stdErr.String()
		results = append(results, error)
	}

	var outputMap map[string]interface{}
	json.Unmarshal(stdOut.Bytes(), &outputMap)

	internalContracts := outputMap["contracts"].(map[string]interface{})

	for key, value := range internalContracts {

		result := ContractJson{}

		result.Filename = strings.Split(key, ":")[1]
		result.Bytecode = "0x" + value.(map[string]interface{})["bin"].(string)

		contractAddress := ethClient.DeployContracts(result.Bytecode, publicKeys, private)

		result.ContractAddress = contractAddress

		result.Interface = value.(map[string]interface{})["abi"].(string)

		var abi interface{}

		json.Unmarshal([]byte(value.(map[string]interface{})["abi"].(string)), &abi)

		jsonData := JsonFormatted{abi, result.Bytecode, contractAddress}
		jsonString, _ := json.Marshal(jsonData)

		result.Json = string(jsonString)

		saveContract(result, strings.Replace(fileName, ".sol", "", -1))

		contNameMap[result.ContractAddress] = result.Filename
		contTimeMap[result.ContractAddress] = strconv.Itoa(int(time.Now().Unix()))
		abiMap[result.ContractAddress] = result.Interface

		results = append(results, result)
	}

	return results
}

func saveContract(contractJson ContractJson, fileName string) {
	path := env.GetAppConfig().ContractsDir + "/" + contractJson.ContractAddress + "_" + fileName

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	}

	filePath := path + "/" + contractJson.Filename + ".json"
	jsByte := []byte(contractJson.Json)
	err := ioutil.WriteFile(filePath, jsByte, 0775)
	if err != nil {
		fmt.Println(err)
	}
	abi := []byte(contractJson.Interface)
	err = ioutil.WriteFile(path+"/ABI", abi, 0775)
	if err != nil {
		fmt.Println(err)
	}
	bin := []byte(contractJson.Bytecode)
	err = ioutil.WriteFile(path+"/BIN", bin, 0775)
	if err != nil {
		fmt.Println(err)
	}
}
