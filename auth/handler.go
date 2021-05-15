package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/soumitra003/go-webrtc/internal/oauth2/google"
	"github.com/soumitra003/go-webrtc/internal/render"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *ModuleAuth) register(writer http.ResponseWriter, req *http.Request) {
	user := NewUser()

	// Assign login times as user will be logged in after this call
	user.FirstLoginAt = time.Now().UTC()
	user.LastLoginAt = time.Now().UTC()

	err := userRepository.Insert(req.Context(), &user)
	if err != nil {
		render.RenderError(writer, err)
		return
	}

	loginUser(writer, user)
	render.RenderBaseResponse(writer, user)
}

func (m *ModuleAuth) oAuth2Init(writer http.ResponseWriter, req *http.Request) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		render.RenderError(writer, err)
		return
	}

	state := uuid.String()
	authServerURL := google.GenerateAuthServerURL(state)

	render.Redirect(writer, authServerURL, http.StatusTemporaryRedirect)
}

func (m *ModuleAuth) oAuth2Login(writer http.ResponseWriter, req *http.Request) {
	// Resolve params
	code := req.URL.Query().Get("code")
	if code == "" {
		render.RenderError(writer, errors.New("required parameter [code, reidrectUri]"))
		return
	}

	respObj, err := google.FetchToken(code)
	if err != nil {
		render.RenderError(writer, err)
		return
	}

	user, err := newUserFromGoogleIDToken(respObj.IDToken)
	if err != nil {
		render.RenderError(writer, err)
		return
	}

	var existingUser User
	err = userRepository.FindByEmail(req.Context(), user.Email, &existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		render.RenderError(writer, err)
		return
	}

	// If user exists then update
	if existingUser.ID != "" {
		existingUser.LastLoginAt = time.Now().UTC()
		existingUser.UpdatedAt = time.Now().UTC()
		err = userRepository.Update(req.Context(), &existingUser)
		if err != nil {
			render.RenderError(writer, err)
			return
		}
		render.RenderBaseResponse(writer, existingUser)
		return
	}

	// If user doesn't exist then insert
	user.FirstLoginAt = time.Now().UTC()
	user.LastLoginAt = time.Now().UTC()
	err = userRepository.Insert(req.Context(), &user)
	if err != nil {
		render.RenderError(writer, err)
		return
	}
	render.RenderBaseResponse(writer, user)
}

func loginUser(writer http.ResponseWriter, user User) {
	// Generate JWT for user
	tokenStr, err := generateTokenWithClaims(user)
	if err != nil {
		render.RenderError(writer, err)
		return
	}
	writer.Header().Add("Set-Cookie", "token="+tokenStr)
}
