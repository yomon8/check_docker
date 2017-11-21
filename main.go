package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
)

// Nagios return codes
const (
	NagiosOk       = 0
	NagiosWarning  = 1
	NagiosCritical = 2
	NagiosUnknown  = 3
)

var version string

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Critical - create docker client connection error:", err)
		os.Exit(NagiosUnknown)
	}

	info, err := cli.Info(context.Background())
	if err != nil {
		fmt.Println("Critical - get docker info error:", err)
		os.Exit(NagiosUnknown)
	}

	msg := fmt.Sprintf("Containers( %d running, %d paused, %d stopped );Images( %d images );Swarm ( Status %s, Error %s )",
		info.ContainersRunning,
		info.ContainersPaused,
		info.ContainersStopped,
		info.Images,
		info.Swarm.LocalNodeState,
		info.Swarm.Error)

	if info.Swarm.Error != "" {
		fmt.Println("Critical - Swarm error message detected:", msg)
		os.Exit(NagiosCritical)
	}

	if info.Swarm.NodeID != "" && info.Swarm.LocalNodeState != "active" {
		fmt.Println("Warning - Swarm State not active:", msg)
		os.Exit(NagiosWarning)
	}

	if info.ContainersRunning == 0 {
		fmt.Println("Warning - No containers runinng:", msg)
		os.Exit(NagiosWarning)
	}

	fmt.Println("OK - ", msg)
	os.Exit(NagiosOk)
}
