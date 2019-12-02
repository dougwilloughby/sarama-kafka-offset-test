package main

import (
	"flag"
	"github.com/Shopify/sarama"
	"github.com/dougwilloughby/sarama-kafka-offset-test/configs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {

	var configFileName = flag.String("config", "", "name of configuration file")
	var logLevel = flag.String("loglevel", "info", "logging level (\"info\",\"debug\")")
	var logPretty = flag.Bool("logpretty", false, "Pretty-print log if true")
	flag.Parse()

	// Set up logging
	switch *logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info": // default...added for completeness and documentation
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
	log.Logger = log.With().Str("service", os.Args[0]).Int("pid", os.Getpid()).Logger()
	if *logPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Info().Msgf("started, logging level: %v", zerolog.GlobalLevel().String())

	// Load the configuration. Configuration can be specified from multiple sources. The order
	// of loading is important. The order is defaults, default configuration file, specified
	//configuration file, environment

	props := configs.NewAppConf()
	if err := props.ReadFile("./offset_test.conf"); err != nil {
		log.Fatal().
			Err(err).
			Msg("Error reading default config file ")
	}
	if len(*configFileName) > 0 {
		if err := props.ReadFile(*configFileName); err != nil {
			log.Fatal().
				Err(err).
				Msg("Error reading specified config file")
		}
	}

	log.Debug().Msgf("Configuration:\n%+v", *props)

	// Set up to handle CTRL-C or kill and shut down gracefully
	//var wg sync.WaitGroup                    // to keep goroutine alive until kafka is done draining
	interruptChan := make(chan os.Signal, 1) // channel for handling interruptChan (ctrl-C)
	signal.Notify(interruptChan, os.Interrupt)

	// Initialize history transaction log Kafka consumer, which is used to ingest
	// history records
	consumer, err := sarama.NewConsumer(props.Brokers, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka consumer")
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error().Err(err).Msg("Problem closing the Kafka consumer")
		}
	}()

	// Listen to just the partition that is producing history records for
	// the customer. Customers are only producing to a single partition
	partitionConsumer, err := consumer.ConsumePartition(props.SrcTopic, props.SrcPartition, sarama.OffsetOldest)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka partition listener")
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Error().Err(err).Msg("Problem closing Kafka comsumer's partition")
		}
	}()

	// Consume messages until interrupted with a control-C
	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-interruptChan:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
