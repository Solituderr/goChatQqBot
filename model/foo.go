package model

import "time"

type Setting struct {
	Summary string `json:"summary"`
	Content string `json:"content"`
}

// 天天宝用户
type UserInfo struct {
	UserId   string `json:"userId"`
	SignTime int    `json:"signInTime" gorm:"default:1"`
	GiftTime int    `json:"rcvGiftTime" gorm:"default:3"`
	LoveRate int    `json:"loveRate" gorm:"default:0"`
	Name     string `json:"name" gorm:"default:name"`
}

// 周边推送
type PushMsg struct {
	MsgTime time.Time `json:"msgTime"`
	Message string    `json:"message"`
	GroupId string    `json:"groupId"`
}

//消息上报结构体

type GroupMsg struct {
	MessageType string `json:"message_type" form:"message_type"`
	MessageId   int32  `json:"message_id" form:"message_id"`
	UserId      int64  `json:"user_id" form:"user_id"`
	Message     string `json:"message" form:"message"`
	RawMessage  string `json:"raw_message" form:"raw_message"`
	GroupId     int64  `json:"group_id" form:"group_id"`
	Sender      struct {
		UserId   int64  `json:"user_id"`
		NickName string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int32  `json:"age"`
	} `json:"sender"`
}

type PersonMsg struct {
	MessageType string `json:"message_type" form:"message_type"`
	MessageId   int32  `json:"message_id" form:"message_id"`
	UserId      int64  `json:"user_id" form:"user_id"`
	Message     string `json:"message" form:"message"`
	RawMessage  string `json:"raw_message" form:"raw_message"`
}

type MutePerson struct {
	GroupId int64 `json:"groupId"`
	UserId  int64 `json:"userId"`
}

type TypeMsg struct {
	MessageType string `json:"message_type" form:"message_type"`
}

type CommonReq struct {
	Time     int64  `json:"time" form:"time"`
	SelfId   int64  `json:"self_id" form:"self_id"`
	PostType string `json:"post_type" form:"post_type"`
}

// 封装通用消息
type CommonMsg struct {
	UserId  string `json:"user_id" form:"user_id"`
	GroupId string `json:"group_id" form:"group_id"`
	Message string `json:"message" form:"message"`
	Sender  struct {
		UserId   int64  `json:"user_id"`
		NickName string `json:"nickname"`
		Sex      string `json:"sex"`
		Age      int32  `json:"age"`
	} `json:"sender"`
}

type AddFriMsg struct {
	RequestType string `json:"request_type" form:"request_type"`
	Comment     string `json:"comment" form:"comment"`
	Flag        string `json:"flag" form:"flag"`
	UserId      int64  `json:"user_id" form:"user_id"`
}

type EnterGroup struct {
	NoticeType string `json:"notice_type"`
	GroupId    int64  `json:"group_id"`
	UserId     int64  `json:"user_id"`
}
