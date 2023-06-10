package trace

import "testing"

func TestNextSpanID(t *testing.T) {
	traceID := NewID()
	for i := 0; i < 100; i++ {
		spanID := NextSpanID(traceID)
		t.Logf("NextSpanID() = %s", spanID)
	}
}

func TestStartSpanID(t *testing.T) {
	traceID := NewID()
	for i := 0; i < 10; i++ {
		spanID := StartSpanID(traceID, "0")
		t.Logf("StartSpanID() = %s", spanID)
	}
}

func TestAll(t *testing.T) {
	for i := 0; i < 10; i++ {
		traceID := NewID()
		spanID := StartSpanID(traceID, "0")
		t.Logf("StartSpanID() = %s", spanID)
		for i := 0; i < 10; i++ {
			spanID = NextSpanID(traceID)
			t.Logf("NextSpanID() = %s", spanID)
		}
	}
}
