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
var statsMapMutex sync.RWMutex
var acceptMap map[string] string
var counter StatsCounter
var startTime time.Time
var reachFlag flag


func initStats() {
	statsMap = make(map[string] time.Time)
	acceptMap = make(map[string] string)
	statsMapMutex =  sync.RWMutex{}
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
	//go latencyCalculator(statsChan)
	//go counterUpDate(statsChan)
	//go throughputCalculator()
}


/*
Latency Part
 */

func counterUpDate(statsChan chan HRBAlgorithm.Message) {
	for {
		data := <- statsChan
		identifier := strconv.Itoa(data.GetRound())

		//fmt.Println("Stats_Counter: " + identifier)

		if _, e:= acceptMap[identifier];!e {
			//end := time.Now()
			//diff := fmt.Sprintf("%f", end.Sub(start).Seconds())
			//fmt.Println(end.String(), start.String())

			acceptMap[identifier] = identifier
			//writeLatencyFile(identifier, diff)
			counter.increment()
		}

		//If equal to the total Round flush to a file
		if counter.counter == round {
			writeAllSuccess()
		}
	}
}

func latencyCalculator(statsChan chan HRBAlgorithm.Message) {
	for {
		data := <- statsChan
		identifier := strconv.Itoa(data.GetRound())

		//fmt.Println("Stats_Counter: " + identifier)

		statsMapMutex.RLock()
		start, recorded := statsMap[identifier]
		statsMapMutex.RUnlock()

		if recorded {

			_, e:= acceptMap[identifier]

			if !e {
				end := time.Now()
				diff := fmt.Sprintf("%f", end.Sub(start).Seconds())
				//fmt.Println(end.String(), start.String())

				acceptMap[identifier] = diff
				//writeLatencyFile(identifier, diff)

				counter.increment()
			}
		} else {
			counter.increment();
		}

		//If equal to the total Round flush to a file
		if counter.counter == round {
			writeAllSuccess()
		}
	}
}

func writeLatencyFile(round, latency string) {
	fileName := strconv.Itoa(algorithm) +":Latency" + "|" + MyId +"|" + startTime.String()+".txt"
	file, err := os.OpenFile(filepath.Join("./Data", fileName), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(round+ ":" + latency )
	fmt.Fprintf(file, round+ ":" + latency +"\n")

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
	time.Sleep(65* time.Second)
	//Write to File
	count := counter.getCount()
	writeThroughPut(count)
}


func latencyCalculator1Min() {
	time.Sleep(1* time.Minute)
}