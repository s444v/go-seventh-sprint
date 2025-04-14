package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestsTable struct {
	request      string
	wantStatus   int
	wantResponse string
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []TestsTable{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range tests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, v.wantStatus, response.Code)
		assert.Equal(t, v.wantResponse, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	tests := []TestsTable{
		{"/cafe?count=2&city=moscow", http.StatusOK, ""},
		{"/cafe?city=tula", http.StatusOK, ""},
		{"/cafe?city=moscow&search=ложка", http.StatusOK, ""},
	}
	for _, v := range tests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, v.wantStatus, response.Code)
	}
}

func TestCafeCountWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	path := "/cafe?city=moscow&count="
	cafe, ok := cafeList["moscow"]
	require.True(t, ok)
	tests := []struct {
		count      int // передаваемое значение count
		wantStatus int // ожидаемый статус ответа
		wantCount  int // ожидаемое количество кафе в ответе
	}{
		{0, http.StatusOK, 0},
		{1, http.StatusOK, 1},
		{2, http.StatusOK, 2},
		{100, http.StatusOK, min(100, len(cafe))},
	}
	for _, v := range tests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, v.wantStatus, response.Code)
		if response.Body.String() == "" {
			assert.Equal(t, v.wantCount, 0)
		} else {
			responseBody := strings.Split(response.Body.String(), ",")
			assert.Equal(t, v.wantCount, len(responseBody))
		}
	}
}

func TestCafeCountNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	path := "/cafe?city=moscow&count="
	tests := []struct {
		count        string // передаваемое значение count
		wantStatus   int    // ожидаемый статус ответа
		wantResponse string // ожидаемое в ответе
	}{
		{"na", http.StatusBadRequest, "incorrect count"},
		{"-1", http.StatusBadRequest, "incorrect count"},
		{"10000000000000000000000000000000000000000000000000000", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range tests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path+v.count, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, v.wantStatus, response.Code)
		assert.Equal(t, v.wantResponse, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	path := "/cafe?city=moscow&search="
	tests := []struct {
		search     string // передаваемое значение search
		wantStatus int    // ожидаемый статус ответа
		wantCount  int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", http.StatusOK, 0},
		{"кофе", http.StatusOK, 2},
		{"вилка", http.StatusOK, 1},
		{"", http.StatusOK, 5},
	}
	for _, v := range tests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path+v.search, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, v.wantStatus, response.Code)
		if response.Body.String() == "" {
			assert.Equal(t, v.wantCount, 0)
		} else {
			responseBody := strings.Split(response.Body.String(), ",")
			assert.Equal(t, v.wantCount, len(responseBody))
		}
	}
}
