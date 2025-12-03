package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ItemI interface {
	// Process 這是一個耗時操作
	Process(context.Context) bool
	TryLock() bool
	Unlock()
	GetKey() string
}

type BaseItem struct {
	Mu            sync.Mutex    `json:"-"`
	FinishedQueue chan struct{} `json:"-"`
	Name          string        `json:"name"`
	Series        int           `json:"series"`
	ProcessTime   time.Duration `json:"process_time"`
	Done          bool          `json:"done"`

	// manager
	MangerMu       *sync.Mutex    `json:"-"`
	RandM          *sync.Map      `json:"-"`
	TotalTimeSpent *time.Duration `json:"total_time_spent"`
}

func NewBaseItem(name string, series int, processTime time.Duration, randM *sync.Map, totalTimeSpent *time.Duration, managerMutex *sync.Mutex, finishedQueue chan struct{}) ItemI {
	return &BaseItem{
		Mu:             sync.Mutex{},
		Done:           false,
		RandM:          randM,
		MangerMu:       managerMutex,
		TotalTimeSpent: totalTimeSpent,
		Name:           name,
		Series:         series,
		ProcessTime:    processTime,
		FinishedQueue:  finishedQueue,
	}
}

func (i *BaseItem) Process(ctx context.Context) bool {
	if i.TryLock() {
		defer i.Unlock()
		if i.Done {
			return false
		}
		doneChan := make(chan struct{})
		var spendTime time.Duration
		go func() {
			defer close(doneChan)
			spendTime = i.doSomething()
		}()
		select {
		case <-ctx.Done():
			// 逾時處理
			i.RandM.Store(i.GetKey(), i)
		case <-doneChan:
			// 正常結束
			i.Done = true
			i.MangerMu.Lock()
			*i.TotalTimeSpent += spendTime
			i.MangerMu.Unlock()
			i.FinishedQueue <- struct{}{}
			return true
		}
	}
	return false
}

func (i *BaseItem) doSomething() (spendTime time.Duration) {
	now := time.Now()
	fmt.Printf("Start processing item..., nowItem: %s, nowSerie: %d, nowTime: %s\n", i.Name, i.Series, now.String())
	time.Sleep(i.ProcessTime)

	//// 模擬失敗
	//if rand.IntN(2) == 1 {
	//	fmt.Printf("item processing failed!, nowItem: %s, nowSerie: %d\n", i.Name, i.Series)
	//	time.Sleep(ExceedTime + 1*time.Second)
	//	return
	//}

	since := time.Since(now)
	fmt.Printf("Finishing item..., nowItem: %s, nowSerie: %d, cost: %d ms\n", i.Name, i.Series, since.Milliseconds())
	spendTime += since

	return
}

func (i *BaseItem) TryLock() bool {
	return i.Mu.TryLock()
}

func (i *BaseItem) Unlock() {
	i.Mu.Unlock()
}

func (i *BaseItem) GetKey() string {
	return fmt.Sprintf("%s_%d", i.Name, i.Series)
}
