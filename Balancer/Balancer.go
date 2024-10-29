package Balancer

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

type Balancer struct {
	ThreadList   []*Thread             // 存放所有线路
	WorkList     []*Thread             // 存放所有工作线路
	DeadList     []*Thread             // 用于存放所有的死亡线路
	SleepList    []*SleepThread        // 用于存放所有休眠的线路
	WordOrderMap map[string]*WorkOrder // 通过uuid获取到工作账单
}

// WorkOrder 配合UUIDMap去统计工作
type WorkOrder struct {
	Th         *Thread
	StartedAt  time.Time
	StatusCode int
}

// Thread 维护一条线路的一些参数
type Thread struct {
	ThreadID        uint
	FailureRate     float64
	AverageDuration int64
	URL             string
}

// SleepThread 休眠的线路, 并记录唤醒时间
type SleepThread struct {
	Th       *Thread
	WakeupAt time.Time
}

// AddThread 用于注册线路,并维护一个自增主键
// return: 线路id
func (b *Balancer) AddThread(url string) uint {
	th := Thread{
		ThreadID: uint(len(b.ThreadList)), // 自增主键
		URL:      url,
	}
	b.ThreadList = append(b.ThreadList, &th)
	return th.ThreadID
}

// KillThread 用于把一个线路加入死亡队列
func (b *Balancer) KillThread(threadID uint) {
	for i := 0; i < len(b.WorkList); i++ {
		if b.WorkList[i].ThreadID == threadID {
			// 1. 添加到死亡队列
			b.DeadList = append(b.DeadList, b.WorkList[i])
			// 2. 踢出工作队列
			b.WorkList = append(b.WorkList[:i], b.WorkList[i+1:]...)
			break
		}
	}
}

type testFunc func(string)

// WakeupThread 用于从死亡队列里唤醒一个线路,
// 需要传入一个测试函数, 用于测试线路
// 需要起一个goroutine去唤醒
func (b *Balancer) WakeupThread(threadID uint, f testFunc, baseline int64) {
	th := b.ThreadList[threadID]

	// 测试, 统计耗时
	startedAt := time.Now()
	f(th.URL)
	duration := time.Now().Sub(startedAt).Milliseconds()

	// 不合格直接鬼
	if duration > baseline {
		return
	}

	for i := 0; i < len(b.DeadList); i++ {
		if b.WorkList[i].ThreadID != threadID {
			// 1. 添加到工作队列
			b.WorkList = append(b.WorkList, b.DeadList[i])
			// 2. 踢出死亡队列
			b.DeadList = append(b.DeadList[:i], b.DeadList[i+1:]...)
			break
		}
	}

}

// Run 用于管理Balancer的运行 执行一些定时任务
func (b *Balancer) Run() {
	// TODO
}

// GetThread 用于获取线路
func (b *Balancer) GetThread() (Thread, string) {
	var th Thread
	// TODO 获取线路的策略
	OrderID := uuid.New().String()
	b.WordOrderMap[OrderID] = &WorkOrder{
		Th:        &th,
		StartedAt: time.Now(),
	}
	return th, OrderID
}

// ReportData 用于维护需要汇报的参数
type ReportData struct {
	ThreadID   uint // 用于表示是哪个线路
	EndedAt    time.Time
	StatusCode int // 请求状态码
}

// Report 用于汇报工作, 并写入redis
func (b *Balancer) Report(redisDB *redis.Client, ctx context.Context, orderID string, data ReportData) error {
	workOrder := b.WordOrderMap[orderID]
	endedAt := data.EndedAt
	th := workOrder.Th

	duration := endedAt.Sub(workOrder.StartedAt).Milliseconds()

	err := redisDB.HMSet(ctx, "user", map[string]interface{}{
		"ThreadID":   th.ThreadID,
		"Duration":   duration,
		"StatusCode": data.StatusCode,
	}).Err()
	return err
}

// Summary 用于总结工作,并更新 工作队列 死亡队列 (定时任务,在Balancer.Run里起一个GoRoutine去跑)
func (b *Balancer) Summary() {
	// TODO
}
