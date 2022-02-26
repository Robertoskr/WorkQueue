package main

import (
	"context"
	"sync"
	"fmt"
	"time"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type task func([]arg) error

// arguments to async tasks 
type arg interface{}

type ServerCall struct {
	response http.ResponseWriter
	request *http.Request
	wg *sync.WaitGroup
}

func dummy_task(n []arg) error {
	new_n, _ := n[1].(int)
	_, _ = n[0].(context.Context)
	time.Sleep(time.Duration(new_n) * time.Second)
	fmt.Println(n)
	return nil
}

type Node struct {
	val *ServerCall
	next *Node
	prev *Node
}


type Queue struct {
	sync.Mutex
	head *Node
	back *Node
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) enqueue(c *ServerCall){
	q.Lock()
	new_node := Node{
		val: c,
		next: q.head,
		prev: nil,
	}
	if q.head != nil {
		q.head.prev = &new_node
	}
	q.head = &new_node
	if q.back == nil{
		q.back = q.head
	}
	q.Unlock()
}

func (q *Queue) dequeue(ctx context.Context) (*ServerCall, error){
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		ctx.Err();
		return nil, nil
	}

	result_node := q.back
	q.back = q.back.prev
	if q.back == nil{
		q.head = nil
	}

	return result_node.val, nil
}

func (q *Queue) empty() bool {
	return q.head == nil
}


type WorkQueue struct {
	tasks Queue
	wg *sync.WaitGroup
}

func newWorkQueue(s *sync.WaitGroup) *WorkQueue {
	return &WorkQueue{wg: s}
}

func (w *WorkQueue) addTask(c *ServerCall) {
	w.tasks.enqueue(c)
}

func (w *WorkQueue) run(ctx context.Context, proxy *httputil.ReverseProxy) {
	for {
		time.Sleep(time.Second / time.Duration(4))
		if !w.tasks.empty(){
			serverCall, err := w.tasks.dequeue(ctx)
			if err == nil {
				proxy.ServeHTTP(
					serverCall.response,
					serverCall.request,
				)
				serverCall.wg.Done()
			}
		}
	}
}

func newDirector(port string) func(*http.Request) {
	origin, _ := url.Parse("http://localhost" + port)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	return director
}

type Router struct {}

func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request){
	var wg sync.WaitGroup
	wg.Add(1)
	worker.addTask(&ServerCall{response, request, &wg})
	wg.Wait()
}

var worker *WorkQueue

func main(){
	ctx := context.TODO()
	var wg sync.WaitGroup
	worker = newWorkQueue(&wg)

	proxy1 := &httputil.ReverseProxy{Director: newDirector(":8000")}
	proxy2 := &httputil.ReverseProxy{Director: newDirector(":9000")}
	proxy3 := &httputil.ReverseProxy{Director: newDirector(":10000")}

	go worker.run(ctx, proxy1)
	go worker.run(ctx, proxy2)
	go worker.run(ctx, proxy3)
	go worker.run(ctx, proxy3)
	go worker.run(ctx, proxy1)
	go worker.run(ctx, proxy2)
	go worker.run(ctx, proxy3)
	go worker.run(ctx, proxy3)

	router := &Router{}

	log.Fatal(http.ListenAndServe(":5000", router))

}