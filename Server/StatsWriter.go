package Server

import (
	"HRB/HRBAlgorithm"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type StatsCounter struct {
	counter int
	m       * sync.Mutex
}

func (counter *StatsCounter) increment(){
	counter.counter += 1
}

func (counter *StatsCounter) getCount() int{
	return counter.counter
}

type flag struct {
	flag bool
	m       * sync.Mutex
}

var statsMap map[string] time.Time
var statsMapMutex sync.RWMutex
var acceptMap map[string] string
var counter StatsCounter
var startTime time.Time
var reachFlag flag
var latencyMap map[string] time.Time


func initStats() {
	statsMap = make(map[string] time.Time)
	acceptMap = make(map[string] string)
	statsMapMutex =  sync.RWMutex{}
	counter = StatsCounter{counter:0, m: &sync.Mutex{}}
	reachFlag = flag{flag:false, m: &sync.Mutex{}}
	startTime = time.Now()
	latencyMap = make(map[string] time.Time)
}

func (f *flag)setFlag() {
	f.m.Lock()
	f.flag = true
	f.m.Lock()
}

func (f *flag)getFlag() bool{
	var reach bool
	f.m.Lock()
	reach = f.flag
	f.m.Unlock()
	return reach
}

func statsCalculate(statsChan chan HRBAlgorithm.Message) {
	go latencyCalculator(statsChan)
	//go counterUpDate(statsChan)
	//go throughputCalculator()
}


/*
Latency Part
 */

func latencyCalculator(statsChan chan HRBAlgorithm.Message) {
	for {
		data := <- statsChan
		identifier := strconv.Itoa(data.GetRound())

		if source {
			if data.GetHeaderType() == HRBAlgorithm.MSG {
				latencyMap[identifier] = time.Now()
			}
		} else {
			if data.GetHeaderType() == HRBAlgorithm.Stat {
				latencyMap[identifier] = time.Now()
			}
		}

		if data.GetHeaderType() == HRBAlgorithm.Stat {
			counter.increment()

			//If equal to the total Round flush to a file
			if counter.getCount() == round {
				lapse := fmt.Sprintf("%f", time.Now().Sub(throughPutBeginTime).Seconds())
				timeLapse, _ := strconv.ParseFloat(lapse, 32)
				//fmt.Println(round, timeLapse, startTime.String())
				throughPut := float64(round) / timeLapse
				writeThroughPut(throughPut)
				writeLatencyFile()
				fmt.Println("Successful receive all message")
			}
		}
	}
}

func writeLatencyFile() {
	fileName := strconv.Itoa(algorithm) +":Latency" + "|" + MyId +"|" + startTime.String()+".txt"
	file, err := os.OpenFile(filepath.Join("./Data", fileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(round+ ":" + latency )
	for round, time := range latencyMap {
		if source {
			fmt.Fprintf(file, "Start^" + round + "^" + time.String() +"\n")
		} else {
			fmt.Fprintf(file, "End^" + round + "^" + time.String() +"\n")
		}
	}

	file.Close()
}

//func writeAllSuccess(throughput string) {
//	fileName := strconv.Itoa(algorithm) +":Latency" + "|" + MyId +"|" + startTime.String()+".txt"
//	file, err := os.OpenFile(filepath.Join("./Data", fileName), os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("Successful receive all message")
//}

func writeThroughPut(throuput float64) {
	fileName := strconv.Itoa(algorithm) +":Throuput" + "|" + MyId +"|" + startTime.String()+".txt"
	file, err := os.OpenFile(filepath.Join("./Data", fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(file, fmt.Sprintf("%f",throuput)+"\n")
	file.Close()
}


/*
Throughput Part
 */
//func throughputCalculator() {
//	//time.Sleep(65* time.Second)
//	//Write to File
//	count := counter.getCount()
//	writeThroughPut(count)
//}


func latencyCalculator1Min() {
	time.Sleep(1* time.Minute)
}

func counterUpDate(statsChan chan HRBAlgorithm.Message) {
	for {
		data := <- statsChan
		identifier := strconv.Itoa(data.GetRound())

		//fmt.Println("Stats_Counter: " + identifier)

		if _, e:= acceptMap[identifier];!e {
			//end := time.Now()
			//diff := fmt.Sprintf("%f", end.Sub(start).Seconds())
			//fmt.Println(end.String(), start.String())

			//acceptMap[identifier] = identifier
			//writeLatencyFile(identifier, diff)
			counter.increment()
		}

		//If equal to the total Round flush to a file
		if counter.counter == round {
			//writeAllSuccess()
		}
	}
}