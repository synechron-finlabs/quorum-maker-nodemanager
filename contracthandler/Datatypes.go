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

	mdt[regexp.MustCompile(`^u?int(([1-9]|[1-5][0-9])|(6[0-4]))$`)] = Uint{}
	mdt[regexp.MustCompile(`^bool$`)] = Bool{}
	mdt[regexp.MustCompile(`^u?int(([1-9]|[1-5][0-9])|(6[0-4]))\[[0-9]+\]$`)] = UintFA{}
	mdt[regexp.MustCompile(`^bytes$`)] = Bytes{}
	mdt[regexp.MustCompile(`^u?int(([1-9]|[1-5][0-9])|(6[0-4]))\[\]$`)] = UintDA{}
	mdt[regexp.MustCompile(`^string$`)] = String{}
	mdt[regexp.MustCompile(`^bytes([1-9]|1[0-9]|2[0-9]|3[0-2])\[\]$`)] = Bytes32DA{}
	mdt[regexp.MustCompile(`^bytes([1-9]|1[0-9]|2[0-9]|3[0-2])\[[0-9]+\]$`)] = Bytes32FA{}
	mdt[regexp.MustCompile(`^bytes([1-9]|1[0-9]|2[0-9]|3[0-2])$`)] = BytesFixed{}
	mdt[regexp.MustCompile(`^u?int(6[5-9]|[7-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-6])$`)] = UintLarge{}
	mdt[regexp.MustCompile(`^u?int(6[5-9]|[7-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-6])\[[0-9]+\]$`)] = UintLargeFA{}
	mdt[regexp.MustCompile(`^u?int(6[5-9]|[7-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-6])\[\]$`)] = UintLargeDA{}
	mdt[regexp.MustCompile(`^address$`)] = Address{}
	mdt[regexp.MustCompile(`^address\[[0-9]+\]$`)] = AddressFA{}
	mdt[regexp.MustCompile(`^address\[\]$`)] = AddressDA{}

}

