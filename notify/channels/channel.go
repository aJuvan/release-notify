package channels

type Channel interface {
	SendMessages(messages []string) error
}
