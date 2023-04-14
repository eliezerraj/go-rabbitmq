package main

import(
	"os"
	"time"
	"os/signal"
	"syscall"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-rabbitmq/internal/service/producer"
	"github.com/go-rabbitmq/internal/core"
)

var (
	logLevel =	zerolog.DebugLevel // InfoLevel DebugLevel
	version	=	"go-rabbitmq producer version 1.0"
	configRabbitMQ core.ConfigRabbitMQ
	duration = 1000
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	configRabbitMQ.User = "guest"
	configRabbitMQ.Password = "guest"
	configRabbitMQ.Port = "localhost:5672/"
	configRabbitMQ.QueueName = "queue_person_quorum"

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
	if os.Getenv("DURATION") !=  "" {
		duration, _ = strconv.Atoi(os.Getenv("DURATION"))
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

	producer, err := producer.NewProducerService(&configRabbitMQ)
	if err != nil {
		log.Error().Err(err).Msg("error message")
		os.Exit(3)
	}
	//log.Debug().Interface("producer",producer).Msg("main producer")

	done := make(chan string)
	go producerMessage(producer, done)

	// Shut down main function
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
		
	log.Debug().Msg("Stopping main producer...")
	defer func() {
		log.Debug().Msg("main producer stopped !!!!")
	}()
	
}

func producerMessage(producer *producer.ProducerService,done chan string){
	for i := 0 ; i < 36000; i++ {
		producer.ProducerQueue(i)
		//producer.ProducerExchange(i)
		time.Sleep(time.Millisecond * time.Duration(duration))
	}
	done <- "END"
}