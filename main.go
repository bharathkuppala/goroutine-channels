package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var (
	result  float64
	logger  *log.Logger
	result1 float64
)

// Numbers ...
type Numbers struct {
	l *log.Logger
}

// NewNumber ...
func NewNumber(l *log.Logger) *Numbers {
	return &Numbers{l}
}

func main() {
	router := mux.NewRouter()
	logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
	numberOp := NewNumber(logger)
	var wg sync.WaitGroup
	wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	router.HandleFunc("/api/v1/add-number", addNumbers).Methods("POST")
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	router.HandleFunc("/api/v1/get-numbers", getNumbers).Methods("GET")
	// }()

	go func() {
		defer wg.Done()
		router.Handle("/api/v1/number", numberOp).Methods("POST")
	}()

	go func() {
		defer wg.Done()
		router.Handle("/api/v1/getNum", numberOp).Methods("GET")
	}()

	wg.Wait()

	server := &http.Server{
		Addr:    ":6000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("something went wrong while starting server", err)
		return
	}
}

// ServeHTTP ...
func (n *Numbers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch := make(chan float64, 1)
	if r.Method == http.MethodPost {
		if r.URL.Path == "/api/v1/check" {
			fmt.Println("Am here!!")
			result1 = n.add(w, r, ch)
			fmt.Println(result1)
			return
		}
	}

	if r.Method == http.MethodGet {
		if r.URL.Path == "/api/v1/getNum" {
			n.getNumber(w, r, result1)
			return
		}
	}
}

// with chan
func (n *Numbers) add(w http.ResponseWriter, r *http.Request, ch chan float64) float64 {
	var requestBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("error in decoding the request body")
		return 0
	}

	if requestBody["firstNumber"] == nil && requestBody["secondNumber"] == nil {
		log.Println("values cannot be emty")
		return 0
	}

	// get values from requestBody
	firstNum := requestBody["firstNumber"].(float64)
	secondNum := requestBody["secondNumber"].(float64)

	log.Printf("firstVal is %f and secondVal is %f\n", firstNum, secondNum)

	result = firstNum + secondNum
	ch <- result
	return <-ch
}

func (n *Numbers) getNumber(w http.ResponseWriter, r *http.Request, result1 float64) {
	s := strconv.Itoa(int(result1))
	w.Write([]byte("the added numbers value is " + s))
}

// without chan
func addNumbers(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("error in decoding the request body")
		return
	}

	if requestBody["firstNumber"] == nil && requestBody["secondNumber"] == nil {
		log.Println("values cannot be emty")
		return
	}

	// get values from requestBody
	firstNum := requestBody["firstNumber"].(float64)
	secondNum := requestBody["secondNumber"].(float64)

	log.Printf("firstVal is %f and secondVal is %f\n", firstNum, secondNum)

	result = firstNum + secondNum
	log.Printf("result: %f\n", result)
}
