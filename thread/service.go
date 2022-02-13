package thread

import (
	"database/sql"
	"errors"
)

// Service serves information about threads.
type Service interface {
	PostComment(threadID ThreadID, body string, author UserID, parentComment CommentID) (CommentID, error)
	GetComment(threadID ThreadID, id CommentID) (Comment, error)
	CreateThread(threadID ThreadID, body string, author UserID) error
	DeleteThread(id ThreadID) error
}

type ServiceMiddleware func(Service) Service

var (
	ErrNotFound = errors.New("not found")
)

const ROOT_COMMENT_ID = 0

type PostgresService struct {
	db *sql.DB
}

func (p PostgresService) createThread(threadID ThreadID, body string, author UserID) error {
	db := p.db
	stmtCreateRow := `INSERT INTO threads (threadID, nextID) VALUES $1, $2`
	_, err := db.Exec(stmtCreateRow, threadID, ROOT_COMMENT_ID+1)
	if err != nil {
		return err
	}

	stmtInsertComment := `
		INSERT INTO comments (threadID, commentID, author, body)
		VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(stmtInsertComment, threadID, ROOT_COMMENT_ID, author, body)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresService) deleteThread(id ThreadID) error {
	db := p.db
	stmtDeleteComments := `DELETE FROM comments WHERE threadID = $1`
	_, err := db.Exec(stmtDeleteComments, id)
	if err != nil {
		return err
	}
	stmtDeleteThreads := `DELETE FROM threads WHERE threadID = $1`
	_, err = db.Exec(stmtDeleteThreads, id)
	return err
}

type DBComment struct {
	Body   string `db:"body"`
	Author string `db:"author"`
	Id     int    `db:"commentID"`
}

func (p PostgresService) PostComment(threadID ThreadID, body string, author UserID, parentComment CommentID) (CommentID, error) {
	db := p.db

	stmtUpdateNextID := `UPDATE threads SET nextID = nextID + 1
	WHERE threadID = $1 
	RETURNING nextID`

	newID := 0
	err := db.QueryRow(stmtUpdateNextID, threadID).Scan(&newID)
	if err != nil {
		return 0, err
	}

	stmtInsertComment := `
		INSERT INTO comments (threadID, commentID, author, parentComment, body)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(stmtInsertComment, threadID, newID, author, parentComment, body)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (p PostgresService) GetComment(threadID ThreadID, commentID CommentID) (Comment, error) {
	db := p.db

	comment := DBComment{}
	stmtGetCommentData := `SELECT (body, author, CommentID) FROM comments WHERE threadID = $1 AND CommentID = $2`
	err := db.QueryRow(stmtGetCommentData, threadID, commentID).Scan(&comment)
	if err != nil {
		return Comment{}, ErrNotFound
	}

	// It's also possible to store list of children in comment row.
	// For now I use the easier version.
	stmtGetChildren := `SELECT CommentID FROM comments WHERE threadID = $1 AND parentComment = $2`
	var ids []int
	rows, err := db.Query(stmtGetChildren, threadID, commentID)
	if err != nil {
		return Comment{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return Comment{}, err
		}
		ids = append(ids, id)
	}
	return Comment{comment.Body, comment.Author, comment.Id, ids}, nil
}

func NewPostgresService(psqlInfo string) *PostgresService {
	db, err := ConnectToDB(psqlInfo)
	if err != nil {
		panic("Unable to connect to database")
	}
	return &PostgresService{
		db: db,
	}
}
