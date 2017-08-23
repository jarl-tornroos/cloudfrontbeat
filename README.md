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

To build the binary for Cloudfrontbeat run the command below. This will generate a binary
in the same directory with the name cloudfrontbeat.

```
make
```

### Run

To run Cloudfrontbeat with debugging output enabled, run:

```
./cloudfrontbeat -c cloudfrontbeat.yml -e -d "*"
```
