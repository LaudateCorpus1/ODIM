//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package chassis

import (
	"bytes"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
)

func Test_sourceProviderImpl_findSwitchChassis(t *testing.T) {
	config.SetUpMockConfig(t)
	col := sresponse.NewChassisCollection()
	type args struct {
		collection *sresponse.Collection
	}
	tests := []struct {
		name string
		c    *sourceProviderImpl
		args args
	}{
		{
			name: "multiple switch chassis collection available for multiple plugins",
			c: &sourceProviderImpl{
				getSwitchFactory: getSwitchFactoryMock,
			},
			args: args{
				collection: &col,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.findSwitchChassis(tt.args.collection)
			assert.Equal(t, 2, tt.args.collection.MembersCount)
		})
	}
}

func getSwitchFactoryMock(collection *sresponse.Collection) *switchFactory {
	chassisMap := make(map[string]bool)
	return &switchFactory{
		collection:        collection,
		chassisMap:        chassisMap,
		wg:                &sync.WaitGroup{},
		mu:                &sync.Mutex{},
		getFabricManagers: getFabricManagersMock,
		contactClient:     contactClientMock,
	}
}

func getFabricManagersMock() ([]smodel.Plugin, error) {
	return []smodel.Plugin{
		{
			ID:                "1",
			PreferredAuthType: "XAuthToken",
		},
		{
			ID:                "2",
			PreferredAuthType: "BasicAuth",
			Username:          "someUser",
			Password:          []byte("password"),
		},
	}, nil
}

func contactClientMock(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	tokenBody := `{"Members": [{"@odata.id":"/ODIM/v1/Chassis/1"}]}`
	basicAuthBody := `{"Members": [{"@odata.id":"/ODIM/v1/Chassis/2"}]}`
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"X-Auth-Token": []string{
				"token",
			},
		},
	}
	if token != "" {
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(tokenBody))
	} else {
		resp.Body = ioutil.NopCloser(bytes.NewBufferString(basicAuthBody))
	}
	return resp, nil
}
