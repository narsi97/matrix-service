package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

func main() {
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/multiply", multiplyHandler)
	http.HandleFunc("/flatten", flattenHandler)
	http.ListenAndServe(":8080", nil)
}

func processCSV(w http.ResponseWriter, r *http.Request) ([][]string, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	records, err := processCSV(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
	}
	var response string
	for _, row := range records {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}
	fmt.Fprint(w, response)
}
func addHandler(w http.ResponseWriter, r *http.Request) {
	records, err := processCSV(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
	}
	//var response string
	sum := 0
	for _, row := range records {
		for _, col := range row {
			val, err := strconv.Atoi(col)
			if err == nil {
				sum += val
			}
		}
	}
	fmt.Fprintf(w, "%d\n", sum)
}
func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	records, err := processCSV(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
	}
	//var response string
	product := 1
	for _, row := range records {
		for _, col := range row {
			val, err := strconv.Atoi(col)
			if err == nil {
				product *= val
			}
		}
	}
	fmt.Fprintf(w, "%d\n", product)
}

// func invertHandler(w http.ResponseWriter, r *http.Request) {
// 	records, err := processCSV(w, r)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
// 	}
// 	//var response string
// 	for _, row := range records {
// 		for _, col := range row {
// 			val, err := strconv.Atoi(col)
// 			if err == nil {
// 				product *= val
// 			}
// 		}
// 	}
// 	fmt.Fprintf(w, "%d\n", product)
// }

func flattenHandler(w http.ResponseWriter, r *http.Request) {
	records, err := processCSV(w, r)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
	}
	//var flatted []int;
	var response string
	for _, row := range records {
		//for _, col := range row {
		//val, err := strconv.Atoi(col)
		//if err == nil {
		//	flatted = append(flatted, val)
		//}
		response = response + strings.Join(row, ",") + ","
		//}
	}
	fmt.Sprintf("%s \n", response.)
	//fmt.Fprintf(w, "%s", response)
}
