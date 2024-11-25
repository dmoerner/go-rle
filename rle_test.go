package rle

import "testing"

func Test_rle(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"singleton", "a", "a1"},
		{"numbers", "444445555", "4554"},
		{"mixed test", "aaaaa5663bbb", "a5516231b3"},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := RunLengthEncoding(d.input)
			if result != d.expected {
				t.Errorf("rle: expected %s, got %s", d.expected, result)
			}
		})
	}
}
