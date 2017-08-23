package workers

import (
	"github.com/jarl-tornroos/cloudfrontbeat/awsfacade"
	"github.com/jarl-tornroos/cloudfrontbeat/cflib"
	"github.com/jarl-tornroos/cloudfrontbeat/config"
	"github.com/elastic/beats/libbeat/logp"
	"path"
)

// Worker publish CloudFront events from S3 files
type Worker struct {
	s3          *awsfacade.S3Facade
	sqs         *awsfacade.SqsFacade
	file        *cflib.FileHandling
	cfLogs      *CloudFrontLogs
	date        *cflib.Date
	geoLog      *GeoLog
	cfPublisher *CfPublisher
	config      *config.Config
	timer       cflib.Timer
}

// Publish the CloudFront events
func (w *Worker) Publish(ch chan WorkerResponse, workerNr int) {

	logp.Debug("worker", "goroutine for worker %d started", workerNr)

	// Receive message from SQS
	message, err := w.sqs.ReceiveMessage()
	if w.responseHasError(err, ch, workerNr) {
		return
	}

	sqsMessage := cflib.NewSqsMessage(message.Body)

	// Get the S3 file location from the SQS message
	files := sqsMessage.GetFiles()

	// Get the S3 buckets for where files are stored in
	buckets := sqsMessage.GetS3Buckets()

	// Should only be one file but the type is an array so lets loop
	for i, srcFile := range files {

		// Set bucket where file is stored
		w.s3.SetBucket(buckets[i])

		fileName := path.Base(srcFile)
		localFile := "/tmp/" + fileName

		// Download the log Cloudfront log file from S3 bucket
		_, err := w.s3.Download(srcFile, localFile)
		if w.responseHasError(err, ch, workerNr) {
			return
		}

		// Read the compressed content of the file and delete the file from local disk
		fileContent, err := w.file.SetFile(localFile).ReadZipContent()
		w.file.Delete()
		if w.responseHasError(err, ch, workerNr) {
			return
		}

		// Set the file content to the CloudFront log object
		w.cfLogs.SetContent(fileContent)

		// Loop trough the struct formatted log lines
		for _, log := range w.cfLogs.Logs {
			// Add the basic date to publish
			timeObject, err := w.date.SetDate(log.Date, log.Time)
			if err != nil {
				logp.Warn("%s", err.Error())
				continue
			}

			timeStamp := timeObject.GetTimeInstance()
			w.cfPublisher.SetMapStr(timeStamp, w.config.Type)

			// Get the IP location information
			err = w.geoLog.SetIP(log.CIp)
			if err != nil {
				logp.Warn("%s", err.Error())
				continue
			}

			geoData := w.geoLog.GetGeoData()

			// Add the CloudFront log data and the geo data before publishing
			w.cfPublisher.AddToMapStr(log)
			w.cfPublisher.AddToMapStr(geoData)

			// publish the event
			w.cfPublisher.PublishEvent()
		}
	}

	err = w.sqs.DeleteMessage(message.ReceiptHandle)
	if w.responseHasError(err, ch, workerNr) {
		return
	}

	// Sleep if queue was empty, we don't want to hit the queue too often
	if len(files) == 0 {
		logp.Debug(
			"worker",
			"Nothing in queue, worker %d will sleep for %d seconds",
			workerNr,
			w.config.Sleep,
		)
		w.timer.SleepSec(w.config.Sleep)
	}

	// Inform channel that the worker has finished
	ch <- WorkerResponse{
		WorkerNr: workerNr,
		Err:      &err,
	}
}

// Check if response has error and write to channel if there is an error
func (w *Worker) responseHasError(err error, ch chan WorkerResponse, workerNr int) bool {
	if err != nil {
		// Inform channel about the error
		ch <- WorkerResponse{
			WorkerNr: workerNr,
			Err:      &err,
		}
		return true
	} else {
		return false
	}
}
