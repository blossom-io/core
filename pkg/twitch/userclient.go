package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"core/internal/entity"

	"github.com/nicklaw5/helix/v2"
)

type UserTwitcher interface {
	GetRefreshToken() string
	GetUserAccessToken() string
	UserInfo(ctx context.Context) (user entity.User, token entity.Token, err error)
	CheckUserSubscription(ctx context.Context, userID, broadcasterID string) (isSubscribed bool, err error)
}

type userTwitch struct {
	UserClient   *helix.Client
	HTTPClient   *http.Client
	RefreshToken string
}

type UserInfoResponse struct {
	Aud               string `json:"aud"`
	Exp               int    `json:"exp"`
	Iat               int    `json:"iat"`
	Iss               string `json:"iss"`
	Sub               string `json:"sub"`
	Azp               string `json:"azp"`
	PreferredUsername string `json:"preferred_username"`
}

func (u *UserInfoResponse) ToUser() entity.User {
	id, _ := strconv.ParseInt(u.Sub, 10, 64)

	user := entity.User{
		TwitchID:       id,
		TwitchUsername: u.PreferredUsername,
	}

	return user
}

func (u *UserInfoResponse) ToToken() entity.Token {
	expiresAt := time.Unix(int64(u.Exp), 0)

	token := entity.Token{
		TwitchBearerExpiresAt: expiresAt,
	}

	return token
}

func NewUserClient(ctx context.Context, authCode string) (UserTwitcher, error) {
	HTTPClient := &http.Client{}

	UserClient, err := helix.NewClientWithContext(ctx, &helix.Options{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		UserAgent:    UserAgent,
		HTTPClient:   HTTPClient,
		RedirectURI:  RedirectURI,
	})
	if err != nil {
		return nil, err
	}

	resp, err := UserClient.RequestUserAccessToken(authCode)
	if err != nil {
		return nil, err
	}

	UserClient.SetUserAccessToken(resp.Data.AccessToken)

	return &userTwitch{
		HTTPClient:   HTTPClient,
		UserClient:   UserClient,
		RefreshToken: resp.Data.RefreshToken,
	}, nil
}

func (ut *userTwitch) CheckUserSubscription(ctx context.Context, userID, broadcasterID string) (isSubscribed bool, err error) {
	params := helix.UserSubscriptionsParams{
		BroadcasterID: broadcasterID,
		UserID:        userID,
	}

	res, err := ut.UserClient.CheckUserSubscription(&params)
	if err != nil {
		return false, err
	}

	if res.ResponseCommon.StatusCode != http.StatusOK {
		return false, nil
	}

	for _, sub := range res.Data.UserSubscriptions {
		if broadcasterID == sub.BroadcasterID {
			return true, nil
		}
	}

	return false, nil
}

func (ut *userTwitch) UserInfo(ctx context.Context) (user entity.User, token entity.Token, err error) {
	var userInfo UserInfoResponse

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, UserInfoURL, http.NoBody)
	if err != nil {
		return user, token, err
	}

	userToken := ut.UserClient.GetUserAccessToken()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userToken))

	res, err := ut.HTTPClient.Do(req)
	if err != nil {
		return user, token, err
	}

	if res.StatusCode != http.StatusOK {
		return user, token, fmt.Errorf("error: %s", res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return user, token, err
	}

	user = userInfo.ToUser()
	token = userInfo.ToToken()

	token.TwitchBearer = userToken
	token.TwitchID = user.TwitchID
	token.TwitchRefreshToken = ut.RefreshToken

	return user, token, err
}

func (ut *userTwitch) GetRefreshToken() string {
	return ut.RefreshToken
}

func (ut *userTwitch) GetUserAccessToken() string {
	return ut.UserClient.GetUserAccessToken()
}
