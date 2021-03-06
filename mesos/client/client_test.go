// +build unit

/*
Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("Should return a new client", t, func() {
		c := NewClient("foo.example.com", "/bar", time.Duration(1))
		So(c, ShouldHaveSameTypeAs, &Client{})
	})
}

func TestClient_Fetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td, err := json.Marshal(map[string]string{
			"foo": "bar",
		})
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(td)
	}))
	defer ts.Close()

	Convey("Should fetch data from the Mesos API and return", t, func() {
		data := map[string]string{}
		host, err := extractHostFromURL(ts.URL)
		if err != nil {
			panic(err)
		}

		c := NewClient(host, "/", time.Duration(1))
		err = c.Fetch(&data)

		So(data["foo"], ShouldEqual, "bar")
		So(err, ShouldBeNil)
	})
}

func TestClient_URL(t *testing.T) {
	Convey("Should return the URL as a string", t, func() {
		c := NewClient("foo.example.com", "/bar", time.Duration(1))
		So(c.URL(), ShouldEqual, "http://foo.example.com/bar")
	})
}

func extractHostFromURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return parsed.Host, nil
}
