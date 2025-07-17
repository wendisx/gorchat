package model

// entity for message table
type Message struct {
	// Topic [string|int64] `json:"tupic"` // 消息kafka主题
	MessageId int64  `json:"messageId"` // 消息id
	Sender    int64  `json:"sender"`    // 消息发送者
	Receiver  int64  `json:"receiver"`  // 消息接收者
	Type      string `json:"type"`      // 消息类型
	Text      string `json:"text"`      // 文本消息映射
	Content   []byte `json:"content"`   // 二进制消息映射
	Status    string `json:"status"`    // 消息状态
	Deleted   int    `json:"deleted"`   // 消息逻辑删除
}
