package grpcmod

import "testing"

func TestDemoServiceInterface(t *testing.T) {
	DemoServiceInterface()
}

func TestDemoErrorHandling(t *testing.T) {
	DemoErrorHandling()
}

// DemoLoggingInterceptor returns UnaryServerInterceptor — skipped.

func TestDemoMetadata(t *testing.T) {
	DemoMetadata()
}

func TestDemoDeadlines(t *testing.T) {
	DemoDeadlines()
}

func TestDemoHealthCheck(t *testing.T) {
	DemoHealthCheck()
}

func TestDemoRetryLogic(t *testing.T) {
	DemoRetryLogic()
}
