package worker

func BuildChannelMessageReceiverForwarder(channel chan Message) func(message Message) {
	return func(message Message) {
		channel <- message
	}
}
