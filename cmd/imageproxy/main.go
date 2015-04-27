// Copyright 2013 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// imageproxy starts an HTTP server that proxies requests for remote images.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/etcinit/imageproxy"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
)

// goxc values
var (
	// VERSION is the version string for imageproxy.
	VERSION = "HEAD"

	// BUILD_DATE is the timestamp of when imageproxy was built.
	BUILD_DATE string
)

func main() {
	// Setup the command line application.
	app := cli.NewApp()
	app.Name = "imageproxy"
	app.Usage = "imageproxy is a caching image proxy server"

	// Set version and authorship info.
	app.Version = VERSION
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Will Norris",
			Email: "will@willnorris.com",
		},
		cli.Author{
			Name:  "Eduardo Trujillo",
			Email: "ed@chromabits.com",
		},
	}

	// Define application flags.
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: ":8080",
			Usage: "TCP address to listen on",
		},
		cli.StringFlag{
			Name:  "whitelist",
			Value: "",
			Usage: "comma separated list of allowed remote hosts",
		},
		cli.StringFlag{
			Name:  "cacheDir",
			Value: "",
			Usage: "directory to use for the file cache",
		},
		cli.IntFlag{
			Name:  "cacheSize",
			Value: 100,
			Usage: "maximum size of the file cache (in MB)",
		},
	}

	// Setup the default action. This action will be triggered when no
	// subcommand is provided as an argument.
	app.Action = func(c *cli.Context) {
		// Collect flags.
		cacheDir := c.String("cacheDir")
		cacheSize := c.Int("cacheSize")
		whitelist := c.String("whitelist")
		addr := c.String("addr")

		var cache httpcache.Cache
		if cacheDir != "" {
			d := diskv.New(diskv.Options{
				BasePath:     cacheDir,
				CacheSizeMax: uint64(cacheSize) * 1024 * 1024,
			})

			cache = diskcache.NewWithDiskv(d)
		} else {
			cache = httpcache.NewMemoryCache()
		}

		p := imageproxy.NewProxy(nil, cache)
		if whitelist != "" {
			p.Whitelist = strings.Split(whitelist, ",")
		}

		// Create the server.
		server := &http.Server{
			Addr:    addr,
			Handler: p,
		}

		// Begin listening.
		fmt.Printf("imageproxy (version %v) listening on %s\n", VERSION, server.Addr)

		if err := server.ListenAndServe(); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}

	// Begin
	app.Run(os.Args)
}
