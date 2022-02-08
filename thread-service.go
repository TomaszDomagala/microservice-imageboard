package main

import (
	"app/thread-service/utils"
	"errors"
	"sync"
)

// ThreadService serves information about threads.
type ThreadService interface {
	PostComment(body string, author UserID, parentComment CommentID) (CommentID, error)
	GetComment(id CommentID) (Comment, error)
}

type ThreadServiceMiddleware func(ThreadService) ThreadService

type InMemoryThreadService struct {
	mtx sync.RWMutex
	m   map[CommentID]Comment
}

func NewInMemoryService() ThreadService {
	return &InMemoryThreadService{
		m: map[CommentID]Comment{},
	}
}

var (
	ErrNotFound = errors.New("not found")
)

var commentIDGenerator = utils.NewAutoInc(1)

func (s *InMemoryThreadService) PostComment(body string, author UserID, parentID CommentID) (CommentID, error) {
	newID := commentIDGenerator.ID()

	cmt := Comment{Body: body, Author: author, Id: newID, Children: []CommentID{}}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[newID] = cmt

	s.m[parentID] = s.m[parentID].addChild(newID)

	return newID, nil
}

func (s *InMemoryThreadService) GetComment(id CommentID) (Comment, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	comment, ok := s.m[id]
	if !ok {
		return Comment{}, ErrNotFound
	}
	return comment, nil
}
