package thread

import (
	"errors"
	"github.com/TomaszDomagala/microservice-imageboard/utils"
	"sync"
)

// Service serves information about threads.
type Service interface {
	PostComment(body string, author UserID, parentComment CommentID) (CommentID, error)
	GetComment(id CommentID) (Comment, error)
}

type ServiceMiddleware func(Service) Service

type InMemoryService struct {
	mtx sync.RWMutex
	m   map[CommentID]Comment
}

func NewInMemoryService() Service {
	return &InMemoryService{
		m: map[CommentID]Comment{},
	}
}

var (
	ErrNotFound = errors.New("not found")
)

var commentIDGenerator = utils.NewAutoInc(1)

func (s *InMemoryService) PostComment(body string, author UserID, parentID CommentID) (CommentID, error) {
	newID := commentIDGenerator.ID()

	cmt := Comment{Body: body, Author: author, Id: newID, Children: []CommentID{}}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[newID] = cmt

	s.m[parentID] = s.m[parentID].addChild(newID)

	return newID, nil
}

func (s *InMemoryService) GetComment(id CommentID) (Comment, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	comment, ok := s.m[id]
	if !ok {
		return Comment{}, ErrNotFound
	}
	return comment, nil
}
