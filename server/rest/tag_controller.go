//Copyright 2017 Huawei Technologies Co., Ltd
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package rest

import (
	pb "github.com/servicecomb/service-center/server/core/proto"
	"github.com/servicecomb/service-center/util"
	"github.com/servicecomb/service-center/util/rest"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TagService struct {
	//
}

func (this *TagService) URLPatterns() []rest.Route {
	return []rest.Route{
		{rest.HTTP_METHOD_POST, "/registry/v3/microservices/:serviceId/tags", this.AddTags},
		{rest.HTTP_METHOD_PUT, "/registry/v3/microservices/:serviceId/tags/:key", this.UpdateTag},
		{rest.HTTP_METHOD_GET, "/registry/v3/microservices/:serviceId/tags", this.GetTags},
		{rest.HTTP_METHOD_DELETE, "/registry/v3/microservices/:serviceId/tags/:key", this.DeleteTags},
	}
}

func (this *TagService) AddTags(w http.ResponseWriter, r *http.Request) {
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.LOGGER.Error("bory err", err)
		WriteText(http.StatusInternalServerError, fmt.Sprintf("body error %s", err.Error()), w)
		return
	}
	var tags map[string]map[string]string
	err = json.Unmarshal(message, &tags)
	if err != nil {
		util.LOGGER.Error("Unmarshal error", err)
		WriteText(http.StatusBadRequest, "Unmarshal error", w)
		return
	}

	resp, err := ServiceAPI.AddTags(r.Context(), &pb.AddServiceTagsRequest{
		ServiceId: r.URL.Query().Get(":serviceId"),
		Tags:      tags["tags"],
	})
	WriteTextResponse(resp.GetResponse(), err, "", w)
}

func (this *TagService) UpdateTag(w http.ResponseWriter, r *http.Request) {
	resp, err := ServiceAPI.UpdateTag(r.Context(), &pb.UpdateServiceTagRequest{
		ServiceId: r.URL.Query().Get(":serviceId"),
		Key:       r.URL.Query().Get(":key"),
		Value:     r.URL.Query().Get("value"),
	})
	WriteTextResponse(resp.GetResponse(), err, "", w)
}

func (this *TagService) GetTags(w http.ResponseWriter, r *http.Request) {
	resp, err := ServiceAPI.GetTags(r.Context(), &pb.GetServiceTagsRequest{
		ServiceId: r.URL.Query().Get(":serviceId"),
	})
	respInternal := resp.Response
	resp.Response = nil
	WriteJsonResponse(respInternal, resp, err, w)
}

func (this *TagService) DeleteTags(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query().Get(":key")
	ids := strings.Split(keys, ",")

	resp, err := ServiceAPI.DeleteTags(r.Context(), &pb.DeleteServiceTagsRequest{
		ServiceId: r.URL.Query().Get(":serviceId"),
		Keys:      ids,
	})
	WriteTextResponse(resp.GetResponse(), err, "", w)
}
