package mqtt2

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/flow/support"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-mqtt")

// todo: switch to use endpoint registration

// MqttTrigger is simple MQTT trigger
type MqttTrigger struct {
	metadata          *trigger.Metadata
	runner            action.Runner
	client            mqtt.Client
	settings          map[string]string
	config            *trigger.Config
	topicToActionURI  map[string]string
	topicToActionType map[string]string
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MqttTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *MqttTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *MqttTrigger) Init(config *trigger.Config, runner action.Runner) {

	t.config = config
	t.settings = config.Settings
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *MqttTrigger) Start() error {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(t.settings["broker"])
	opts.SetClientID(t.settings["id"])
	opts.SetUsername(t.settings["user"])
	opts.SetPassword(t.settings["password"])
	b, err := strconv.ParseBool(t.settings["cleansess"])
	if err != nil {
		log.Error("Error converting \"cleansess\" to a boolean ", err.Error())
		return err
	}
	opts.SetCleanSession(b)
	if t.settings["store"] != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(t.settings["store"]))
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {

		topic := msg.Topic()
		//TODO we should handle other types, since mqtt message format are data-agnostic
		payload := string(msg.Payload())
		log.Debug("Received msg:", payload)
		log.Debug("Actual topic: ", topic)

		// try topic without wildcards
		actionType, found := t.topicToActionType[topic]
		actionURI, _ := t.topicToActionURI[topic]

		if found {
			t.RunAction(actionType, actionURI, payload, topic)
		} else {
		// search for wildcards

			for _, endpoint := range t.config.Endpoints {
				eptopic := endpoint.Settings["topic"]
				if strings.HasSuffix(eptopic,"/#") {
					// is wildcard, now check actual topic starts with wildcard
					if strings.HasPrefix(topic,strings.TrimSuffix(eptopic,"/#")) {
						// Got a match, now get the action for the wildcard topic
						actionType, found := t.topicToActionType[eptopic]
						actionURI, _ := t.topicToActionURI[eptopic]
						if found {
							t.RunAction(actionType, actionURI, payload, topic)
						}
					}
				}
			}
		}

	})


	client := mqtt.NewClient(opts)
	t.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	i, err := strconv.Atoi(t.settings["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return err
	}

	t.topicToActionType = make(map[string]string)
	t.topicToActionURI = make(map[string]string)

	for _, endpoint := range t.config.Endpoints {
		if token := t.client.Subscribe(endpoint.Settings["topic"], byte(i), nil); token.Wait() && token.Error() != nil {
			log.Errorf("Error subscribing to topic %s: %s", endpoint.Settings["topic"], token.Error())
			panic(token.Error())
		} else {
			t.topicToActionURI[endpoint.Settings["topic"]] = endpoint.ActionURI
			t.topicToActionType[endpoint.Settings["topic"]] = endpoint.ActionType
		}
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *MqttTrigger) Stop() error {
	//unsubscribe from topic
	log.Debug("Unsubcribing from topic: ", t.settings["topic"])
	for _, endpoint := range t.config.Endpoints {
		if token := t.client.Unsubscribe(endpoint.Settings["topic"]); token.Wait() && token.Error() != nil {
			log.Errorf("Error unsubscribing from topic %s: %s", endpoint.Settings["topic"], token.Error())
		}
	}

	t.client.Disconnect(250)

	return nil
}

// RunAction starts a new Process Instance
func (t *MqttTrigger) RunAction(actionType string, actionURI string, payload string, topic string) {

	log.Debug("Starting new Process Instance")
	log.Debug("Action Type: ", actionType)
	log.Debug("Action URI: ", actionURI)
	log.Debug("Payload: ", payload)
	log.Debug("Actual Topic: ", topic)



	req := t.constructStartRequest(payload, topic)
	//err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	//if err != nil {
	//	//http.Error(w, err.Error(), http.StatusBadRequest)
	//	log.Error("Error Starting action ", err.Error())
	//	return
	//}

	//todo handle error
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	action := action.Get(actionType)

	context := trigger.NewContext(context.Background(), startAttrs)

	//todo handle error
	_, replyData, err := t.runner.Run(context, action, actionURI, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s-%s]", actionType, actionURI)

	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			log.Error(err)
		} else {
			t.publishMessage(req.ReplyTo, string(data))
		}
	}
}

func (t *MqttTrigger) publishMessage(topic string, message string) {

	log.Debug("ReplyTo topic: ", topic)
	log.Debug("Publishing message: ", message)

	qos, err := strconv.Atoi(t.settings["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}
	token := t.client.Publish(topic, byte(qos), false, message)
	token.Wait()
}

func (t *MqttTrigger) constructStartRequest(message string, topic string) *StartRequest {

	log.Debug("Received contstruct start request")

	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	data["actualtopic"] = topic
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}
