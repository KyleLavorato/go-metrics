package metrics

import "testing"

func createTestReg() Registry {
	r := NewRegistry()
	NewRegisteredCounter("bar", r).Inc(randomInt64())
	return r
}

func BenchmarkSlice(b *testing.B) {
	s := NewSlice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Append(createTestReg())
	}
}

func TestSliceClear(t *testing.T) {
	s := NewSlice()
	s.Append(createTestReg())

	s.Clear()
	if len(s.GetAll()) != 0 {
		t.Errorf("Slice elements not cleared")
	}
}

func TestSliceAppend(t *testing.T) {
	s := NewSlice()
	s.Append(createTestReg())
	if len(s.GetAll()) == 0 {
		t.Errorf("Slice elements not appended")
	}
}

func TestSliceZero(t *testing.T) {
	s := NewSlice()
	if len(s.GetAll()) != 0 {
		t.Errorf("Did not generate empty slice")
	}
}
