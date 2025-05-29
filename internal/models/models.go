package models

import "time"

// everywhere is SQLx style db tags

type UserPosition struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type User struct {
	UID        string     `db:"uid"`
	Name       string     `db:"name"`
	Surname    string     `db:"surname"`
	PositionID int        `db:"position_id"`
	Email      string     `db:"email"`
	Password   string     `db:"password"`
	Phone      string     `db:"phone"`
	HireDate   *time.Time `db:"hire_date"`
	LastOnline *time.Time `db:"last_online"`
	PfpURL     *string    `db:"pfp_url"`
	IsAdmin    bool       `db:"is_admin"`
	ReAuth     bool       `db:"re_auth"`
}

type IdeaCategory struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type IdeaStatus struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Idea struct {
	IdeaUID      string    `db:"idea_uid"`
	Name         string    `db:"name"`
	Text         string    `db:"text"`
	Author       string    `db:"author"`
	CreationDate time.Time `db:"creation_date"`
	StatusID     int       `db:"status_id"`
	CategoryID   int       `db:"category_id"`
	LikeCount    int       `db:"like_count"`
	DislikeCount int       `db:"dislike_count"`
}

type IdeaComment struct {
	Idea           Idea
	CommentReplies []CommentReply
}

type Comment struct {
	CommentUID  string    `db:"comment_uid"`
	IdeaUID     string    `db:"idea_uid"`
	AuthorID    string    `db:"author_uid"`
	CommentText string    `db:"comment_text"`
	Timestamp   time.Time `db:"timestamp"`
}

type CommentReply struct {
	Comment Comment
	Replies []Reply
}

type Reply struct {
	ReplyUID   string    `db:"reply_uid"`
	CommentUID string    `db:"comment_uid"`
	AuthorID   string    `db:"author_uid"`
	Timestamp  time.Time `db:"timestamp"`
	ReplyText  string    `db:"reply_text"`
}

type BrowseHistory struct {
	VisitorID string `db:"visitor_uid"`
	IdeaID    string `db:"idea_uid"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	PositionID int    `json:"positionID"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
}

type InsertIdeaRequest struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	Author   string `json:"author"`
	Status   int    `json:"status"`
	Category int    `json:"category"`
}

type InsertCommentRequest struct {
	IdeaUID     string `json:"ideaUID"`
	CommentText string `json:"commentText"`
}

type InsertReplyRequest struct {
	CommentUID string `json:"commentUID"`
	ReplyText  string `json:"replyText"`
}
