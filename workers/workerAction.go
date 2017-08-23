package workers

import (
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/jarl-tornroos/cloudfrontbeat/awsfacade"
	"github.com/jarl-tornroos/cloudfrontbeat/cflib"
	"github.com/jarl-tornroos/cloudfrontbeat/config"
)

// WorkerAction
type WorkerAction struct {
	done       bool
	ch         chan WorkerResponse
	Config     *config.Config
	Client     publisher.Client
	Sqs        *awsfacade.SqsFacade
	S3         *awsfacade.S3Facade
	Stopper    *cflib.StopPublisher
	workerList []Worker
}

// WorkerResponse is a response struct from goroutines
type WorkerResponse struct {
	WorkerNr int
	Err      *error
}

// Name returns the name of the action
func (wa *WorkerAction) Name() string {
	return "worker"
}

// Do fetch the log information and publish them
func (wa *WorkerAction) Do() error {
	logp.Info("cloudfrontbeat worker(s) is running! Hit CTRL-C to stop it.")

	var err error

	// Add observer
	wa.Stopper.Add(wa)

	// Create the concurrency communication channel
	wa.ch = make(chan WorkerResponse)

	// Create the needed amount of workers
	for i := 0; i < wa.Config.Workers; i++ {
		file := &cflib.FileHandling{}
		cfLogs := &CloudFrontLogs{}
		date := &cflib.Date{Location: "UTC"}

		geo, err := GetGeoProvider(wa.Config)
		if err != nil {
			return err
		}
		defer geo.CloseDb()
		geoLog := &GeoLog{
			provider: geo,
		}
		cfPublisher := NewCfPublisher(wa.Client)

		// Timer handle the sleeping and interruptions for the sleep
		timer := &cflib.Time{}

		// Add observer
		wa.Stopper.Add(timer)

		worker := Worker{
			sqs:         wa.Sqs,
			s3:          wa.S3,
			file:        file,
			cfLogs:      cfLogs,
			date:        date,
			geoLog:      geoLog,
			cfPublisher: cfPublisher,
			config:      wa.Config,
			timer:       timer,
		}

		wa.workerList = append(wa.workerList, worker)

		// Start worker in goroutine
		go worker.Publish(wa.ch, i)
	}

	wa.handleWorkers()

	return err
}

// handleWorkers take care of the goroutines / workers after send to channel
func (wa *WorkerAction) handleWorkers() {
	// If CTRL-C was hit, we want to gracefully wait for all workers to finnish
	// This counter keep track on the number of workers finished
	workerListDone := 0

	// Receive from channel
	for workerResp := range wa.ch {

		// If Ctrl + C was hit and we want to exit
		if wa.done {
			logp.Debug("worker", "worker %d done", workerResp.WorkerNr)
			workerListDone++
			// If all workers are done
			if workerListDone == wa.Config.Workers {
				wa.Client.Close()
				close(wa.ch)
			}
			continue
		}

		// If an error occurred
		if *workerResp.Err != nil {
			err := *workerResp.Err
			logp.Err("%s", err.Error())
		}

		// Start new goroutine for worker
		go wa.workerList[workerResp.WorkerNr].Publish(wa.ch, workerResp.WorkerNr)
	}

}

// Stop the action
func (wa *WorkerAction) Stop() {
	wa.done = true
}
