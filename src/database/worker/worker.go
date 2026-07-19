package worker

import (
	"diyd/src/database"
	"diyd/src/utils"
	"sync"
)

type Params struct {
	WC      int
	KVStore database.KeyValueStore
}

func NewKVWorker(p Params) *KVWorker {
	return &KVWorker{
		wc:    p.WC,
		taskQ: utils.NewQueue[*Task](),
		kvs:   p.KVStore,
	}
}

type KVWorker struct {
	wc    int
	taskQ utils.Queue[*Task]
	kvs   database.KeyValueStore
	mu    sync.RWMutex
	wg    sync.WaitGroup
}

func (w *KVWorker) Start() {
	for i := 0; i < w.wc; i++ {
		w.wg.Add(1)
		go w.runWorker()
	}
}

func (w *KVWorker) Submit(req *Request) (*Response, error) {
	task := toTask(req)
	task.done = make(chan struct{})

	w.taskQ.WaitAndPush(task)
	<-task.done

	if task.err != nil {
		return &Response{}, task.err
	}

	return toResponse(task), nil
}

func (w *KVWorker) Stop() error {
	w.wg.Wait()
	return nil
}

func (w *KVWorker) runWorker() {
	for true {
		task := w.taskQ.WaitAndPop()

		// TODO: use operation prioritisation strategy
		if task.tType == TaskTypeGet {
			w.performGet(task)
		} else if task.tType == TaskTypeSet {
			w.performSet(task)
		} else if task.tType == TaskTypeDelete {
			w.performDelete(task)
		}

		close(task.done)
	}
}

func (w *KVWorker) performGet(task *Task) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	task.val, task.err = w.kvs.Get(task.key)
}

func (w *KVWorker) performSet(task *Task) {
	w.mu.Lock()
	defer w.mu.Unlock()

	task.err = w.kvs.Set(task.key, task.val)
}

func (w *KVWorker) performDelete(task *Task) {
	w.mu.Lock()
	defer w.mu.Unlock()

	task.err = w.kvs.Delete(task.key)
}
