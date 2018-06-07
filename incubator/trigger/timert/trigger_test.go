package timert

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJSONMetadata()

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

// Run Once, Start Immediately
const testConfig1 string = `{
  "name": "timert",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "false",
				"startImmediate": "true"
      }
    }
  ]
}`

//Run Every 5 seconds, start Immediately
const testConfig2 string = `{
"name": "timert",
"settings": {
},
"handlers": [
	{
		"actionId": "local://testFlow2",
		"settings": {
			"repeating": "true",
			"seconds": "5",
			"startImmediate": "true"
		}
	}
]
}`

// Run Once, Start Delayed at given datetime
const testConfig3 string = `{
  "name": "timert",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "false",
				"startImmediate": "false",
				"startDate" : "2018-03-12T14:44:00Z02:00"
      }
    }
  ]
}`

//Run Every 5 seconds, start Delayed at given time
const testConfig4 string = `{
"name": "timert",
"settings": {
},
"handlers": [
	{
		"actionId": "local://testFlow2",
		"settings": {
			"repeating": "true",
			"seconds": "5",
			"startImmediate": "false",
			"startDate" : "2018-03-12T14:43:00Z02:00"
		}
	}
]
}`

// Multiple timer configurations
const testConfig5 string = `{
  "name": "timert",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "local://testFlow",
      "settings": {
        "repeating": "false",
				"startImmediate": "true"
      }
    },
		{
      "actionId": "local://testFlow2",
      "settings": {
        "repeating": "false",
				"startImmediate": "false",
        "startDate" : "2018-03-12T14:49:00Z02:00"
      }
    },
    {
      "actionId": "local://testFlow3",
      "settings": {
        "repeating": "true",
				"seconds": "5",
				"startImmediate": "false",
        "startDate" : "2018-03-12T14:49:00Z02:00"
      }
    }
  ]
}`

type TestRunner struct {
}

var Test action.Runner

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action (Run): %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action (RunAction): %v", act)
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action (Execute): %v", act)
	return nil, nil
}

func TestTimer(t *testing.T) {
	log.Info("Testing Timer")
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig1), &config)
	f := &TimerFactory{}
	f.metadata = trigger.NewMetadata(jsonMetadata)
	tgr := f.New(&config)
	//runner := &TestRunner{}
	//tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()
	log.Infof("Press CTRL-C to quit")
	for {
	}
}
