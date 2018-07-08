package contractclient

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
	"strings"
)

type ParamTableRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var abiMap = map[string]string{}
var funcSigMap = map[string]string{}
var funcParamNameMap = map[string]string{}
var funcNameMap = map[string]string{}

func ABIParser(contractAdd string, abiContent string, payload string) ([]ParamTableRow, string) {
	var abi abi.ABI
	if abiMap[contractAdd] == "" {
		abiMap[contractAdd] = abiContent
		abi.UnmarshalJSON([]byte(abiMap[contractAdd]))
		methodMap := abi.Methods

		var functionSigs []string
		var paramNames []string
		var keccakHashes []string
		i := 1
		for key := range methodMap {
			if !methodMap[key].Const {
				var funcSig string
				var params string
				for _, elem := range methodMap[key].Inputs {
					funcSig = funcSig + elem.Type.String() + ","
					params = params + elem.Name + ","
				}
				paramNames = append(paramNames, strings.TrimSuffix(params, ","))
				functionSigs = append(functionSigs, key+"("+strings.TrimSuffix(funcSig, ",")+")")
				keccakHashes = append(keccakHashes, hex.EncodeToString(crypto.Keccak256([]byte(functionSigs[i-1]))[:4]))
				funcParamNameMap[contractAdd + ":" + keccakHashes[i-1]] = paramNames[i-1]
				functionSigs[i-1] = strings.TrimSuffix(funcSig, ",")
				funcSigMap[contractAdd + ":" + keccakHashes[i-1]] = functionSigs[i-1]

				isSupported := contracthandler.IsSupported(functionSigs[i-1])

				if !isSupported {
					funcSigMap[contractAdd + ":" + keccakHashes[i-1]] = "unsupported"
				}

				funcNameMap[contractAdd + ":" + keccakHashes[i-1]] = key + "(" + strings.TrimSuffix(funcSig, ",") + ")"
				i++
			}
		}
	}
	return Decode(payload, contractAdd)
}

func Decode(r string, contractAdd string) ([]ParamTableRow, string) {
	keccakHash := r[2:10]
	if funcSigMap[contractAdd + ":" + keccakHash] == "" {
		abiMismatch := make([]ParamTableRow, 1)
		abiMismatch[0].Key = "decodeFailed"
		abiMismatch[0].Value = "ABI Mismatch"
		abiMap[contractAdd] = ""
		return abiMismatch, ""
	}
	if funcSigMap[contractAdd + ":" + keccakHash] == "unsupported" {
		abiMismatch := make([]ParamTableRow, 1)
		abiMismatch[0].Key = "decodeFailed"
		abiMismatch[0].Value = "Unsupported Datatype"
		return abiMismatch, ""
	}
	encodedParams := r[10:]
	params := strings.Split(funcSigMap[contractAdd + ":" + keccakHash], ",")
	paramTable := make([]ParamTableRow, len(params))
	if r == "" || len(r) < 1 {
		return paramTable, ""
	}
	paramNamesArr := strings.Split(funcParamNameMap[contractAdd + ":" + keccakHash], ",")
	resultArray := contracthandler.FunctionProcessor{funcSigMap[contractAdd + ":" + keccakHash]}.Decode(encodedParams)
	for i := 0; i < len(params); i++ {
		paramTable[i].Key = paramNamesArr[i]
		paramTable[i].Value = fmt.Sprint(resultArray[i])
	}

	return paramTable, funcNameMap[contractAdd + ":" + keccakHash]
}
