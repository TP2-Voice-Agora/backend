package models

import "time"

// everywhere is SQLx style db tags

type UserPosition struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
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
}

type IdeaCategory struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type IdeaStatus struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Idea struct {
	IdeaUID      string    `db:"idea_uid"`
	Author       string    `db:"author"`
	CreationDate time.Time `db:"creation_date"`
	StatusID     int       `db:"status_id"`
	CategoryID   int       `db:"category_id"`
	LikeCount    int       `db:"like_count"`
	DislikeCount int       `db:"dislike_count"`
}

type Comment struct {
	CommentID   int       `db:"comment_id"`
	IdeaUID     string    `db:"idea_uid"`
	AuthorID    string    `db:"author_id"`
	Timestamp   time.Time `db:"timestamp"`
	CommentText string    `db:"comment_text"`
}

type Reply struct {
	CommentID int       `db:"comment_id"`
	AuthorID  string    `db:"author_id"`
	Timestamp time.Time `db:"timestamp"`
	ReplyText string    `db:"reply_text"`
}

type BrowseHistory struct {
	VisitorID string `db:"visitor_id"`
	IdeaID    string `db:"idea_id"`
}
