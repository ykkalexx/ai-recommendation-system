package models

type UserBehavior struct {
	UserID    string `json:"user_id"`
	ItemID    string `json:"item_id"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}