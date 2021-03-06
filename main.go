/**
 * Copyright (c) 2017 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net/http"
	"os"

	"github.com/mainflux/mainflux-mongodb-reader/api"
	"github.com/mainflux/mainflux-mongodb-reader/db"

	"github.com/cenkalti/backoff"
)

const (
	help string = `
Usage: mainflux-influxdb [options]
Options:
	-a, --host	Host address
	-p, --port	Port
	-m, --nats	MongoDB host
	-q, --nport	MongoDB port
	-d, --db	MongoDB database
	-h, --help	Prints this message end exits`
)

type (
	Opts struct {
		HTTPHost string
		HTTPPort string

		MongoHost     string
		MongoPort     string
		MongoDatabase string

		Help bool
	}
)

var (
	opts Opts
)

func tryMongoInit() error {
	var err error

	log.Print("Connecting to MongoDB... ")
	err = db.InitMongo(opts.MongoHost, opts.MongoPort, opts.MongoDatabase)
	return err
}

func main() {
	flag.StringVar(&opts.HTTPHost, "a", "localhost", "HTTP server address.")
	flag.StringVar(&opts.HTTPPort, "p", "7071", "HTTP server port.")
	flag.StringVar(&opts.MongoHost, "m", "localhost", "MongoDB host.")
	flag.StringVar(&opts.MongoPort, "q", "27017", "MongoDB port.")
	flag.StringVar(&opts.MongoDatabase, "d", "mainflux", "MongoDB database name.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	// MongoDb
	// Connect to MongoDB
	if err := backoff.Retry(tryMongoInit, backoff.NewExponentialBackOff()); err != nil {
		log.Fatalf("MongoDd: Can't connect: %v\n", err)
	} else {
		log.Println("OK")
	}

	// Print banner
	color.Cyan(banner)

	// Serve HTTP
	httpHost := fmt.Sprintf("%s:%s", opts.HTTPHost, opts.HTTPPort)
	http.ListenAndServe(httpHost, api.HTTPServer())
}

var banner = `
       MAINFLUX Mongo Reader 
                                      
    == Industrial IoT System ==

    Made with <3 by Mainflux Team
[w] http://mainflux.io
[t] @mainflux

`
