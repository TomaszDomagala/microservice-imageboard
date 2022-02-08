package main

import (
	"errors"
	"sync"
)

type UserID = string
type CommentID = int

// ThreadService serves information about threads.
type ThreadService interface {

	// I don't know if context is necessary here, for now I deleted it.
	PostComment(body string, author UserID, parentComment CommentID) (CommentID, error)
	GetComment(id CommentID) (Comment, error)
}

type inmemThreadService struct {
	mtx sync.RWMutex
	m   map[CommentID]Comment
}

func NewInmemService() ThreadService {
	return &inmemThreadService{
		m: map[CommentID]Comment{},
	}
}

func (c Comment) addChild(id CommentID) Comment {
	c.Children = append(c.Children, id)
	return c
}

type Comment struct {
	Body     string      `json:"body,omitempty"`
	Author   UserID      `json:"author,omitempty"`
	Id       CommentID   `json:"id,omitempty"`
	Children []CommentID `json:"children,omitempty"`
}

var (
	ErrNotFound = errors.New("not found")
)

// https://stackoverflow.com/questions/64631848/how-to-create-an-autoincrement-id-field
type autoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id         int
}

func (a *autoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}

var commentIDGenerator autoInc // global instance

func (s *inmemThreadService) PostComment(body string, author UserID, parentID CommentID) (CommentID, error) {
	newID := commentIDGenerator.ID()

	cmt := Comment{Body: body, Author: author, Id: newID, Children: []CommentID{}}

	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[newID] = cmt

	s.m[parentID] = s.m[parentID].addChild(newID)

	return newID, nil
}

func (s *inmemThreadService) GetComment(id CommentID) (Comment, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	comment, ok := s.m[id]
	if !ok {
		return Comment{}, ErrNotFound
	}
	return comment, nil
}
