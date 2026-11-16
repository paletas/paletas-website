package chat

type IChatOperation interface {
	Perform(user *ChatUser) error
}

type BaseOperation struct {
	Operation string `json:"operation"`
}
