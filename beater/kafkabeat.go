package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/vrecan/kafkabeat/config"
)

// Kafkabeat configuration.
type Kafkabeat struct {
	done        chan struct{}
	config      config.Config
	clusterConf *cluster.Config
	client      beat.Client
}

// New creates an instance of kafkabeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	bt := &Kafkabeat{
		done:        make(chan struct{}),
		config:      c,
		clusterConf: config,
	}
	return bt, nil
}

// Run starts kafkabeat.
func (bt *Kafkabeat) Run(b *beat.Beat) error {
	logp.Info("kafkabeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	consumer, err := cluster.NewConsumer(bt.config.Brokers, bt.config.Group, bt.config.Topics, bt.clusterConf)
	if err != nil {
		return err
	}
	defer consumer.Close()
	errChan := consumer.Errors()
	notifyChan := consumer.Notifications()
	kafkaMsgs := consumer.Messages()
	cnt := 1
	for {
		select {

		case err := <-errChan:
			fmt.Println("Kafka Error: ", err)
		case notify := <-notifyChan:
			fmt.Println("Kafka Notification: ", notify)
		case msg := <-kafkaMsgs:
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type": b.Info.Name,
					"cnt":  cnt,
					"msg":  msg,
				},
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
			cnt++
		case <-bt.done:
			return nil
		}

	}
}

// Stop stops kafkabeat.
func (bt *Kafkabeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