/*
* supports below formats
* functionName(datatype1, datatype2,...)
* datatype1, datatype2,...
* datatype1, datatype2,.., (with an extra comma and the end
*
*/
func IsSupported(sig string) bool {

	rex := regexp.MustCompile(`.*\((.*)\)|(.*)`)
	var datatypes string

	if match := rex.FindStringSubmatch(sig); match[1] != "" {
		datatypes = match[1]
	}else{
		datatypes = strings.TrimSuffix(match[2], ",")
	}

	if datatypes == "" {
		return true
	}

	for _, param := range strings.Split(datatypes, ",") {
		var found bool
		for k := range mdt {

			if k.MatchString(param) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}


type DataType interface {
	IsDynamic() bool
	Length() int
	Encode() []string
	New(interface{}, string) DataType
	Decode([]string, int) (int, interface{})
}

type BaseDataType struct {
	value     interface{}
	signature string
}


type Uint struct {
	BaseDataType
}


func (t Uint) New(i interface{}, sig string) DataType {

	return Uint{BaseDataType{i, sig}}
}

func (t Uint) IsDynamic() bool {
	return false
}

func (t Uint) Length() int {
	return 1
}


func (t Uint) Decode(data []string, index int) (int, interface{}) {

	return 1, util.StringToInt(data[index])
}

func (t Uint) Encode() []string {

	i := t.value.(int)

	return []string{util.IntToString(i)}
}



type Bool struct {
	Uint
}

func (t Bool) New(i interface{}, sig string) DataType {

	if i != nil && i.(bool) {
		return Bool{Uint{BaseDataType{1, sig}}}
	}
	return Bool{Uint{BaseDataType{0, sig}}}
}

func (t Bool) Decode(data []string, index int) (int, interface{}) {

	_, t.value = t.Uint.Decode(data, index)

	return 1, t.value == 1
}

type UintDA struct {
	BaseDataType
}

func (t UintDA) New(i interface{}, sig string) DataType {

	return UintDA{BaseDataType{i, sig}}
}

func (t UintDA) IsDynamic() bool {
	return true
}

func (t UintDA) Length() int {
	i := t.value.([]int)

	return len(i) + 1
}

func (t UintDA) Decode(data []string, index int) (int, interface{}) {

	offset := util.StringToInt(data[index])

	length := util.StringToInt(data[offset/32])

	var a = make([]int, length)

	for i, j := offset/32+1, 0; i < offset/32+1+length; i++ {
		a[j] = util.StringToInt(data[i])

		j++
	}

	return 1, a
}

func (t UintDA) Encode() []string {

	i := t.value.([]int)

	r := make([]string, len(i)+1)

	r[0] = util.IntToString(len(i))

	for j := 1; j <= len(i); j++ {
		r[j] = util.IntToString(i[j-1])
	}

	return r
}

type UintFA struct {
	UintDA
}

func (t UintFA) New(i interface{}, sig string) DataType {

	return UintFA{UintDA{BaseDataType{i, sig}}}
}

func (t UintFA) IsDynamic() bool {
	return false
}

func (t UintFA) Length() int {
	i := t.value.([]int)

	return len(i)
}

func (t UintFA) Encode() []string {

	i := t.value.([]int)

	var output []string
	for _, v := range i {
		output = append(output, util.IntToString(v))
	}

	return output
}

func (t UintFA) Decode(data []string, index int) (int, interface{}) {

	length := util.StringToInt(util.Between(t.signature, "[", "]"))

	var a = make([]int, length)

	for i, j := index, 0; j < length; i++ {
		a[j] = util.StringToInt(data[i])

		j++
	}

	return length, a
}

type UintLarge struct {
	BaseDataType
}

func (t UintLarge) New(i interface{}, sig string) DataType {

	return UintLarge{BaseDataType{i, sig}}
}

func (t UintLarge) IsDynamic() bool {
	return false
}

func (t UintLarge) Length() int {
	return 1
}

func (t UintLarge) Decode(data []string, index int) (int, interface{}) {
	return 1, util.DecodeLargeInt(data[index])
}

func (t UintLarge) Encode() []string {

	i := t.value.(string)

	return []string{util.EncodeLargeInt(i)}
}

type UintLargeDA struct {
	BaseDataType
}

func (t UintLargeDA) New(i interface{}, sig string) DataType {

	return UintLargeDA{BaseDataType{i, sig}}
}

func (t UintLargeDA) IsDynamic() bool {
	return true
}

func (t UintLargeDA) Length() int {
	i := t.value.([]string)

	return len(i) + 1
}

func (t UintLargeDA) Decode(data []string, index int) (int, interface{}) {

	offset := util.StringToInt(data[index])

	length := util.StringToInt(data[offset/32])

	var a = make([]string, length)

	for i, j := offset/32+1, 0; i < offset/32+1+length; i++ {
		a[j] = util.DecodeLargeInt(data[i])

		j++
	}

	return 1, a
}

func (t UintLargeDA) Encode() []string {

	i := t.value.([]string)

	r := make([]string, len(i)+1)

	r[0] = util.IntToString(len(i))

	for j := 1; j <= len(i); j++ {
		r[j] = util.EncodeLargeInt(i[j-1])
	}

	return r
}

type UintLargeFA struct {
	UintLargeDA
}

func (t UintLargeFA) New(i interface{}, sig string) DataType {

	return UintLargeFA{UintLargeDA{BaseDataType{i, sig}}}
}

func (t UintLargeFA) IsDynamic() bool {
	return false
}

func (t UintLargeFA) Length() int {
	i := t.value.([]string)

	return len(i)
}

func (t UintLargeFA) Encode() []string {

	i := t.value.([]string)

	var output []string
	for _, v := range i {
		output = append(output, util.EncodeLargeInt(v))
	}

	return output
}

func (t UintLargeFA) Decode(data []string, index int) (int, interface{}) {

	length := util.StringToInt(util.Between(t.signature, "[", "]"))

	var a = make([]string, length)

	for i, j := index, 0; j < length; i++ {
		a[j] = util.DecodeLargeInt(data[i])

		j++
	}

	return length, a
}

type Address struct {
	UintLarge
}

func (t Address) New(i interface{}, sig string) DataType {

	return Address{UintLarge{BaseDataType{i, sig}}}
}


func (t Address) Decode(data []string, index int) (int, interface{}) {

	return 1, "0x" + strings.TrimLeft(data[index], "0")
}

type AddressDA struct {
	UintLargeDA
}

func (t AddressDA) New(i interface{}, sig string) DataType {

	return AddressDA{UintLargeDA{BaseDataType{i, sig}}}
}

func (t AddressDA) Decode(data []string, index int) (int, interface{}) {

	offset := util.StringToInt(data[index])

	length := util.StringToInt(data[offset/32])

	var a = make([]string, length)

	for i, j := offset/32+1, 0; i < offset/32+1+length; i++ {
		a[j] = "0x" + strings.TrimLeft(data[i], "0")

		j++
	}

	return 1, a
}


type AddressFA struct {
	UintLargeFA
}

func (t AddressFA) New(i interface{}, sig string) DataType {

	return AddressFA{UintLargeFA{UintLargeDA{BaseDataType{i, sig}}}}
}


func (t AddressFA) Decode(data []string, index int) (int, interface{}) {

	length := util.StringToInt(util.Between(t.signature, "[", "]"))

	var a = make([]string, length)

	for i, j := index, 0; j < length; i++ {
		a[j] = "0x" + strings.TrimLeft(data[i], "0")

		j++
	}

	return length, a
}

type Bytes struct {
	BaseDataType
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

type Bytes32DA struct {
	BaseDataType
}

func (t Bytes32DA) New(i interface{}, sig string) DataType {

	return Bytes32DA{BaseDataType{i, sig}}
}

func (t Bytes32DA) IsDynamic() bool {
	return true
}

func (t Bytes32DA) Length() int {
	i := t.value.([]int)

	return len(i) + 1
}


func (t Bytes32DA) Decode(data []string, index int) (int, interface{}) {

	offset := util.StringToInt(data[index])

	length := util.StringToInt(data[offset/32])

	var a = make([][]byte, length)

	for i, j := offset/32+1, 0; j < length; i++ {

		a[j], _ = hex.DecodeString(strings.Replace(data[i], "00", "",-1))

		j++

	}

	t.value = a

	return 1, a
}

func (t Bytes32DA) Encode() []string {

	i := t.value.([][]byte)

	r := make([]string, len(i)+1)

	r[0] = util.IntToString(len(i))

	for j := 1; j <= len(i); j++ {
		b := make([]byte, 32)

		copy(b, i[j-1])
		r[j] = hex.EncodeToString(b)

	}

	return r
}


type Bytes32FA struct {
	BaseDataType
}

func (t Bytes32FA) New(i interface{}, sig string) DataType {

	return Bytes32FA{BaseDataType{i, sig}}
}

func (t Bytes32FA) IsDynamic() bool {
	return false
}

func (t Bytes32FA) Length() int {
	i := t.value.([][]byte)

	return len(i)
}

func (t Bytes32FA) Decode(data []string, index int) (int, interface{}) {

	length := util.StringToInt(util.Between(t.signature, "[", "]"))

	var a = make([][]byte, length)

	for i, j := index, 0; j < length; i++ {

		a[j], _ = hex.DecodeString(strings.Replace(data[i], "00", "",-1))

		j++

	}

	t.value = a

	return length, a
}

func (t Bytes32FA) Encode() []string {

	i := t.value.([][]byte)

	r := make([]string, len(i))


	for j := 0; j < len(i); j++ {
		b := make([]byte, 32)

		copy(b, i[j])
		r[j] = hex.EncodeToString(b)

	}

	return r
}

type BytesFixed struct {
	Uint
}

func (t BytesFixed) New(i interface{}, sig string) DataType {

	return BytesFixed{Uint{BaseDataType{i, sig}}}
}

func (t BytesFixed) Decode(data []string, index int) (int, interface{}) {

	t.value, _ = hex.DecodeString(strings.Replace(data[index], "00", "",-1))

	return 1, t.value
}

func (t BytesFixed) Encode() []string {
	i := t.value.([]byte)

	b := make([]byte, 32)

	copy(b, i)

	return []string{hex.EncodeToString(b)}

}