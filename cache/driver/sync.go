package driver

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
)

type syncDriver struct {
	items      map[string]syncEntry
	maxEntries int
	mu         sync.Mutex
}

type syncEntry struct {
	expiresAt time.Time
	value     string
}

func newSyncDriver(maxEntries int) *syncDriver {
	return &syncDriver{
		items:      make(map[string]syncEntry),
		maxEntries: max(maxEntries, 1),
	}
}

func (d *syncDriver) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.deleteEntry(key)

	return nil
}

func (d *syncDriver) Fetch(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	entry, ok := d.items[key]
	if !ok {
		return "", ErrMissing
	}

	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		d.deleteEntry(key)

		return "", ErrExpired
	}

	return entry.value, nil
}

func (d *syncDriver) Flush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	clear(d.items)

	return nil
}

func (d *syncDriver) Save(ctx context.Context, key, value string, lifetime time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.deleteExpired()
	if _, ok := d.items[key]; !ok {
		d.evictIfFull()
	}

	d.items[key] = syncEntry{
		expiresAt: expiresAt(lifetime),
		value:     value,
	}

	return nil
}

func (d *syncDriver) deleteExpired() {
	now := time.Now()
	for key, entry := range d.items {
		if !entry.expiresAt.IsZero() && now.After(entry.expiresAt) {
			d.deleteEntry(key)
		}
	}
}

func (d *syncDriver) evictIfFull() {
	for len(d.items) >= d.maxEntries {
		for key := range d.items {
			d.deleteEntry(key)

			break
		}
	}
}

func (d *syncDriver) deleteEntry(key string) {
	delete(d.items, key)
}

func expiresAt(lifetime time.Duration) time.Time {
	if lifetime <= 0 {
		return time.Time{}
	}

	return time.Now().Add(lifetime.Duration())
}
