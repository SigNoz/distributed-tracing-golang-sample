package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type errResponse struct {
	Message string `json:"message"`
}

func ReadBody(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return fmt.Errorf("read body error: %w", err)
	}

	// unmarshal into object
	if err := json.Unmarshal(body, obj); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return fmt.Errorf("json unmarshal error: %w", err)
	}

	// validate object
	if err := validator.New().Struct(obj); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return fmt.Errorf("validate object error: %w", err)
	}

	return nil
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	WriteResponse(w, statusCode, errResponse{err.Error()})
}

func WriteResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("encode response error: %v", err)
	}
}

func SendRequest(ctx context.Context, method string, url string, data []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	client := http.Client{
		// Wrap the Transport with one that starts a span and injects the span context
		// into the outbound request headers.
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	return client.Do(request)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	log.Printf("status: %d", statusCode)
	rw.statusCode = statusCode
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	// WriteHeader() is not claled if our response implicitly returns 200 OK, so
	// we default to that status code
	return &responseWriter{w, http.StatusOK}
}

func LoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// wrap the response writer to capture the response
		rw := newResponseWriter(w)
		// Once the body is read, it cannot be re-read. Hence, use the TeeReader
		// to write the r.Body to buf as it is being read.
		// This buf is later used for logging.
		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		r.Body = ioutil.NopCloser(tee)
		next.ServeHTTP(rw, r)

		log.Printf("%s %s %d from %s (query:%s body:%s)",
			r.Method, r.URL.Path, rw.statusCode, r.RemoteAddr, r.URL.Query(), buf.String())
	})
}
