package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

func threadServiceloggingMiddleware(logger log.Logger) ThreadServiceMiddleware {
	return func(next ThreadService) ThreadService {
		return &logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	next   ThreadService
}

func (l *logmw) PostComment(body string, author UserID, parentComment CommentID) (newId CommentID, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "PostComment", "body", body, "author",
			author, "parentComment", parentComment, "generated", newId, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.PostComment(body, author, parentComment)
}

func (l *logmw) GetComment(id CommentID) (comm Comment, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "GetComment", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.GetComment(id)
}
