package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/tamnguyenwizeline/ondemand-go-bootcamp/model"
)

func ReadCSVFile(filePath string) ([]model.Pokemon, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []model.Pokemon
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %v", record[0])
		}

		pokemon := model.Pokemon{
			ID:   id,
			Name: record[1],
		}
		result = append(result, pokemon)
	}

	return result, nil
}

func worker(id int, jobs <-chan int, results chan<- model.Pokemon, file *os.File, itemsPerWorker int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)

		reader := csv.NewReader(file)

		// Read and process the specified number of items
		// Read the record
		record, err := reader.Read()
		if err != nil {
			// Handle any error in reading the CSV file
			fmt.Println("Error reading CSV file:", err)
			return
		}
		fmt.Println("Record", record)

		// Process the record
		id, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Printf("Invalid ID: %v\n", record[0])
			return
		}

		pokemon := model.Pokemon{
			ID:   id,
			Name: record[1],
		}
		results <- pokemon
	}
}

func ReadCSVDataWithWorkers(filePath string, items, itemsPerWorker int) ([]*model.Pokemon, error) {
	file, err := os.Open(filePath)
	response := make([]*model.Pokemon, 0)

	if err != nil {
		return nil, err
	}
	defer file.Close()
	fcsv := csv.NewReader(file)

	totalWorkers := items / itemsPerWorker
	// Create a channel to communicate with workers
	jobs := make(chan []string, totalWorkers)
	result := make(chan *model.Pokemon)

	var wg sync.WaitGroup
	worker := func(jobs <-chan []string, results chan<- *model.Pokemon) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}
				results <- parseData(job)
			}
		}
	}

	for w := 0; w < totalWorkers; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output at line 107 (func worker: line 71)
			defer wg.Done()
			worker(jobs, result)
		}()
	}

	go func() {
		for {
			rStr, err := fcsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				break
			}
			jobs <- rStr
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(result) // when you close(res) it breaks the below loop.
	}()

	for r := range result {
		response = append(response, r)
	}

	return response, nil
}

func parseData(data []string) *model.Pokemon {
	id, _ := strconv.Atoi(data[0])
	return &model.Pokemon{
		ID:   id,
		Name: data[1],
	}
}
