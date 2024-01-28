package channel

func testChannel() <-chan string {

	result := make(chan string)

	go func() {
		defer close(result)

		result <- "Hello"
		result <- "There"
		result <- "General"
		result <- "Kenobi"
	}()

	return result
}
