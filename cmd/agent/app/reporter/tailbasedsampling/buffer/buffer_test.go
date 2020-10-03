package buffer

import (
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"sync"
	"testing"
	"time"
)

func TestNewWindowBuffer(t *testing.T) {
	buffer := NewWindowBuffer(10)
	g := 100
	go func() {
		for i := 0; i < g; i++ {
			buffer.Put("test"+fmt.Sprintf("_%d", i), &model.Span{
				TraceID: model.TraceID{
					High: uint64(i),
					Low:  uint64(i),
				},
				SpanID:               0,
				OperationName:        "",
				References:           nil,
				Flags:                0,
				StartTime:            time.Time{},
				Duration:             0,
				Tags:                 nil,
				Logs:                 nil,
				Process:              nil,
				ProcessID:            "",
				Warnings:             nil,
			})
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 10)
	for j := 0; j < g; j++ {
		if _, ok := buffer.Get("test" + fmt.Sprintf("_%d", j)); ok {
			fmt.Println(j, "ok")
		} else {
			fmt.Println(j, "not ok")
		}
		time.Sleep(time.Second)
	}
}

func TestWindowBuffer_Put(t *testing.T) {
	buffer := NewWindowBuffer(60)
	g := 10000
	w := sync.WaitGroup{}
	w.Add(g)

	startTime := time.Now()
	for i := 0; i < g; i++ {
		j := i
		go func() {
			buffer.Put("test"+fmt.Sprintf("_%d", j), &model.Span{
				TraceID: model.TraceID{
					High: uint64(j),
					Low:  uint64(j),
				},
				SpanID:               0,
				OperationName:        "",
				References:           nil,
				Flags:                0,
				StartTime:            time.Time{},
				Duration:             0,
				Tags:                 nil,
				Logs:                 nil,
				Process:              nil,
				ProcessID:            "",
				Warnings:             nil,
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     nil,
				XXX_sizecache:        0,
			})
		}()
		w.Done()
	}

	w.Wait()
	t.Log(time.Since(startTime).String())
}

func TestWindowBuffer_Get(t *testing.T) {
	buffer := NewWindowBuffer(60)
	g := 10000
	w := sync.WaitGroup{}
	w.Add(g)

	startTime := time.Now()
	for i := 0; i < g; i++ {
		j := i
		go func() {
			buffer.Put("test"+fmt.Sprintf("_%d", j), &model.Span{
				TraceID: model.TraceID{
					High: uint64(j),
					Low:  uint64(j),
				},
				SpanID:               0,
				OperationName:        "",
				References:           nil,
				Flags:                0,
				StartTime:            time.Time{},
				Duration:             0,
				Tags:                 nil,
				Logs:                 nil,
				Process:              nil,
				ProcessID:            "",
				Warnings:             nil,
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     nil,
				XXX_sizecache:        0,
			})
		}()
		w.Done()
	}

	w.Wait()
	t.Log(time.Since(startTime).String())

	w = sync.WaitGroup{}
	w.Add(g)
	startTime = time.Now()
	for i := 0; i < g; i++ {
		j := i
		go func() {
			buffer.Get("test" + fmt.Sprintf("_%d", j))
			w.Done()
		}()
	}
	w.Wait()

	t.Log(time.Since(startTime).String())
}