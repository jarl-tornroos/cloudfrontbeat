// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Type          string   `config:"type"`
	Sleep         int      `config:"sleep"`
	Environment   string   `config:"environment"`
	Workers       int      `config:"workers"`
	Region        string   `config:"region"`
	S3Bucket      string   `config:"s3_bucket"`
	Distributions []string `config:"distributions"`
	QueueManager  string   `config:"queue_manager"`
	QueueName     string   `config:"queue_name"`
	GeoManager    string   `config:"geo_manager"`
	MaxMindDb     string   `config:"max_mind_db"`
	Action        string   `config:"action"`
	StartDate     string   `config:"start_date"`
	EndDate       string   `config:"end_date"`
}

var DefaultConfig = Config{
	Type:          "cloudfrontbeat",
	Sleep:         60,
	Environment:   "staging",
	Workers:       2,
	Region:        "eu-west-1",
	S3Bucket:      "cloudfrontlogs",
	Distributions: []string{"prefix/CF_DISTRIBUTION"},
	QueueManager:  "SQS",
	QueueName:     "cloudfrontbeat",
	GeoManager:    "maxmind",
	MaxMindDb:     "GeoIP2-City.mmdb",
	Action:        "worker",
	StartDate:     "2017-01-01",
	EndDate:       "2017-01-01",
}
