package auth

import (
	"context"
	"time"

	"github.com/soumitra003/go-webrtc/internal/mongo"
	"github.com/soumitra003/go-webrtc/internal/oauth2/google"
	"github.com/soumitra003/goframework/config"
)

var (
	userRepository UserRepository
)

// ModuleAuth main auth module
type ModuleAuth struct {
	config *config.Config
}

//New creates module instance
func New(config config.Config) *ModuleAuth {
	md := &ModuleAuth{config: &config}
	return md
}

// Init initializes auth module
func (h *ModuleAuth) Init(ctx context.Context, config config.Config) {
	mongo.InitMongo(config)
	google.InitGoogleOauth()
	coll := mongo.GetClient().Database("dev").Collection("user")
	timeoutContext := func() context.Context {
		ct, _ := context.WithTimeout(context.Background(), 10*time.Second)
		return ct
	}
	mongoRepo := UserRepositoryMongoImpl{
		collection:     coll,
		timeoutContext: timeoutContext,
	}
	userRepository = UserRepository(&mongoRepo)
}
