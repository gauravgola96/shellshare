package storage

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Status string

const (
	UserDatabase       string = "user"
	CollectionDownload string = "download"
	CollectionUser     string = "user_info"
	CollectionUsage    string = "usage"
)

const (
	Success      Status = "successful"
	Timeout      Status = "timeout"
	Failed       Status = "failed"
	SessionClose Status = "session_closed"
)

type Download struct {
	SSHKeys      string    `bson:"ssh_keys,omitempty"`
	BytesWritten int64     `bson:"bytes_written,omitempty"`
	Status       Status    `bson:"Status,omitempty"`
	UpdatedAt    time.Time `bson:"updated_at,omitempty"`
}

type User struct {
	UserId    string    `bson:"user_id,omitempty"`
	SSHKeys   string    `bson:"ssh_keys,omitempty"`
	JoinedAt  time.Time `bson:"joined_at,omitempty"`
	LastLogin time.Time `bson:"last_login,omitempty"`
}

func UpdateDownloadDetail(ctx context.Context, data Download) error {
	ctx, cancle := context.WithTimeout(ctx, 5*time.Second)
	defer cancle()
	_, err := S.Mongo.Database(UserDatabase).Collection(CollectionDownload).InsertOne(ctx, Download{
		data.SSHKeys,
		data.BytesWritten,
		data.Status,
		time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func RegisterUserData(ctx context.Context, data User) error {
	subLogger := log.With().Str("module", "mongo_helper.RegisterUserData").Logger()

	ctx, cancle := context.WithTimeout(ctx, 5*time.Second)
	defer cancle()

	user := User{}
	filter := bson.D{{"user_id", data.UserId}}
	err := S.Mongo.Database(UserDatabase).Collection(CollectionUser).FindOne(ctx, filter).Decode(&user)
	if err != mongo.ErrNoDocuments {
		subLogger.Error().Err(err).Msgf("User %+v already registered", user)
		return nil
	}
	_, err = S.Mongo.Database(UserDatabase).Collection(CollectionUser).InsertOne(ctx, User{
		UserId:    data.UserId,
		JoinedAt:  time.Now(),
		SSHKeys:   data.SSHKeys,
		LastLogin: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserLastLogin(ctx context.Context, id string) error {
	subLogger := log.With().Str("module", "mongo_helper.RegisterUserData").Logger()

	ctx, cancle := context.WithTimeout(ctx, 5*time.Second)
	defer cancle()

	filter := bson.D{{"user_id", id}}
	update := bson.D{{"$set", bson.D{{"last_login", time.Now()}}}}
	_, err := S.Mongo.Database(UserDatabase).Collection(CollectionUser).UpdateOne(ctx, filter, update)
	if err == mongo.ErrNoDocuments {
		subLogger.Error().Err(err).Msgf("UserId %s not registered", id)
		return nil
	}

	if err != nil {
		return err
	}
	return nil
}

// GetUsers : returns list of users. set limit -1 to get all users limits
func GetUsers(ctx context.Context, limit int64) ([]User, error) {
	subLogger := log.With().Str("module", "mongo_helper.GetUsers").Logger()

	ctx, cancle := context.WithTimeout(ctx, 5*time.Second)
	defer cancle()
	filter := bson.D{{}}

	var (
		res  []User
		opts *options.FindOptions
	)

	if limit != -1 {
		opts = options.Find().SetLimit(limit)
	}

	cursor, err := S.Mongo.Database(UserDatabase).Collection(CollectionUser).Find(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			subLogger.Error().Err(err).Msg("No user found")
			return nil, nil
		} else {
			return nil, err
		}
	}

	if err = cursor.All(context.TODO(), &res); err != nil {
		return nil, err
	}

	return res, nil
}
