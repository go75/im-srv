package model

type AddGroupMessage struct {
	// 发送方id
	SenderId uint32
	// 被添加的群聊id
	GroupId uint32
}

func (m *AddGroupMessage) TableName() string {
	return "add_group_message"
}