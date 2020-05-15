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
	"github.com/iwita/kube-scheduler/customcache"
	"k8s.io/klog"
)

var (
	customResourcePriority = &CustomAllocationPriority{"CustomResourceAllocation", customResourceScorer}
	//customResourcePriority = &CustomAllocationPriority{"CustomRequestedPriority", customResourceScorer}
	// LeastRequestedPriorityMap is a priority function that favors nodes with fewer requested resources.
	// It calculates the percentage of memory and CPU requested by pods scheduled on the node, and
	// prioritizes based on the minimum of the average of the fraction of requested to capacity.
	//
	// Details:
	// (cpu((capacity-sum(requested))*10/capacity) + memory((capacity-sum(requested))*10/capacity))/2
	CustomRequestedPriorityMap = customResourcePriority.PriorityMap
)

func customScoreFn(si scorerInput) float64 {
	return si.metrics["ipc"] / (si.metrics["mem_read"] + si.metrics["mem_write"])
}

func onlyIPC(metrics map[string]float64) float64 {
	return metrics["ipc"]
}

func onlyL3(metrics map[string]float64) float64 {
	return 1 / metrics["l3m"]
}

func onlyNrg(metrics map[string]float64) float64 {
	return 1 / metrics["procnrg"]
}

func calculateScore(si scorerInput,
	logicFn func(scorerInput) float64) float64 {

	res := logicFn(si)
	//klog.Infof("Has score (in float) %v\n", res)

	return res
}

func calculateWeightedAverage(response *client.Response,
	numberOfRows, numberOfMetrics int) (map[string]float64, error) {
	// initialize the metrics map with a constant size
	metrics := make(map[string]float64, numberOfMetrics)
	rows := response.Results[0].Series[0]
	for i := 1; i < len(rows.Columns); i++ {
		for j := 0; j < numberOfRows; j++ {
			val, err := rows.Values[j][i].(json.Number).Float64()
			if err != nil {
				klog.Infof("Error while calculating %v", rows.Columns[i])
				return nil, err
			}
			metrics[rows.Columns[i]] += val * float64(numberOfRows-j)
		}
		metrics[rows.Columns[i]] = metrics[rows.Columns[i]] / float64((numberOfRows * (numberOfRows + 1) / 2))
		//klog.Infof("%v : %v", rows.Columns[i], metrics[rows.Columns[i]])
	}
	// TODO better handling for the returning errors
	return metrics, nil
}

func customScoreInfluxDB(metrics []string, uuid string, socket,
	numberOfRows int, cfg Config, c client.Client) (map[string]float64, error) {

	// calculate the number of rows needed
	// i.e. 20sec / 0.5s interval => 40rows
	//numberOfRows := int(float32(time) / cfg.MonitoringSpecs.TimeInterval)
	// merge all the required columns
	columns := strings.Join(metrics, ", ")
	// build the coommand
	var command strings.Builder
	fmt.Fprintf(&command, "SELECT %s from socket_metrics where uuid = '%s' and socket_id='%d' order by time desc limit %d", columns, uuid, socket, numberOfRows)
	//klog.Infof("%s", command.String())
	//q := client.NewQuery("select ipc from system_metrics", "evolve", "")
	q := client.NewQuery(command.String(), cfg.Database.Name, "")
	response, err := c.Query(q)
	if err != nil {
		klog.Infof("Error while executing the query: %v", err.Error())
		return nil, err
	}

	// Calculate the average for the metrics provided
	return calculateWeightedAverage(response, numberOfRows, len(metrics))
}

func InvalidateCache() {
	// Check if the cache needs update
	select {
	// clean the cache if 10 seconds are passed
	case <-customcache.LabCache.Timeout.C:
		klog.Infof("Time to erase")
		klog.Infof("Cache: %v", customcache.LabCache.Cache)
		//customcache.LabCache.Timeout.Stop()
		customcache.LabCache.CleanCache()

	default:
	}
}

