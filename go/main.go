package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type worker struct {
	num int
}

type multiThreadedHandler struct {
	workerChan    chan worker
	workIntensity int
}

func NewMultiThreadedHandler(workers, workIntensity int) *multiThreadedHandler {
	workerChan := make(chan worker, workers)
	for i := 0; i < workers; i++ {
		w := worker{num: i}
		workerChan <- w
	}
	return &multiThreadedHandler{
		workerChan:    workerChan,
		workIntensity: workIntensity,
	}

}

type ResponseObject struct {
	Response  string `json:"response"`
	TimeTaken int64  `json:"time_taken_microseconds"`
}

func IntenseON4Operation(n int) {
	// CPU intensive algo
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			for k := 1; k <= n; k++ {
				for l := 1; l <= n; l++ {
					// Perform some CPU-intensive computation
					_ = i * j * k * l
				}
			}
		}
	}
}

func (h *multiThreadedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Acquire a worker from the channel (blocking if the pool is empty)
	n := <-h.workerChan
	go func() {
		start := time.Now()
		IntenseON4Operation(h.workIntensity)
		response := fmt.Sprintf("request processed by worker %d", n.num)
		resp := &ResponseObject{
			Response:  response,
			TimeTaken: time.Since(start).Microseconds(),
		}
		b, _ := json.Marshal(resp)
		_, _ = w.Write(b)

		fmt.Println(string(b))

		// Release the worker back to the channel
		<-h.workerChan
	}()
}

type HTTPService struct {
	port   int
	server *http.Server
}

func NewHTTPService(port, maxThreads int, difficulty int) HTTPService {

	return HTTPService{
		port: port,
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           NewMultiThreadedHandler(maxThreads, difficulty),
			ReadHeaderTimeout: 20 * time.Second,
		},
	}
}

func (service *HTTPService) Start() {
	go func() {
		if err := service.server.ListenAndServe(); err != nil {
			fmt.Println("serverTerminated")
		}
	}()
}

func (service *HTTPService) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return service.server.Shutdown(ctx)
}

func main() {
	service := NewHTTPService(8000, 100, 200)
	service.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	service.Stop()
}
