package HRBAlgorithm

import (
	"fmt"
	"testing"
)

//func TestMsgHandler(t *testing.T ) {
//	AlgorithmSetUp()
//
//	val := "abc"
//	targetHash := ConvertBytesToString(Hash([]byte(val)))
//
//	//Send the first Message
//	m := MSGStrcut{Header:MSG, Id:"1", Data:val, Round:0}
//	boo, count, hashStr:= Msghandler(m)
//	if boo != true || count != 1 || targetHash != hashStr{
//		t.Errorf("Wrong awnser")
//	}
//
//	//Send the Repeated Message
//	boo, count, hashStr = Msghandler(m)
//	if boo == true {
//		t.Errorf("Wrong awnser in MsgTest")
//	}
//
//	//send a different message
//	m = MSGStrcut{Header:MSG, Id:"1", Data:val, Round:1}
//	boo, count, hashStr = Msghandler(m)
//	if boo != true || count != 1 || targetHash != hashStr{
//		t.Errorf("Wrong awnser in MsgTest")
//	}
//}

//func TestEchoHandler(t *testing.T) {
//	AlgorithmSetUp()
//
//	val := "abc"
//	targetHash := ConvertBytesToString(Hash([]byte(val)))
//
//	//Send the first Message
//	m := ECHOStruct{Header:ECHO, Id:"1", HashData:targetHash, Round:0, SenderId:"0"}
//	boo, count := EchoHandler(m)
//	if boo != true || count != 1 {
//		t.Errorf("Wrong awnser in MsgTest1")
//	}
// 	//Send a repeated data
//	m = ECHOStruct{Header:ECHO, Id:"1", HashData:targetHash, Round:0, SenderId:"0"}
//	boo, count = EchoHandler(m)
//	if boo == true {
//		t.Errorf("Wrong awnser in MsgTest2")
//	}
//
//	//Send a different Echo from another id
//	m = ECHOStruct{Header:ECHO, Id:"1", HashData:targetHash, Round:0, SenderId:"1"}
//	boo, count = EchoHandler(m)
//	if boo != true || count != 2 {
//		t.Errorf("Wrong awnser in MsgTest3")
//	}
//}

//func TestAccHandler(t *testing.T) {
//	AlgorithmSetUp()
//
//	val := "abc"
//	targetHash := ConvertBytesToString(Hash([]byte(val)))
//
//	//Send the first Message
//	m := ACCStruct{Header:ACC, Id:"1", HashData:targetHash,Round:0, SenderId:"0"}
//	notSeen, count, req := AccHandler(m)
//	if notSeen != true || count != 1 || req != false {
//		t.Errorf("Wrong awnser in AccTest1")
//	}
//
//	notSeen, count, req = AccHandler(m)
//	if notSeen == true {
//		t.Errorf("Wrong awnser in AccTest2 Unique")
//	}
//
//	/*
//	Test Threshold f + 1
//	 */
//	m = ACCStruct{Header:ACC, Id:"1", HashData:targetHash,Round:0, SenderId:"1"}
//	notSeen, count, req = AccHandler(m)
//	if notSeen != true || count != 2 || req != true {
//		fmt.Println(notSeen, count, req)
//		t.Errorf("Wrong awnser in AccTest3\n")
//	}
//}

//func TestReqHandler(t *testing.T) {
//	AlgorithmSetUp()
//
//	val := "abc"
//	targetHash := ConvertBytesToString(Hash([]byte(val)))
//	fmt.Println(targetHash)
//	//Send the first Message
//	msg := MSGStrcut{Header:MSG, Id:"1", Data:val, Round:0, SenderId:"1"}
//	Msghandler(msg)
//
//	//Send the first Message
//	m := REQStruct{Header: REQ, Id: "1", HashData: targetHash, Round: 0, SenderId: "0"}
//	f1, f2 := ReqHandler(m)
//	if f1 != true || f2 != true {
//		t.Errorf("Wrong awnser in ReqTest1\n")
//	}
//
//	f1, f2 = ReqHandler(m)
//	if f1 != false || f2 != false {
//		t.Errorf("Wrong awnser in ReqTest2\n")
//	}
//
//	/*
//	Check the case when we donot have the data
//	 */
//
//	val1 := "abcd"
//	targetHash = ConvertBytesToString(Hash([]byte(val1)))
//	//Send the first Message
//	m = REQStruct{Header: REQ, Id: "1", HashData: targetHash, Round: 0, SenderId: "0"}
//	f1, f2 = ReqHandler(m)
//	if f1 != false || f2 != false {
//		t.Errorf("Wrong awnser in ReqTest1\n")
//	}
//
//	f1, f2 = ReqHandler(m)
//	if f1 != false || f2 != false {
//		t.Errorf("Wrong awnser in ReqTest2\n")
//	}
//
//}

