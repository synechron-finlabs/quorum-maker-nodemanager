package contracthandler

import (
	"bytes"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"strings"		
)


type DataType interface {
	IsDynamic() bool
	Length() int
	Encode() []string
	New(interface{}, string) DataType
	Decode([]string, int) (int, interface{})
	SetSignature(string)
	GetSignature() string
}

type BaseDataType struct {
	value    interface{}
	Signatue string
}

func (bdt BaseDataType) SetSignature(s string) {
	bdt.Signatue = s
}

func (bdt BaseDataType) GetSignature() string {
	return bdt.Signatue
}

type FunctionProcessor struct {
	Signature string
	Param     []interface{}
	Result    string
}

func (fp FunctionProcessor) GetData() string {

	datatypes := ParseParameters(fp)

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

func (fp FunctionProcessor) GetResults() []interface{} {

	length := len(fp.Result) / 64

	data := make([]string, length)

	for i, j := 0, 0; i < length; i++ {
		data[i] = fp.Result[j : j+64]

		j += 64

	}

	datatypes := ParseResults(fp)

	results := make([]interface{}, len(datatypes))

	nextIndex := 0

	for i := range strings.Split(fp.Signature, ",") {
		ni, result := datatypes[i].Decode(data, nextIndex)

		nextIndex += ni

		results[i] = result
	}

	return results
}
