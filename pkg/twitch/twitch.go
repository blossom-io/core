package twitch

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/nicklaw5/helix/v2"
)

const (
	AuthURL     = "https://id.twitch.tv/oauth2/authorize"
	TokenURL    = "https://id.twitch.tv/oauth2/token"
	UserInfoURL = "https://id.twitch.tv/oauth2/userinfo"
	GrantType   = "authorization_code"
	UserAgent   = "blossom-core"
)

var (
	ClientID     = os.Getenv("TWITCH_CLIENT_ID")
	ClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	RedirectURI  = os.Getenv("TWITCH_AUTH_REDIRECT_URL")
)

type Twitcher interface {
	GetAppAccessToken() (AccessToken string, ExpiresIn int, RefreshToken string, err error)
	GetUserAccessToken(ctx context.Context, authCode string) (err error)
	GetUsers(ctx context.Context, IDs []string, logins []string) (user helix.ManyUsers, err error)
	GetUserID(ctx context.Context, login string) (userID string, err error)
}

type twitch struct {
	HTTPClient   *http.Client
	UserClient   *helix.Client
	AppClient    *helix.Client
	ClientID     string
	ClientSecret string
	redirectURI  string
}

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

func New(ctx context.Context, clientID, clientSecret, redirectURI string) (Twitcher, error) {
	HTTPClient := &http.Client{}

	AppClient, err := helix.NewClientWithContext(ctx, &helix.Options{
		HTTPClient:   HTTPClient,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	})
	if err != nil {
		return nil, err
	}

	res, err := AppClient.RequestAppAccessToken([]string{""})
	if err != nil {
		return nil, fmt.Errorf(res.ResponseCommon.ErrorMessage)
	}

	AppClient.SetAppAccessToken(res.Data.AccessToken)

	return &twitch{
		HTTPClient:   HTTPClient,
		AppClient:    AppClient,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		redirectURI:  redirectURI,
	}, nil
}

func (t *twitch) GetUserAccessToken(ctx context.Context, authCode string) (err error) {
	t.UserClient, err = helix.NewClientWithContext(ctx, &helix.Options{
		ClientID:     t.ClientID,
		ClientSecret: t.ClientSecret,
		RedirectURI:  t.redirectURI,
	})
	if err != nil {
		return err
	}

	resp, err := t.UserClient.RequestUserAccessToken(authCode)
	if err != nil {
		return err
	}

	t.UserClient.SetUserAccessToken(resp.Data.AccessToken)

	return nil
}

func (t *twitch) RefreshUserAccessToken(ctx context.Context, authCode string) (err error) {
	resp, err := t.UserClient.RequestUserAccessToken(authCode)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	// Set the access token on the client
	t.UserClient.SetUserAccessToken(resp.Data.AccessToken)

	return nil
}

func (t *twitch) GetAppAccessToken() (AccessToken string, ExpiresIn int, RefreshToken string, err error) {
	res, err := t.AppClient.RequestAppAccessToken([]string{""}) // []string{"user:read:email"}
	if err != nil {
		return "", 0, "", err
	}

	t.AppClient.SetAppAccessToken(res.Data.AccessToken)

	return res.Data.AccessToken, res.Data.ExpiresIn, res.Data.RefreshToken, nil
}

func (t *twitch) GetChannelInformation(ctx context.Context, broadcastersIDs []string) (err error) {
	params := helix.GetChannelInformationParams{
		BroadcasterIDs: broadcastersIDs,
	}

	res, err := t.AppClient.GetChannelInformation(&params)
	if err != nil {
		return err
	}

	if res.ResponseCommon.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to get channel information")
	}

	return nil
}

func (t *twitch) GetUsers(ctx context.Context, IDs []string, logins []string) (user helix.ManyUsers, err error) {
	params := helix.UsersParams{
		IDs:    IDs,
		Logins: logins,
	}

	res, err := t.AppClient.GetUsers(&params)
	if err != nil {
		return user, err
	}

	if res.ResponseCommon.StatusCode == http.StatusOK {
		return res.Data, nil
	}

	return user, fmt.Errorf("%s", res.ResponseCommon.Error)
}

func (t *twitch) GetUserID(ctx context.Context, login string) (userID string, err error) {
	params := helix.UsersParams{
		Logins: []string{login},
	}

	res, err := t.AppClient.GetUsers(&params)
	if err != nil {
		return userID, err
	}

	if res.ResponseCommon.StatusCode != http.StatusOK {
		return userID, fmt.Errorf(res.ResponseCommon.ErrorMessage)
	}

	for _, u := range res.Data.Users {
		if u.Login == login {
			return u.ID, nil
		}
	}

	return userID, fmt.Errorf("twitch user id not found")
}
