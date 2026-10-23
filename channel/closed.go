package channel

// Closed returns TRUE if a receive from the channel succeeds at once: the channel is closed, or
// it holds a value, which this call consumes and discards. It returns FALSE otherwise.
func Closed[T any](channel <-chan T) bool {
	select {
	case <-channel:
		return true
	default:
		return false
	}
}
