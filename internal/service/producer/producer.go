package producer

import(
	"context"
	"time"
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/go-rabbitmq/internal/core"
)

var childLogger = log.With().Str("service", "Producer").Logger()

type ProducerService struct{
	producer *amqp.Connection
	configRabbitMQ *core.ConfigRabbitMQ
}

func NewProducerService(configRabbitMQ *core.ConfigRabbitMQ) (*ProducerService, error){
	childLogger.Debug().Msg("NewProducerService")

	rabbitmqURL := "amqp://" + configRabbitMQ.User + ":" + configRabbitMQ.Password + "@" + configRabbitMQ.Port
	childLogger.Debug().Str("rabbitmqURL :", rabbitmqURL).Msg("Rabbitmq URI")
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		childLogger.Error().Err(err).Msg("error connect to server message") 
		return nil, err
	}
	return &ProducerService{
		producer: conn,
		configRabbitMQ: configRabbitMQ,
	}, nil
}

func (p *ProducerService) Producer(i int) error {
	childLogger.Debug().Msg("Producer")

	ch, err := p.producer.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		p.configRabbitMQ.QueueName, // name
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	person_mock := p.CreateDataMock(i)
	body, _ := json.Marshal(person_mock)

	payloadMsg := amqp.Publishing{	
							ContentType:  "text/plain",
							//DeliveryMode: amqp.Persistent,
							//Priority:     0,
							Body:         []byte(body),
	}

	//childLogger.Error().Interface("",payloadMsg).Msg("error publish message") 
	err = ch.PublishWithContext(ctx,
								"",     // exchange
								q.Name, // routing key
								false,  // mandatory
								false,  // immediate
								payloadMsg)
	if err != nil {
		childLogger.Error().Err(err).Msg("error publish message") 
		return err
	}

	childLogger.Debug().Str("msg :", string(body)).Msg("Success Publish a message")

	return nil	
}

func (p *ProducerService) CreateDataMock(i int) *core.Message{
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000
	salt := rand.Intn(max-min+1) + min
	key_person := "PERSON-"+ strconv.Itoa(salt)
	key_msg := "key-"+ strconv.Itoa(i)

	person := core.NewPerson(key_person, key_person,"Mr Luigi","F")
	message := core.NewMessage(key_msg, key_msg, *person)
	return message
}