package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// UserRepository interface
type UserRepository interface {
	FindByID(context.Context, string, *User) error
	FindByEmail(context.Context, string, *User) error
	Insert(context.Context, *User) error
	Update(context.Context, *User) error
	Save(context.Context, *User) error
	DeleteByID(context.Context, string) error
}

type newTimeoutContext func() context.Context

// UserRepositoryMongoImpl mongo impl
type UserRepositoryMongoImpl struct {
	collection     *mongo.Collection
	timeoutContext newTimeoutContext
}

// FindByID find user by id
func (urmi UserRepositoryMongoImpl) FindByID(ctx context.Context, uuid string, user *User) error {
	if err := urmi.collection.FindOne(urmi.timeoutContext(), bson.M{"_id": uuid}).Decode(user); err != nil {
		return err
	}
	return nil
}

// FindByEmail find user by id
func (urmi UserRepositoryMongoImpl) FindByEmail(ctx context.Context, email string, user *User) error {
	if err := urmi.collection.FindOne(urmi.timeoutContext(), bson.M{"email": email}).Decode(user); err != nil {
		return err
	}
	return nil
}

// Insert inserts a user into mongo db
func (urmi UserRepositoryMongoImpl) Insert(ctx context.Context, user *User) error {
	_, err := urmi.collection.InsertOne(urmi.timeoutContext(), *user)
	return err
}

// Update update user
func (urmi UserRepositoryMongoImpl) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now().UTC()
	map1 := make(map[string]interface{})
	map1["$set"] = user
	_, err := urmi.collection.UpdateOne(
		urmi.timeoutContext(),
		bson.M{"_id": user.ID},
		map1,
	)
	return err
}

// Save save a user
func (urmi UserRepositoryMongoImpl) Save(ctx context.Context, user *User) error {
	existingUser := &User{}
	// Look for existing user with ID
	err := urmi.FindByID(ctx, user.ID, existingUser)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		// If User not found then create
		err = urmi.Insert(ctx, user)
		return err
	}
	// If User fond then update
	err = urmi.Update(ctx, existingUser)
	return err
}

// DeleteByID delete user by ID
func (urmi UserRepositoryMongoImpl) DeleteByID(ctx context.Context, id string) error {
	return nil
}
