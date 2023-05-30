// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/bidirectional/multer"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type Multer struct {
}

type data struct {
	Value int64
}

func (k *Multer) Put(key string, value int64, a multer.AddHelper) error {
	v, _ := k.Get(key)

	r, err := a.Sum(v, value)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(&data{r})
	if err != nil {
		return err
	}

	return ioutil.WriteFile("kv2_"+key, buf, 0644)
}

func (k *Multer) Get(key string) (int64, error) {
	dataRaw, err := ioutil.ReadFile("kv2_" + key)
	if err != nil {
		return 0, err
	}

	data := &data{}
	err = json.Unmarshal(dataRaw, data)
	if err != nil {
		return 0, err
	}

	return data.Value, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: multer.Handshake,
		Plugins: map[string]plugin.Plugin{
			"counter": &multer.MulterPlugin{Impl: &Multer{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
