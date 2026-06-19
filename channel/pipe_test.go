package channel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {

	input := FromSlice([]int{1, 2, 3, 4})
	output := make(chan int)

	Pipe(input, output)

	require.Equal(t, []int{1, 2, 3, 4}, Slice(output))
}

func TestPipeWithCancel_Complete(t *testing.T) {

	input := FromSlice([]int{1, 2, 3})
	output := make(chan int)
	done := make(chan struct{})

	PipeWithCancel(input, output, done)

	require.Equal(t, []int{1, 2, 3}, Slice(output))
}

func TestPipeWithCancel_Cancelled(t *testing.T) {

	input := make(chan int) // never sends
	output := make(chan int)
	done := make(chan struct{})

	PipeWithCancel(input, output, done)

	// Closing "done" stops the pipe and closes the output channel
	close(done)

	// Draining the output channel must complete without blocking
	require.Equal(t, []int{}, Slice(output))
}
