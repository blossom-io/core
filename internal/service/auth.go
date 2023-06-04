package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"core/internal/entity"
	"core/internal/infrastructure/repository"
	"core/pkg/key"
	"core/pkg/logger"
	"core/pkg/twitch"
)

var (
	ErrIndexerNotFound          = fmt.Errorf("service - FindIndexerAndShowIDByURL - unknown indexer")
	ErrURLIsEmpty               = fmt.Errorf("service - AddShow - URL is empty")
	ErrUserNotSubscribed        = errors.New("not_subscribed")
	ErrSubchatNotFoundNotActive = errors.New("subchat not found or not active")
)

type auth struct {
	log    logger.Logger
	repo   repository.Repository
	twitch twitch.Twitcher
}

type Auther interface {
	AuthTwitchSubchat(ctx context.Context, authCode, state string) (inviteKey string, err error)
	Test(ctx context.Context) error
}

// NewShow injects repository and returns show service.
func NewAuth(log logger.Logger, repo repository.Repository, twitch twitch.Twitcher) Auther {
	return &auth{
		log:    log,
		repo:   repo,
		twitch: twitch,
	}
}

func (au *auth) Test(ctx context.Context) error {
	// user := entity.User{
	// 	TwitchID:       68247475,
	// 	TwitchUsername: "evan_64",
	// }

	token := entity.Token{
		TwitchID:              68247475,
		TwitchAuthCode:        "authcode",
		TwitchBearer:          "bearer",
		TwitchBearerExpiresAt: time.Now().Add(time.Hour * 1),
		TwitchRefreshToken:    "refresh",
	}

	funcs := []func(context.Context) error{
		// func(txCtx context.Context) error {
		// 	return au.repo.AddUser(txCtx, user)
		// },
		func(txCtx context.Context) error {
			return au.repo.AddToken(txCtx, token)
		},
	}

	err := au.repo.InTX(ctx, funcs)
	if err != nil {
		return err
	}

	return nil
}

// AuthTwitchSubchat
func (au *auth) AuthTwitchSubchat(ctx context.Context, authCode, state string) (inviteKey string, err error) {
	userClient, err := twitch.NewUserClient(ctx, authCode)
	if err != nil {
		return "", err
	}

	user, token, err := userClient.UserInfo(ctx)
	if err != nil {
		return "", err
	}

	broadcasterID, err := au.twitch.GetUserID(ctx, state)
	if err != nil {
		return "", err
	}

	token.TwitchAuthCode = authCode
	token.InviteKey = key.WrapKey(broadcasterID)

	ownerTwitchID, _ := strconv.ParseInt(broadcasterID, 10, 64)

	if isSubchatActive, err := au.repo.IsSubchatExistsAndActive(ctx, ownerTwitchID); !isSubchatActive || err != nil {
		return "", ErrSubchatNotFoundNotActive
	}

	isSubscribed, err := userClient.CheckUserSubscription(ctx, fmt.Sprint(user.TwitchID), broadcasterID)
	if err != nil {
		return "", ErrUserNotSubscribed
	}

	if !isSubscribed {
		au.log.Error("svc - AuthTwitchSubchat: user is not subscribed", "twitch_username", user.TwitchUsername, "twitch_user_id", user.TwitchID, "broadcaster_id", broadcasterID)
		return "", ErrUserNotSubscribed
	}

	funcs := []func(context.Context) error{
		func(txCtx context.Context) error {
			return au.repo.AddUser(txCtx, user)
		},
		func(txCtx context.Context) error {
			return au.repo.AddToken(txCtx, token)
		},
	}

	err = au.repo.InTX(ctx, funcs)
	if err != nil {
		return "", err
	}

	return token.InviteKey, nil
}
