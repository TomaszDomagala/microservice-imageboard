package board

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

func (l *logmw) CreateThread(boardID BoardID, owner UserID) (id ThreadID, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "CreateThread", "board", boardID, "threadID", id, "owner",
			owner, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.CreateThread(boardID, owner)
}

func (l *logmw) DeleteThread(boardID BoardID, threadID ThreadID) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "CreateThread", "board", boardID, "threadID", threadID,
			"took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.DeleteThread(boardID, threadID)
}

func (l *logmw) GetThreads(boardID BoardID) (threads []ThreadID, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "CreateThread", "board", boardID, "returned ids", threads,
			"took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.GetThreads(boardID)
}
