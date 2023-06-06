package conversion

import (
	"im-srv/ws-srv/model"
	"im-srv/ws-srv/pb"
)

func UserMessage(msg *model.UserMessage) *pb.UserMessage {
	return &pb.UserMessage{
		ID: msg.ID,
		SenderId: msg.SenderId,
		ReceiverId: msg.ReceiverId,
		Type: uint32(msg.Type),
		Content: msg.Content,
	}
}

func UserMessages(msgs []*model.UserMessage) *pb.UserMessages {
	res := &pb.UserMessages{
		Msgs: make([]*pb.UserMessage, len(msgs)),
	}
	
	for i, msg := range msgs {
		res.Msgs[i] = UserMessage(msg)
	}

	return res
}

func UserMessage2(msg *pb.UserMessage) *model.UserMessage {
	return &model.UserMessage{
		ID: msg.ID,
		SenderId: msg.SenderId,
		ReceiverId: msg.ReceiverId,
		Type: uint8(msg.Type),
		Content: msg.Content,
	}
}

func UserMessages2(msgs *pb.UserMessages) []*model.UserMessage {
	res := make([]*model.UserMessage, len(msgs.Msgs))

	for i, m := range msgs.Msgs {
		res[i] = UserMessage2(m)
	}

	return res
}

func User(user *model.User) *pb.User {
	return &pb.User{
		ID: user.ID,
		Name: user.Name,
		Introduce: user.Introduce,
	}
}

func Users(users []*model.User) *pb.Users {
	res := &pb.Users{
		Users: make([]*pb.User, len(users)),
	}

	for i, u := range users {
		res.Users[i] = User(u)
	}

	return res
}

