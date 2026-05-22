package driver

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	sync "github.com/alexfalkowski/go-sync"
)

type syncDriver struct {
	items sync.Map[string, syncEntry]
}

type syncEntry struct {
	expiresAt time.Time
	value     string
}

func (d *syncDriver) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.items.Delete(key)

	return nil
}

func (d *syncDriver) Fetch(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	value, ok := d.items.Load(key)
	if !ok {
		return "", ErrMissing
	}

	if !value.expiresAt.IsZero() && time.Now().After(value.expiresAt) {
		d.items.Delete(key)

		return "", ErrExpired
	}

	return value.value, nil
}

func (d *syncDriver) Flush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.items.Clear()

	return nil
}

func (d *syncDriver) Save(ctx context.Context, key, value string, lifetime time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.items.Store(key, syncEntry{
		expiresAt: expiresAt(lifetime),
		value:     value,
	})

	return nil
}

func expiresAt(lifetime time.Duration) time.Time {
	if lifetime <= 0 {
		return time.Time{}
	}

	return time.Now().Add(lifetime.Duration())
}
