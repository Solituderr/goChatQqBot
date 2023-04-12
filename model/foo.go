package model


type Setting struct {
	Summary string `json:"summary"`
	Content string `json:"content"`
}

// 两个表
type UserInfo struct {
	UserId  string `json:"user_id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type UserAccess struct {
	UserId string `json:"user_id"`
	Acc    bool   `json:"acc"`
}

//消息上报结构体

type GroupMsg struct {
	MessageType string `json:"message_type" form:"message_type"`
	MessageId   int32  `json:"message_id" form:"message_id"`
	UserId      int64  `json:"user_id" form:"user_id"`
	Message     string `json:"message" form:"message"`
	RawMessage  string `json:"raw_message" form:"raw_message"`
	GroupId     int64  `json:"group_id" form:"group_id"`
}

type PersonMsg struct {
	MessageType string `json:"message_type" form:"message_type"`
	MessageId   int32  `json:"message_id" form:"message_id"`
	UserId      int64  `json:"user_id" form:"user_id"`
	Message     string `json:"message" form:"message"`
	RawMessage  string `json:"raw_message" form:"raw_message"`
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
}

type AddFriMsg struct {
	RequestType string `json:"request_type" form:"request_type"`
	Comment     string `json:"comment" form:"comment"`
	Flag        string `json:"flag" form:"flag"`
	UserId      int64  `json:"user_id" form:"user_id"`
}
