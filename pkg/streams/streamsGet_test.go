/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package streams

import (
	"net/http"
	"testing"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestMultipleStreamsGetFormatsResultsOk(t *testing.T) {

	//Given..

	body := `
[
  {
    "apiVersion": "galasa-dev/v1alpha1",
    "kind": "GalasaStream",
    "metadata": {
      "name": "mystream",
      "url": "http://localhost:8080/api/streams/myStream",
      "description": "A stream which I use to..."
    },
    "data": {
      "isEnabled": true,
      "repository": {
        "url": "http://mymavenrepo.host/testmaterial"
      },
      "obrs": [
        {
          "group-id": "com.ibm.zosadk.k8s",
          "artifact-id": "com.ibm.zosadk.k8s.obr",
          "version": "0.1.0-SNAPSHOT"
        }
      ],
      "testCatalog": {
        "url": "http://mymavenrepo.host/testmaterial/com.ibm.zosadk.k8s/com.ibm.zosadk.k8s.obr/0.1.0-SNAPSHOT/testcatalog.yaml"
      }
    }
  },
  {
    "apiVersion": "galasa-dev/v1alpha1",
    "kind": "GalasaStream",
    "metadata": {
      "name": "mystream2",
      "url": "http://localhost:8080/api/streams/myStream",
      "description": "Another stream to..."
    },
    "data": {
      "isEnabled": true,
      "repository": {
        "url": "http://mymavenrepo.host/testmaterial"
      },
      "obrs": [
        {
          "group-id": "com.ibm.zosadk.k8s",
          "artifact-id": "com.ibm.zosadk.k8s.obr",
          "version": "0.1.0-SNAPSHOT"
        }
      ],
      "testCatalog": {
        "url": "http://mymavenrepo.host/testmaterial/com.ibm.zosadk.k8s/com.ibm.zosadk.k8s.obr/0.1.0-SNAPSHOT/testcatalog.yaml"
      }
    }
  }
  
]
`

	getStreamsInteraction := utils.NewHttpInteraction("/streams", http.MethodGet)
	getStreamsInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("ClientApiVersion", "myVersion")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(body))
	}

	interactions := []utils.HttpInteraction{
		getStreamsInteraction,
	}

	server := utils.NewMockHttpServer(t, interactions)
	defer server.Server.Close()

	apiClient := api.InitialiseAPI(server.Server.URL)
	console := utils.NewMockConsole()

	expectedOutput := `name      state   description
mystream  enabled A stream which I use to...
mystream2 enabled Another stream to...

Total:2
`

	err := GetStreams("", apiClient, console)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())

}

func TestMissingStreamNameFlagReturnsBadRequest(t *testing.T) {
	//Given...

	body := `{"error_code": 2505,"error_message": "GAL2505I: The stream name provided by the --name field cannot be an empty string."}`

	getStreamsInteraction := utils.NewHttpInteraction("/streams", http.MethodGet)
	getStreamsInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("ClientApiVersion", "myVersion")
		writer.WriteHeader(http.StatusBadRequest) // It's going to fail with an error on purpose !
		writer.Write([]byte(body))
	}

	interactions := []utils.HttpInteraction{
		getStreamsInteraction,
	}

	server := utils.NewMockHttpServer(t, interactions)
	defer server.Server.Close()

	apiClient := api.InitialiseAPI(server.Server.URL)

	console := utils.NewMockConsole()
	expectedOutput := `GAL2505I: The stream name provided by the --name field cannot be an empty string.`

	//When
	err := GetStreams("     ", apiClient, console)

	//Then
	assert.NotNil(t, err)
	assert.Equal(t, expectedOutput, err.Error())
}

func TestMultipleStreamByNameGetFormatsResultsOk(t *testing.T) {

	//Given..
	var streamName = "mystream"

	body := `
{
  "apiVersion": "galasa-dev/v1alpha1",
  "kind": "GalasaStream",
  "metadata": {
    "name": "mystream",
    "url": "http://localhost:8080/api/streams/myStream",
    "description": "A stream which I use to..."
  },
  "data": {
    "isEnabled": true,
    "repository": {
      "url": "http://mymavenrepo.host/testmaterial"
    },
    "obrs": [
      {
        "groupId": "com.ibm.zosadk.k8s",
        "artifactId": "com.ibm.zosadk.k8s.obr",
        "version": "0.1.0-SNAPSHOT"
      }
    ],
    "testCatalog": {
      "url": "http://mymavenrepo.host/testmaterial/com.ibm.zosadk.k8s/com.ibm.zosadk.k8s.obr/0.1.0-SNAPSHOT/testcatalog.yaml"
    }
  }
}
`

	getStreamsInteraction := utils.NewHttpInteraction("/streams/"+streamName, http.MethodGet)
	getStreamsInteraction.WriteHttpResponseFunc = func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("ClientApiVersion", "myVersion")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(body))
	}

	interactions := []utils.HttpInteraction{
		getStreamsInteraction,
	}

	server := utils.NewMockHttpServer(t, interactions)
	defer server.Server.Close()

	apiClient := api.InitialiseAPI(server.Server.URL)
	console := utils.NewMockConsole()

	expectedOutput := `name     state   description
mystream enabled A stream which I use to...

Total:1
`

	err := GetStreams(streamName, apiClient, console)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, console.ReadText())

}
