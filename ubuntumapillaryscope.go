// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 JaquerEspeis
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

// Package main implements an Ubuntu scope for Mapillary street level photos.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/JaquerEspeis/mapillary"
	mapillaryapi "github.com/JaquerEspeis/mapillary/v2"
	"launchpad.net/go-unityscopes/v2"
)

// UbuntuMapillaryScope is a scope that shows Mapillary street level photos.
type UbuntuMapillaryScope struct {
}

type conf struct {
	ClientID string
}

type imgInfo struct {
	*mapillaryapi.GetSearchImRandomSelect
	ImageURL string
}

func clientID() (string, error) {
	confFile, err := os.Open("conf.json")
	if err != nil {
		return "", err
	}
	defer confFile.Close()
	var config conf
	confFileContents, err := ioutil.ReadAll(confFile)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(confFileContents, &config)
	if err != nil {
		return "", err
	}
	return config.ClientID, nil
}

// SetScopeBase is defined just to conform with the interface.
func (s *UbuntuMapillaryScope) SetScopeBase(base *scopes.ScopeBase) {
}

// Search pushes results to the scope.
func (s *UbuntuMapillaryScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	cat := reply.RegisterCategory("curated", "Curated Images", "", "")

	result := scopes.NewCategorisedResult(cat)
	info, err := s.getRandomCuratedImageInfo()
	if err != nil {
		return err
	}
	result.SetURI(info.ImageURL)
	result.SetTitle(info.Location)
	result.SetArt(info.ImageURL)
	return reply.Push(result)
}

func (s *UbuntuMapillaryScope) getRandomCuratedImageInfo() (info imgInfo, err error) {
	var response imgInfo
	id, err := clientID()
	if err != nil {
		return response, err
	}
	client := mapillaryapi.NewClient(id)
	if err := client.Get("search/im/randomselected", url.Values{}, &response); err != nil {
		return response, err
	}
	response.ImageURL, err = mapillary.GetImageURL(response.Key, 320)
	return response, err
}

// Preview shows more information about a result.
func (s *UbuntuMapillaryScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	header := scopes.NewPreviewWidget("header", "header")
	header.AddAttributeMapping("title", "title")
	header.AddAttributeMapping("subtitle", "subtitle")

	image := scopes.NewPreviewWidget("image", "image")
	image.AddAttributeMapping("source", "art")

	return reply.PushWidgets(header, image)
}

func main() {
	if err := scopes.Run(&UbuntuMapillaryScope{}); err != nil {
		log.Fatalln(err)
	}
}
