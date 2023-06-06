package model

type AddUserMessage struct {
	// 发送方id
	SenderId uint32
	// 接收方id
	ReceiverId uint32
}

func (m *AddUserMessage) TableName() string {
	return "add_user_message"
}