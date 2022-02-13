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

func (l *logmw) PostComment(threadID ThreadID, body string, author UserID, parentComment CommentID) (newId CommentID, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "PostComment", "threadID", threadID, "body", body, "author",
			author, "parentComment", parentComment, "generated", newId, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.PostComment(threadID, body, author, parentComment)
}

func (l *logmw) GetComment(threadID ThreadID, id CommentID) (comm Comment, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "GetComment", "threadID", threadID, "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.GetComment(threadID, id)
}
