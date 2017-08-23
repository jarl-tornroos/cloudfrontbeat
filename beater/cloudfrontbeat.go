package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/jarl-tornroos/cloudfrontbeat/config"
	"github.com/jarl-tornroos/cloudfrontbeat/workers"
	"github.com/jarl-tornroos/cloudfrontbeat/backfill"
	"github.com/jarl-tornroos/cloudfrontbeat/awsfacade"
	"github.com/jarl-tornroos/cloudfrontbeat/cflib"
)

// Cloudfrontbeat contains configurations and the action we want to run
type Cloudfrontbeat struct {
	config  *config.Config
	action  Action
	stopper *cflib.StopPublisher
	done    bool
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := &config.DefaultConfig
	if err := cfg.Unpack(config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Cloudfrontbeat{
		config: config,
	}
	return bt, nil
}

// Run the specified action
func (bt *Cloudfrontbeat) Run(b *beat.Beat) error {

	// TODO: logging in all areas
	var err error

	// Create stop object to add and publish to observers
	bt.stopper = &cflib.StopPublisher{}

	client := b.Publisher.Connect()

	// Common authentication object for S3 and SQS
	auth := awsfacade.NewAuth(bt.config.Region, bt.config.Environment)

	// SQS Aws facade
	sqs, err := awsfacade.NewSqs(auth, bt.config.QueueName)
	if err != nil {
		return err
	}

	// S3 Aws facade
	s3 := awsfacade.NewS3(auth)

	// Worker action for publishing events
	worker := &workers.WorkerAction{
		Config:  bt.config,
		Client:  client,
		Sqs:     sqs,
		S3:      s3,
		Stopper: bt.stopper,
	}
	// Back fill old events to the queue
	backFill := &backfill.BackFillAction{
		Config:     bt.config,
		Sqs:        sqs,
		S3:         s3,
		Date:       &cflib.Date{Location: "UTC"},
		SqsMessage: &cflib.SqsMessage{},
		Stopper:    bt.stopper,
	}

	// Get the preferred action
	bt.action, err = NewActionFactory().
		AddAction(worker).
		AddAction(backFill).
		GetAction(&bt.config.Action)
	if err != nil {
		return err
	}

	// Run the action unless stop has been initiated
	if bt.done != true {
		err = bt.action.Do()
	}

	return err
}

// Stop cloudfrontbeat with the running action
func (bt *Cloudfrontbeat) Stop() {
	logp.Info("Graceful stop has been initiated, please wait for some processes to finnish.")
	bt.done = true
	bt.stopper.NotifySubscribers()
}
