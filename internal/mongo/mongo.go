package mongo

import (
	"context"
	"time"

	"github.com/soumitra003/goframework/config"
	"github.com/soumitra003/goframework/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	primaryMongoClient *mongo.Client
)

// InitMongo Initialize primary mongo client
func InitMongo(config config.Config) error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(getCredentials()))
	if err != nil {
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return err
	}
	primaryMongoClient = client
	if err := client.Ping(ctx, options.Client().ReadPreference); err != nil {
		logging.GetLogger().Info("Unable to connect to mongo client")
		logging.GetLogger().Fatal(err.Error())
		panic(err)
	}
	return nil
}

// GetClient returns primary mongo client
func GetClient() *mongo.Client {
	return primaryMongoClient
}

func getCredentials() options.Credential {
	creds := options.Credential{
		Username: "root",
		Password: "123",
	}
	return creds
}
