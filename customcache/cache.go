package customcache

import (
	"strconv"
	"sync"
	"time"
)

var LabCache *MlabCache

// duration of caching (seconds)
var duration int = 10

type MlabCache struct {
	Cache   map[string]map[string]float64
	Mux     sync.Mutex
	Timeout *time.Ticker
}

func (c *MlabCache) init() {
	LabCache = &MlabCache{
		Cache:   make(map[string]map[string]float64, 0),
		Timeout: time.NewTicker(time.Duration(duration) * time.Second),
	}
}

func (c *MlabCache) CleanCache() {
	c.Mux.Lock()

	for k, v := range c.Cache {
		for key, _ := range v {
			c.Cache[k][key] = -1
		}
	}

	c.Mux.Unlock()
}

func (c *MlabCache) UpdateCache(input map[string]float64, c6res float64, socketId int32, nodename string) error {
	c.Mux.Lock()
	nn := nodename + strconv.Itoa(int(socketId))
	c.Cache[nn]["ipc"] = input["ipc"]
	c.Cache[nn]["mem_read"] = input["mem_read"]
	c.Cache[nn]["mem_write"] = input["mem_write"]
	c.Cache[nn]["c6res"] = c6res

	// Reset the ticker
	c.Timeout = time.NewTicker(time.Duration(duration) * time.Second)
	c.Mux.Unlock()
	c.printCached(nodename)
	return nil
}

func (c *MlabCache) AddAppMetrics(app map[string]float64, nodename string, socketId int32, numCores int) {
	c.Mux.Lock()
	nn := nodename + strconv.Itoa(int(socketId))
	c.Cache[nn]["mem_read"] += app["mem_read"]
	c.Cache[nn]["mem_write"] += app["mem_write"]
	c.Cache[nn]["c6res"] -= (100 - app["c6res"]) / float64(100*numCores)
	if c.Cache[nn]["c6res"] <= 0 {
		c.Cache[nn]["c6res"] = 0.00000001
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
