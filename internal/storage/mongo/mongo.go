package mongo

import (
	"Rsbot_only/pkg/clientDB/mongodb"
	"github.com/mentalisit/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	s   *mongo.Client
	log *logger.Logger
}

func InitMongoDB(log *logger.Logger) *DB {
	client, err := mongodb.NewMongoClient()
	if err != nil {
		log.ErrorErr(err)
		return nil
	}

	d := &DB{
		s:   client,
		log: log,
	}
	return d
}
