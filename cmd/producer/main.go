package main

import(
	"os"
	"time"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/go-rabbitmq/internal/service/producer"
	"github.com/go-rabbitmq/internal/core"
)

var (
	logLevel =	zerolog.DebugLevel // InfoLevel DebugLevel
	version	=	"go-rabbitmq producer version 1.0"
	configRabbitMQ core.ConfigRabbitMQ
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	configRabbitMQ.User = "default_user_bz0ey1Nu-TKmjL-BQvz"
	configRabbitMQ.Password = "wbFpBiXKxXtj3T3L1LOtPwkosaM1uzZD"
	configRabbitMQ.Port = "localhost:5672/"
	configRabbitMQ.QueueName = "task_queue"

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
}

func main () {
	log.Debug().Msg("main producer")
	log.Debug().Msg("-------------------")
	log.Debug().Str("version", version).
				Msg("Enviroment Variables")
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
		producer.Producer(i)
		time.Sleep(time.Millisecond * time.Duration(300))
	}
	done <- "END"
}