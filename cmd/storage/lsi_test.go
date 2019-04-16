package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (

	// storcli show J
	outputNoController = `{
	"Controllers":[
	{
			"Command Status" : {
					"CLI Version" : "007.0813.0000.0000 Dec 14, 2018",
					"Operating system" : "Linux 4.18.0-16-generic",
					"Status Code" : 0,
					"Status" : "Success",
					"Description" : "None"
			},
			"Response Data" : {
					"Number of Controllers" : 0,
					"Host Name" : "cengalo",
					"Operating System " : "Linux 4.18.0-16-generic",
					"StoreLib IT Version" : "07.0900.0200.0200",
					"StoreLib IR3 Version" : "16.04-0"
			}
	}
	]
	}`

	// storcli show J
	outputWithController = `{
		"Controllers":[
		{
				"Command Status" : {
						"CLI Version" : "007.0813.0000.0000 Dec 14, 2018",
						"Operating system" : "Linux 4.18.0-15-generic",
						"Status Code" : 0,
						"Status" : "Success",
						"Description" : "None"
				},
				"Response Data" : {
						"Number of Controllers" : 1,
						"Host Name" : "machine",
						"Operating System " : "Linux 4.18.0-15-generic",
						"StoreLib IT Version" : "07.0900.0200.0200",
						"StoreLib IR3 Version" : "16.04-0",
						"System Overview" : [
								{
										"Ctl" : 0,
										"Model" : "AVAGO3108MegaRAID",
										"Ports" : 8,
										"PDs" : 45,
										"DGs" : 0,
										"DNOpt" : 0,
										"VDs" : 0,
										"VNOpt" : 0,
										"BBU" : "Msng",
										"sPR" : "On",
										"DS" : "1&2",
										"EHS" : "Y",
										"ASOs" : 3,
										"Hlth" : "Opt"
								}
						]
				}
		}
		]
		}
		`
)

func TestCheckControllerPresent(t *testing.T) {

	var result LSIController
	err := json.Unmarshal([]byte(outputNoController), &result)
	if err != nil {
		t.Error("unable to unmarshal data")
	}
	ctrl := result.Controllers[0]
	controllerCount := ctrl.ResponseData["Number of Controllers"]

	assert.Equal(t, controllerCount.(float64), float64(0), "controllercount is expected to be zero")

	err = json.Unmarshal([]byte(outputWithController), &result)
	if err != nil {
		t.Error("unable to unmarshal data")
	}
	ctrl = result.Controllers[0]
	controllerCount = ctrl.ResponseData["Number of Controllers"]

	assert.Equal(t, controllerCount.(float64), float64(1), "controllercount is expected to be zero")
}
