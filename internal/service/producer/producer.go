package producer

import(
	"context"
	"time"
	"encoding/json"
	"math/rand"
	"strconv"
	"net"

	"github.com/rs/zerolog/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/go-rabbitmq/internal/core"
)

var childLogger = log.With().Str("service", "Producer").Logger()
var my_ip = "x.x.x.x"

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

func (p *ProducerService) ProducerQueue(i int) error {
	childLogger.Debug().Msg("ProducerQueue")

	// Get IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		childLogger.Error().Err(err).Msg("Error to get the POD IP addresd") 
		return err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				my_ip = ipnet.IP.String()
			}
		}
	}

	ch, err := p.producer.Channel()
	if err != nil {
		childLogger.Error().Err(err).Msg("error channel the server message") 
		return err
	}
	defer ch.Close()

	args := amqp.Table{ // queue args
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	}
	q, err := ch.QueueDeclare(p.configRabbitMQ.QueueName, // name
								true,         // durable
								false,        // delete when unused
								false,        // exclusive
								false,        // no-wait
								args,          // arguments
	)
	if err != nil {
		childLogger.Error().Err(err).Msg("error declare queue !!!") 
		return err
	}

	person_mock := p.CreateDataMock(i,my_ip)
	body, _ := json.Marshal(person_mock)

	payloadMsg := amqp.Publishing{	
							ContentType:  "application/json",
							Timestamp:    time.Now(),
							DeliveryMode: amqp.Persistent,
							Body:         []byte(body),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
								"", // exchange
								q.Name, // routing key
								false,  // mandatory
								false,  // immediate
								payloadMsg)
	if err != nil {
		childLogger.Error().Err(err).Msg("error publish message") 
		return err
	}

	childLogger.Debug().Str("msg :", string(body)).Msg("Success Publish a message (ProducerQueue)")

	return nil	
}

func (p *ProducerService) ProducerExchange(i int) error {
	childLogger.Debug().Msg("ProducerExchange")

	// Get IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		childLogger.Error().Err(err).Msg("Error to get the POD IP addresd") 
		return err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				my_ip = ipnet.IP.String()
			}
		}
	}

	ch, err := p.producer.Channel()
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

	person_mock := p.CreateDataMock(i,my_ip)
	body, _ := json.Marshal(person_mock)

	payloadMsg := amqp.Publishing{	
							ContentType:  "application/json",
							Timestamp:    time.Now(),
							DeliveryMode: amqp.Persistent,
							Body:         []byte(body),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
								topic_exchange, // exchange
								"info11", // routing key
								false,  // mandatory
								false,  // immediate
								payloadMsg)
	if err != nil {
		childLogger.Error().Err(err).Msg("error publish message") 
		return err
	}

	childLogger.Debug().Str("msg :", string(body)).Msg("Success Publish a message (ProducerExchange)")

	return nil	
}

func (p *ProducerService) CreateDataMock(i int, my_ip string) *core.Message{
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000
	salt := rand.Intn(max-min+1) + min
	key_person := "PERSON-"+ strconv.Itoa(salt)
	key_msg := "key-"+ strconv.Itoa(i)

	person := core.NewPerson(key_person, key_person ,"Mr Luigi","F")
	message := core.NewMessage(key_msg, key_msg, my_ip ,*person)
	return message
}