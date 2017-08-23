package workers

import (
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/elastic/beats/libbeat/common"
	"time"
	"reflect"
)

// CfPublisher is an easy way to publish events
type CfPublisher struct {
	client publisher.Client
	event  *common.MapStr
}

// NewCfPublisher is a construct function
func NewCfPublisher(client publisher.Client) *CfPublisher {
	return &CfPublisher{
		client: client,
	}
}

// PublishEvent publish the content in event
func (cfp *CfPublisher) PublishEvent() bool {
	return cfp.client.PublishEvent(*cfp.event)
}

// Close the publisher client
func (cfp *CfPublisher) Close() error {
	return cfp.client.Close()
}

// SetMapStr create a new common.MapStr to publish
func (cfp *CfPublisher) SetMapStr(eventTime time.Time, eventType string) {
	cfp.event = &common.MapStr{
		"@timestamp": common.Time(eventTime),
		"type":       eventType,
	}
}

// AddToMapStr appends content to common.MapStr
func (cfp *CfPublisher) AddToMapStr(fromEvent interface{}) {
	modelReflect := reflect.ValueOf(fromEvent)
	kind := modelReflect.Kind()
	if kind == reflect.Struct {
		fieldsCount := modelReflect.NumField()
		modelRefType := modelReflect.Type()
		for i := 0; i < fieldsCount; i++ {
			value := modelReflect.Field(i).Interface()
			key := modelRefType.Field(i).Tag.Get("workers")
			//key := modelRefType.Field(i).Name
			(*cfp.event)[key] = value
		}
	}
}
