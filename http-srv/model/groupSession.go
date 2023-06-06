package model

// 群聊会话
type GroupSession struct {
	// 群聊id
	GroupId uint32
	// 用户id
	UserId uint32
}

func (s *GroupSession) TableName() string {
	return "group_session"
}