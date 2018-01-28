package client

import (
	"encoding/hex"
	"fmt"
)

type RequestHandler interface {
	Encode() string
}

type ResponseHandler interface {
	Decode(r string)
}

type ContractParam struct {
	From    string
	To      string
	Passwd  string
	Parties []string
}

func Field(b [] byte, i int) []byte {
	s, e := 0, 0

	s = 64*i + 2

	e = s + 64

	return trim(b[s:e])
}

func Decode(b []byte) string {

	dst := make([]byte, hex.DecodedLen(len(b)))

	hex.Decode(dst, b)

	return string(dst)
}

func trim(b []byte) []byte {
	i := 0

	for ; i < len(b); i += 2 {
		if b[i] == 48 && b[i+1] == 48 {
			break
		}
	}

	return b[:i]

}

func EncodeAndPad(f interface{}) string {

	var ret string

	switch v := f.(type) {
	case int:
		i := f.(int)

		b := []byte{byte(i)}
		ret = hex.EncodeToString(b)
		ret = fmt.Sprintf("%064s", ret)
	default:
		a := []byte(fmt.Sprintf("%v", v))
		b := make([]byte, 32)

		copy(b, a)
		ret = hex.EncodeToString(b)
	}


	return ret
}
