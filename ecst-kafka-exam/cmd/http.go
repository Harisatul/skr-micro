package cmd

import (
	"context"
	"ecst-kafka-exam/feature/shared"
	"ecst-kafka-exam/feature/tryout"
	"ecst-kafka-exam/pkg"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

func runHTTPServer(ctx context.Context) {
	// Load configuration
	cfg := shared.LoadConfig("config/app.yaml")

	dbCfg, err := pgxpool.ParseConfig(cfg.DBConfig.ConnStr())
	if err != nil {
		log.Fatalln("unable to parse database config", err)
	}

	// Set needed dependencies
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

	tryout.SetDBPool(pool)
	tryout.SetKafkaProducer(producer)

	// Create a new server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tryout.HttpRoute(mux)

	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("unable to start server", err)
		}
	}()

	log.Println("server started")

	// Wait for signal to shut down
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalln("unable to shutdown server", err)
	}

	log.Println("server shutdown")
}
