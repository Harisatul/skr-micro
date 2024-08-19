package cmd

import (
	"context"
	"ecst-kafka-leaderboard/feature/leaderboard"
	"ecst-kafka-leaderboard/feature/shared"
	"ecst-kafka-leaderboard/pkg"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func createGradeConsumer(ctx context.Context) {
	cfg := shared.LoadConfig("config/create_grade_consumer.yaml")
	kafkaCfg := pkg.NewKafkaConsumerConfig()

	consumer, err := sarama.NewConsumerGroup([]string{cfg.Kafka.Broker}, leaderboard.GradeSubmissionConsumerGroup, kafkaCfg)
	if err != nil {
		log.Fatalln("unable to create consumer group", err)
	}

	defer consumer.Close()

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	// Set needed dependencies
	newCtx, cancel := context.WithCancel(ctx)

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		log.Fatalln("unable to create database connection pool", err)
	}
	defer pool.Close()

	producer, err := sarama.NewSyncProducer([]string{cfg.Kafka.Broker}, pkg.NewKafkaProducerConfig())
	if err != nil {
		log.Fatalln("unable to create kafka producer", err)
	}

	defer producer.Close()

	leaderboard.SetDBPool(pool)
	leaderboard.SetKafkaProducer(producer)

	go func() {
		for err = range consumer.Errors() {
			log.Printf("consumer error, topic %s, error %s", leaderboard.CreateGradeTopic, err.Error())
		}
	}()

	go func() {
		for {
			select {
			case <-newCtx.Done():
				log.Println("consumer stopped")
				return
			default:
				err = consumer.Consume(newCtx, []string{leaderboard.CreateGradeTopic},
					pkg.NewKafkaConsumer(&leaderboard.CreateGradeHandler{}, 1000),
				)
				if err != nil {
					log.Printf("consume message error, topic %s, error %s", leaderboard.CreateGradeTopic, err.Error())
					return
				}
			}
		}
	}()

	log.Printf("consumer up and running, topic %s, group: %s", leaderboard.CreateGradeTopic, leaderboard.GradeSubmissionConsumerGroup)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	cancel()
	log.Println("cancelled message without marking offsets")
}
