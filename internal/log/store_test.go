package log

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	write = []byte("hello world")
	width = uint64(len(write)) + lenWidth
)

func TestStoreAppendRead(t *testing.T) {
	f, err := ioutil.TempFile("", "store_append_read_test")
	assertNoError(t, err)
	defer os.Remove(f.Name())

	_, err = newStore(f)
	assertNoError(t, err)
}

func assertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()

	if got != want {
		var format string

		switch any(want).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			format = "got %d, want %d"
		case float32:
			format = "got %.2f, want %.2f"
		case float64:
			format = "got %g, want %g"
		case string:
			format = "got %q, want %q"
		default:
			format = "got %v, want %v"
		}

		t.Errorf(format, got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()

	if got != nil {
		t.Fatal("got an error but didn't want one: ", got)
	}
}
