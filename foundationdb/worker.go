package foundationdb

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

var WorkerPool FDBWorkerPool

type FDBWorkerPool struct {
	Concurrency    int
	contextChannel chan *FDBWorker
}

func init() {
	MaxProc := runtime.NumCPU()
	fmt.Println("Initiated SECP256K1 workers: " + strconv.Itoa(MaxProc))
	c := make(chan *FDBWorker, MaxProc)
	for i := 0; i < MaxProc; i++ {
		proc := NewFDBWorker()
		c <- proc
	}
	newPool := FDBWorkerPool{Concurrency: MaxProc, contextChannel: c}
	WorkerPool = newPool
}

func Reinit(numProc int) {
	MaxProc := runtime.NumCPU()
	fmt.Println("Initiated SECP256K1 workers: " + strconv.Itoa(MaxProc))
	c := make(chan *FDBWorker, MaxProc)
	for i := 0; i < MaxProc; i++ {
		proc := NewFDBWorker()
		c <- proc
	}
	newPool := FDBWorkerPool{Concurrency: MaxProc, contextChannel: c}
	WorkerPool = newPool
}

type FDBWorker struct {
	db fdb.Database
}

func NewFDBWorker() *FDBWorker {
	db := fdb.MustOpenDefault()
	new := &FDBWorker{db}
	return new
}

func Transact(tx func(tr fdb.Transaction) (interface{}, error)) (interface{}, error) {
	boundContext := <-WorkerPool.contextChannel
	defer func() { WorkerPool.contextChannel <- boundContext }()
	return boundContext.db.Transact(tx)
}

func ReadTransact(tx func(tr fdb.ReadTransaction) (interface{}, error)) (interface{}, error) {
	boundContext := <-WorkerPool.contextChannel
	defer func() { WorkerPool.contextChannel <- boundContext }()
	return boundContext.db.ReadTransact(tx)
}
