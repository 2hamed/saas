package jobq

import (
	"fmt"
	"os"
	"strconv"
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

	workersPerInstance int
)

func initConfig() {
	rabbitMQHost = os.Getenv("RABBITMQ_HOST")
	rabbitMQPort = os.Getenv("RABBITMQ_PORT")
	rabbitMQUser = os.Getenv("RABBITMQ_USER")
	rabbitMQPass = os.Getenv("RABBITMQ_PASS")

	var err error
	workersPerInstance, err = strconv.Atoi(os.Getenv("WORKERS_PER_INSTANCE"))

	// if it's not set or invalid, set it to 3
	if err != nil {
		workersPerInstance = 3
	}

	// if it's less than 1 set it to 1
	if workersPerInstance < 0 {
		workersPerInstance = 1
	}

	// if it's more than 10 set it to 10
	// we don't want to overload the instance
	if workersPerInstance > 10 {
		workersPerInstance = 10
	}
}

func createRabbitMQConnection() (*amqp.Connection, error) {

	initConfig()

	waitfor.WaitForServices([]string{
		fmt.Sprintf("%s:%s", rabbitMQHost, rabbitMQPort),
	}, 60*time.Second)

	log.Infof("Conncting to RabbitMQ on %s:%s", rabbitMQHost, rabbitMQPort)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPass, rabbitMQHost, rabbitMQPort))
	if err != nil {
		return nil, fmt.Errorf("failed connecting to rabbit: %w", err)
	}
	return conn, nil
}

func createRabbitMQ(wc webCapture, conn *amqp.Connection) (*rabbitMQManager, error) {
	pubChan, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed creating pub channel: %w", err)
	}

	consumerChan, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed creating consumer channel: %w", err)
	}

	// tell RabbitMQ to buffer this much messages
	err = consumerChan.Qos(workersPerInstance, 0, true)
	if err != nil {
		return nil, fmt.Errorf("failed setting Qos: %w", err)
	}

	// declaring job queue
	_, err = pubChan.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating job queue: %w", err)
	}

	rmq := &rabbitMQManager{
		qCon:         conn,
		pubChan:      pubChan,
		consumerChan: consumerChan,
		wc:           wc,
		finishChan:   make(chan []string, 100),
		failChan:     make(chan []string, 100),

		stopChan: make(chan struct{}),
	}

	for i := 0; i < workersPerInstance; i++ {
		rmq.startConsumer()
	}

	return rmq, nil
}

type rabbitMQManager struct {
	qCon *amqp.Connection

	pubChan      *amqp.Channel
	consumerChan *amqp.Channel

	wc webCapture

	finishChan chan []string
	failChan   chan []string

	stopChan chan struct{}
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
	// There could more than one consumers (worker) for jobs to
	// to process jobs
	go func() {
		log.Debug("Starting a worker to process jobs")
		for {
			select {
			case d := <-consumChan:
				m.processJob(d)
			case <-m.stopChan:
				// this is it guys, stop listening on channel
				return
			}
		}
	}()
}
func (m *rabbitMQManager) processJob(d amqp.Delivery) {

	urlPathStr := string(d.Body)

	log.Debugf("Received job: %s", urlPathStr)

	urlPath := strings.Split(urlPathStr, "::")

	err := m.wc.Save(urlPath[0], urlPath[1])

	if err != nil {
		log.Errorf("Saving screenshot failed: %v", err)

		// TODO: retry this or push to a failed jobs queue

		m.failChan <- urlPath
		d.Nack(false, false)
	} else {
		m.finishChan <- urlPath
		d.Ack(false)
	}

}

func (m *rabbitMQManager) CleanUp() {
	m.stopChan <- struct{}{}
	m.pubChan.Close()
	m.consumerChan.Close()
	m.qCon.Close()
}
