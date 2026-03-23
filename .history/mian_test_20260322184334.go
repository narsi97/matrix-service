package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/////////////////////////
// Helpers

func createMultipartRequest(t *testing.T, endpoint string, content string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "matrix.csv")
	if err != nil {
		t.Fatal(err)
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", endpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

/////////////////////////
// Business Logic Tests
/////////////////////////

func TestSum(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}

	result := sum(matrix)
	expected := 10

	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestMultiply(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}

	result := multiply(matrix)
	expected := 24

	if result != expected {
		t.Errorf("expected %d, got %d", expected, result)
	}
}

func TestFlatten(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}

	result := flatten(matrix)
	expected := "1,2,3,4"

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestInvert(t *testing.T) {
	matrix := [][]int{
		{1, 2},
		{3, 4},
	}

	result := invert(matrix)

	expected := [][]int{
		{1, 3},
		{2, 4},
	}

	for i := range expected {
		for j := range expected[i] {
			if result[i][j] != expected[i][j] {
				t.Errorf("expected %v, got %v", expected, result)
			}
		}
	}
}

/////////////////////////
// Handler Tests
/////////////////////////

func TestEchoHandler(t *testing.T) {
	req := createMultipartRequest(t, "/echo", "1,2\n3,4")
	rr := httptest.NewRecorder()

	echoHandler(rr, req)

	expected := "1,2\n3,4\n"

	if rr.Body.String() != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}

func TestSumHandler(t *testing.T) {
	req := createMultipartRequest(t, "/add", "1,2\n3,4")
	rr := httptest.NewRecorder()

	addHandler(rr, req)

	expected := "10"

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}

func TestMultiplyHandler(t *testing.T) {
	req := createMultipartRequest(t, "/multiply", "1,2\n3,4")
	rr := httptest.NewRecorder()

	multiplyHandler(rr, req)

	expected := "24"

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}

func TestFlattenHandler(t *testing.T) {
	req := createMultipartRequest(t, "/flatten", "1,2\n3,4")
	rr := httptest.NewRecorder()

	flattenHandler(rr, req)

	expected := "1,2,3,4"

	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}

func TestInvertHandler(t *testing.T) {
	req := createMultipartRequest(t, "/invert", "1,2\n3,4")
	rr := httptest.NewRecorder()

	invertHandler(rr, req)

	expected := "1,3\n2,4\n"

	if rr.Body.String() != expected {
		t.Errorf("expected %s, got %s", expected, rr.Body.String())
	}
}

/////////////////////////
// Error Handling Tests
/////////////////////////

func TestInvalidMatrix(t *testing.T) {
	req := createMultipartRequest(t, "/add", "1,2,3\n4,5,6") // not square
	rr := httptest.NewRecorder()

	addHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestInvalidData(t *testing.T) {
	req := createMultipartRequest(t, "/add", "1,2\n3,x")
	rr := httptest.NewRecorder()

	addHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}