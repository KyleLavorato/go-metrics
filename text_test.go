package metrics

import (
	"math/rand"
	"testing"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func BenchmarkText(b *testing.B) {
	tx := NewText()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tx.Set(RandString(100))
	}
}

func TestTextClear(t *testing.T) {
	tx := NewText()
	tx.Set(RandString(100))
	tx.Clear()
	if text := tx.Text(); "" != text {
		t.Errorf("tx.Text() is not empty after clear")
	}
}

func TestTextSet(t *testing.T) {
	tx := NewText()
	str := RandString(100)
	tx.Set(str)
	if text := tx.Text(); text != str {
		t.Errorf("tx.Text(): %s != %s", text, str)
	}
}

func TestTextAppend(t *testing.T) {
	tx := NewText()
	str1 := RandString(100)
	str2 := RandString(25)
	tx.Set(str1)
	tx.Append(str2)
	if text := tx.Text(); str1+str2 != text {
		t.Errorf("tx.Text(): %s != %s", text, str1+str2)
	}
}

func TestGetText(t *testing.T) {
	r := NewRegistry()
	str := RandString(100)
	NewRegisteredText("foo", r).Set(str)
	if tx := GetText("foo", r); str != tx.Text() {
		t.Fatal(tx)
	}
}
