package metrics

import (
	"encoding/json"
	"testing"
)

type JsonData struct {
	ValueOne int
	ValueTwo int
}

func getRawJson(t *testing.T) json.RawMessage {
	data := JsonData{
		ValueOne: 1,
		ValueTwo: 2,
	}
	raw, err := json.Marshal(&data)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func BenchmarkJson(b *testing.B) {
	js := NewJson()
	data := JsonData{
		ValueOne: 1,
		ValueTwo: 2,
	}
	raw, err := json.Marshal(&data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		js.Set(raw)
	}
}

func TestJsonClear(t *testing.T) {
	js := NewJson()
	js.Set(getRawJson(t))
	js.Clear()
	if Json := js.Json(); nil != Json {
		t.Errorf("js.Json() is not empty after clear")
	}
}

func TestJsonSet(t *testing.T) {
	js := NewJson()
	raw := getRawJson(t)
	js.Set(raw)
	if Json := js.Json(); string(Json) != string(raw) {
		t.Errorf("js.Json(): %s != %s", string(Json), string(raw))
	}
}

func TestGetJson(t *testing.T) {
	r := NewRegistry()
	raw := getRawJson(t)
	NewRegisteredJson("foo", r).Set(raw)
	if js := GetJson("foo", r); string(raw) != string(js.Json()) {
		t.Fatal(js)
	}
}
