package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"

	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
)

type OAuth2UserInfo struct {
	Id    string
	Name  string
	Email string
}

type OAuth2Service interface {
	GenerateAuthURL() string
	GetToken(state, code string) (string, error)
	GetUserInfo(token string) (*OAuth2UserInfo, error)
}

type oauthService struct {
	config       oauth2.Config
	stateString  string
	infoQueryURL string
	logger       log.Logger
}

func NewOauth2(config oauth2.Config, infoQueryURL string, logger log.Logger) OAuth2Service {
	return &oauthService{
		config:       config,
		stateString:  uuid.NewV4().String(),
		infoQueryURL: infoQueryURL,
		logger:       logger,
	}
}

func (os *oauthService) GenerateAuthURL() string {
	Url, _ := url.Parse(os.config.Endpoint.AuthURL)

	parameters := url.Values{}
	parameters.Add("client_id", os.config.ClientID)
	parameters.Add("scope", strings.Join(os.config.Scopes, " "))
	parameters.Add("redirect_uri", os.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", os.stateString)
	Url.RawQuery = parameters.Encode()

	return Url.String()
}

func (os *oauthService) GetToken(state, code string) (string, error) {
	if state != os.stateString {
		return "", errors.New("invalid state")
	}

	token, err := os.config.Exchange(context.Background(), code)
	if err != nil {
		return "", errors.New("could not exchange code for token")
	}

	return token.AccessToken, nil
}

func (os *oauthService) GetUserInfo(token string) (*OAuth2UserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s?fields=%s&access_token=%s", os.infoQueryURL, strings.Join(os.config.Scopes, ","), url.QueryEscape(token)))
	if err != nil {
		return nil, errors.Wrap(err, "could not make get request")
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			os.logger.WithError(err).Warnln("Could not close response body")
		}
	}()

	userInfo := new(OAuth2UserInfo)
	if err = json.NewDecoder(resp.Body).Decode(userInfo); err != nil {
		return nil, errors.Wrap(err, "could not decode response body")
	}

	return userInfo, nil
}
