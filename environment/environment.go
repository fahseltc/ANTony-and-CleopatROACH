package environment

import (
	"log/slog"
	"os"
	"time"
)

type Env struct {
	*slog.Logger
	*Config
}

var ENV *Env

func NewEnv() *Env {
	env := &Env{
		Logger: setupLogger(),
		Config: newConfig(),
	}
	ENV = env
	return env
}

func setupLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.Attr{
						Key:   slog.TimeKey,
						Value: slog.Int64Value(t.Unix()),
					}
				}
			}
			return a
		},
	})

	return slog.New(handler)
}
