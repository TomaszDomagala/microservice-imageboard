package thread

import (
	"github.com/go-kit/kit/log"
	"time"
)

func ServiceLoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return &logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	next   Service
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

