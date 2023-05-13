package auth

import (
	"errors"
	"fmt"
	"net/http"

	"core/internal/config"
	"core/internal/infrastructure/controller/http/v1/response"
	"core/internal/service"
	"core/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	AuthCodeLength         = 30
	ScopeUserSubscriptions = "user_subscriptions"
)

type authRoutes struct {
	cfg  *config.Config
	auth service.Auther
	log  logger.Logger
}

type AuthTwitchSubchatRequest struct {
	Code  string `json:"code"`
	Scope string `json:"scope"`
	State string `json:"state"`
}

// Bind request to struct and validate it.
func (p *AuthTwitchSubchatRequest) Bind(r *http.Request) error {
	p.Code = r.URL.Query().Get("code")
	p.Scope = r.URL.Query().Get("scope")
	p.State = r.URL.Query().Get("state")

	if len(p.Code) != AuthCodeLength {
		return fmt.Errorf("Bind - code length is not %d", AuthCodeLength)
	}

	if p.Scope != ScopeUserSubscriptions {
		return fmt.Errorf("Bind - scope is not %s", ScopeUserSubscriptions)
	}

	if p.State == "" {
		return fmt.Errorf("Bind - state is empty")
	}

	return nil
}

func New(r chi.Router, cfg *config.Config, authSvc service.Auther, log logger.Logger) {
	au := &authRoutes{
		cfg:  cfg,
		auth: authSvc,
		log:  log,
	}
	r.Route("/auth", func(r chi.Router) {
		r.Route("/twitch", func(r chi.Router) {
			r.Get("/subchat", au.AuthTwitchSubchat)
			r.Get("/test", au.Test)
		})
	})
}

func (au *authRoutes) AuthTwitchSubchat(w http.ResponseWriter, r *http.Request) {
	var in AuthTwitchSubchatRequest
	in.Bind(r)

	au.log.Info("auth - AuthTwitchSubchat", in)

	inviteKey, err := au.auth.AuthTwitchSubchat(r.Context(), in.Code, in.State)
	if errors.Is(err, service.ErrUserNotSubscribed) {
		au.log.Error("auth - AuthTwitchSubchat: %w", err)

		invFailedRedirectURL := fmt.Sprintf("%s/failed?err=%s&state=%s", au.cfg.SubchatInviteRedirectURL, err.Error(), in.State)
		http.Redirect(w, r, invFailedRedirectURL, http.StatusFound)

		return
	}
	if err != nil {
		au.log.Error("auth - AuthTwitchSubchat: %w", err)

		invFailedRedirectURL := fmt.Sprintf("%s/failed?err=%s&state=%s", au.cfg.SubchatInviteRedirectURL, err.Error(), in.State)
		http.Redirect(w, r, invFailedRedirectURL, http.StatusFound)

		return
	}

	invSucceedRedirectURL := fmt.Sprintf("%s/succeed?code=%s&state=%s", au.cfg.SubchatInviteRedirectURL, inviteKey, in.State)
	http.Redirect(w, r, invSucceedRedirectURL, http.StatusFound)
}

func (au *authRoutes) Test(w http.ResponseWriter, r *http.Request) {
	err := au.auth.Test(r.Context())
	if err != nil {
		au.log.Error("auth - AuthTwitchSubchat: %w", err)

		render.JSON(w, r, response.Response{Error: response.Error{Code: http.StatusUnauthorized, Message: err.Error()}})
	}

	render.JSON(w, r, response.Response{Data: "ok"})
}
