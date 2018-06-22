package contractclient

import (
	"strings"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"encoding/hex"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/contracthandler"
	"fmt"
	"regexp"
)

type ParamTableRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var SupportedDatatypes = []*regexp.Regexp{regexp.MustCompile(`uint256`), regexp.MustCompile(`bool`), regexp.MustCompile(`int256\[.+\]`), regexp.MustCompile(`bytes`), regexp.MustCompile(`uint32\[\]`), regexp.MustCompile(`string`)}

var abiMap = map[string]string{}
var funcSigMap = map[string]string{}
var funcParamNameMap = map[string]string{}

func ABIParser(contractAdd string, abiContent string, payload string) []ParamTableRow {
	var abiVal abi.ABI
	if abiMap[contractAdd] == "" {
		abiMap[contractAdd] = abiContent
		methodMap := abiVal.UnmarshalJSON([]byte(abiMap[contractAdd]))
		size := 0
		for key := range methodMap {
			if methodMap[key].Const == false {
				size++
			}
		}
		functionSigs := make([]string, size)
		paramNames := make([]string, size)
		keccakHashes := make([]string, size)
		i := 1

		for key := range methodMap {
			if methodMap[key].Const == false {
				var funcSig string
				var params string
				for _, elem := range methodMap[key].Inputs {
					i := 0
					for _, v := range SupportedDatatypes {
						if v.MatchString(elem.Type.String()) != true {
							i++
						}
					}
					if (i == len(SupportedDatatypes)) {
						abiMap[contractAdd] = "Unsupported"
						decodeUnsupported := make([]ParamTableRow, 1)
						decodeUnsupported[0].Key = "Unsupported"
						decodeUnsupported[0].Value = "Some input parameter datatypes are unsupported"
						return decodeUnsupported
					}
					funcSig = funcSig + elem.Type.String() + ","
					params = params + elem.Name + ","
				}
				paramNames[i-1] = strings.TrimSuffix(params, ",")
				functionSigs[i-1] = key + "(" + strings.TrimSuffix(funcSig, ",") + ")"
				keccakHashes[i-1] = hex.EncodeToString(crypto.Keccak256([]byte(functionSigs[i-1]))[:4])
				funcParamNameMap[contractAdd+":"+keccakHashes[i-1]] = paramNames[i-1]
				functionSigs[i-1] = strings.TrimSuffix(funcSig, ",")
				funcSigMap[contractAdd+":"+keccakHashes[i-1]] = functionSigs[i-1]
				i++
			}
		}

	} else if abiMap[contractAdd] == "Unsupported" {
		decodeUnsupported := make([]ParamTableRow, 1)
		decodeUnsupported[0].Key = "Unsupported"
		decodeUnsupported[0].Value = "Some input parameter datatypes are unsupported"
		return decodeUnsupported
	}

	return Decode(payload, contractAdd)
}

func Decode(r string, contractAdd string) []ParamTableRow {
	keccakHash := r[2:10]
	encodedParams := r[10:]
	params := strings.Split(funcSigMap[contractAdd+":"+keccakHash], ",")
	paramTable := make([]ParamTableRow, len(params))
	if r == "" || len(r) < 1 {
		return paramTable
	}
	paramNamesArr := strings.Split(funcParamNameMap[contractAdd+":"+keccakHash], ",")
	resultArray := contracthandler.FunctionProcessor{funcSigMap[contractAdd+":"+keccakHash], nil, encodedParams}.GetResults()
	for i := 0; i < len(params); i++ {
		paramTable[i].Key = paramNamesArr[i]
		paramTable[i].Value = fmt.Sprint(resultArray[i])
	}
	return paramTable
}
