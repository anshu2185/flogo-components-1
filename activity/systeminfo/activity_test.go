package systeminfo

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
  "github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"io/ioutil"
	"testing"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}
		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}
	return activityMetadata
}

func TestCreate(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}
func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	fmt.Println("Retrieving Sytem Information, IpAddress without netmask")

	tc.SetInput("includenetmask", false)

	act.Eval(tc)

	//check result attr

	hostname := tc.GetOutput("hostname")
	ipaddress := tc.GetOutput("ipaddress")
	ip6address := tc.GetOutput("ip6address")
	macaddress := tc.GetOutput("macaddress")

	fmt.Println("hostname: ", hostname)
	fmt.Println("ipaddress: ", ipaddress)
	fmt.Println("ip6address: ", ip6address)
	fmt.Println("macaddress: ", macaddress)


	fmt.Println("Retrieving Sytem Information, IpAddress with netmask")

	tc.SetInput("includenetmask", true)

	act.Eval(tc)

	//check result attr

	hostname = tc.GetOutput("hostname")
	ipaddress = tc.GetOutput("ipaddress")
	ip6address = tc.GetOutput("ip6address")
	macaddress = tc.GetOutput("macaddress")

	fmt.Println("hostname: ", hostname)
	fmt.Println("ipaddress: ", ipaddress)
	fmt.Println("ip6address: ", ip6address)
	fmt.Println("macaddress: ", macaddress)


	if hostname == nil {
		t.Fail()
	}

}
