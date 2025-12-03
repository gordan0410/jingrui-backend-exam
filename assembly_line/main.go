package main

import (
	"context"
	"fmt"
	"time"
)

const (
	ExceedTime            = 10 * time.Second
	ProcessTimeMSNum      = 5000
	TotalWorker           = 5
	TotalItem             = 3
	TotalPieceForEachItem = 10
)

func main() {
	taskQueue := make(chan ItemI, TotalItem*TotalPieceForEachItem)
	w := NewEmployee(taskQueue, ExceedTime, TotalWorker)
	ctx, cancelFun := context.WithCancel(context.Background())
	defer cancelFun()
	w.Work(ctx)

	m := NewItemManger([]string{"item1", "item2", "item3"}, taskQueue)
	m.GenItems()
	m.Queue(cancelFun)
	<-ctx.Done()
	fmt.Printf("All work done!, total time spent: %d\n", m.GetTotalTimeSpend().Milliseconds())
	time.Sleep(1 * time.Second)
}
