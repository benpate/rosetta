package channel

// Done is just sugar around a standard channel type
// that identifies that a process has completed
type Done chan struct{}
