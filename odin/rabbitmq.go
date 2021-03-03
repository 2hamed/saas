package odin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/2hamed/saas/waitfor"
	"github.com/rs/zerolog/log"
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

}

func createRabbitMQConnection() (*amqp.Connection, error) {

	initConfig()

	waitfor.WaitForServices([]string{
		fmt.Sprintf("%s:%s", rabbitMQHost, rabbitMQPort),
	}, 60*time.Second)

	log.Info().Msgf("Conncting to RabbitMQ on %s:%s", rabbitMQHost, rabbitMQPort)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPass, rabbitMQHost, rabbitMQPort))
	if err != nil {
		return nil, fmt.Errorf("failed connecting to rabbit: %w", err)
	}
	return conn, nil
}

func createRabbitMQ(conn *amqp.Connection) (*rabbitMQManager, error) {
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
		qCon: conn,
	}

	return rmq, nil
}

type rabbitMQManager struct {
	qCon *amqp.Connection
}

func (m *rabbitMQManager) Enqueue(ctx context.Context, job CaptureJob) error {
	qChan, err := m.qCon.Channel()
	if err != nil {
		return fmt.Errorf("failed creating RabbitMQ channel: %w", err)
	}
	defer qChan.Close()

	bytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed marshalling job: %w", err)
	}

	return qChan.Publish("", qName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	})
}
func (m *rabbitMQManager) GetJobChan(ctx context.Context) (<-chan CaptureJob, error) {
	qChan, err := m.qCon.Channel()
	if err != nil {
		return nil, err
	}
	delivery, err := qChan.Consume(qName, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	jobChan := make(chan CaptureJob)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(jobChan)
				qChan.Close()
			case d := <-delivery:
				var job CaptureJob
				err := json.Unmarshal(d.Body, &job)
				if err != nil {
					log.Error().Bytes("body", d.Body).Err(err).Msg("invalid job")
					continue
				}
				job.Ack = func() {
					d.Ack(false)
				}
				job.Nack = func() {
					d.Nack(false, false)
				}
				jobChan <- job
			}
		}
	}()
	return jobChan, nil
}
func (m *rabbitMQManager) CleanUp() {
	m.qCon.Close()
}
