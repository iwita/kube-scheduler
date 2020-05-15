/*
Copyright 2020 Achilleas Tzenetopoulos.

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

package priorities

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	client "github.com/influxdata/influxdb1-client/v2"
	"k8s.io/klog"
)

var (
	nodeSelectionPriority = &CustomAllocationPriority{"NodeSelection", nodeSelectionScorer}
	//customResourcePriority = &CustomAllocationPriority{"CustomRequestedPriority", customResourceScorer}
	// LeastRequestedPriorityMap is a priority function that favors nodes with fewer requested resources.
	// It calculates the percentage of memory and CPU requested by pods scheduled on the node, and
	// prioritizes based on the minimum of the average of the fraction of requested to capacity.
	//
	// Details:
	// (cpu((capacity-sum(requested))*10/capacity) + memory((capacity-sum(requested))*10/capacity))/2
	NodeSelectionPriorityMap = nodeSelectionPriority.PriorityMap
)

func OneScorer(si scorerInput) float64 {
	return si.metrics[si.metricName]
}

func calculateWeightedAverageCores(response *client.Response,
	numberOfRows, numberOfMetrics, numberOfCores int) (map[string]float64, error) {
	// initialize the metrics map with a constant size
	metrics := make(map[string]float64, numberOfMetrics)
	rows := response.Results[0].Series[0]
	for i := 1; i < len(rows.Columns); i++ {
		//klog.Infof("Name of column %v : %v\nrange of values: %v\nnumber of rows: %v\nnumber of cores %v\n", i, rows.Columns[i], len(rows.Values), numberOfRows, numberOfCores)
		for j := 0; j < numberOfRows; j++ {
			avg := 0.0
			for k := 0; k < numberOfCores; k++ {
				val, err := rows.Values[j*numberOfCores+k][i].(json.Number).Float64()
				if err != nil {
					klog.Infof("Error while calculating %v", rows.Columns[i])
					return nil, err
				}
				//metrics[rows.Columns[i]] += val * float64(numberOfRows-j)
				avg += val / float64(numberOfCores)
			}
			metrics[rows.Columns[i]] += avg * float64(numberOfRows-j)
		}
		metrics[rows.Columns[i]] = metrics[rows.Columns[i]] / float64((numberOfRows * (numberOfRows + 1) / 2))
		//klog.Infof("%v : %v", rows.Columns[i], metrics[rows.Columns[i]])
	}
	// TODO better handling for the returning errors
	return metrics, nil
}

// This function does the following:
// 1. Queries the DB with the provided metrics and cores
// 2. Calculates and returns the weighted average of each of those metrics
func queryInfluxDbCores(metrics []string, uuid string, socket,
	numberOfRows int, cfg Config, c client.Client, cores []int) (*client.Response, error) {

	// calculate the number of rows needed
	// i.e. 20sec / 0.2s interval => 100rows
	//numberOfRows := int(float32(time) / cfg.MonitoringSpecs.TimeInterval)
	// EDIT
	// This time we will fetch data for multiple cores
	// so we will need more rows, proportional to the core number
	// merge all the required columns
	columns := strings.Join(metrics, ", ")

	// build the cores part of the command
	var coresPart strings.Builder
	fmt.Fprintf(&coresPart, "core_id='%d'", cores[0])
	for i := 1; i < len(cores); i++ {
		fmt.Fprintf(&coresPart, " or core_id='%d'", cores[i])
	}

	// build the coommand
	var command strings.Builder
	fmt.Fprintf(&command, "SELECT %s from core_metrics where uuid = '%s' and socket_id='%d' and %s order by time desc limit %d", columns, uuid, socket, coresPart.String(), numberOfRows*len(cores))
	//klog.Infof("The query is: %v", command.String())
	q := client.NewQuery(command.String(), cfg.Database.Name, "")
	response, err := c.Query(q)
	if err != nil {
		klog.Infof("Error while executing the query: %v", err.Error())
		return nil, err
	}
	// Calculate the average for the metrics provided
	return response, nil
}

func nodeSelectionScorer(nodeName string) (float64, error) {
	//return (customRequestedScore(requested.MilliCPU, allocable.MilliCPU) +
	//customRequestedScore(requested.Memory, allocable.Memory)) / 2

	//read database information
	var cfg Config
	err := readFile(&cfg, "/etc/kubernetes/scheduler-monitoringDB.yaml")
	if err != nil {
		return 0, err
	}

	/*-------------------------------------
	//TODO read also nodes to uuid mappings for EVOLVE
	-------------------------------------*/

	// InfluxDB
	c, err := connectToInfluxDB(cfg)
	if err != nil {
		return 0, err
	}
	// close the connection in the end of execution
	defer c.Close()

	//Get the uuid of this node in order to query in the database
	curr_uuid, ok := Nodes[nodeName]
	socket, _ := Sockets[nodeName]
	cores, _ := Cores[nodeName]
	if len(cores) == 0 {
		return 0.0, nil
	}

	//klog.Infof("Node %v has %v cores", nodeName, len(cores))
	if ok {

		metrics := []string{"c6res"}
		time := 20

		numberOfRows := int(float32(time) / cfg.MonitoringSpecs.TimeInterval)
		// Select Socket
		r, err := queryInfluxDbCores(metrics, curr_uuid, socket, numberOfRows, cfg, c, cores)
		if err != nil {
			klog.Infof("Error in querying or calculating average: %v", err.Error())
			return 0, nil
		}

		results, err := calculateWeightedAverageCores(r, numberOfRows, len(metrics), len(cores))

		res := calculateScore(scorerInput{metricName: "c6res", metrics: results}, OneScorer) * float64(len(cores))

		// Select Node

		klog.Infof("Node name %s, Score %v\n", nodeName, res)
		return res, nil
	} else {
		klog.Infof("Error finding the uuid: %v", ok)
		return 0, nil
	}
}
