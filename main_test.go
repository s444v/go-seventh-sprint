package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Request struct {
	request string
	status  int
	message string
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []Request{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []Request{
		{"/cafe?count=2&city=moscow", http.StatusOK, ""},
		{"/cafe?city=tula", http.StatusOK, ""},
		{"/cafe?city=moscow&search=ложка", http.StatusOK, ""},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, v.status, response.Code)
	}
}

func TestCafeCountWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	moscowCafeList, ok := cafeList["moscow"]
	if !ok {
		return
	}
	requests := []Request{
		{"/cafe?city=moscow&count=0", http.StatusOK, ""},
		{"/cafe?city=moscow&count=1", http.StatusOK, strings.Join(moscowCafeList[:1], ",")},
		{"/cafe?city=moscow&count=2", http.StatusOK, strings.Join(moscowCafeList[:2], ",")},
		{"/cafe?city=moscow&count=101", http.StatusOK, strings.Join(moscowCafeList, ",")},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCountNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []Request{
		{"/cafe?city=moscow&count=na", http.StatusBadRequest, "incorrect count"},
		{"/cafe?city=moscow&count=blabla", http.StatusBadRequest, "incorrect count"},
		{"/cafe?city=moscow&count=-1", http.StatusBadRequest, "incorrect count"},
		{"/cafe?city=moscow&count=10000000000000000000000000000000000000000000000000000", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []Request{
		{"/cafe?city=moscow&search=фасоль", http.StatusOK, ""},
		{"/cafe?city=moscow&search=кофе", http.StatusOK, "Мир кофе,Кофе и завтраки"},
		{"/cafe?city=moscow&search=вилка", http.StatusOK, "Ложка и вилка"},
		{"/cafe?city=moscow&search=", http.StatusOK, strings.Join(cafeList["moscow"], ",")},
		{"/cafe?city=moscow&search=бредик", http.StatusOK, ""},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}
