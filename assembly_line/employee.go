package main

import (
	"context"
	"fmt"
	"time"
)

type EmployeeI interface {
	Work(ctx context.Context)
}
type Employee struct {
	TaskQueue     chan ItemI    `json:"-"`
	ExceedTime    time.Duration `json:"exceed_time"`
	TotalEmployee int           `json:"total_employee"`
}

func NewEmployee(taskQueue chan ItemI, exceedTime time.Duration, totalEmployee int) EmployeeI {
	return &Employee{
		TaskQueue:     taskQueue,
		ExceedTime:    exceedTime,
		TotalEmployee: totalEmployee,
	}
}

func (e *Employee) Work(ctx context.Context) {
	for i := 0; i < e.TotalEmployee; i++ {
		go e.work(ctx)
	}
}

func (e *Employee) work(ctx context.Context) {
	itemHandleCount := 0
	for {
		select {
		case item := <-e.TaskQueue:
			timeoutCtx, cancelFun := context.WithTimeout(ctx, e.ExceedTime)
			if item.Process(timeoutCtx) {
				itemHandleCount++
			}
			cancelFun()
		case <-ctx.Done():
			fmt.Printf("employee out of work, total handled items: %d\n", itemHandleCount)
			return
		}
	}
}
