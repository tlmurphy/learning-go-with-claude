package goroutines

import "testing"

func TestDemoBasicGoroutine(t *testing.T) {
	DemoBasicGoroutine()
}

func TestDemoMultipleGoroutines(t *testing.T) {
	DemoMultipleGoroutines(5)
}

func TestDemoPingPong(t *testing.T) {
	DemoPingPong(3)
}

func TestDemoBufferedChannel(t *testing.T) {
	DemoBufferedChannel()
}

func TestDemoPipeline(t *testing.T) {
	DemoPipeline(5)
}

func TestDemoRangeOverChannel(t *testing.T) {
	DemoRangeOverChannel()
}

func TestDemoCheckChannelClosed(t *testing.T) {
	DemoCheckChannelClosed()
}

func TestDemoConcurrentFetch(t *testing.T) {
	DemoConcurrentFetch([]string{"https://example.com", "https://go.dev", "https://pkg.go.dev"})
}

func TestDemoWaitGroup(t *testing.T) {
	DemoWaitGroup()
}

func TestDemoSemaphore(t *testing.T) {
	DemoSemaphore(10, 3)
}
