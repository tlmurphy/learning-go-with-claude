package jsonserialization

import "testing"

func TestDemoMarshalUnmarshal(t *testing.T) {
	DemoMarshalUnmarshal()
}

func TestDemoStructTags(t *testing.T) {
	DemoStructTags()
}

func TestDemoNullVsAbsent(t *testing.T) {
	DemoNullVsAbsent()
}

func TestDemoDisallowUnknownFields(t *testing.T) {
	DemoDisallowUnknownFields()
}

func TestDemoCustomMarshal(t *testing.T) {
	DemoCustomMarshal()
}

func TestDemoRawMessage(t *testing.T) {
	DemoRawMessage()
}

func TestDemoStreamingJSON(t *testing.T) {
	DemoStreamingJSON()
}

func TestDemoTimeSerialization(t *testing.T) {
	DemoTimeSerialization()
}

func TestDemoJSONNumber(t *testing.T) {
	DemoJSONNumber()
}

func TestDemoMapJSON(t *testing.T) {
	DemoMapJSON()
}
