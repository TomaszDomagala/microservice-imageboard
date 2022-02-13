package media

import (
	"github.com/go-kit/kit/log"
	"io"
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

func (l logmw) PostMedia(name string, reader io.Reader) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "PostMedia", "name", name, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.PostMedia(name, reader)
}

func (l logmw) GetMedia(name string) (b []byte, err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "GetMedia", "name", name, "size", len(b), "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.GetMedia(name)
}

func (l logmw) DeleteMedia(name string) (err error) {
	defer func(begin time.Time) {
		l.logger.Log("method", "DeleteMedia", "name", name, "took", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.DeleteMedia(name)
}
