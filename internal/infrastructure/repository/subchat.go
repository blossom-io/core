package repository

import (
	"context"
	"fmt"
)

type Subchater interface {
	IsSubchatExistsAndActive(ctx context.Context, ownerTwitchID int64) (bool, error)
}

var _ Subchater = (*repository)(nil)

// IsSubchatExistsAndActive checks if subchat exists and active.
func (r *repository) IsSubchatExistsAndActive(ctx context.Context, ownerTwitchID int64) (exists bool, err error) {
	q, a, err := r.DB.Sq.Select("1").Prefix("SELECT EXISTS (").From("subchat").Where("twitch_id = $1", ownerTwitchID).Suffix("AND NOT disabled)").ToSql()
	if err != nil {
		return false, fmt.Errorf("IsSubchatExistsAndActive - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return false, fmt.Errorf("IsSubchatExistsAndActive - r.prepare: %w", err)
	}

	err = stmt.QueryRowContext(ctx, a...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("IsSubchatExistsAndActive - Exec: %w", err)
	}

	return exists, nil
}
