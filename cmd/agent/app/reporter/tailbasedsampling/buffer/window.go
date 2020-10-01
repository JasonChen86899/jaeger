package buffer

import (
	"container/ring"
	"github.com/jaegertracing/jaeger/model"
	"sync"
	"time"
)

type WindowBuffer struct {
	rwLock     sync.RWMutex
	ringWindow *ring.Ring
	bufferMap  map[string][]*model.Span

	windowSize int
}

type keyList struct {
	creatTime int64
	keys      map[string]struct{}
}

func NewWindowBuffer(windowSize int) *WindowBuffer {
	buffer := &WindowBuffer{
		rwLock:     sync.RWMutex{},
		ringWindow: ring.New(windowSize),
		bufferMap:  make(map[string][]*model.Span, 1024),
		windowSize: windowSize,
	}

	go buffer.loopClean()
	return buffer
}

func (wb *WindowBuffer) Get(key string) ([]*model.Span, bool) {
	wb.rwLock.RLock()
	defer wb.rwLock.RUnlock()

	v, ok := wb.bufferMap[key]
	return v, ok
}

func (wb *WindowBuffer) Put(key string, value *model.Span) {
	wb.rwLock.Lock()
	defer wb.rwLock.Unlock()

	spanList := wb.bufferMap[key]
	spanList = append(spanList, value)
	wb.bufferMap[key] = spanList

	if v := wb.ringWindow.Value; v == nil {
		wb.ringWindow.Value = &keyList{
			creatTime: time.Now().Unix(),
			keys:      make(map[string]struct{}, 1024),
		}
	}
	createTime := wb.ringWindow.Value.(*keyList).creatTime
	keys := wb.ringWindow.Value.(*keyList).keys

	// check clean
	if createTime+int64(wb.windowSize) <= time.Now().Unix() {
		// clean
		for k, _ := range keys {
			delete(wb.bufferMap, k)
		}
		wb.ringWindow.Value = nil
	}
	keys[key] = struct{}{}
}

func (wb *WindowBuffer) loopClean() {
	for {
		wb.rwLock.Lock()
		wb.ringWindow = wb.ringWindow.Next()
		wb.rwLock.Unlock()
		time.Sleep(time.Second)
	}
}
