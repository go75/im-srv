package conversion

import (
	"im-srv/ws-srv/dao"
	"im-srv/ws-srv/model"
	"im-srv/ws-srv/pb"
)

func Group(group *model.Group) *pb.Group {
	return &pb.Group{
		ID: group.ID,
		Name: group.Name,
		MasterId: group.MasterId,
		Introduce: group.Introduce,
	}
}

func Groups(groups []*model.Group) *pb.Groups {
	res := &pb.Groups{
		Groups: make([]*pb.Group, len(groups)),
	}
	for i, g := range groups {
		res.Groups[i] = Group(g)
	}
	return res
}

func GroupMessage(msg *model.GroupMessage) *pb.GroupMessage {
	user := model.User{
		ID: msg.SenderId,
	}
	dao.GetUserNameById(&user)

	return &pb.GroupMessage{
		ID: msg.ID,
		GroupId: msg.GroupId,
		SenderId: msg.SenderId,
		SenderName: user.Name,
		Type: uint32(msg.Type),
		Content: msg.Content,
	}
}

func GroupMessages(msgs []*model.GroupMessage) *pb.GroupMessages {
	res := &pb.GroupMessages{
		Msgs: make([]*pb.GroupMessage, len(msgs)),
	}
	for i, m := range msgs {
		res.Msgs[i] = GroupMessage(m)
	}
	return res
}

func GroupMessage2(msg *pb.GroupMessage) *model.GroupMessage {
	return &model.GroupMessage{
		ID: msg.ID,
		GroupId: msg.GroupId,
		SenderId: msg.SenderId,
		Type: uint8(msg.Type),
		Content: msg.Content,
	}
}

func GroupMessages2(msgs *pb.GroupMessages) []*model.GroupMessage {
	res := make([]*model.GroupMessage, len(msgs.Msgs))
	for i, m := range msgs.Msgs {
		res[i] = GroupMessage2(m)
	}
	return res
}

func NewGroupMessage(msg *model.AddGroupMessage) *pb.NewGroupMessage {
	group := &model.Group{
		ID: msg.GroupId,
	}
	dao.QueryGroupNameByGroupId(group)
	user := &model.User{
		ID: msg.SenderId,
	}
	dao.GetUserNameById(user)
	return &pb.NewGroupMessage{
		GroupId: msg.GroupId,
		Username: user.Name,
		GroupName: group.Name,
	}
}

func NewGroupMessages(msgs []*model.AddGroupMessage) *pb.NewGroupMessages {
	res := &pb.NewGroupMessages{
		Msgs: make([]*pb.NewGroupMessage, len(msgs)),
	}

	for i, m := range msgs {
		res.Msgs[i] = NewGroupMessage(m)
	}

	return res
}

func AddGroup(group *model.Group) *pb.AddGroupRes {
	return &pb.AddGroupRes{
		MasterId: group.MasterId,
		GroupId: group.ID,
		GroupName: group.Name,
	}
}
