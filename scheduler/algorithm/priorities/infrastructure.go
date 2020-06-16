package priorities

var NovaNodes InfraConfig

func init() {
	err := readInfra(&NovaNodes, "/etc/kubernetes/infra.yaml")
	if err != nil {
		klog.Infof("Error while reading nodes' info")
	}
}

import (
	"os"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)




type scorerInput struct {
	metricName string
	metrics    map[string]float64
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Type     string `yaml:"type"`
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	MonitoringSpecs struct {
		TimeInterval float32 `yaml:"interval"`
	} `yaml:"monitoring"`
}

type InfraConfig struct {
	Nodes []struct {
		Name           string  `yaml:"name"`
		ThreadsPerCore int     `yaml:"threadsPerCore"`
		MaxGHz         float64 `yaml:"maxGHz"`
		L1DCache       int     `yaml:"l1dCache"`
		L1ICache       int     `yaml:"l1iCache"`
		L2Cache        int     `yaml:"l2Cache"`
		L3Cache        int     `yaml:"l3Cache"`
		Sockets        []struct {
			ID    int `yaml:"id"`
			Cores []struct {
				ID int `yaml:"id"`
			} `yaml:"cores"`
		} `yaml:"sockets"`
	} `yaml:"nodes"`
}

type Application struct {
	Metrics  map[string]float64
	Duration time.Duration
}

var Applications = map[string]Application{
	"scikit-lasso": Application{
		Metrics: map[string]float64{
			"ipc":       1.87,
			"mem_read":  0.1753,
			"mem_write": 0.008856,
			"c6res":     0.003058,
		},
		Duration: 69 * time.Second,
	},
	"scikit-ada": Application{
		Metrics: map[string]float64{
			"ipc":       1.10,
			"mem_read":  0.09868,
			"mem_write": 0.00669,
			"c6res":     0,
		},
		Duration: 138 * time.Second,
	},
	"scikit-rfr": Application{
		Metrics: map[string]float64{
			"ipc":       1.25,
			"mem_read":  0.0228,
			"mem_write": 0.00503,
			"c6res":     0,
		},
		Duration: 115 * time.Second,
	},
	"scikit-rfc": Application{
		Metrics: map[string]float64{
			"ipc":       1.802,
			"mem_read":  0.02423,
			"mem_write": 0.010603,
			"c6res":     0,
		},
		Duration: 38 * time.Second,
	},
	"scikit-linregr": Application{
		Metrics: map[string]float64{
			"ipc":       1.9464,
			"mem_read":  0.040475,
			"mem_write": 0.01974,
			"c6res":     0.00928149,
		},
		Duration: 45 * time.Second,
	},
	"scikit-lda": Application{
		Metrics: map[string]float64{
			"ipc":       1.9162,
			"mem_read":  0.0541,
			"mem_write": 0.029381,
			"c6res":     0.003805,
		},
		Duration: 53 * time.Second,
	},
	"cloudsuite-data-serving-client": Application{
		Metrics: map[string]float64{
			"ipc":       0.6619,
			"mem_read":  0,
			"mem_write": 0,
			"c6res":     44.48,
		},
		Duration: 72 * time.Second,
	},
	"cloudsuite-in-memory-analytics": Application{
		Metrics: map[string]float64{
			"ipc":       1.3399,
			"mem_read":  0.0052142,
			"mem_write": 0.61361,
			"c6res":     3.76196,
		},
		Duration: 60 * time.Second,
	},
	"cloudsuite-web-serving-client": Application{
		Metrics: map[string]float64{
			"ipc":       0.6619,
			"mem_read":  0,
			"mem_write": 0,
			"c6res":     44.48,
		},
		Duration: 203 * time.Second,
	},
	"spec-sphinx": Application{
		Metrics: map[string]float64{
			"ipc":       2.035,
			"mem_read":  0.0042372,
			"mem_write": 0.0021131,
			"c6res":     0.07497,
		},
		Duration: 592 * time.Second,
	},
	"spec-cactus": Application{
		Metrics: map[string]float64{
			"ipc":       1.353,
			"mem_read":  0.07105,
			"mem_write": 0.0273161,
			"c6res":     0.0532267,
		},
		Duration: 780 * time.Second,
	},
	"spec-astar": Application{
		Metrics: map[string]float64{
			"ipc":       0.86314,
			"mem_read":  0.0063,
			"mem_write": 0.0032874,
			"c6res":     0.09115,
		},
		Duration: 468 * time.Second,
	},
	"spec-leslie": Application{
		Metrics: map[string]float64{
			"ipc":       1.5225,
			"mem_read":  0.3221,
			"mem_write": 0.1532,
			"c6res":     0.1215,
		},
		Duration: 378 * time.Second,
	},
}

var NodesToUuid = map[string]string{
	"ns50": "871e6bc6-af0a-11ea-81b7-0800383a77be",
	"ns51": "8745a5ba-af0a-11ea-bfab-0800383a77b5",
	"ns54": "8761281c-af0a-11ea-a8eb-0800383e2631",
	"ns55": "877ee97e-af0a-11ea-8d62-0800383e22d1",
	"ns56": "8795d7c4-af0a-11ea-8fac-0800383e2ec5",
	"ns57": "87aedf94-af0a-11ea-912d-080038b27c19",
	"ns58": "87ca3f00-af0a-11ea-83d1-0800383e2dbd",
	"ns59": "87e5a38a-af0a-11ea-94e7-0800383e2e59",
	"ns60": "8803abbe-af0a-11ea-8f4a-0800383e1e39",
	"ns61": "8820131c-af0a-11ea-b716-0800383e1e45",
	"ns62": "883ba30c-af0a-11ea-817b-0800383e20e5",
	"ns63": "8859ae88-af0a-11ea-8cbf-0800383e206d",
	"ns64": "887e06e8-af0a-11ea-9db3-0025909c1c9c",
	"ns65": "889dd770-af0a-11ea-87ec-0025909c1cac",
	"ns66": "88b8b0ae-af0a-11ea-b3c9-080038b65204",
}

