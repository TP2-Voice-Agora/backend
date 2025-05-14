package postgres

import (
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ictisagora/backend/internal/models"
	"log/slog"
)

// PostgresRepository - implements Repository interface for PostgreSQL
type PostgresRepository struct {
	db  *sqlx.DB
	log slog.Logger
}

// use ConnectDB before query, and CloseConnectDB when a query is finished
func (pg *PostgresRepository) ConnectDB(sourceURL string, log slog.Logger) error {
	var err error
	// sourceURL := "postgres://username:password@localhost:5432/database_name"
	pg.db, err = sqlx.Connect("pgx", sourceURL)
	pg.log = log
	pg.log.Debug("connecting to database")

	if err != nil {
		log.Error("failed to connect to database" + err.Error())
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
			"uid", "password", "name", "surname", "position_id", "email",
			"phone", "hire_date", "last_online", "pfp_url", "is_admin",
		).Values(
		user.UID, user.Password, user.Name, user.Surname, user.PositionID, user.Email, user.Phone, user.HireDate, user.LastOnline, user.PfpURL, user.IsAdmin,
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

	err = pg.db.QueryRowx(q, args...).StructScan(&user)

	return user, err
}

func (pg *PostgresRepository) SelectUserByUID(uid string) (models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("users").Where(sq.Eq{"uid": uid}).ToSql()
	if err != nil {
		return models.User{}, err
	}
	var user models.User

	err = pg.db.QueryRowx(q, args...).StructScan(&user)

	return user, err
}

// SelectPositions selects all job positions into a slice, then returns it
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

// Probably final
func (pg *PostgresRepository) InsertIdea(idea models.Idea) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, _, err := psql.Insert("ideas").
		Columns(
			"idea_uid", "name", "text", "author", "status_id",
			"category_id",
		).
		Values(
			sq.Expr(":idea_uid"), sq.Expr(":name"), sq.Expr(":text"), sq.Expr(":author"), sq.Expr(":status_id"),
			sq.Expr(":category_id"),
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

func (pg *PostgresRepository) SelectIdeaByUID(uid string) (models.Idea, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q, args, err := psql.Select("*").From("ideas").Where(sq.Eq{"idea_uid": uid}).ToSql()
	if err != nil {
		return models.Idea{}, err
	}
	var idea models.Idea

	err = pg.db.QueryRowx(q, args...).StructScan(&idea)
	if err != nil {
		return models.Idea{}, err
	}
	return idea, nil
}

func (pg *PostgresRepository) InsertIdeaComment(comment models.Comment) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	// expect potential problems with inserting time.Time into timestamp
	q, _, err := psql.Insert("comments").
		Columns(
			"comment_uid", "idea_uid",
			"author_uid", "comment_text",
		).
		Values(
			sq.Expr(":comment_uid"), sq.Expr(":idea_uid"),
			sq.Expr(":author_uid"), sq.Expr(":comment_text"),
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
			"reply_uid", "comment_uid", "author_uid", "reply_text",
		).
		Values(
			sq.Expr(":reply_uid"), sq.Expr(":comment_uid"),
			sq.Expr(":author_uid"), sq.Expr(":reply_text"),
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

func (pg *PostgresRepository) SelectCommentReplies(uid string) ([]models.Reply, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Select("*").From("replies").Where(sq.Eq{"comment_uid": uid}).ToSql()
	if err != nil {
		return []models.Reply{}, err
	}
	var replies []models.Reply

	rows, err := pg.db.Queryx(q, args...)
	defer rows.Close()
	if err != nil {
		return []models.Reply{}, err
	}

	for rows.Next() {
		var reply models.Reply
		err = rows.StructScan(&reply)
		if err != nil {
			return []models.Reply{}, err
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

func (pg *PostgresRepository) UpdateUserPfpURL(uid string, url string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q, args, err := psql.Update("users").
		Set("pfp_url", url).
		Where(sq.Eq{"uid": uid}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = pg.db.Exec(q, args...)
	return err
}
