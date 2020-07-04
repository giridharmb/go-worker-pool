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
	for j := range jobs {
		hashValue := getMD5Hash(j.uuid)
		results <- OutputDataType{hash: hashValue}
	}
}

// ---- Helper Func ---
func getMD5Hash(text string) string {
	time.Sleep(2 * time.Second)
	hasher := md5.New()
	hasher.Write([]byte(text))
	md5Value := hex.EncodeToString(hasher.Sum(nil))
	return md5Value
}

// ---- Helper Func ---
func generateUUID4() string {
	id := guuid.New()
	generatedUUID := id.String()
	return generatedUUID
}

func main() {
	hashList := executeAllJobs()
	fmt.Printf("%v", hashList)
}
