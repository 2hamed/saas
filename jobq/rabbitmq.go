package jobq

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/2hamed/saas/waitfor"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	qName = "jobs"
)

var (
	rabbitMQHost string
	rabbitMQPort string
	rabbitMQUser string
	rabbitMQPass string
)

func loadConfig() {
	rabbitMQHost = os.Getenv("RABBITMQ_HOST")
	rabbitMQPort = os.Getenv("RABBITMQ_PORT")
	rabbitMQUser = os.Getenv("RABBITMQ_USER")
	rabbitMQPass = os.Getenv("RABBITMQ_PASS")
}

func createRabbitMQ(wc webCapture) (*rabbitMQManager, error) {
	loadConfig()

	waitfor.WaitForServices([]string{
		fmt.Sprintf("%s:%s", rabbitMQHost, rabbitMQPort),
	}, 10*time.Second)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPass, rabbitMQHost, rabbitMQPort))
	if err != nil {
		return nil, fmt.Errorf("failed connecting to rabbit: %v", err)
	}

	log.Infof("Conncted to RabbitMQ on %s:%s", rabbitMQHost, rabbitMQPort)

	pubChan, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed creating pub channel: %v", err)
	}

	consumerChan, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed creating consumer channel: %v", err)
	}

	err = consumerChan.Qos(1, 0, true)
	if err != nil {
		return nil, fmt.Errorf("failed setting Qos: %v", err)
	}

	// declaring job queue
	_, err = pubChan.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating job queue: %v", err)
	}

	rmq := &rabbitMQManager{
		qCon:         conn,
		pubChan:      pubChan,
		consumerChan: consumerChan,
		wc:           wc,
		finishChan:   make(chan []string, 100),
		failChan:     make(chan []string, 100),
	}

	rmq.startConsumer()

	return rmq, nil
}

type rabbitMQManager struct {
	qCon *amqp.Connection

	pubChan      *amqp.Channel
	consumerChan *amqp.Channel

	wc webCapture

	finishChan chan []string
	failChan   chan []string
}

func (m *rabbitMQManager) Enqueue(url string, destination string) error {
	return m.pubChan.Publish("", qName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(fmt.Sprintf("%s::%s", url, destination)),
	})
}

func (m *rabbitMQManager) FinishChan() <-chan []string {
	return m.finishChan
}

func (m *rabbitMQManager) FailChan() <-chan []string {
	return m.failChan
}

func (m *rabbitMQManager) startConsumer() {
	consumChan, err := m.consumerChan.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			m.processJob(<-consumChan)
		}
	}()
}
func (m *rabbitMQManager) processJob(d amqp.Delivery) {

	urlPathStr := string(d.Body)

	log.Debug("received job:", urlPathStr)

	urlPath := strings.Split(urlPathStr, "::")

	err := m.wc.Save(urlPath[0], urlPath[1])

	if err != nil {
		log.Error("Saving screenshot failed", err)
		// TODO: retry this or push to a failed jobs queue

		m.failChan <- urlPath
		d.Nack(false, false)
	} else {
		m.finishChan <- urlPath
		d.Ack(false)
	}

}
