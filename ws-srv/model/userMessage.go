package model

type UserMessage struct {
	ID uint64 `gorm:"primary_key"`
	// 发送方id
	SenderId uint32
	// 接收方id
	ReceiverId uint32
	// 消息类型 文字为0, 图片为1, 音频为2
	Type uint8
	// 消息内容
	Content string
}

func (m *UserMessage) TableName() string {
	return "user_message"
}