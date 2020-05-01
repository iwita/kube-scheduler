package customcache

import (
	"sync"
	"time"
)

var Τimeout *time.Ticker

type MlabCache struct {
	Cache map[string]map[string]float64
	Mux   sync.Mutex
}

var LabCache MlabCache = MlabCache{
	Cache: map[string]map[string]float64{
		"kube-01": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-02": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-03": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-04": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-05": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-06": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-07": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-08": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
	},
}

func (c *MlabCache) CleanCache() {
	c.Mux.Lock()
	c.Cache = map[string]map[string]float64{
		"kube-01": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-02": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-03": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-04": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-05": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-06": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-07": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
		"kube-08": map[string]float64{
			"ipc":       -1,
			"mem_read":  -1,
			"mem_write": -1,
			"c6res":     -1,
		},
	}
	c.Mux.Unlock()
}

func (c *MlabCache) UpdateCache(input map[string]float64, c6res float64, nodename string) error {
	c.Mux.Lock()
	c.Cache[nodename] = map[string]float64{
		"ipc":       input["ipc"],
		"mem_read":  input["mem_read"],
		"mem_write": input["mem_write"],
		"c6res":     c6res,
	}
	c.Mux.Unlock()

	//klog.Infof("After cache update")
	c.printCached(nodename)
	return nil
}

func (c *MlabCache) AddAppMetrics(app map[string]float64, nodename string, numCores int) {
	c.Mux.Lock()
	c.Cache[nodename]["mem_read"] += app["mem_read"]
	c.Cache[nodename]["mem_write"] += app["mem_write"]
	//TODO
	// handle c6res addition
	c.Cache[nodename]["c6res"] -= (1 - app["c6res"]) / float64(100*numCores)
	if c.Cache[nodename]["c6res"] <= 0 {
		c.Cache[nodename]["c6res"] = 0.00000001
	}

	//TODO
	// handle ipc addition
	c.Mux.Unlock()
	//klog.Infof("After application metrics addition")
	c.printCached(nodename)
}

func (c *MlabCache) printCached(nodename string) {
	//klog.Infof("IPC: %v, Reads: %v,  Writes: %v, C6res: %v", c.Cache[nodename]["ipc"], c.Cache[nodename]["mem_read"],
	//c.Cache[nodename]["mem_write"], c.Cache[nodename]["c6res"])
}
