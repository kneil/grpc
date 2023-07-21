// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/bidirectional/frontend"
	"github.com/hashicorp/go-plugin/examples/bidirectional/multer"
	"github.com/hashicorp/go-plugin/examples/bidirectional/shared"
)

type addHelper struct{}

type configHelper struct{}

func (*addHelper) Sum(a, b int64) (int64, error) {
	return a + b, nil
}

func (*configHelper) GetConfig(s string) (string, error) {
	return "We got some config", nil
}

func main() {
	// We don't want to see the plugin logs.
	log.SetOutput(ioutil.Discard)

	// 	// We're a host. Start by launching the plugin process.
	// 	client := plugin.NewClient(&plugin.ClientConfig{
	// 		HandshakeConfig: shared.Handshake,
	// 		Plugins:         shared.PluginMap,
	// 		Cmd:             exec.Command("sh", "-c", os.Getenv("COUNTER_PLUGIN")),
	// 		AllowedProtocols: []plugin.Protocol{
	// 			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	// 	})

	// 	clientMulter := plugin.NewClient(&plugin.ClientConfig{
	// 		HandshakeConfig: multer.Handshake,
	// 		Plugins:         multer.PluginMap,
	// 		Cmd:             exec.Command("sh", "-c", os.Getenv("MULTER_PLUGIN")),
	// 		AllowedProtocols: []plugin.Protocol{
	// 			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	// 	})

	//    /* Client required to connect to plugin */
	// 	clientFrontend := plugin.NewClient(&plugin.ClientConfig{
	// 		HandshakeConfig: frontend.Handshake,
	// 		Plugins:         frontend.PluginMap,
	// 		Cmd:             exec.Command("sh", "-c", os.Getenv("FRONTEND_PLUGIN")),
	// 		AllowedProtocols: []plugin.Protocol{
	// 			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	// 	})
	// 	defer client.Kill()
	// 	defer clientMulter.Kill()
	// 	defer clientFrontend.Kill()

	// // Connect via RPC
	// rpcClient, err := client.Client()
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// // Request the plugin
	// raw, err := rpcClient.Dispense("counter")
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// For multer:

	// Connect via RPC
	// rpcClientMulter, err := clientMulter.Client()
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// // Request the plugin
	// rawMult, err := rpcClientMulter.Dispense("multer")
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// For frontend

	// rpcClientFrontend, err := clientFrontend.Client()
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// Request the plugin
	// rawFrontend, err := rpcClientFrontend.Dispense("frontend")
	// if err != nil {
	// 	fmt.Println("Error:", err.Error())
	// 	os.Exit(1)
	// }

	// We should have a Counter store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	//counter := raw.(shared.Counter)

	counter := LoadPlugin(shared.Handshake, shared.PluginMap, "COUNTER_PLUGIN").(shared.Counter)
	//multer := rawMult.(multer.Multer)

	multer := LoadPlugin(multer.Handshake, multer.PluginMap, "MULTER_PLUGIN").(multer.Multer)
	//frontend := rawFrontend.(frontend.Frontend)

	frontend := LoadPlugin(frontend.Handshake, frontend.PluginMap, "FRONTEND_PLUGIN").(frontend.Frontend)

	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "get":
		result, err := counter.Get(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		fmt.Println(result)

	case "put":
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		err = counter.Put(os.Args[1], int64(i), &addHelper{})
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

	case "mget":
		result, err := multer.Get(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		fmt.Println(result)

	case "mput":
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		err = multer.Put(os.Args[1], int64(i), &addHelper{})
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		/* Frontend */
	case "frontend":
		t, errt := frontend.Build("33")

		if errt != nil {
			fmt.Println("Error", errt.Error())
			os.Exit(3)
		}

		fmt.Println(t)

		s, err := frontend.Compile("22", &configHelper{})
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		fmt.Println(s)

	default:
		fmt.Println("Please only use 'frontend', 'get' or 'put'")
		os.Exit(1)
	}
}

func LoadPlugins() {

	plugin := LoadPlugin(shared.Handshake, shared.PluginMap, "COUNTER_PLUGIN").(shared.Counter)
	frontend := LoadPlugin(frontend.Handshake, frontend.PluginMap, "FRONTEND_PLUGIN").(frontend.Frontend)
	plugin.Get("Hello")
	frontend.Build("123")

}
func LoadPlugin(
	handshake plugin.HandshakeConfig,
	pluginMap map[string]plugin.Plugin,
	// pluginStruct interface{},
	PathVariable string) interface{} {

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv(PathVariable)),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})

	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("counter")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// pluginType := reflect.TypeOf (pluginStruct )
	// returnType := raw.(pluginStruct)
	return raw

}
