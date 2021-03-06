/*
 * Copyright 2017 Baidu, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions
 * and limitations under the License.
 */

// media.go - the media APIs definition supported by the VCR service

// Package api defines all APIs supported by the VCR service of BCE.
package api

import (
	"encoding/json"
	"github.com/kougazhang/bce-sdk-go/bce"
	"github.com/kougazhang/bce-sdk-go/http"
)

func PutMedia(cli bce.Client, args *PutMediaArgs) error {
	jsonBytes, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return jsonErr
	}
	body, err := bce.NewBodyFromBytes(jsonBytes)
	if err != nil {
		return err
	}

	req := &bce.BceRequest{}
	req.SetUri(URI_PREFIX + MEDIA_URI)
	req.SetMethod(http.PUT)
	req.SetBody(body)

	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return err
	}
	if resp.IsFail() {
		return resp.ServiceError()
	}
	defer func() { resp.Body().Close() }()
	return nil
}

func GetMedia(cli bce.Client, source string) (*GetMediaResult, error) {
	req := &bce.BceRequest{}
	req.SetUri(URI_PREFIX + MEDIA_URI)
	req.SetMethod(http.GET)
	req.SetParam("source", source)

	// Send request and get response
	resp := &bce.BceResponse{}
	if err := cli.SendRequest(req, resp); err != nil {
		return nil, err
	}
	if resp.IsFail() {
		return nil, resp.ServiceError()
	}
	jsonBody := &GetMediaResult{}
	if err := resp.ParseJsonBody(jsonBody); err != nil {
		return nil, err
	}
	return jsonBody, nil
}
