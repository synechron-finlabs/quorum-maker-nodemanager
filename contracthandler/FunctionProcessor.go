package contracthandler

import (
	"bytes"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"strings"		
)



type FunctionProcessor struct {
	Signature string
}

func (fp FunctionProcessor) Encode(paramValues []interface{}) string {

	datatypes := fp.getDataTypes(paramValues)

	var output []string

	for _, dt0 := range datatypes {

		if dt0.IsDynamic() {
			output = append(output, util.IntToString(0))
		} else {
			output = append(output, dt0.Encode()...)
		}

	}

	index, loc := 0, 0

	for i, dt0 := range datatypes {

		if dt0.IsDynamic() {

			for j, dt1 := range datatypes {

				if !dt1.IsDynamic() || j < i {
					index += dt1.Length()

					if dt1.IsDynamic() {
						index++
					}
				} else {
					index++
				}
			}

			output[loc] = util.IntToString(index * 32)
			output = append(output, dt0.Encode()...)

			index = 0

			loc++
		} else {
			loc += dt0.Length()

		}

	}

	var buffer bytes.Buffer

	for _, v := range output {
		buffer.WriteString(v)
	}

	return buffer.String()

}

func (fp FunctionProcessor) Decode(encodedString string) []interface{} {

	length := len(encodedString) / 64

	data := make([]string, length)

	for i, j := 0, 0; i < length; i++ {
		data[i] = encodedString[j : j+64]

		j += 64

	}

	datatypes := fp.getDataTypes(nil)

	results := make([]interface{}, len(datatypes))

	if len(datatypes) == 0 {
		return results
	}

	nextIndex := 0

	for i := range strings.Split(fp.Signature, ",") {
		ni, result := datatypes[i].Decode(data, nextIndex)

		nextIndex += ni

		results[i] = result
	}

	return results
}


func (fp FunctionProcessor) getDataTypes(paramValues []interface{}) []DataType {

	var dt []DataType
	for i, paramType := range strings.Split(fp.Signature, ",") {
		for k, v := range mdt {

			if k.MatchString(paramType) {
				if paramValues == nil {
					dt = append(dt, v.New(nil, paramType))
				} else {
					dt = append(dt, v.New(paramValues[i], paramType))
				}

			}

		}

	}

	return dt
}