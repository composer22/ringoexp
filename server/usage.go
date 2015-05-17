package server

import (
	"fmt"
	"os"
)

const usageText = `
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
`

// end help text

// PrintUsageAndExit is used to print out command line options.
func PrintUsageAndExit() {
	fmt.Printf("%s\n", usageText)
	os.Exit(0)
}
