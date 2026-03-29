package errorhandling

import "testing"

func TestDemoCreatingErrors(t *testing.T) {
	DemoCreatingErrors()
}

func TestDemoSentinelErrors(t *testing.T) {
	DemoSentinelErrors()
}

func TestDemoCustomErrors(t *testing.T) {
	DemoCustomErrors()
}

func TestDemoErrorWrapping(t *testing.T) {
	DemoErrorWrapping()
}

func TestDemoRecover(t *testing.T) {
	DemoRecover()
}
