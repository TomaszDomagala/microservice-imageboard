package board

import (
	"database/sql"
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

type ThreadID int
type UserID string
type BoardID string

// Service serves information about threads.
type Service interface {
	CreateThread(boardID BoardID, owner UserID) (ThreadID, error)
	DeleteThread(boardID BoardID, threadID ThreadID) error
	GetThreads(boardID BoardID) ([]ThreadID, error)
}

type ServiceMiddleware func(Service) Service

var (
	ErrNotFound = errors.New("not found")
)

type PostgresService struct {
	db    *sql.DB
	cache *cache.Cache
}

func (p *PostgresService) CreateThread(boardID BoardID, owner UserID) (ThreadID, error) {
	db := p.db
	stmtCreateThread := `INSERT INTO threads (boardID, owner) VALUES ($1, $2) RETURNING threadID`
	var threadID ThreadID

	err := db.QueryRow(stmtCreateThread, boardID, owner).Scan(&threadID)
	if err != nil {
		return 0, err
	}

	return threadID, nil
}

func (p *PostgresService) DeleteThread(boardID BoardID, threadID ThreadID) error {
	db := p.db
	stmtDeleteThread := `DELETE FROM threads WHERE boardID = $1 AND threadID = $2`
	_, err := db.Exec(stmtDeleteThread, boardID, threadID)
	if err != nil {
		return err
	}
	return err
}

type DBComment struct {
	Body   string `db:"body"`
	Author string `db:"author"`
	Id     int    `db:"commentID"`
}

func (p *PostgresService) GetThreads(boardID BoardID) ([]ThreadID, error) {
	cacheKey := string(boardID)
	if val, found := p.cache.Get(cacheKey); found {
		return val.([]ThreadID), nil
	}

	db := p.db

	stmtGetThreads := `SELECT threadID FROM threads WHERE boardID = $1`

	var ids []ThreadID
	rows, err := db.Query(stmtGetThreads, boardID)
	if err != nil {
		return nil, ErrNotFound
	}

	defer rows.Close()

	for rows.Next() {
		var id ThreadID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	p.cache.Set(cacheKey, ids, cache.DefaultExpiration)
	return ids, nil
}

func NewPostgresService(psqlInfo string) Service {
	db, err := ConnectToDB(psqlInfo)
	if err != nil {
		panic("Unable to connect to database")
	}
	return &PostgresService{
		db: db, cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}
