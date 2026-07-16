package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: 0},
	}

	totalCafes := len(cafeList["moscow"])

	requests[3].want = totalCafes
	if totalCafes > 100 {
		requests[3].want = 100
	}

	for _, req := range requests {
		response := httptest.NewRecorder()
		respUrl := "/cafe?city=moscow&count=" + strconv.Itoa(req.count)

		httpReq := httptest.NewRequest("GET", respUrl, nil)
		handler.ServeHTTP(response, httpReq)

		assert.Equal(t, http.StatusOK, response.Code)

		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")

		gotCafes := len(cafes)
		if response.Body.String() == "" || response.Body.String() == "/n" {
			gotCafes = 0
		}

		assert.Equal(t, req.want, gotCafes, "expected: %d, got: %d", req.want, gotCafes)
	}
}
