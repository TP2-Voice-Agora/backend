package pgmodels

import "time"

type UserPosition struct {
	ID   int    `db:"id"` // SQLx style tag
	Name string `db:"name"`
}

type User struct {
	UID        string     `db:"uid"` // UUID
	Name       string     `db:"name"`
	Surname    string     `db:"surname"`
	PositionID int        `db:"position_id"`
	Email      string     `db:"email"`
	Phone      string     `db:"phone"`
	HireDate   *time.Time `db:"hire_date"`   // DATE
	LastOnline *time.Time `db:"last_online"` // TIMESTAMP
	PfpURL     *string    `db:"pfp_url"`     // TEXT
	IsAdmin    bool       `db:"is_admin"`    // BOOL
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
	IdeaUID      string    `db:"idea_uid"`      // UUID
	Author       string    `db:"author"`        // UUID
	CreationDate time.Time `db:"creation_date"` // consider using time.Time with TIMESTAMP
	StatusID     int       `db:"status_id"`
	CategoryID   int       `db:"category_id"`
	LikeCount    int       `db:"like_count"`
	DislikeCount int       `db:"dislike_count"`
}

type Comment struct {
	CommentID   int       `db:"comment_id"`
	IdeaUID     string    `db:"idea_uid"`     // UUID
	AuthorID    string    `db:"author_id"`    // UUID
	Timestamp   time.Time `db:"timestamp"`    // TIMESTAMP
	CommentText string    `db:"comment_text"` // TEXT
}

type Reply struct {
	CommentID int       `db:"comment_id"`
	AuthorID  string    `db:"author_id"`  // UUID
	Timestamp time.Time `db:"timestamp"`  // TIMESTAMP
	ReplyText string    `db:"reply_text"` // TEXT
}
