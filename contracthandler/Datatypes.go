package contracthandler

import (
	"bytes"
	"encoding/hex"
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"regexp"
	"strings"
)

var mdt map[*regexp.Regexp]DataType

func init() {

	mdt = make(map[*regexp.Regexp]DataType)

	mdt[regexp.MustCompile(`uint256`)] = Uint256{}
	mdt[regexp.MustCompile(`bool`)] = Bool{}
	mdt[regexp.MustCompile(`int256\[.+\]`)] = Int256_FA{}
	mdt[regexp.MustCompile(`bytes`)] = Bytes{}
	mdt[regexp.MustCompile(`uint32\[\]`)] = Uint32_DA{}
	mdt[regexp.MustCompile(`string`)] = String{}
}

func ParseParameters(fp FunctionProcessor) []DataType {

	var dt []DataType
	for i, param := range strings.Split(fp.Signature, ",") {
		for k, v := range mdt {

			if k.MatchString(param) {
				dt = append(dt, v.New(fp.Param[i], ""))
			}
		}
	}

	return dt
}

func ParseResults(fp FunctionProcessor) []DataType {

	var dt []DataType
	for _, param := range strings.Split(fp.Signature, ",") {
		for k, v := range mdt {

			if k.MatchString(param) {
				dt = append(dt, v.New(nil, param))
			}

		}

	}

	return dt
}

type Uint256 struct {
	BaseDataType
}

func (t Uint256) Decode(data []string, index int) (int, interface{}) {

	return 1, util.StringToInt(data[index])
}

func (t Uint256) New(i interface{}, sig string) DataType {

	return Uint256{BaseDataType{i, sig}}
}

func (t Uint256) IsDynamic() bool {
	return false
}

func (t Uint256) Length() int {
	return 1
}

func (t Uint256) Encode() []string {

	i := t.value.(int)

	return []string{util.IntToString(i)}
}

type Bool struct {
	Uint256
}

func (t Bool) New(i interface{}, sig string) DataType {

	if i != nil && i.(bool) {
		return Bool{Uint256{BaseDataType{1, sig}}}
	}
	return Bool{Uint256{BaseDataType{0, sig}}}
}

func (t Bool) Decode(data []string, index int) (int, interface{}) {

	_, t.value = t.Uint256.Decode(data, index)

	return 1, t.value == 1
}

type Uint32_DA struct {
	BaseDataType
}

func (t Uint32_DA) Decode(data []string, index int) (int, interface{}) {

	offset := util.StringToInt(data[index])

	length := util.StringToInt(data[offset/32])

	var a = make([]int, length)

	for i, j := offset/32+1, 0; i < offset/32+1+length; i++ {
		a[j] = util.StringToInt(data[i])

		j++
	}

	return 1, a
}

func (t Uint32_DA) New(i interface{}, sig string) DataType {

	return Uint32_DA{BaseDataType{i, sig}}
}

func (t Uint32_DA) IsDynamic() bool {
	return true
}

func (t Uint32_DA) Length() int {
	i := t.value.([]int)

	return len(i) + 1
}

func (t Uint32_DA) Encode() []string {

	i := t.value.([]int)

	r := make([]string, len(i)+1)

	r[0] = util.IntToString(len(i))

	for j := 1; j <= len(i); j++ {
		r[j] = util.IntToString(i[j-1])
	}

	return r
}

type Int256_FA struct {
	Uint32_DA
}

func (t Int256_FA) New(i interface{}, sig string) DataType {

	return Int256_FA{Uint32_DA{BaseDataType{i, sig}}}
}

func (t Int256_FA) IsDynamic() bool {
	return false
}

func (t Int256_FA) Length() int {
	i := t.value.([]int)

	return len(i)
}

func (t Int256_FA) Encode() []string {

	i := t.value.([]int)

	var output []string
	for _, v := range i {
		output = append(output, util.IntToString(v))
	}

	return output
}

func (t Int256_FA) Decode(data []string, index int) (int, interface{}) {

	length := util.StringToInt(util.Between(t.GetSignature(), "[", "]"))

	var a = make([]int, length)

	for i, j := index, 0; j < length; i++ {
		a[j] = util.StringToInt(data[i])

		j++
	}

	return length, a
}

type Bytes struct {
	BaseDataType
}

func (t Bytes) Decode(data []string, index int) (int, interface{}) {
	offset := util.StringToInt(data[index])

	length := util.StringToInt(data[offset/32])

	var buffer bytes.Buffer

	for i, c := offset/32+1, 0; c < length; i++ {

		buffer.WriteString(data[i])

		c += 32

	}

	t.value, _ = hex.DecodeString(buffer.String()[:length*2])

	return 1, t.value

}

func (t Bytes) New(i interface{}, sig string) DataType {

	return Bytes{BaseDataType{i, sig}}
}

func (t Bytes) IsDynamic() bool {
	return true
}

func (t Bytes) Length() int {
	i := t.value.([]byte)

	if len(i)%32 == 0 {
		return len(i)/32 + 1
	}
	return len(i)/32 + 2

}

func (t Bytes) Encode() []string {

	s := t.value.([]byte)

	var d []string

	d = append(d, util.IntToString(len(s)))

	limit := 0

	if len(s)%32 == 0 {
		limit = len(s) / 32
	} else {
		limit = len(s)/32 + 1
	}

	for i := 0; i < limit; i++ {

		j := i * 32

		k := j + 32

		var b []byte

		if k > len(s) {
			b = make([]byte, 32)

			copy(b, s[j:])
		} else {
			b = s[j:k]
		}
		d = append(d, hex.EncodeToString(b))
	}

	return d
}

type String struct {
	Bytes
}

func (t String) New(i interface{}, sig string) DataType {

	if i == nil {
		i = ""
	}

	return String{Bytes{BaseDataType{[]byte(i.(string)), sig}}}
}

func (t String) Decode(data []string, index int) (int, interface{}) {

	_, t.value = t.Bytes.Decode(data, index)

	return 1, string(t.value.([]byte))
}
