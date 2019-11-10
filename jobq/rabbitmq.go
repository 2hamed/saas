package jobq

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

const (
	qName = "jobs"
)

var (
	rabbitMQHost = os.Getenv("RABBITMQ_HOST")
	rabbitMQPort = os.Getenv("RABBITMQ_PORT")
	rabbitMQUser = os.Getenv("RABBITMQ_USER")
	rabbitMQPass = os.Getenv("RABBITMQ_PASS")
)

func createRabbitMQ(webCapture WebCapture) (*rabbitMQManager, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPass, rabbitMQHost, rabbitMQPort))
	if err != nil {
		return nil, err
	}

	pubChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	consumerChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	consumerChan.Qos(10, 10, true)

	// declaring job queue
	pubChan.QueueDeclare(qName, true, false, true, true, nil)

	return &rabbitMQManager{conn, pubChan, consumerChan, webCapture}, nil
}

type rabbitMQManager struct {
	qCon *amqp.Connection

	pubChan      *amqp.Channel
	consumerChan *amqp.Channel

	webCapture WebCapture
}

func (m *rabbitMQManager) Enqueue(url string) error {
	return m.pubChan.Publish("", qName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(url),
	})
}

func (m *rabbitMQManager) startConsumer() {
	consumChan, err := m.consumerChan.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func(ch <-chan amqp.Delivery) {
		for {
			d := <-ch

			err := m.webCapture.Save(string(d.Body), "")
			if err != nil {
				// retry this or push to a failed jobs queue
				d.Nack(false, false)
			} else {
				d.Ack(false)
			}
		}
	}(consumChan)
}
