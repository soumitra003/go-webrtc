package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/soumitra003/goframework/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	ResponseType string
	Scope        string
}

type oauthTokenRequestBody struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	Code         string `json:"code"`
}

type OAuthTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	IDToken          string `json:"id_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type IDTokenPayload struct {
	Issuer          string
	AuthorizedParty string
	Audience        string
	Subject         string
	Email           string
	EmailVerified   bool
	IssuedAt        int
	ExpiresAt       int
}

const (
	authAPIHost  = "accounts.google.com"
	authAPIPath  = "/o/oauth2/v2/auth"
	tokenAPIHost = "www.googleapis.com"
	tokenAPIPath = "/oauth2/v4/token"
)

var config Config
var logger *zap.Logger

// InitGoogleOauth load config from file
func InitGoogleOauth() {
	logger = logging.GetLogger()
	logger.Info(viper.GetString("OAuth2.Google.ClientID"))
	viper.UnmarshalKey("Oauth2.Google", &config)
}

// GenerateAuthServerURL returns URL string for google auth server
func GenerateAuthServerURL(state string) string {
	query := url.Values{}
	query.Add("response_type", config.ResponseType)
	query.Add("client_id", config.ClientID)
	query.Add("redirect_uri", config.RedirectURI)
	query.Add("scope", config.Scope)
	query.Add("state", state)
	authServerURL := url.URL{
		Scheme:   "https",
		Host:     authAPIHost,
		Path:     authAPIPath,
		RawQuery: query.Encode(),
	}
	return authServerURL.String()
}

func FetchToken(code string) (OAuthTokenResponse, error) {
	var respObj OAuthTokenResponse
	reqBody := oauthTokenRequestBody{
		GrantType:    "authorization_code",
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURI:  config.RedirectURI,
		Code:         code,
	}

	url := url.URL{
		Scheme: "https",
		Host:   tokenAPIHost,
		Path:   tokenAPIPath,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		logging.GetLogger().Error(err.Error())
		return respObj, err
	}

	response, err := http.Post(url.String(), "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		logging.GetLogger().Error(err.Error())
		return respObj, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		data, _ := ioutil.ReadAll(response.Body)
		logging.GetLogger().Error(string(data))
		return respObj, errors.New(string(data))
	}

	err = json.NewDecoder(response.Body).Decode(&respObj)
	if err != nil {
		logging.GetLogger().Error(err.Error())
	}

	return respObj, err
}

func DecodeIDToken(idToken string) (IDTokenPayload, error) {
	var idTokenPayload IDTokenPayload
	segmentArr := strings.Split(idToken, ".")
	data, err := jwt.DecodeSegment(segmentArr[1])
	if err != nil {
		logging.GetLogger().Error(err.Error())
		return idTokenPayload, err
	}

	err = json.Unmarshal(data, &idTokenPayload)
	if err != nil {
		logging.GetLogger().Error(err.Error())
	}

	return idTokenPayload, err
}
