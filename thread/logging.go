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

func (l *logmw) CreateThread(ip, board, body string, hasMedia bool) (threadID ThreadID, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "CreateThread",
			"ip", ip,
			"board", board,
			"body", body,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.CreateThread(ip, board, body, hasMedia)
}

func (l *logmw) DeleteThread(threadID ThreadID) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "DeleteThread", "threadID", threadID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.DeleteThread(threadID)
}

func (l *logmw) PostComment(ip string, threadID ThreadID, body string, parentComment CommentID, hasMedia bool) (newId CommentID, err error) {
	defer func(begin time.Time) {
		l.logger.Log(
			"method", "PostComment",
			"ip", ip,
			"threadID", threadID,
			"body", body,
			"parentComment", parentComment,
			"newId", newId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return l.next.PostComment(ip, threadID, body, parentComment, hasMedia)
}

func (l *logmw) GetComment(threadID ThreadID, id CommentID) (comm Comment, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "GetComment", "threadID", threadID, "Id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.GetComment(threadID, id)
}

func (l *logmw) GetChildren(threadID ThreadID, id CommentID) (comm []Comment, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "GetCommentChildren", "threadID", threadID, "Id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.GetChildren(threadID, id)
}
