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
	"math/big"
	"bufio"
	"io"
	"io/ioutil"
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

func HexStringtoLargeInt64(hexVal string) string {
	hexVal = strings.TrimSuffix(hexVal, "\n")
	hexVal = strings.TrimPrefix(hexVal, "0x")
	intVal := new(big.Int)
	intVal.SetString(hexVal, 16)
	return intVal.String()
}

func MustGetString(key string, filename *properties.Properties) (val string) {
	val = filename.MustGetString(key)
	val = strings.TrimSuffix(val, "\n")
	return val
}

func ComposeJSON(intrfc string, bytecode string, contractAdd string) (json string) {
	json = "{\n\"interface\": " + intrfc + ",\n\"bytecode\": \"" + bytecode + "\",\n\"address\": \"" + contractAdd + "\"\n}"
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

func DecodeLargeInt(s string) string {
	i := new(big.Int)
	i.SetString(s, 16)

	return i.String()
}

func EncodeLargeInt(i string) string {

	j := new(big.Int)
	j.SetString(i, 0)

	return fmt.Sprintf("%064s", j.Text(16))
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
	out, _ := exec.Command("bash", "-c", command).Output()
	return string(out)
}

func DeleteProperty(property string, filepath string) {
	command := fmt.Sprint("sed -i '0,/", property, "/ s///' ", filepath)
	cmd := exec.Command("bash", "-c", command)
	cmd.Run()
}

type Uint256 struct {
	part1 int64
	part2 int64
	part3 int64
	part4 int64
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

func CreateFile(path string) {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}
}

func WriteFile(path string, content string) {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if isError(err) {
		return
	}
	err = file.Sync()
	if isError(err) {
		return
	}
}

func DeleteFile(path string) {
	var err = os.Remove(path)
	if isError(err) {
		return
	}
}

func isError(err error) bool {
	if err != nil {
	}
	return (err != nil)
}
