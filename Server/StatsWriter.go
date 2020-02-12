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
	counter.m.Lock()
	counter.counter += 1
	counter.m.Unlock()
}

func (counter *StatsCounter) getCount() int{
	var count int
	counter.m.Lock()
	count = counter.counter
	counter.m.Unlock()
	return count
}

type flag struct {
	flag bool
	m       * sync.Mutex
}

var statsMap map[string] time.Time
var acceptMap map[string] string
var counter StatsCounter
var startTime time.Time
var reachFlag flag


func initStats() {
	statsMap = make(map[string] time.Time)
	acceptMap = make(map[string] string)
	counter = StatsCounter{counter:0, m: &sync.Mutex{}}
	reachFlag = flag{flag:false, m: &sync.Mutex{}}
	startTime = time.Now()
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
	go throughputCalculator()
}


/*
Latency Part
 */
func latencyCalculator(statsChan chan HRBAlgorithm.Message) {
	initStats()
	for {
		data := <- statsChan
		identifier := strconv.Itoa(data.GetRound())

		//fmt.Println("Stats_Counter: " + identifier)

		start, recorded := statsMap[identifier]
		if recorded {
			if _, e := acceptMap[identifier]; !e {
				end := time.Now()
				diff := fmt.Sprintf("%f", end.Sub(start).Seconds())
				//fmt.Println(end.String(), start.String())
				acceptMap[identifier] = diff
				counter.increment()
			}
		} else {
			counter.increment();
		}

		//If equal to the total Round flush to a file
		if counter.counter == round {
			writeLatencyFile()
			time.Sleep(30*time.Second)
			writeAllSuccess()
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
	for round, latency := range acceptMap {
		fmt.Fprintf(file, round+ ":" + latency +"\n")
	}
}

func writeAllSuccess() {
	//fileName := strconv.Itoa(algorithm) +":Latency" + "|" + MyId +"|" + startTime.String()+".txt"
	//file, err := os.OpenFile(filepath.Join("./Data", fileName), os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	fmt.Println("Successful receive all message")
}

func writeThroughPut(throuput int) {
	fileName := strconv.Itoa(algorithm) +":Throuput" + "|" + MyId +"|" + startTime.String()+".txt"
	file, err := os.OpenFile(filepath.Join("./Data", fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(file, strconv.Itoa(throuput)+"\n")
}


/*
Throughput Part
 */
func throughputCalculator() {
	time.Sleep(1* time.Minute)
	//Write to File
	writeThroughPut(counter.getCount())
}


func latencyCalculator1Min() {
	time.Sleep(1* time.Minute)
}