package thread

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Service serves information about threads.
type Service interface {
	PostComment(ip string, threadID ThreadID, body string, parentComment CommentID) (CommentID, error)
	GetComment(threadID ThreadID, id CommentID) (Comment, error)
	CreateThread(ip, board, body string) error
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

func identify(ip string) (string, error) {
	req := fmt.Sprintf("{ip: %s}", ip)
	resp, err := http.Post("http://identification/identify", "application/json", strings.NewReader(req))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("identification failed: %s", resp.Status)
	}
	defer resp.Body.Close()
	var idResp struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&idResp)
	if err != nil {
		return "", err
	}

	return idResp.ID, nil
}

func requestNewThread(board, author string) (ThreadID, error) {
	req := fmt.Sprintf("{boardID: %s, owner: %s}", board, author)
	resp, err := http.Post("http://board/createThread", "application/json", strings.NewReader(req))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to create thread: %s", resp.Status)
	}
	defer resp.Body.Close()
	var idResp struct {
		ID ThreadID `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&idResp)
	if err != nil {
		return 0, err
	}

	return idResp.ID, nil

}

func (p *PostgresService) CreateThread(ip, board, body string) error {
	author, err := identify(ip)
	if err != nil {
		return err
	}
	threadID, err := requestNewThread(board, author)
	if err != nil {
		return err
	}
	db := p.db
	stmtCreateRow := `INSERT INTO threads (threadID, nextID) VALUES ($1, $2)`
	_, err = db.Exec(stmtCreateRow, threadID, ROOT_COMMENT_ID+1)
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

func (p *PostgresService) DeleteThread(id ThreadID) error {
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

func (p *PostgresService) PostComment(ip string, threadID ThreadID, body string, parentComment CommentID) (CommentID, error) {
	db := p.db
	author, err := identify(ip)
	if err != nil {
		return 0, err
	}

	stmtUpdateNextID := `UPDATE threads SET nextID = nextID + 1
	WHERE threadID = $1 
	RETURNING nextID`

	newID := 0
	err = db.QueryRow(stmtUpdateNextID, threadID).Scan(&newID)
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

func (p *PostgresService) GetComment(threadID ThreadID, commentID CommentID) (Comment, error) {
	db := p.db

	comment := DBComment{}
	stmtGetCommentData := `SELECT (body, author, commentID) FROM comments WHERE threadID = $1 AND commentID = $2`
	err := db.QueryRow(stmtGetCommentData, threadID, commentID).Scan(&comment)
	if err != nil {
		return Comment{}, ErrNotFound
	}

	// It's also possible to store list of children in comment row.
	// For now I use the easier version.
	stmtGetChildren := `SELECT commentID FROM comments WHERE threadID = $1 AND parentComment = $2`
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

func NewPostgresService(psqlInfo string) Service {
	db, err := ConnectToDB(psqlInfo)
	if err != nil {
		panic("Unable to connect to database")
	}
	return &PostgresService{
		db: db,
	}
}
