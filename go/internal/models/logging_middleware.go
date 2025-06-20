package models

import (
	"bytes"
	"net/http"
)

// Сохранение статуса кода и сообщения
type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode    int
	StatusMessage string
}

// Конструктор для loggingResponseWriter
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, "OK"}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.StatusMessage = http.StatusText(code)
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	var errDto ExceptionDto
	JsonToStruct(&errDto, bytes.NewReader(b))
	lrw.StatusMessage += ": " + errDto.ErrorMessage
	return lrw.ResponseWriter.Write(b)
}
