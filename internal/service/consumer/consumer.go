package consumer

import(

	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/go-rabbitmq/internal/core"
)

var childLogger = log.With().Str("service", "Consumer").Logger()

type ConsumerService struct{
	consumer *amqp.Connection
	configRabbitMQ *core.ConfigRabbitMQ
}

func NewConsumerService(configRabbitMQ *core.ConfigRabbitMQ) (*ConsumerService, error){
	childLogger.Debug().Msg("NewConsumerService")

	rabbitmqURL := "amqp://" + configRabbitMQ.User + ":" + configRabbitMQ.Password + "@" + configRabbitMQ.Port
	childLogger.Debug().Str("rabbitmqURL :", rabbitmqURL).Msg("Rabbitmq URI")
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		childLogger.Error().Err(err).Msg("error connect to server message") 
		return nil, err
	}
	return &ConsumerService{
		consumer: conn,
		configRabbitMQ: configRabbitMQ,
	}, nil
}

func (c *ConsumerService) Consumer() error {
	childLogger.Debug().Msg("Consumer")

	ch, err := c.consumer.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		c.configRabbitMQ.QueueName, // name
		false,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error queue declare message") 
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error consume message") 
		return err
	}
	
	var forever chan struct{}

	childLogger.Debug().Msg("Starting Consumer...")
	go func() {
		for d := range msgs {
			childLogger.Debug().Msg("++++++++++++++++++++++++++++")
			childLogger.Debug().Str("msg.Body:", string(d.Body)).Msg(" *** ") 
		}
	}()
	<-forever

	return nil
}