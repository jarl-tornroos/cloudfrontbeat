package cflib

import (
	"encoding/json"
)

// SqsMessage handling the S3 notification
type SqsMessage struct {
	notification S3Notification
}

// S3Notification follows the S3 notification structure
// http://docs.aws.amazon.com/AmazonS3/latest/dev/notification-content-structure.html
type S3Notification struct {
	Records []S3NotificationRecords `json:"Records"`
}

type S3NotificationRecords struct {
	EventVersion string `json:"eventVersion"`
	EventSource  string `json:"eventSource"`
	AwsRegion    string `json:"awsRegion"`
	EventTime    string `json:"eventTime"`
	EventName    string `json:"eventName"`
	UserIdentity struct {
		PrincipalID string `json:"principalId"`
	} `json:"userIdentity"`
	RequestParameters struct {
		SourceIPAddress string `json:"sourceIPAddress"`
	} `json:"requestParameters"`
	ResponseElements struct {
		XAmzRequestID string `json:"x-amz-request-id"`
		XAmzID2       string `json:"x-amz-id-2"`
	} `json:"responseElements"`
	S3 struct {
		S3SchemaVersion string `json:"s3SchemaVersion"`
		ConfigurationID string `json:"configurationId"`
		Bucket struct {
			Name string `json:"name"`
			OwnerIdentity struct {
				PrincipalID string `json:"principalId"`
			} `json:"ownerIdentity"`
			Arn string `json:"arn"`
		} `json:"bucket"`
		Object struct {
			Key       string `json:"key"`
			Size      int `json:"size"`
			ETag      string `json:"eTag"`
			Sequencer string `json:"sequencer"`
		} `json:"object"`
	} `json:"s3"`
}

// NewSqsMessage construction function, take the SQS message as argument
func NewSqsMessage(message *string) *SqsMessage {
	sqsmessage := SqsMessage{}
	if message != nil {
		json.Unmarshal([]byte(*message), &sqsmessage.notification)
	}
	return &sqsmessage
}

// Set file/key for s3 notification.
func (s *SqsMessage) SetFile(file string, bucket string) *SqsMessage {
	records := make([]S3NotificationRecords, 1)
	s.notification.Records = records
	s.notification.Records[0].S3.Object.Key = file
	s.notification.Records[0].S3.Bucket.Name = bucket
	return s
}

// GetFiles return the S3 object(s)
func (s *SqsMessage) GetFiles() []string {
	var files []string
	for _, record := range s.notification.Records {
		if record.S3.Object.Key != "" {
			files = append(files, record.S3.Object.Key)
		}
	}

	return files
}

// GetS3Bucket return the name of the bucket where S3 object are stored in
func (s *SqsMessage) GetS3Buckets() []string {
	var buckets []string
	for _, record := range s.notification.Records {
		if record.S3.Bucket.Name != "" {
			buckets = append(buckets, record.S3.Bucket.Name)
		}
	}

	return buckets
}

// GetNotificationJson get notification in json formant
func (s *SqsMessage) GetNotificationJson() string {
	notificationJson, _ := json.Marshal(s.notification)
	return string(notificationJson)
}
