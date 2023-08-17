package consumer

import(
	"time"
	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/go-rabbitmq/internal/core"
)

var (
	childLogger = log.With().Str("service", "Consumer").Logger()
	consumer_name = "consumer.person"
)


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

func (c *ConsumerService) ConsumerQueue() error {
	childLogger.Debug().Msg("ConsumerQueue")

	ch, err := c.consumer.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	args := amqp.Table{ // queue args
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	}
	q, err := ch.QueueDeclare(	c.configRabbitMQ.QueueName, // name
								true,         // durable
								false,        // delete when unused
								false,        // exclusive
								false,        // no-wait
								args,          // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error declare queue !!!!") 
		return err
	}

	// Adjust qos for Keda purposes
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error declare Qos !!!!") 
		return err
	}

	msgs, err := ch.Consume(	q.Name, // queue
								consumer_name,    // consumer
								false,   // auto-ack
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
			childLogger.Debug().Str("msg.Body:", string(d.Body)).Msg(" Success Receive a message (ConsumerQueue) ") 
			time.Sleep(time.Duration(c.configRabbitMQ.TimeDelayQueue) * time.Millisecond)
			d.Ack(false)
		}
	}()
	<-forever

	return nil
}

func (c *ConsumerService) ConsumerExchange() error {
	childLogger.Debug().Msg("ConsumerExchange")

	ch, err := c.consumer.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	// declare exchange if not exist
	topic_exchange :=  "personCreated" //"personCreated" 
	err = ch.ExchangeDeclare(	topic_exchange, // name
								"direct",      // type
								true,          // durable
								false,         // auto-deleted
								false,         // internal
								false,         // no-wait
								nil,           // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error exchange queue !!!") 
		return err
	}

	q, err := ch.QueueDeclare(	"", // name
								true,         // durable
								false,        // delete when unused
								false,        // exclusive
								false,        // no-wait
								nil,          // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error declare queue !!!!") 
		return err
	}

	err = ch.QueueBind(q.Name, 
						"info11", 		// routing key
						topic_exchange,  // exchange
						false, 
						nil,
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error binding queue !!!!") 
		return err
	}

	err = ch.Qos(2, 0, false)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Qos queue !!!!") 
		return err
	}

	msgs, err := ch.Consume(	q.Name, // queue
								consumer_name,    // consumer
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
			childLogger.Debug().Str("msg.Body:", string(d.Body)).Msg(" Success Receive a message (ConsumerExchange) ") 
			time.Sleep(time.Duration(c.configRabbitMQ.TimeDelayQueue) * time.Millisecond)
		}
	}()
	<-forever

	return nil
}