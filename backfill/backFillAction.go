package backfill

import (
	"github.com/jarl-tornroos/cloudfrontbeat/awsfacade"
	"github.com/jarl-tornroos/cloudfrontbeat/cflib"
	"github.com/jarl-tornroos/cloudfrontbeat/config"
	"github.com/elastic/beats/libbeat/logp"
	"fmt"
)

type BackFillAction struct {
	done       bool
	Config     *config.Config
	Sqs        *awsfacade.SqsFacade
	S3         *awsfacade.S3Facade
	Date       *cflib.Date
	SqsMessage *cflib.SqsMessage
	Stopper    *cflib.StopPublisher
}

const dateFormat = "2006-01-02"
const timeString = "00:00:00"

// Name returns the name of the action
func (bf *BackFillAction) Name() string {
	return "backfill"
}

// Do list old log files in the S3 bucket and publish them to the queue
func (bf *BackFillAction) Do() error {
	logp.Info("cloudfrontbeat backfill is running! Hit CTRL-C to stop it.")

	// Add observer
	bf.Stopper.Add(bf)

	// Set bucket to read from
	bf.S3.SetBucket(bf.Config.S3Bucket)

	// Get dates from where we want to publish log files
	fromDate, endDate, err := bf.getDates()
	if err != nil {
		return err
	}

	logp.Info("Backfilling logs from %s to %s", fromDate, endDate)

	// Loop trough day by day and publish log files to the queue
	for fromDate <= endDate {
		// If Ctrl + C was hit and we want to exit
		if bf.done {
			break
		}
		err := bf.addLogFilesToQueue(fromDate)
		if err != nil {
			break
		}
		fromDate = bf.nextDay(fromDate)
	}

	return err
}

// getDAtes return start and end dated
func (bf *BackFillAction) getDates() (string, string, error) {
	var err error

	startDateObject, err := bf.Date.SetDate(bf.Config.StartDate, timeString)
	if err != nil {
		return "", "", err
	}
	startDate := startDateObject.Format(dateFormat)

	endDateObject, err := bf.Date.SetDate(bf.Config.EndDate, timeString)
	if err != nil {
		return "", "", err
	}
	endDate := endDateObject.Format(dateFormat)

	if startDate > endDate {
		err = fmt.Errorf("Start date can't be more reacent than end date!")
	}

	return startDate, endDate, err
}

// nextDay increases the date by a day
func (bf *BackFillAction) nextDay(date string) string {
	dateObject, _ := bf.Date.SetDate(date, timeString)
	return dateObject.IncDays(1).Format(dateFormat)
}

// addLogFilesToQueue list all log files for the given date and insert
// references to them in the queue
func (bf *BackFillAction) addLogFilesToQueue(date string) error {
	var err error
	var filesList []string

	distributions := bf.Config.Distributions
	for _, distribution := range distributions {
		files, err := bf.S3.ListFiles(distribution + "." + date)
		if err != nil {
			return err
		}
		filesList = append(filesList, *files...)
	}

	if len(filesList) > 0 {
		logp.Info("Backfilling logs for %s", date)
	} else {
		logp.Info("Found no logs for %s", date)
	}

	for _, file := range filesList {
		logp.Debug("backfill", "Notifying queue about file %s", file)

		// Create Sqs message that have the same structure as S3 notification
		bf.SqsMessage.SetFile(file, bf.Config.S3Bucket)

		message := bf.SqsMessage.GetNotificationJson()
		// Insert file message to the queue
		err = bf.Sqs.SendMessage(&message)
		if err != nil {
			break
		}
	}

	return err
}

// Stop the action
func (bf *BackFillAction) Stop() {
	bf.done = true
}
