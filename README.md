# Ringo Exp
[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/composer22/ringoexp.svg?branch=master)](http://travis-ci.org/composer22/ringoexp)
[![Current Release](https://img.shields.io/badge/release-v0.1.3-brightgreen.svg)](https://github.com/composer22/ringoexp/releases/tag/v0.1.3)
[![Coverage Status](https://coveralls.io/repos/composer22/chattypantz/badge.svg?branch=master)](https://coveralls.io/r/composer22/ringoexp?branch=master)

A ring-buffer experiment written in [Go.](http://golang.org)

## About

This is a small server that demonstrates ring-buffer usage as a means to optimize memory and processing usage.

Some key objectives in this demonstration:

* Publisher Server: as a stand-alone publisher it receives work requests from an external websocket source and places this work into a ring-buffer queue. Background workers are in parallel reading the ring-buffer queue and forward these items to a consumer server.
* Consumer Server: as a stand-alone consumer server, it receives connections from a Publisher Server and stores these objects into an internal database. This is basically a proxy over a database that can distribute I/O extensive writes.

Issues:

* TBD

## Usage

```
Description: is a small server that demonstrates ring-buffer usage as a means to optimize memory and processing usage.

Usage: ringoexp [options...]

Basic Server options:
    -N, --name NAME					NAME of the server (default: empty).
    -H, --hostname HOSTNAME         	HOSTNAME of the server (default: localhost).
    -p, --port PORT					PORT to listen on (default: 6660).
	-n, --connections MAX				MAX incoming connections allowed (default: unlimited).
	-I, --is_publisher				   	Is the server a publisher? (default: true).

Publisher Server Mode - additional options (is_publisher = true):
    -r, --ring_size SIZE			    SIZE of the incoming ring buffer (default: 4096).
    -U, --consumer_hostname HOSTNAME	HOSTNAME of the remote consumer server (default: localhost).
    -T, --consumer_port PORT			PORT of the remote consumer server (default: 6660).
    -W, --workers MAX         			MAX worker connections to the consumer (default: 1024).

System level options:
	-X, --procs MAX                  *MAX processor cores to use from the machine.
	-L, --profiler_port PORT         *PORT the profiler is listening on (default: off).
    -d, --debug                      Enable debugging output (default: false)

     *  Anything <= 0 is no change to the environment (default: 0).

Common options:
    -h, --help                       Show this message
    -V, --version                    Show version

Examples:

	# Publisher Mode:
	#
	# Name: "San Francisco"
	# Host: 0.0.0.0
	# Port: 6661
	# Consumer Port: 6662
	# Procs: 2

    ringoexp -N "San Francisco" -H 0.0.0.0 -p 6661 -T 6662 -X 2

	# Consumer Mode:
	#
	# Name: "San Francisco"
	# Host: 0.0.0.0
	# Port: 6662
	# Connections: 1280
	# Is Publisher: false

	ringoexp -N "San Francisco" -H 0.0.0.0 -p 6662 -n 1280 -I false

```
## Server Connection Specifications

The websocket connection endpoint is:
```
ws://{host:port}/v1.0/ingest
```

Note: for performance sake this is a binary packet in the socket.  Json is not used.

The fixed format of the raw message in the socket or http call is

int - 4 byte integer being sent as payload

8 bytes total

## HTTP API for Alive and Stats

Two additional API routes are provided:

* http://localhost:6660/v1.0/alive - GET: Is the server alive?
* http://localhost:6660/v1.0/stats - GET: Returns information about the server state.

For these calls, json headers are required:

* Accept: application/json
* Content-Type: application/json

Example cURL:

```
$ curl -i -H "Accept: application/json" \
-H "Content-Type: application/json" \
-X GET "http://0.0.0.0:6660/v1.0/alive"

HTTP/1.1 200 OK
Content-Type: application/json;charset=utf-8
Date: Fri, 03 Apr 2015 17:29:17 +0000
Server: San Francisco
X-Request-Id: DC8D9C2E-8161-4FC0-937F-4CA7037970D5
Content-Length: 0
```
## Building

This code currently requires version 1.42 or higher of Go.

Information on Golang installation, including pre-built binaries, is available at
<http://golang.org/doc/install>.

Run `go version` to see the version of Go which you have installed.

Run `go build` inside the directory to build.

Run `go test ./...` to run the unit regression tests.

A successful build run produces no messages and creates an executable called `ringoexp` in this
directory.

Run `go help` for more guidance, and visit <http://golang.org/> for tutorials, presentations, references and more.

## Docker Images

A prebuilt docker image is available at (http://www.docker.com) [ringoexp](https://registry.hub.docker.com/u/composer22/ringoexp/)

If you have docker installed, run:
```
docker pull composer22/ringoexp:latest

or

docker pull composer22/ringoexp:<version>

if available.
```
See /docker directory README for more information on how to run it.

## License

(The MIT License)

Copyright (c) 2015 Pyxxel Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to
deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.
