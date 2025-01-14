package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var hosts = []string{
	"http://172.232.132.21",
	"http://172.232.132.130",
	"http://172.232.158.39",
	"http://172.232.158.78",
	"http://172.232.158.51",
	"http://172.232.158.61",
	"http://172.232.148.207",
	"http://172.232.148.51",
}

type Result struct {
	ID         string `json:"id"`
	TestID     string `json:"test_id"`
	Second     int    `json:"second"`
	Requests   int    `json:"requests"`
	ErrorCodes string `json:"error_codes"`
}

type Response struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Method      string   `json:"method"`
	ContentType string   `json:"content_type"`
	Status      string   `json:"status"`
	Body        string   `json:"body"`
	CreatedAt   string   `json:"created_at"`
	FinishedAt  string   `json:"finished_at"`
	Results     []Result `json:"results"`
}

func main() {
	var wg sync.WaitGroup

	ch := make(chan []Result, 1000)

	for _, url := range hosts {
		wg.Add(1)

		go func(remoteUrl string) {
			defer wg.Done()

			payload := strings.NewReader(`{
				"method": "GET",
				"tasks": 100,
				"seconds": 400,
				"start_at": "2023-09-17T10:16:34.675Z",
				"url": "http://172.232.156.13:3000/json", 
				"content_type": "application/json",
				"body": ""
			}`)
			// TODO:: change ip
			res, err := http.Post(remoteUrl, "application/json", payload)
			if err != nil {
				return
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(body))

			var response Response
			err = json.Unmarshal(body, &response)
			if err != nil {
				fmt.Println("Unmarshal post request")
				fmt.Println(err)
				return
			}

			for {
				getTest, err := http.Get(fmt.Sprintf("%s/%s", remoteUrl, response.ID))
				if err != nil {
					fmt.Println(err)
					return
				}

				body, err := io.ReadAll(getTest.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(body))

				var getTestResponse Response
				err = json.Unmarshal(body, &getTestResponse)
				if err != nil {
					return
				}

				if getTestResponse.Status == "PROCESSING" {
					time.Sleep(1 * time.Second)
					continue
				}

				if getTestResponse.Status == "FINISHED" {
					ch <- getTestResponse.Results

					break
				}

				fmt.Println(fmt.Sprintf("unexpected status :%s", getTestResponse.Status))
				break
			}
		}(url)
	}

	wg.Wait()

	close(ch)

	var finalResult []Result

	for r := range ch {
		for i, res := range r {
			if len(finalResult) < (i + 1) {
				finalResult = append(finalResult, res)
			} else {
				finalResult[i].Requests += res.Requests
				finalResult[i].ErrorCodes += res.ErrorCodes
			}
		}
	}

	data := [][]string{}

	for _, row := range finalResult {
		s := strconv.Itoa(row.Second)
		r := strconv.Itoa(row.Requests)
		l := strconv.Itoa(len(strings.Split(row.ErrorCodes, ",")))
		data = append(data, []string{s, r, l})
	}

	file, err := os.Create("results/node.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	writer.Write([]string{"Second", "Requests", "Error codes"})
	writer.WriteAll(data)
}
