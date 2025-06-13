package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Message(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func GetLogger(r *http.Request) *logrus.Entry {
	return r.Context().Value("logger").(*logrus.Entry)
}

func SetLogger(r *http.Request, id string) *http.Request {
	entry := logrus.WithField("id", id)
	ctx := context.WithValue(r.Context(), "logger", entry)
	return r.WithContext(ctx)
}
