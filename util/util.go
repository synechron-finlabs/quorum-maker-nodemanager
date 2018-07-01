package util

import (
	"encoding/hex"
	"fmt"
	"github.com/magiconair/properties"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"os/exec"
)

func TakeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := takeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func HexStringtoInt64(hexVal string) (intVal int64) {
	hexVal = strings.TrimSuffix(hexVal, "\n")
	hexVal = strings.TrimPrefix(hexVal, "0x")
	intVal, err := strconv.ParseInt(hexVal, 16, 64)
	if err != nil {
	}
	return intVal
}

func MustGetString(key string, filename *properties.Properties) (val string) {
	val = filename.MustGetString(key)
	val = strings.TrimSuffix(val, "\n")
	return val
}

func ComposeJSON(intrfc string, bytecode string, contractAdd string) (json string) {
	json = "{\n\"interface\" : " + intrfc + ",\n\"bytecode\" : \"" + bytecode + "\",\n\"address\" : \"" + contractAdd + "\"\n}"
	return json
}

func IntToString(i int) string {

	return fmt.Sprintf("%064s", fmt.Sprintf("%x", i))
}
func StringToInt(s string) int {

	s = strings.TrimLeft(s, "0")

	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		//panic(err)
		return 0
	}

	return int(n)
}

func ByteToString(a []byte) string {

	b := make([]byte, 32)

	copy(b, a)
	return hex.EncodeToString(b)

}

func Between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func TotalTime(start int) {
	fmt.Println("Total time took = ", time.Now().Nanosecond()-start)
}

func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

func PropertyExists(property string, filepath string) string {
	command := fmt.Sprint("grep -R ", "\"", property, "\" ", "\"", filepath, "\"")
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}

func DeleteProperty(property string, filepath string) {
	command := fmt.Sprint("sed -i '0,/", property, "/ s///' ", filepath)
	cmd := exec.Command("bash", "-c", command)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

type Uint256 struct {
	part1 int64
	part2 int64
	part3 int64
	part4 int64
}