var links = map[string][]float32{
	"e77467ad-636e-4e7e-8bc9-53e46ae51da1": []float32{3, 10.4},
	"c4766d29-4dc1-11ea-9d98-0242ac110002": []float32{2, 9.6},
}

var maxSpeed = map[string]float32{
	"e77467ad-636e-4e7e-8bc9-53e46ae51da1": 2.2,
	"c4766d29-4dc1-11ea-9d98-0242ac110002": 2.0,
}

var Sockets = map[string]int{
	"kube-01": 1,
	"kube-02": 0,
	"kube-03": 0,
	"kube-04": 1,
	"kube-05": 0,
	"kube-06": 1,
	"kube-07": 0,
	"kube-08": 1,

}

type Core struct {
	ServerName string
	SocketId int
	CoreId int
}

type Socket struct {
	ServerName string
	SocketId int
	Cores []Core
}

type Node struct {
	Name string
	L1dCache int
	L1iCache int
	L2Cache int
	L3Cache int
	UUid string
	Sockets []Socket
	ThreadsPerCore int
}

// var Nodes = []Node{
// 	&Node{
// 		Name: "ns50",
// 		UUid: "",
// 		L1dCache: 32,
// 		L1iCache: 32,
// 		L2Cache: 256,
// 		L3Cache: 20480,
// 		ThreadsPerCore: 1,
// 		Sockets: []Socket {
// 			&Socket{
// 				SocketId: 0,
// 				Cores: []Core {&Core {CoreId: 0},&Core {CoreId: 1},&Core {CoreId: 2},&Core {CoreId: 3},&Core {CoreId: 4},&Core {CoreId: 5},&Core {CoreId: 6},&Core {CoreId: 7}},
// 			},
// 			&Socket {
// 				SocketId: 1,
// 				Cores: []Core {&Core {CoreId: 8},&Core {CoreId: 9},&Core {CoreId: 10},&Core {CoreId: 11},&Core {CoreId: 12},&Core {CoreId: 13},&Core {CoreId: 14},&Core {CoreId: 15}},
// 			},
// 		},
// 	},
// 	&Node{
// 		Name: "ns51",
// 		UUid: "",
// 		L1dCache: 32,
// 		L1iCache: 32,
// 		L2Cache: 256,
// 		L3Cache: 20480,
// 		ThreadsPerCore: 1,
// 		Sockets: []Socket {
// 			&Socket{
// 				SocketId: 0,
// 				Cores: []Core {&Core {CoreId: 0},&Core {CoreId: 1},&Core {CoreId: 2},&Core {CoreId: 3},&Core {CoreId: 4},&Core {CoreId: 5},&Core {CoreId: 6},&Core {CoreId: 7}},
// 			},
// 			&Socket {
// 				SocketId: 1,
// 				Cores: []Core {&Core {CoreId: 8},&Core {CoreId: 9},&Core {CoreId: 10},&Core {CoreId: 11},&Core {CoreId: 12},&Core {CoreId: 13},&Core {CoreId: 14},&Core {CoreId: 15}},
// 			},
// 		},
// 	}
// 	&Node{
// 		Name: "ns50",
// 		UUid: "",
// 		Sockets: []Socket {
// 			&Socket{
// 				SocketId: 0,
// 				Cores: []Core {&Core {CoreId: 0},&Core {CoreId: 1},&Core {CoreId: 2},&Core {CoreId: 3},&Core {CoreId: 4},&Core {CoreId: 5},&Core {CoreId: 6},&Core {CoreId: 7}},
// 			},
// 			&Socket {
// 				SocketId: 1,
// 				Cores: []Core {&Core {CoreId: 8},&Core {CoreId: 9},&Core {CoreId: 10},&Core {CoreId: 11},&Core {CoreId: 12},&Core {CoreId: 13},&Core {CoreId: 14},&Core {CoreId: 15}},
// 			},
// 		},
// 	}
// }


var Cores = map[string][]int{
	"kube-01": []int{20, 21, 22, 23},
	"kube-02": []int{2, 3, 4, 5, 6, 7, 8, 9},
	"kube-03": []int{40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55},
	"kube-04": []int{24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75},
	"kube-05": []int{0, 1, 2, 3},
	"kube-06": []int{12, 13, 14, 15, 16, 17, 18, 19},
	"kube-07": []int{4, 5, 6, 7, 8, 9, 10, 11, 24, 25, 26, 27, 28, 29, 30, 31},
	"kube-08": []int{20, 21, 22, 23, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
}

func readFile(cfg *Config, file string) error {
	f, err := os.Open(file)
	if err != nil {
		klog.Infof("Config file for scheduler not found. Error: %v", err)
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		klog.Infof("Unable to decode the config file. Error: %v", err)
		return err
	}
	return nil
}
func readInfra(cfg *InfraConfig, file string) error {
	f, err := os.Open(file)
	if err != nil {
		klog.Infof("Config file for scheduler not found. Error: %v", err)
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		klog.Infof("Unable to decode the config file. Error: %v", err)
		return err
	}
	return nil
}

func connectToInfluxDB(cfg Config) (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://" + cfg.Server.Host + ":" + cfg.Server.Port + "",
	})
	if err != nil {
		klog.Infof("Error while connecting to InfluxDB: %v ", err.Error())
		return nil, err
	}
	klog.Infof("Connected Successfully to InfluxDB")
	return c, nil

}
