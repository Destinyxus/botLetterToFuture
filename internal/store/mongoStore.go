package store

import (
	"context"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Destinyxus/botLetterToFuture/internal/model"
	"github.com/Destinyxus/botLetterToFuture/pkg/config"
)

type Letter struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Email           string             `bson:"email,omitempty"`
	Date            string             `bson:"date,omitempty"`
	EncryptedLetter string             `bson:"encrypted_letter,omitempty"`
	Sent            string             `bson:"sent,omitempty,default:false"`
}

type MongoDB struct {
	client *mongo.Client
	log    *logrus.Entry
}

func NewStoreMongo(log *logrus.Entry) *MongoDB {
	return &MongoDB{
		log: log.WithFields(logrus.Fields{
			"package": "store",
		}),
	}
}
func (s *MongoDB) Open(cfg *config.Config) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.StoreURL))
	if err != nil {
		s.log.Error("db connection error")
	}

	err = client.Ping(context.TODO(), nil)
	s.log.Info("Db opened")

	s.client = client
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *MongoDB) CloseConnection() error {
	if err := s.client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
	return nil
}

func (s *MongoDB) CreateALetter(m *model.User, userInfo *model.User) error {
	coll := s.client.Database("letterToFuture").Collection("letters")

	letter := &Letter{
		Email:           m.Email,
		Date:            m.Date,
		EncryptedLetter: m.EncryptedLetter,
		Sent:            "false",
	}

	_, err := coll.InsertOne(context.Background(), letter)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"UserName": userInfo.UserName,
		}).Error("Creating letter error")

	}

	s.log.WithFields(logrus.Fields{
		"UserName": userInfo.UserName,
	}).Info("Letter successfully created")
	return nil
}

func (s *MongoDB) GetLetter() ([]*Letter, error) {
	coll := s.client.Database("letterToFuture").Collection("letters")

	currentDate := time.Now().Format("2006-01-02")
	filter := bson.M{"date": currentDate, "sent": "false"}
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	var letters []*Letter

	for cursor.Next(context.Background()) {
		var letter Letter
		if err := cursor.Decode(&letter); err != nil {
			return nil, err
		}
		letters = append(letters, &letter)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return letters, nil
}

func (s *MongoDB) IsSent(id primitive.ObjectID) error {
	coll := s.client.Database("letterToFuture").Collection("letters")

	_, err := coll.UpdateByID(context.Background(), id, bson.M{"$set": bson.M{
		"sent": "true",
	}})
	if err != nil {
		s.log.Error("err")
	}

	return nil
}
