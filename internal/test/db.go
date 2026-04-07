package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/errors"
)

func (w *World) registerDatabase() {
	pg.Register()

	w.Append(di.Hook{
		OnStart: func(_ context.Context) error {
			if w.PG == nil || !w.PG.IsEnabled() {
				return nil
			}

			return w.openDatabase()
		},
		OnStop: func(_ context.Context) error {
			if w.DB == nil {
				return nil
			}

			return errors.Join(w.DB.Destroy()...)
		},
	})
}

func (w *World) openDatabase() error {
	db, err := pg.Connect(FS, w.PG)
	if err != nil {
		return err
	}

	w.DB = db

	return nil
}