//func TestFwdHandler(t *testing.T) {
//	AlgorithmSetUp()
//
//	val := "abc"
//
//	//Suppose have not sent a request before
//	msg := FWDStruct{Header:FWD, Id:"1", Data:val, Round:0, SenderId:"1"}
//	f1, f2 := FwdHandler(msg)
//	if f1 != false || f2 != false {
//		t.Errorf("Wrong awnser in FWDTest1\n")
//	}
//
//	//Sent a request before
//	targetHash := ConvertBytesToString(Hash([]byte(val)))
//	req := REQStruct{Header:REQ, Id:"1", HashData:targetHash, Round:0}
//	l:= []string{"1"}
//	ReqSentSet[req] = l
//	f1, f2 = FwdHandler(msg)
//	if f1 != true || f2 != true {
//		t.Errorf("Wrong awnser in FWDTest2\n")
//	}
//
//	//Sent a request before but have FWD
//	f1, f2 = FwdHandler(msg)
//	if f1 != true || f2 != false {
//		t.Errorf("Wrong awnser in FWDTest3\n")
//	}
//}

/*
Simple Testing without Byzantinne
 */
func TestSimple1(t *testing.T) {
	AlgorithmSetUp()

	/*
	Broadcast phase
	 */
	val := "abc"
	m := MSGStrcut{Header:MSG, Id:"0", Data:val, Round:0, SenderId:"0"}
	hash := ConvertBytesToString(Hash([]byte(val)))
	/*
	Accept Phase
	 */

	//MSG receive:
	Msghandler(m)
	fmt.Println("****************************")
	//Echo Receive:
	echoM := ECHOStruct{Header:ECHO, Id:"0", HashData:hash, Round:0, SenderId:"1"}
	f, count, flags := EchoHandler(echoM)
	if f != true || count != 2 || flags[0] != false {
		t.Errorf("Wrong awnser in Simple1\n")
	}

	fmt.Println("****************************")
	echoM = ECHOStruct{Header:ECHO, Id:"0", HashData:hash, Round:0, SenderId:"2"}
	f, count, flags = EchoHandler(echoM)
	if f != true || count != 3 || flags[0] != false || flags[1] != true {
		t.Errorf("Wrong awnser in Simple1\n")
	}

	/*
	Receive ACC
	 */
	accM := ACCStruct{Header:ACC, Id:"0", HashData:hash, Round:0, SenderId:"0"}
	AccHandler(accM)
	accM = ACCStruct{Header:ACC, Id:"0", HashData:hash, Round:0, SenderId:"1"}
	AccHandler(accM)
	accM = ACCStruct{Header:ACC, Id:"0", HashData:hash, Round:0, SenderId:"2"}
	AccHandler(accM)

	flags = check(accM)
	if flags[0] != false || flags[1] != false || flags[2] != false || flags[3] != true {
		t.Errorf("Wrong awnser in Simple1\n")
	}

}

//func TestSimple2(t *testing.T) {
//	AlgorithmSetUp()
//
//	//Broadcast Phase:
//	val := "abc"
//	m := MSGStrcut{Header:MSG, Id:"0", Data:val, Round:0, SenderId:"0"}
//
//	//Echo Receive
//}
//
//func TestSimple3(t *testing.T) {
//	AlgorithmSetUp()
//	//Broadcast Phase:
//	val := "abc"
//	m := MSGStrcut{Header:MSG, Id:"0", Data:val, Round:0, SenderId:"0"}
//
//	//ACC Phase
//}

/*
Simple Testing with Byzantine
 */
