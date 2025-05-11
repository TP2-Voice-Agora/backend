package postgres

import (
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ictisagora/backend/internal/models"
)

// PostgresRepository - implements Repository interface for PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// use ConnectDB before query, and CloseConnectDB when query is finished
func (pg *PostgresRepository) ConnectDB(sourceURL string) error {
	var err error
	// sourceURL := "postgres://username:password@localhost:5432/database_name"
	pg.db, err = sqlx.Connect("pgx", sourceURL)

	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresRepository) CloseConnectDB() error {
	err := pg.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresRepository) InsertUser(user models.User) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Insert("users").
		Columns(
			"uid", "name", "surname", "position_id", "email",
			"phone", "hire_date", "last_online", "pfp_url", "is_admin",
		).Values(
		user.UID, user.Name, user.Surname, user.PositionID, user.Email, user.Phone, user.HireDate, user.LastOnline, user.PfpURL, user.IsAdmin,
	).ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(q, args...)

	return err
}

// SelectUserByEmail selects user by email from table users, returns User struct
func (pg *PostgresRepository) SelectUserByEmail(email string) (models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("users").Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return models.User{}, err
	}
	var user models.User

	err = pg.db.QueryRowx(q, args...).Scan(&user)

	return user, err
}

// SelectPositions selects all job positions into slice, then returns it
func (pg *PostgresRepository) SelectPositions() ([]models.UserPosition, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("positions").ToSql()
	if err != nil {
		return nil, err
	}
	var positions []models.UserPosition

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var position models.UserPosition
		err = rows.StructScan(&position)
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	return positions, nil
}

func (pg *PostgresRepository) InsertIdea(idea models.Idea) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, _, err := psql.Insert("ideas").
		Columns(
			"idea_uid", "author", "creation_date", "status_id",
			"category_id", "like_count", "dislike_count",
		).
		Values(
			sq.Expr(":idea_uid"), sq.Expr(":author"), sq.Expr(":creation_date"),
			sq.Expr(":status_id"), sq.Expr(":category_id"), sq.Expr(":like_count"), sq.Expr(":dislike_count"),
		).ToSql()
	if err != nil {
		return err
	}
	_, err = pg.db.NamedExec(q, idea)
	return err
}

func (pg *PostgresRepository) SelectIdeas() ([]models.Idea, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q, args, err := psql.Select("*").From("ideas").ToSql()
	if err != nil {
		return nil, err
	}
	var ideas []models.Idea

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var idea models.Idea
		err = rows.StructScan(&idea)
		if err != nil {
			return nil, err
		}
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

func (pg *PostgresRepository) SelectUserIdeas(uid string, limit int) ([]models.Idea, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.
		Select("*").
		From("ideas").
		Where(sq.Eq{"author": uid}).
		OrderBy("created_at DESC")

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	}

	q, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var ideas []models.Idea
	rows, err := pg.db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idea models.Idea
		if err := rows.StructScan(&idea); err != nil {
			return nil, err
		}
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

func (pg *PostgresRepository) InsertIdeaComment(comment models.Comment) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	// expect potential problems with inserting time.Time into timestamp
	q, _, err := psql.Insert("comments").
		Columns(
			"idea_uid", "author", "creation_date", "status_id",
			"category_id", "like_count", "dislike_count",
		).
		Values(
			sq.Expr(":idea_uid"), sq.Expr(":author"), sq.Expr(":creation_date"),
			sq.Expr(":status_id"), sq.Expr(":category_id"), sq.Expr(":like_count"), sq.Expr(":dislike_count"),
		).ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.NamedExec(q, comment)
	return err
}

func (pg *PostgresRepository) InsertCommentReply(reply models.Reply) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// expect potential problems with inserting time.Time into timestamp
	q, _, err := psql.Insert("replies").
		Columns(
			"comment_id", "author_id", "timestamp", "reply_text",
		).
		Values(
			sq.Expr(":comment_id"), sq.Expr(":author_id"),
			sq.Expr(":timestamp"), sq.Expr(":reply_text"),
		).ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.NamedExec(q, reply)
	return err
}

func (pg *PostgresRepository) SelectIdeaComments(uid string) ([]models.Comment, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("comments").Where(sq.Eq{"idea_uid": uid}).ToSql()
	if err != nil {
		return nil, err
	}
	var comments []models.Comment

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment models.Comment
		err = rows.StructScan(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (pg *PostgresRepository) SelectCommentReplies(id int) ([]models.Reply, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("replies").Where(sq.Eq{"comment_id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var replies []models.Reply

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var reply models.Reply
		err = rows.StructScan(&reply)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}

	return replies, nil
}

func (pg *PostgresRepository) SelectIdeaCategories() ([]models.IdeaCategory, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("idea_categories").ToSql()
	if err != nil {
		return nil, err
	}
	var ideaCategories []models.IdeaCategory

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ideaCategory models.IdeaCategory
		err = rows.StructScan(&ideaCategory)
		if err != nil {
			return nil, err
		}
		ideaCategories = append(ideaCategories, ideaCategory)
	}

	return ideaCategories, nil
}

func (pg *PostgresRepository) SelectIdeaStatuses() ([]models.IdeaStatus, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("idea_statuses").ToSql()
	if err != nil {
		return nil, err
	}
	var ideaStatuses []models.IdeaStatus

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ideaStatus models.IdeaStatus
		err = rows.StructScan(&ideaStatus)
		if err != nil {
			return nil, err
		}
		ideaStatuses = append(ideaStatuses, ideaStatus)
	}

	return ideaStatuses, nil
}
