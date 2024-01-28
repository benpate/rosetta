package channel

import "testing"

func TestBeep(t *testing.T) {

	in := testChannel()

	for value := range Beep(in) {
		t.Log(value)
	}
}
