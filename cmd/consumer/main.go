package main

import(
	"os"
	"strconv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-rabbitmq/internal/service/consumer"
	"github.com/go-rabbitmq/internal/core"
)

var (
	logLevel =	zerolog.DebugLevel // InfoLevel DebugLevel
	version	=	"go-rabbitmq consumer version 1.0"
	configRabbitMQ core.ConfigRabbitMQ
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	configRabbitMQ.User = "guest"
	configRabbitMQ.Password = "guest"
	configRabbitMQ.Port = "localhost:5672/"
	configRabbitMQ.QueueName = "queue_person_quorum"
	configRabbitMQ.TimeDeleyQueue = 500

	getEnv()
}

func getEnv() {
	if os.Getenv("LOG_LEVEL") !=  "" {
		if (os.Getenv("LOG_LEVEL") == "DEBUG"){
			logLevel = zerolog.DebugLevel
		}else if (os.Getenv("LOG_LEVEL") == "INFO"){
			logLevel = zerolog.InfoLevel
		}else if (os.Getenv("LOG_LEVEL") == "ERROR"){
				logLevel = zerolog.ErrorLevel
		}else {
			logLevel = zerolog.DebugLevel
		}
	}
	if os.Getenv("VERSION") !=  "" {
		version = os.Getenv("VERSION")
	}
	if os.Getenv("RMQ_USER") !=  "" {
		configRabbitMQ.User = os.Getenv("RMQ_USER")
	}
	if os.Getenv("RMQ_PASS") !=  "" {
		configRabbitMQ.Password = os.Getenv("RMQ_PASS")
	}
	if os.Getenv("RMQ_PORT") !=  "" {
		configRabbitMQ.Port = os.Getenv("RMQ_PORT")
	}
	if os.Getenv("RMQ_QUEUE") !=  "" {
		configRabbitMQ.QueueName = os.Getenv("RMQ_QUEUE")
	}
	if os.Getenv("TIME_DELAY_QUEUE") !=  "" {
		intVar, _ := strconv.Atoi(os.Getenv("TIME_DELAY_QUEUE"))
		configRabbitMQ.TimeDeleyQueue = intVar
	}
}

func main () {
	log.Debug().Msg("main consumer")
	log.Debug().Msg("-------------------")
	log.Debug().Str("version", version).
				Msg("Enviroment Variables")
	log.Debug().Str("configRabbitMQ.User: ", configRabbitMQ.User).
				Msg("-----")
	log.Debug().Str("configRabbitMQ.Password: ", configRabbitMQ.Password).
				Msg("-----")
	log.Debug().Str("configRabbitMQ.Port :", configRabbitMQ.Port).
				Msg("----")
	log.Debug().Str("configRabbitMQ.QueueName :", configRabbitMQ.QueueName).
				Msg("----")
	log.Debug().Msg("--------------------")
	
	consumer, err := consumer.NewConsumerService(&configRabbitMQ)
	if err != nil {
		log.Error().Err(err).Msg("error message")
		os.Exit(3)
	}

	log.Debug().Interface("consumer",consumer).Msg("main consumer")
	
	consumer.ConsumerQueue()
	//consumer.ConsumerExchange()
}