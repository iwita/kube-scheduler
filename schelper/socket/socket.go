package socket

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type Server struct {
}

var ShellToUse = "/bin/sh"

func init() {
	space := regexp.MustCompile(`\s+`)
	nodesString := space.ReplaceAllString(getNumaNodes(), " ")
	for _, numaNode := range strings.Split(nodesString, " ") {
		c1 := "NUMA " + numaNode + " CPU(s)"
		command := "lscpu | grep '" + c1 + " ' | awk '{print $1}'"
		res, err1, err2 := ShellCommand(command)
		if err1 != "" || err2 != nil {
			//fmt.Errorf(err1, err2)
		}
		numaNodes[numaNode] = res
	}
}

func ShellCommand(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getNumaNodes() string {
	res, err1, err2 := ShellCommand("numastat | head -n 1 | cut -f2-")
	if err1 != "" || err2 != nil {
		//fmt.Errorf(err1, err2)
	}
	return res
}

var numaNodes map[string]string = make(map[string]string, 0)

// Get sockets nodes
// Initialize a map like the following:
// Numa-node-0: "0-7".
// Numa-node-1: "8-15"

func updateContainerAffinity(containerName string, socketId string) error {
	_, _, err2 := ShellCommand("sudo docker update" + containerName + "cpuset-cpus:\"" + numaNodes[socketId] + "\"")
	return err2
}

func (s *Server) HandleCpuAffinity(ctx context.Context, message *Socket) (*Response, error) {
	log.Printf("Received Socket: %v, Container: %v", message.Id, message.ContainerName)

	// Further operations
	// socketName := "node" + socket.Id
	// err := updateContainerAffinity(message.ContainerName, socketName)

	return &Response{
		Ok: true,
	}, nil
}
