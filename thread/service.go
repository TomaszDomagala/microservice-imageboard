package thread

import (
	"database/sql"
	"errors"
)

// Service serves information about threads.
type Service interface {
	PostComment(threadID ThreadID, body string, author UserID, parentComment CommentID) (CommentID, error)
	GetComment(threadID ThreadID, id CommentID) (Comment, error)
}

type ServiceMiddleware func(Service) Service

//type InMemoryService struct {
//	mtx sync.RWMutex
//	m   map[CommentID]Comment
//}
//
//func NewInMemoryService() Service {
//	return &InMemoryService{
//		m: map[CommentID]Comment{},
//	}
//}
//

//
//var commentIDGenerator = utils.NewAutoInc(1)
//
//func (s *InMemoryService) PostComment(body string, author UserID, parentID CommentID, id ThreadID) (CommentID, error) {
//	newID := commentIDGenerator.ID()
//
//	cmt := Comment{Body: body, Author: author, CommentID: newID, Children: []CommentID{}}
//
//	s.mtx.Lock()
//	defer s.mtx.Unlock()
//	s.m[newID] = cmt
//
//	s.m[parentID] = s.m[parentID].addChild(newID)
//
//	return newID, nil
//}
//
//func (s *InMemoryService) GetComment(id CommentID) (Comment, error) {
//	s.mtx.RLock()
//	defer s.mtx.RUnlock()
//	comment, ok := s.m[id]
//	if !ok {
//		return Comment{}, ErrNotFound
//	}
//	return comment, nil
//}

var (
	ErrNotFound = errors.New("not found")
	ErrDB       = errors.New("database error")
)

type PostgresService struct {
	db *sql.DB
}

type DBComment struct {
	Body   string `db:"body"`
	Author string `db:"author"`
	id     int    `db:"CommentID"`
}

func (p PostgresService) PostComment(threadID ThreadID, body string, author UserID, parentComment CommentID) (CommentID, error) {
	db := p.db

	updateMaxComment := `UPDATE threads SET nextID = nextID + 1
	WHERE threadID = $1 
	RETURNING nextID`

	newID := 0
	err := db.QueryRow(updateMaxComment, threadID).Scan(&newID)
	if err != nil {
		return 0, err
	}

	sqlStatement := `
INSERT INTO comments (threadID, CommentID, author, parentComment, body)
VALUES ($1, $2, $3, $4, $5)x
RETURNING id`

	_, err = db.Exec(sqlStatement, threadID, newID, author, parentComment, body)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (p PostgresService) GetComment(threadID ThreadID, commentID CommentID) (Comment, error) {
	db := p.db

	var comment DBComment
	stmtGetCommentData := `SELECT (body, author, CommentID) 
FROM comments 
WHERE threadID = $1 AND CommentID = $2`
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
	return Comment{comment.Body, comment.Author, comment.id, ids}, nil
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
