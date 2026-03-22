package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/multiply", multiplyHandler)
	http.HandleFunc("/flatten", flattenHandler)
	http.HandleFunc("/invert", invertHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

//////////////////////
// Parsing & Validation
//////////////////////

func parseCSV(r *http.Request) ([][]int, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("invalid CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("empty matrix")
	}

	matrix := make([][]int, len(records))

	for i, row := range records {
		matrix[i] = make([]int, len(row))
		for j, val := range row {
			num, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return nil, fmt.Errorf("invalid number at row %d col %d", i, j)
			}
			matrix[i][j] = num
		}
	}

	// Validate square matrix
	n := len(matrix)
	for _, row := range matrix {
		if len(row) != n {
			return nil, errors.New("matrix must be square")
		}
	}

	return matrix, nil
}

//////////////////////
// Business Logic
//////////////////////

func echo(matrix [][]int) string {
	var b strings.Builder
	for _, row := range matrix {
		for i, val := range row {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(strconv.Itoa(val))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func flatten(matrix [][]int) string {
	var b strings.Builder
	first := true
	for _, row := range matrix {
		for _, val := range row {
			if !first {
				b.WriteString(",")
			}
			b.WriteString(strconv.Itoa(val))
			first = false
		}
	}
	return b.String()
}

func sum(matrix [][]int) int {
	total := 0
	for _, row := range matrix {
		for _, val := range row {
			total += val
		}
	}
	return total
}

func multiply(matrix [][]int) int {
	product := 1
	for _, row := range matrix {
		for _, val := range row {
			product *= val
		}
	}
	return product
}

func invert(matrix [][]int) [][]int {
	n := len(matrix)
	result := make([][]int, n)

	for i := range result {
		result[i] = make([]int, n)
		for j := range matrix {
			result[i][j] = matrix[j][i]
		}
	}
	return result
}

func matrixToString(matrix [][]int) string {
	var b strings.Builder
	for _, row := range matrix {
		for i, val := range row {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(strconv.Itoa(val))
		}
		b.WriteString("\n")
	}
	return b.String()
}

//////////////////////
// Handlers
//////////////////////

func handleRequest(w http.ResponseWriter, r *http.Request, fn func([][]int) string) {
	matrix, err := parseCSV(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, fn(matrix))
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, echo)
}

func flattenHandler(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, flatten)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	matrix, err := parseCSV(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, sum(matrix))
}

func multiplyHandler(w http.ResponseWriter, r *http.Request) {
	matrix, err := parseCSV(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, multiply(matrix))
}

func invertHandler(w http.ResponseWriter, r *http.Request) {
	matrix, err := parseCSV(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, matrixToString(invert(matrix)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Started %s %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}