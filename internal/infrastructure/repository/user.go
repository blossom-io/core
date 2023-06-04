package repository

import (
	"context"
	"fmt"
	"time"

	"core/internal/entity"
)

type User interface {
	AddUser(ctx context.Context, user entity.User) error
	AddToken(ctx context.Context, token entity.Token) error
}

var _ User = (*repository)(nil)

// AddUser adds user to db.
func (r *repository) AddUser(ctx context.Context, user entity.User) error {
	q, a, err := r.DB.Sq.Insert("person").
		Columns("twitch_id, twitch_username, updated_at").
		Values(user.TwitchID, user.TwitchUsername, time.Now()).
		Suffix("ON CONFLICT (twitch_id) DO UPDATE SET twitch_username = EXCLUDED.twitch_username, updated_at = EXCLUDED.updated_at").
		ToSql()
	if err != nil {
		return fmt.Errorf("AddUser - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return fmt.Errorf("AddUser - r.prepare: %w", err)
	}

	result, err := stmt.ExecContext(ctx, a...)
	if err != nil {
		return fmt.Errorf("AddUser - Exec: %w", err)
	}

	fmt.Println(result.RowsAffected())

	return nil
}

// AddToken adds user's token to db.
func (r *repository) AddToken(ctx context.Context, token entity.Token) error {
	q, a, err := r.DB.Sq.Insert("token").
		Columns("twitch_id, twitch_auth_code, twitch_bearer, twitch_bearer_expires_at, twitch_refresh_token, invite_key").
		Values(token.TwitchID, token.TwitchAuthCode, token.TwitchBearer, token.TwitchBearerExpiresAt, token.TwitchRefreshToken, token.InviteKey).
		Suffix(`ON CONFLICT (twitch_id) DO UPDATE SET twitch_auth_code = EXCLUDED.twitch_auth_code,
			twitch_bearer = EXCLUDED.twitch_bearer,
			twitch_bearer_expires_at = EXCLUDED.twitch_bearer_expires_at,
			twitch_refresh_token = EXCLUDED.twitch_refresh_token,
			invite_key = EXCLUDED.invite_key
			`).ToSql()
	if err != nil {
		return fmt.Errorf("AddToken - r.Sq: %w", err)
	}

	stmt, err := r.prepare(ctx, q)
	if err != nil {
		return fmt.Errorf("AddToken - r.prepare: %w", err)
	}

	result, err := stmt.ExecContext(ctx, a...)
	if err != nil {
		return fmt.Errorf("AddToken - Exec: %w", err)
	}

	fmt.Println(result.RowsAffected())

	return nil
}
