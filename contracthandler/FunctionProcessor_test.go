package contracthandler

import (
	"github.com/synechron-finlabs/quorum-maker-nodemanager/util"
	"reflect"
	"testing"
	"time"
)

func TestGetData(t *testing.T) {

	defer util.TotalTime(time.Now().Nanosecond())

	for _, table := range getFunctionProcessorTestData() {
		data := FunctionProcessor{table.x}.Encode(table.y)

		if data != table.z {
			t.Errorf("Encoding was incorrect, got: %d, want: %d.", data, table.z)
		}
	}
}

func TestGetResults(t *testing.T) {

	for _, table := range getFunctionProcessorTestData() {
		data := FunctionProcessor{table.x}.Decode(table.z)
		if !checkEquality(table.y, data) {
			t.Errorf("Decoding was incorrect")
		}

	}
}

func getFunctionProcessorTestData() []test_data {

	tables := []test_data{
		{"uint64[]", []interface{}{[]int{8, 9, 10, 11, 12}}, "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000000c"},
		{"uint256,int256[1],bytes,bool,int256[5],uint32[],int256[3],string", []interface{}{"999999", []string{"1"}, []byte("Perhaps we need a quick connect meeting on the current state for TD - paperwork, timeline for POC delivery, engagement, resourcing etc."), true, []string{"3", "4", "5", "6", "7"}, []int{8, 9, 10, 11, 12}, []string{"13", "14", "15"}, "Perhaps we need a quick connect meeting on the current state for TD - paperwork, timeline for POC delivery, engagement, resourcing etc. So blockchain gives us the possibility to confidentially share access to a data structure. This will first and foremost drives efficiencies in removing reconciliations and middlemen, but longer term will lead to new business models, often peer-to-peer in nature driven by the reduction in friction - either speed, cost, or accuracy - from what we can do today"}, "00000000000000000000000000000000000000000000000000000000000f423f000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000070000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000000d000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000f0000000000000000000000000000000000000000000000000000000000000340000000000000000000000000000000000000000000000000000000000000008750657268617073207765206e656564206120717569636b20636f6e6e656374206d656574696e67206f6e207468652063757272656e7420737461746520666f72205444202d207061706572776f726b2c2074696d656c696e6520666f7220504f432064656c69766572792c20656e676167656d656e742c207265736f757263696e67206574632e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000b000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000001ef50657268617073207765206e656564206120717569636b20636f6e6e656374206d656574696e67206f6e207468652063757272656e7420737461746520666f72205444202d207061706572776f726b2c2074696d656c696e6520666f7220504f432064656c69766572792c20656e676167656d656e742c207265736f757263696e67206574632e20536f20626c6f636b636861696e2067697665732075732074686520706f73736962696c69747920746f20636f6e666964656e7469616c6c792073686172652061636365737320746f20612064617461207374727563747572652e20546869732077696c6c20666972737420616e6420666f72656d6f73742064726976657320656666696369656e6369657320696e2072656d6f76696e67207265636f6e63696c696174696f6e7320616e64206d6964646c656d656e2c20627574206c6f6e676572207465726d2077696c6c206c65616420746f206e657720627573696e657373206d6f64656c732c206f6674656e20706565722d746f2d7065657220696e206e61747572652064726976656e2062792074686520726564756374696f6e20696e206672696374696f6e202d206569746865722073706565642c20636f73742c206f72206163637572616379202d2066726f6d20776861742077652063616e20646f20746f6461790000000000000000000000000000000000"},
		{"uint256,int256[1],bytes", []interface{}{"999999", []string{"1"}, []byte("Perhaps we need a quick connect meeting on the current state for TD - paperwork, timeline for POC delivery, engagement, resourcing etc.")}, "00000000000000000000000000000000000000000000000000000000000f423f00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000008750657268617073207765206e656564206120717569636b20636f6e6e656374206d656574696e67206f6e207468652063757272656e7420737461746520666f72205444202d207061706572776f726b2c2074696d656c696e6520666f7220504f432064656c69766572792c20656e676167656d656e742c207265736f757263696e67206574632e00000000000000000000000000000000000000000000000000"},
		{"string,string,string", []interface{}{"10d93926bcd78f37dbbb2c95d65e7bc4c723a75e66fe60317525f0607c7111d3d9088c8a33c944c3e2e24b3281115a8f688312cee61573119e47555e8cd31e30", "JPM", "Custodian"}, "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000080313064393339323662636437386633376462626232633935643635653762633463373233613735653636666536303331373532356630363037633731313164336439303838633861333363393434633365326532346233323831313135613866363838333132636565363135373331313965343735353565386364333165333000000000000000000000000000000000000000000000000000000000000000034a504d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009437573746f6469616e0000000000000000000000000000000000000000000000"},
		{"bytes32[],uint256,uint256", []interface{}{[][]byte{[]byte("abc"), []byte("cde")}, "1370", "10000"}, "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000055a0000000000000000000000000000000000000000000000000000000000002710000000000000000000000000000000000000000000000000000000000000000261626300000000000000000000000000000000000000000000000000000000006364650000000000000000000000000000000000000000000000000000000000"},
		{"bytes32[2],uint256,uint256", []interface{}{[][]byte{[]byte("abc"), []byte("cde")}, "1370", "10000"}, "61626300000000000000000000000000000000000000000000000000000000006364650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000055a0000000000000000000000000000000000000000000000000000000000002710"},
		{"bytes32[],uint256,bytes32[3],uint256", []interface{}{[][]byte{[]byte("abc"), []byte("cde")}, "1370", [][]byte{[]byte("fgh"), []byte("ijk"), []byte("lmn")}, "10000"}, "00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000055a6667680000000000000000000000000000000000000000000000000000000000696a6b00000000000000000000000000000000000000000000000000000000006c6d6e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002710000000000000000000000000000000000000000000000000000000000000000261626300000000000000000000000000000000000000000000000000000000006364650000000000000000000000000000000000000000000000000000000000"},
		{"address,bytes32[],uint256,bytes32,uint256", []interface{}{"0xdad324753d1d84ccaad81180e3f6866637cda99b", [][]byte{[]byte("abc"), []byte("cde")}, "1370", []byte("fgh"), "20000"}, "000000000000000000000000dad324753d1d84ccaad81180e3f6866637cda99b00000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000055a66676800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004e20000000000000000000000000000000000000000000000000000000000000000261626300000000000000000000000000000000000000000000000000000000006364650000000000000000000000000000000000000000000000000000000000"},
		{"", []interface{}{}, ""},
	}

	return tables
}

type test_data struct {
	x string
	y []interface{}
	z string
}

func checkEquality(a []interface{}, b []interface{}) bool {

	if len(a) != len(b) {
		return false
	}

	for i, av := range a {

		if !reflect.DeepEqual(av, b[i]) {
			return false
		}

	}

	return true
}
