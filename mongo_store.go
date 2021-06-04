package users27

import (
	"context"
	"github.com/google/uuid"
	"github.com/kirigaikabuto/config27"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type usersStore struct {
	collection *mongo.Collection
}

func NewUsersStore(config config27.MongoConfig) (UsersStore, error) {
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Host + ":" + config.Port)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(config.Database)
	err = db.CreateCollection(context.TODO(), config.CollectionName)
	if err != nil && !strings.Contains(err.Error(), "NamespaceExists") {
		return nil, err
	}
	collection := db.Collection(config.CollectionName)
	return &usersStore{collection: collection}, nil
}

func (u *usersStore) Create(user *User) (*User, error) {
	token := uuid.New()
	user.Id = token.String()
	_, err := u.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usersStore) Get(id string) (*User, error) {
	filter := bson.D{{"id", id}}
	user := &User{}
	err := u.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usersStore) GetByUsernameAndPassword(username, password string) (*User, error) {
	filter := bson.D{{"username", username}, {"password", password}}
	user := &User{}
	err := u.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return nil, ErrNoUser
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
