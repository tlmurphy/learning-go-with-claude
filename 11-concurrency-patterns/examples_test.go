package concurrency

import "testing"

func TestDemoWaitGroup(t *testing.T) {
	DemoWaitGroup()
}

func TestDemoMutex(t *testing.T) {
	DemoMutex(4, 100)
}

func TestDemoSelect(t *testing.T) {
	DemoSelect()
}

func TestDemoContextCancellation(t *testing.T) {
	DemoContextCancellation()
}

func TestDemoContextTimeout(t *testing.T) {
	DemoContextTimeout()
}

func TestDemoWorkerPool(t *testing.T) {
	jobs := []Job{
		{ID: 1, Input: 2},
		{ID: 2, Input: 3},
		{ID: 3, Input: 4},
		{ID: 4, Input: 5},
	}
	DemoWorkerPool(2, jobs)
}

func TestDemoFanOutFanIn(t *testing.T) {
	DemoFanOutFanIn([]int{1, 2, 3, 4, 5}, 3)
}

func TestDemoPipeline(t *testing.T) {
	DemoPipeline([]int{1, 2, 3, 4, 5, 6, 7, 8})
}
