package codesmells

import "testing"

func TestDemoInterfacePollution(t *testing.T) {
	DemoInterfacePollution()
}

func TestDemoGoroutineLeak(t *testing.T) {
	DemoGoroutineLeak()
}

func TestDemoErrorFormatting(t *testing.T) {
	DemoErrorFormatting()
}

func TestDemoContextMisuse(t *testing.T) {
	DemoContextMisuse()
}

func TestDemoSets(t *testing.T) {
	DemoSets()
}

func TestDemoMutexVsChannel(t *testing.T) {
	DemoMutexVsChannel()
}
