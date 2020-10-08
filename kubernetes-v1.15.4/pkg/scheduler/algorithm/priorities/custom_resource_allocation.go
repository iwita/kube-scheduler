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
	customcache "github.com/iwita/kube-scheduler/customcache"
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
	//klog.Infof("Query: %v, Response: %v", q, response)
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

func customResourceScorer(nodeName string) (float64, int, error) {

	//InvalidateCache()
	//klog.Infof("The value of the Ticker: %v", customcache.LabCache.Timeout.C)
	//cores, _ := Cores[nodeName]

	var results map[string]float64

	// // Check the cache
	// // In case it is nil just continue
	// customcache.LabCache.Mux.Lock()
	// ipc, ok_ipc := customcache.LabCache.Cache[nodeName]["ipc"]
	// if !ok_ipc {
	// 	klog.Infof("IPC field is nil")
	// }
	// reads, ok_read := customcache.LabCache.Cache[nodeName]["mem_read"]
	// if !ok_read {
	// 	klog.Infof("Memory Reads field is nil")
	// }
	// writes, ok_write := customcache.LabCache.Cache[nodeName]["mem_write"]
	// if !ok_write {
	// 	klog.Infof("Memory Writes field is nil")
	// }

	// // in case the cache is not nil
	// if ok_write && ok_ipc && ok_read && reads != -1 && writes != -1 && ipc != -1 {
	// 		results := map[string]float64{
	// 			"ipc":       ipc,
	// 			"mem_read":  reads,
	// 			"mem_write": writes,
	// 			//"c6res":     c6res,
	// 		}
	// }

	// socket, _ := Sockets[nodeName]
	// curr_uuid, ok := Nodes[nodeName]

	// If the cache has value use it
	// if ipc != -1 && reads != -1 && writes != -1 {

	// 	var socketNodes []string
	// 	for kubenode, s := range Sockets {
	// 		if s == socket && Nodes[kubenode] == curr_uuid {
	// 			socketNodes = append(socketNodes, kubenode)
	// 		}
	// 	}
	// 	//var float64 currentNodeC6res
	// 	socketSum := 0.0
	// 	socketCores := 0
	// 	sum := 0

	// 	for _, snode := range socketNodes {
	// 		c6res, ok := customcache.LabCache.Cache[snode]["c6res"]
	// 		if !ok {
	// 			klog.Infof("C6 state is nil")
	// 		}
	// 		if c6res*float64(len(Cores[snode])) > 1 {
	// 			sum++
	// 		}
	// 		socketSum += c6res * float64(len(Cores[snode]))
	// 		socketCores += len(Cores[snode])
	// 		// if snode == nodeName {
	// 		// 	currentNodeC6res = c6res
	// 		// }
	// 	}
	// 	customcache.LabCache.Mux.Unlock()

	// 	klog.Infof("Found in the cache: ipc: %v, reads: %v, writes: %v, c6: %v\n", ipc, reads, writes, socketSum/float64(socketCores))
	// 	results["c6res"] = socketSum / float64(socketCores)
	// 	res := calculateScore(scorerInput{metrics: results}, customScoreFn)

	// 	if sum < 1 {
	// 		//klog.Infof("Average C6 is less than 1, so we get: %v", average["c6res"])
	// 		res = res * socketSum / float64(socketCores)
	// 	} else {
	// 		res = res * 1
	// 	}

	// 	//Apply heterogeneity
	// 	speed := links[Nodes[nodeName]][0] * links[Nodes[nodeName]][1]
	// 	maxFreq := maxSpeed[Nodes[nodeName]]
	// 	res = res * float64(speed) * float64(maxFreq)

	// 	// Select Node
	// 	klog.Infof("Using the cached values, Node name %s, has score %v\n", nodeName, res)
	// 	return res, nil
	// }

	// EVOLVE

	// we have the config file in NovaNodes(type InfraConfig) variable

	//read database information
	var cfg Config
	err := readFile(&cfg, "/etc/kubernetes/scheduler-monitoringDB.yaml")
	if err != nil {
		return 0, -1, err
	}
	// Connect to InfluxDB
	c, err := connectToInfluxDB(cfg)
	if err != nil {
		return 0, -1, err
	}
	// close the connection in the end of execution
	defer c.Close()

	uuid := NameToNode[nodeName].Uuid
	klog.Infof("Node: %v, uuid: %v", nodeName, uuid)

	// query the last 'time' seconds
	time := 20
	// calculate the rows of data needed for this interval
	numberOfRows := int(float32(time) / cfg.MonitoringSpecs.TimeInterval)

	nodeSumOfFreeCores := 0
	sumOfFreeCores := 0

	// Calculate the score for each socket of the node
	maxRes := 0.0
	var winningSocket int
	for _, socket := range NameToNode[nodeName].Sockets {
		socketId := socket.Id

		// First check the cache
		nn := nodeName + "" + string(socketId)
		customcache.LabCache.Mux.Lock()
		ipc, ok_ipc := customcache.LabCache.Cache[nn]["ipc"]
		reads, ok_reads := customcache.LabCache.Cache[nn]["mem_read"]
		writes, ok_writes := customcache.LabCache.Cache[nn]["mem_write"]
		c6, ok_c6 := customcache.LabCache.Cache[nn]["c6res"]
		customcache.LabCache.Mux.Unlock()
		if ok_ipc && ok_reads && ok_writes && ipc != -1 && reads != -1 && writes != -1 {
			results = map[string]float64{
				"ipc":       ipc,
				"mem_read":  reads,
				"mem_write": writes,
				//"c6res":     c6res,
			}
		} else { // if nothing exists here read from database
			results, err = customScoreInfluxDB([]string{"ipc", "mem_read", "mem_write"}, uuid, socketId, numberOfRows, cfg, c)
			if err != nil {
				klog.Infof("Error in querying or calculating average for the custom score in the first stage: %v", err.Error())
				return 0, -1, nil
			}
		}

		klog.Infof("Node: %v, Calculating score...", nodeName)
		res := calculateScore(scorerInput{metrics: results}, customScoreFn)
		klog.Infof("Node: %v, Finished calculating score ", nodeName)

		// Check the core availability

		// Get the core ids of this socket
		currCores := []int{}
		for _, c := range socket.Cores {
			currCores = append(currCores, c.Id)
		}
		klog.Infof("Cores of this socket are: %v\n", currCores) 
		// if c6 does not exist in the database
		if ok_c6 && c6 != -1 {
			if c6*float64(len(currCores)) > 1 {
				sumOfFreeCores++
			}
			if sumOfFreeCores < 1 {
				klog.Infof("Less than 1 node is available\nC6contribution: %v", c6)
				res = res * c6
			} else {
				res = res * 1
			}
		} else {
			// Return the metrics of those cores
			metrics := []string{"c6res"} // check the c6res for now
			r, err := queryInfluxDbCores(metrics, uuid, socketId, numberOfRows, cfg, c, currCores)
			//r, err := queryInfluxDbSocket(metrics, curr_uuid, socket, numberOfRows, cfg, c)
			if err != nil {
				klog.Infof("Error in querying or calculating core availability in the first stage: %v", err.Error())
			}

			// Calculate the average of those metrics
			average, err := calculateWeightedAverageCores(r, numberOfRows, len(metrics), len(currCores))
			if err != nil {
				klog.Infof("Error defining core availability")
			}

			// Finally check the availability of the socket
			if average["c6res"]*float64(len(currCores)) > 1 {
				klog.Infof("Node's %v, socket %v has C6 sum: %v", nodeName, socketId, average["c6res"]*float64(len(currCores)))
				sumOfFreeCores++
			}
			if sumOfFreeCores < 1 {
				klog.Infof("Less than 1 node is available\nC6contribution: %v", average["c6res"])
				res = res * average["c6res"]
			} else {
				res = res * 1
			}

			// Update the cache
			klog.Infof("Node/Socket: %v/%v, Updating the cache ", nodeName, socketId)
			err = customcache.LabCache.UpdateCache(results, average["c6res"], int32(socketId), nodeName)
			klog.Infof("Node: %v, Finishing updating the cache.... ", nodeName)
			if err != nil {
				klog.Infof(err.Error())
			} else {
				klog.Infof("Cache updated successfully for node: %v, socket: %v", nodeName, socketId)
			}
		}

		if res > maxRes {
			maxRes = res
			winningSocket = socketId
		}
		nodeSumOfFreeCores += sumOfFreeCores

		// Update the cache with the new metrics

		// klog.Infof("Node/Socket: %v/%v, Updating the cache ", nodeName, socketId)
		// err = customcache.LabCache.UpdateCache(results, currentNodeC6res, nodeName)
		// klog.Infof("Node: %v, Finishing updating the cache.... ", nodeName)

		// if err != nil {
		// 	klog.Infof(err.Error())
		// } else {
		// 	klog.Infof("Cache updated successfully for %v", nodeName)
		// }

		// here we need to return the score of the best socket and its id
		// TODO
		// **************************************************************

		// Send it on a grpc server
	}

	// EVOLVE

	//TODO Apply heterogeneity
	// speed := links[Nodes[nodeName]][0] * links[Nodes[nodeName]][1]
	// maxFreq := maxSpeed[Nodes[nodeName]]
	// res = res * float64(speed) * float64(maxFreq)

	// Select Node
	klog.Infof("Node/Socket %s/%v, has score %v\n", nodeName, winningSocket, maxRes)
	return maxRes, winningSocket, nil
}