func customResourceScorer(nodeName string) (float64, error) {

	//InvalidateCache()
	//klog.Infof("The value of the Ticker: %v", customcache.LabCache.Timeout.C)
	cores, _ := Cores[nodeName]

	var results map[string]float64
	// Check the cache
	customcache.LabCache.Mux.Lock()
	ipc, ok := customcache.LabCache.Cache[nodeName]["ipc"]
	if !ok {
		klog.Infof("IPC is nil")
	}
	reads, ok := customcache.LabCache.Cache[nodeName]["mem_read"]
	if !ok {
		klog.Infof("Memory Reads is nil")
	}
	writes, ok := customcache.LabCache.Cache[nodeName]["mem_write"]
	if !ok {
		klog.Infof("Memory Writes is nil")
	}
	c6res, ok := customcache.LabCache.Cache[nodeName]["c6res"]
	if !ok {
		klog.Infof("C6 state is nil")
	}
	customcache.LabCache.Mux.Unlock()
	// If the cache has value use it
	if ipc != -1 && reads != -1 && writes != -1 && c6res != -1 {
		results := map[string]float64{
			"ipc":       ipc,
			"mem_read":  reads,
			"mem_write": writes,
			"c6res":     c6res,
		}

		klog.Infof("Found in the cache: ipc: %v, reads: %v, writes: %v", ipc, reads, writes)
		res := calculateScore(scorerInput{metrics: results}, customScoreFn)

		if sum := c6res * float64(len(cores)); sum < 1 {
			//klog.Infof("Average C6 is less than 1, so we get: %v", average["c6res"])
			res = res * c6res
		} else {
			res = res * 1
		}

		//Apply heterogeneity
		speed := links[Nodes[nodeName]][0] * links[Nodes[nodeName]][1]
		maxFreq := maxSpeed[Nodes[nodeName]]
		res = res * float64(speed) * float64(maxFreq)

		// Select Node

		klog.Infof("Using the cached values, Node name %s, has score %v\n", nodeName, res)
		return res, nil
	}

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
	// cores, _ := Cores[nodeName]
	var socketNodes []string

	if ok {

		metrics := []string{"c6res"}
		time := 20

		numberOfRows := int(float32(time) / cfg.MonitoringSpecs.TimeInterval)

		// Define Core availability
		// Select all the nodes belonging to the current socket
		for kubenode, s := range Sockets {
			if s == socket && Nodes[kubenode] == curr_uuid {
				socketNodes = append(socketNodes, kubenode)
			}
		}
		sum := 0
		socketSum := 0.0
		socketCores := 0.0
		for _, snode := range socketNodes {
			currCores, _ := Cores[snode]
			r, err := queryInfluxDbCores(metrics, curr_uuid, socket, numberOfRows, cfg, c, currCores)
			//r, err := queryInfluxDbSocket(metrics, curr_uuid, socket, numberOfRows, cfg, c)
			if err != nil {
				klog.Infof("Error in querying or calculating core availability in the first stage: %v", err.Error())
			}
			average, err := calculateWeightedAverageCores(r, numberOfRows, len(metrics), len(currCores))
			if err != nil {
				klog.Infof("Error defining core availability")
			}
			if average["c6res"]*float64(len(currCores)) >= 1 {
				klog.Infof("Node %v has C6 sum: %v", snode, average["c6res"]*float64(len(currCores)))
				sum++
			}
			socketSum += average["c6res"] * float64(len(currCores))
			socketCores += float64(len(currCores))
		}

		// Select Socket
		results, err = customScoreInfluxDB([]string{"ipc", "mem_read", "mem_write"}, curr_uuid, socket, numberOfRows, cfg, c)
		if err != nil {
			klog.Infof("Error in querying or calculating average for the custom score in the first stage: %v", err.Error())
			return 0, nil
		}

		res := calculateScore(scorerInput{metrics: results}, customScoreFn)

		//klog.Infof("Node: %v\t res before: %v", nodeName, res)

		if sum < 1 {
			klog.Infof("Less than 1 node is available\nC6contribution: %v", socketSum/socketCores)
			res = res * socketSum / socketCores
		} else {
			res = res * 1
		}

		//Update the cache with the new metrics
		err = customcache.LabCache.UpdateCache(results, socketSum/socketCores, nodeName)
		if err != nil {
			klog.Infof(err.Error())
		} else {
			klog.Infof("Cache updated successfully for %v", nodeName)
		}

		//Apply heterogeneity
		speed := links[Nodes[nodeName]][0] * links[Nodes[nodeName]][1]
		maxFreq := maxSpeed[Nodes[nodeName]]
		res = res * float64(speed) * float64(maxFreq)

		// Select Node

		klog.Infof("Node name %s, has score %v\n", nodeName, res)
		return res, nil
	} else {
		klog.Infof("Error finding the uuid: %v", ok)
		return 0, nil
	}

}

// WARNING
// c6res is not a dependable metric for isnpecting core availability
// Some Systems use higher core states (e.g c7res)
// func findAvailability(response *client.Response, numberOfMetrics, numberOfRows, numberOfCores int, floor float64) (map[string]float64, error) {
// 	// initialize the metrics map with a constant size
// 	metrics := make(map[string]float64, numberOfMetrics)
// 	rows := response.Results[0].Series[0]
// 	for i := 1; i < len(rows.Columns); i++ {
// 		//klog.Infof("Name of column %v : %v\nrange of values: %v\nnumber of rows: %v\nnumber of cores %v\n", i, rows.Columns[i], len(rows.Values), numberOfRows, numberOfCores)
// 		for j := 0; j < numberOfRows; j++ {
// 			//avg, max := 0.0, 0.0
// 			for k := 0; k < numberOfCores; k++ {
// 				val, err := rows.Values[j*numberOfCores+k][i].(json.Number).Float64()
// 				if err != nil {
// 					klog.Infof("Error while calculating %v", rows.Columns[i])
// 					return false, err
// 				}
// 				// if val > floor {
// 				// 	return true, nil
// 				// }
// 				// sum += val

// 				//avg += val / float64(numberOfCores)
// 				avg += val
// 			}
// 			metrics[rows.Columns[i]] += avg * float64(numberOfRows-j)
// 		}
// 		metrics[rows.Columns[i]] = metrics[rows.Columns[i]] / float64((numberOfRows * (numberOfRows + 1) / 2))
// 		if metrics[row.Columns[i]] > 1 {
// 			return true, nil
// 		}
// 		//klog.Infof("%v : %v", rows.Columns[i], metrics[rows.Columns[i]])
// 	}
// 	// TODO better handling for the returning errors
// 	return false, nil
// }
