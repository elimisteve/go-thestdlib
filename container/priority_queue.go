package main

import (
    "container/heap"
    "fmt"
    "log"
    "math/rand"
    "time"
)

const (
    MaxQueueLength = 100
    MaxRequesters  = 10
    Seconds        = 2e9
)

func requester(work chan Request) {
    for {
        time.Sleep(time.Duration(rand.Int63n(MaxRequesters * Seconds)))
        work <- func() {
            time.Sleep(time.Duration(rand.Int63n(MaxRequesters*Seconds) + 10))
        }
    }
}

type Request func()

type Worker struct {
    id       int
    pending  int
    requests chan Request
    index    int
}

func (w *Worker) work(done chan *Worker) {
    for {
        req := <-w.requests
        req()
        done <- w
    }
}

func (w *Worker) String() string {
    return fmt.Sprintf("W%d{pending: %d}", w.id, w.pending)
}

type Pool []*Worker

func (p Pool) Len() int {
    return len(p)
}

func (p Pool) Less(i, j int) bool {
    return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
    a := *p
    a[i], a[j] = a[j], a[i]
    a[i].index = i
    a[j].index = j
}

func (p *Pool) Push(i interface{}) {
    w := i.(*Worker)
    a := *p
    n := len(a)
    w.index = n
    a = append(a, w)
    *p = a
}

func (p *Pool) Pop() interface{} {
    a := *p
    n := len(a)
    w := a[n-1]
    w.index = -1
    *p = a[0 : n-1]
    return w
}

type Balancer struct {
    pool Pool
    done chan *Worker
}

func NewBalancer(size int) *Balancer {
    done := make(chan *Worker, size)
    b := &Balancer{
        pool: make(Pool, 0, size),
        done: done,
    }
    for i := 0; i < size; i++ {
        w := &Worker{id: i, requests: make(chan Request, MaxQueueLength)}
        heap.Push(&b.pool, w)
        go w.work(done)
    }
    return b
}

func (b *Balancer) Balance(requests chan Request) {
    for {
        select {
        case req := <-requests:
            b.dispatch(req)
            log.Printf("New request, %s", b.pool)
        case w := <-b.done:
            b.completed(w)
            log.Printf("Request finished, %s", b.pool)
        }
    }
}

func (b *Balancer) dispatch(req Request) {
    w := heap.Pop(&b.pool).(*Worker)
    w.requests <- req
    w.pending++
    heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
    w.pending--
    heap.Remove(&b.pool, w.index)
    heap.Push(&b.pool, w)
}

func main() {
    requests := make(chan Request)
    for i := 0; i < MaxRequesters; i++ {
        go requester(requests)
    }
    NewBalancer(4).Balance(requests)
}
