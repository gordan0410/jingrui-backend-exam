package main

import (
	"math/rand/v2"
	"sync"
	"time"
)

type ItemManagerI interface {
	GenItems()
	Queue(cancelFun func())
	GetTotalTimeSpend() time.Duration
}
type ItemManager struct {
	Mu             *sync.Mutex    `json:"-"`
	TaskQueue      chan ItemI     `json:"-"`
	FinishedQueue  chan struct{}  `json:"-"`
	RandM          *sync.Map      `json:"-"`
	TotalTimeSpent *time.Duration `json:"total_time_spent"`
	ItemTypes      []string       `json:"item_types"`
}

func NewItemManger(itemTypes []string, taskQueue chan ItemI) ItemManagerI {
	m := &sync.Map{}
	mu := &sync.Mutex{}
	totalTimeSpent := time.Duration(0)

	return &ItemManager{
		Mu:             mu,
		TotalTimeSpent: &totalTimeSpent,
		ItemTypes:      itemTypes,
		RandM:          m,
		TaskQueue:      taskQueue,
		FinishedQueue:  make(chan struct{}, TotalItem*TotalPieceForEachItem),
	}
}

func (im *ItemManager) GenItems() {
	items := make([]ItemI, 0, TotalPieceForEachItem*TotalItem)
	for _, itemTypeName := range im.ItemTypes {
		r := rand.IntN(ProcessTimeMSNum-1) + 1
		// 測試用
		//r = 1000
		for i := 0; i < TotalPieceForEachItem; i++ {
			items = append(items, NewBaseItem(itemTypeName, i, time.Duration(r)*time.Millisecond, im.RandM, im.TotalTimeSpent, im.Mu, im.FinishedQueue))
		}
	}

	for _, item := range items {
		im.RandM.Store(item.GetKey(), item)
	}
}

func (im *ItemManager) Queue(cancelFun func()) {
	for {
		im.RandM.Range(func(key, value interface{}) bool {
			im.TaskQueue <- value.(ItemI)
			im.RandM.Delete(key)
			return true
		})
		if len(im.TaskQueue) == 0 && len(im.FinishedQueue) == TotalItem*TotalPieceForEachItem {
			cancelFun()
			break
		}
	}

}

func (im *ItemManager) GetTotalTimeSpend() time.Duration {
	return *im.TotalTimeSpent
}
