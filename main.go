package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	guuid "github.com/google/uuid"
)

/*
InputDataType ...
*/
type InputDataType struct {
	uuid string
}

/*
OutputDataType ..
*/
type OutputDataType struct {
	hash                    string
	timeTakenForComputation int64
}

func generateInputDataList() []InputDataType {
	myDataList := make([]InputDataType, 0)
	for i := 0; i < 50; i++ {
		myUUID := generateUUID4()
		data := InputDataType{uuid: myUUID}
		myDataList = append(myDataList, data)
	}
	return myDataList
}

func executeAllJobs() []interface{} {
	jobs := make(chan InputDataType, 1)
	results := make(chan OutputDataType, 1)
	resultList := make([]interface{}, 0)
	maxNumberOfWorkers := 5

	myListOfData := generateInputDataList()

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for w := 1; w <= maxNumberOfWorkers; w++ {
			go workerFunc(jobs, results)
		}

		fmt.Printf("\n@ Putting data into jobs...\n")
		for _, jobData := range myListOfData {
			jobs <- jobData
		}
		fmt.Printf("\n@ jobs channel got all data !\n")
		close(jobs)
		fmt.Printf("\n@ jobs channel closed !\n")
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 1; i <= len(myListOfData); i++ {
			x := <-results
			resultList = append(resultList, x)
		}
	}(&wg)

	fmt.Printf("\n BEFORE wg.Wait()... ! \n")

	wg.Wait()

	fmt.Printf("\n AFTER wg.Wait()... ! \n")

	fmt.Printf("\n resultList : %v \n", resultList)

	fmt.Printf("\n executeAllJobs() Done ! \n")

	return resultList
}

func workerFunc(jobs <-chan InputDataType, results chan<- OutputDataType) {
	for job := range jobs {
		startTime := GetMillis()
		hashValue := getMD5Hash(job)
		fmt.Printf("\n workerFunc() : hashValue : %v \n", hashValue)
		endTime := GetMillis()
		timeTakenForCalculating := endTime - startTime
		results <- OutputDataType{hash: hashValue, timeTakenForComputation: timeTakenForCalculating}
	}
}

// ---- Helper Func ---
func getMD5Hash(data InputDataType) string {
	time.Sleep(2 * time.Second)
	hasher := md5.New()
	hasher.Write([]byte(data.uuid))
	md5Value := hex.EncodeToString(hasher.Sum(nil))
	return md5Value
}

/*
GetMillis ...
// ---- Helper Func ---
*/
func GetMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

// ---- Helper Func ---
func generateUUID4() string {
	id := guuid.New()
	generatedUUID := id.String()
	return generatedUUID
}

func main() {
	resultData := executeAllJobs()
	for _, item := range resultData {
		data := item.(OutputDataType)
		fmt.Printf("\ndata.hash : (%v) , data.timeTakenForComputation (%v) milli-secs\n", data.hash, data.timeTakenForComputation)
	}
}
