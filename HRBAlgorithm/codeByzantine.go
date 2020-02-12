package HRBAlgorithm

import (
	"strconv"
	"time"
)

var ByzCodeCounter map[string] int
var ByzCodeElement map[string] [][]byte

func InitByzCode() {
	ByzCodeCounter = make (map[string] int)
	ByzCodeElement = make (map[string] [][]byte)
	dataFromSrc = make(map[string] string)
}

func ECByzBroadCast(length, round int) {
	//need to make sure that coded element > f
	time.Sleep(3*time.Second)
	var shards[][] byte

	for r := 0; r < round; r++ {
		s := RandStringBytes(length)
		if faulty == 0 {
			shards = Encode(s, total, total)
		} else {
			shards = Encode(s, total - faulty, 2*(total) - (total - faulty))
		}
		//fmt.Println("Shards are ", shards)


		for i := 0; i < total; i++ {
			code1 := ConvertBytesToString(shards[i])
			code2 := ConvertBytesToString(shards[i + total])
			m := MSGStruct{Header:MSG, Id:MyID, SenderId:MyID, Data: code1, Round: r, HashData: code2}
			sendReq := PrepareSend{M: m, SendTo: serverList[i]}
			SendReqChan <- sendReq
		}
	}
}

func ByzRecMsg(m Message) {
	identifier := identifierCreate(m.GetId(), m.GetRound())
	count, exist := ByzCodeCounter[identifier]

	stats := Stats{}
	stats.Start = time.Now()
	statsRecord[identifier] = stats

	//fmt.Printf("Begin Stats: %+v\n",stats)

	if exist {
		ByzCodeCounter[identifier] = count + 2
	} else {
		ByzCodeCounter[identifier] = 2
		ByzCodeElement[identifier] = make([][]byte, 2 * total)
	}

	index1 := serverMap[MyID]
	index2 := index1 + total
	code1, _ := ConvertStringToBytes(m.GetData())
	code2, _ := ConvertStringToBytes(m.GetHashData())
	//fmt.Println("code1: ", code1, " code2: ", code2)
	ByzCodeElement[identifier][index1] = code1
	ByzCodeElement[identifier][index2] = code2

	code := ConvertBytesToString(code1)
	id := m.GetId();
	round := m.GetRound()
	//Send Echo
	for i := 0; i < total; i++ {
		message := ECHOStruct{Header:ECHO, Id:id, SenderId:MyID, Data: code, Round: round}
		sendReq := PrepareSend{M: message, SendTo: serverList[i]}
		SendReqChan <- sendReq
	}
}

func ByzRecEcho(m Message) {
	if MyID != m.GetSenderId() {
		identifier := identifierCreate(m.GetId(), m.GetRound())
		count, exist := ByzCodeCounter[identifier]

		if exist {
			ByzCodeCounter[identifier] = count + 1
		} else {
			ByzCodeCounter[identifier] = 1
			ByzCodeElement[identifier] = make([][]byte, 2 * total)
		}

		index := serverMap[m.GetSenderId()]
		code, _ := ConvertStringToBytes(m.GetData())
		ByzCodeElement[identifier][index] = code

		if ByzCodeCounter[identifier] == total + 1 {
			var vals []string
			//fmt.Println("ByzCodes:",ByzCodeElement[identifier])
			if faulty == 0 {
				vals = permutation(ByzCodeElement[identifier], total, total)
			} else {
				vals = permutation(ByzCodeElement[identifier], total - faulty, 2*(total) - (total - faulty))
			}
			detected := validateByzCode(vals)
			if !detected {
				//fmt.Println("Vals", vals)
				dataFromSrc[identifier] = vals[0]
			}
			broadcastBinary(detected, m.GetId(), m.GetRound())
		}
	}
}

func ByzRecBin(m Message) {
	if m.GetSenderId() != MyID {
		identifier := serverList[0] + ":" + strconv.Itoa(m.GetRound())
		if l, ok := binarySet[identifier]; !ok {
			firstL := []Message{m}
			binarySet[identifier] = firstL
		} else {
			l = append(l, m)
			binarySet[identifier] = l
			//fmt.Println(binarySet[identifier])
			if len(l) == total - 1 {
				detect := checkDetect(l)

				if detect {
					if _, e:= acceptData[identifier]; !e {
						acceptData[identifier] = true
						//stat := statsRecord[identifier]
						//stat.End = time.Now()
						//fmt.Printf("Stats: %+v\n",stat)
						//diff := fmt.Sprintf("%f",stat.End.Sub(stat.Start).Seconds())
						//fmt.Println()
						//fmt.Println(stat.Start.String(), stat.End.String())
						//fmt.Println("Reliable Accept Failure"  + strconv.Itoa(m.GetRound()) + " " + diff)
						//fmt.Println()

						stats := StatStruct{Id:m.GetId(), Round: m.GetRound(), Header:Stat}
						statInfo :=PrepareSend{M:stats, SendTo:MyID}
						SendReqChan <- statInfo
					}

				} else {
					if _, e:= acceptData[identifier]; !e {
						acceptData[identifier] = true
						//stat := statsRecord[identifier]
						//stat.End = time.Now()
						//fmt.Printf("Stats: %+v\n",stat)
						//diff := fmt.Sprintf("%f",stat.End.Sub(stat.Start).Seconds())
						//fmt.Println()
						//fmt.Println("Reliable Accept "  + strconv.Itoa(m.GetRound()) + " " + diff, dataFromSrc[identifier])
						//fmt.Println()

						stats := StatStruct{Id:m.GetId(), Round: m.GetRound(), Header:Stat}
						statInfo :=PrepareSend{M:stats, SendTo:MyID}
						SendReqChan <- statInfo
					}
				}

			}
		}
	}
}

func validateByzCode(vals []string) bool{
	//fmt.Println(vals)
	data := vals[0]
	for _, val := range vals {
		if val != data {
			//fmt.Println("Wrong", val, data)
			return true
		}
	}
	return false
}



