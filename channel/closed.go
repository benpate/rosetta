package channel

// Closed returns TRUE if this channel is closed, and FALSE otherwise.
// It does not read from the channel, and is just a simple wrapper around a select statement.
func Closed[T any](channel <-chan T) bool {
	select {
	case <-channel:
		return true
	default:
		return false
	}
}
