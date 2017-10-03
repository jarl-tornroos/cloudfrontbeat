# Cloudfrontbeat

Cloudfrontbet is an Elastic [beat](https://www.elastic.co/products/beats) that fetch [CloudFront](https://aws.amazon.com/cloudfront/) logs and publish them to Elasticsearch.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/jarl-tornroos/cloudfrontbeat`

## Getting Started with Cloudfrontbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.8
* [AWS SDK](https://aws.amazon.com/sdk-for-go/)
* [GeoIP2 Reader for Go](https://github.com/oschwald/geoip2-golang)

### Build

To build the binary for Cloudfrontbeat run the command below. This will generate a binary in the same directory with the name cloudfrontbeat.

```bash
make
```

It is also possible to compile the binary in a Docker container without the need for installing the dependencies on the host. This command will generate a binary for linux/amd64 in the same directory with the name cloudfrontbeat.

```bash
docker run --rm \
--mount type=bind,source="$(pwd)",target=/go/src/github.com/jarl-tornroos/cloudfrontbeat \
-w /go/src/github.com/jarl-tornroos/cloudfrontbeat golang:1.9 \
bash -c "go get github.com/aws/aws-sdk-go/... && \
go get github.com/oschwald/geoip2-golang/... && \
make"
```

### AWS Configuration

You can configure CloudFront to create log files that contain detailed information about every user request that CloudFront receives. If you enable logging for your CloudFront distribution, you can specify an Amazon S3 bucket that you want CloudFront to save files in.

For this beat to work you need to enable the logging and store files in an S3 bucket. Additionally you need to create an SQS queue and add notification event on your S3 bucket to notify the SQS queue that new files has appeared in the S3 bucket. 

Note: the SQS queue need to have "Default Visibility Timeout" set to at least 5 minutes.

You can follow this [example](http://docs.aws.amazon.com/AmazonS3/latest/dev/ways-to-add-notification-config-to-bucket.html) to configure a bucket for notifications (with SQS queue as message destination)

You can also use CloudFormation to provision your environment, you can use the [provided templates](/cloudformation) which will provision the S3 bucket and the SQS queue.

#### Access control configuration

Cloudfrontbeat supports usage of IAM roles, API keys in environment variables and API key in a credential file. It is preferable to use IAM roles if Cloudfrontbeat is run from an EC2 instance.

### IP Geolocation

Cloudfrontbeat support MixMind's geolocation database. [MaxMind's GeoLite2 City](https://dev.maxmind.com/geoip/geoip2/geolite2/) binary version is available for free. However, more accurate version is available for purchase.

### Running Cloudfrontbeat

1. Enable logging to an S3 bucket from your CloudFormation distribution
2. Create the SQS queue and modify the permissions so that the S3 bucket is allowed to write to the queue
3. Add notification on the S3 bucket to notify the queue of new files
4. Download the IP Geolocation database
5. Build Cloudfrontbeat using the instructions above
6. Modify the included cloudfrontbeat.yml for your application

You can also use the [provided CloudFormation templates](/cloudformation) to provision the S3 bucket and the SQS queue.

To run Cloudfrontbeat, run:

```bash
./cloudfrontbeat -c cloudfrontbeat.yml -e
```

#### Backfilling

If you would like to backfill events that are available in the S3 bucket you can modify the following values in cloudfrontbeat.yml: action, start_date,end_date, s3_bucket and distributions.

The action value need to be backfill.

You can also overwrite the values from the command line with e.g.

```bash
./cloudfrontbeat -e -E cloudfrontbeat='{action:backfill,start_date:2017-08-22,end_date:2017-08-23}'
```

### Running Cloudfrontbeat on Docker

Docker images for Cloudfrontbeat are available from the Docker Hub. You can retrieve an image with a docker pull command.

```bash
docker pull jallu/cloudfrontbeat:0.1
```

#### Configuring Cloudfrontbeat on Docker

The Docker image provides several methods for configuring Cloudfrontbeat. The conventional approach is to provide a configuration file and the IP Geolocation database via bind-mounted volumes, but it’s also possible to create a custom image with your configuration and IP Geolocation database included.

#### Bind-Mounted Configuration

One way to configure Cloudfrontbeat on Docker is to provide cloudfrontbeat.yml and GeoLite2-City.mmdb via bind-mounting. Note that the owner of the files has to be root in order work.

In this example we'll pass the AWS credentials as environment variables from a file aws-credentials.list. The content of the file should look like this (Copy paste your keys into the file):

```bash
AWS_ACCESS_KEY_ID=YOUR_KEY_ID
AWS_SECRET_ACCESS_KEY=YOUR_SECRET
```

And then run the container with:

```bash
docker run \
--env-file "$(pwd)"/aws-credentials.list \
--mount type=bind,source="$(pwd)"/cloudfrontbeat.yml,target=/cloudfrontbeat/cloudfrontbeat.yml \
--mount type=bind,source="$(pwd)"/GeoLite2-City.mmdb,target=/cloudfrontbeat/GeoLite2-City.mmdb \
jallu/cloudfrontbeat:0.1
```

#### Custom Image Configuration

It’s possible to embed your Cloudfrontbeat configuration and IP Geolocation database in a custom image. Here is an example Dockerfile to achieve this:

```
FROM jallu/cloudfrontbeat:0.1
ADD cloudfrontbeat.yml /cloudfrontbeat
ADD GeoLite2-City.mmdb /cloudfrontbeat
```

#### Custom Image

The parent image for the Cloudfrontbeat Docker images is Alpine. It is easy to build your own image with help of the provided Dockerfile. However, the cloudfrontbeat binary has to be compiled before creating the image. Here is an example how to compile for Alpine:

```bash
docker run --rm \
--mount type=bind,source="$(pwd)",target=/go/src/github.com/jarl-tornroos/cloudfrontbeat \
-w /go/src/github.com/jarl-tornroos/cloudfrontbeat golang:1.9-alpine \
sh -c "apk add --update gcc musl-dev git && \
go get github.com/aws/aws-sdk-go/... && \
go get github.com/oschwald/geoip2-golang/... && \
go build"
```
