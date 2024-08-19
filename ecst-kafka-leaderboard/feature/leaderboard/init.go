package leaderboard

import (
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
	kp sarama.SyncProducer
)

func SetDBPool(dbPool *pgxpool.Pool) {
	if dbPool == nil {
		panic("cannot assign nil db pool")
	}

	db = dbPool
}

func SetKafkaProducer(producer sarama.SyncProducer) {
	if producer == nil {
		panic("cannot assign nil kafka producer")
	}

	kp = producer
}
