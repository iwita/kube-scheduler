package priorities

import (
	"os"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

var NovaNodes InfraConfig
var NameToNode map[string]Node

func init() {
	err := readInfra(&NovaNodes, "/etc/kubernetes/infra.yaml")
	if err != nil {
		klog.Infof("Error while reading nodes' info")
	}

	// Build a map [name : Node]
	NameToNode = make(map[string]Node)
	klog.Infof("Nodes: %v", NovaNodes.Nodes)
	for _, node := range NovaNodes.Nodes {
		name := node.Name
		temp := Node{
			Uuid:    node.Uuid,
			Sockets: make([]Socket, 0),
		}
		for _, socket := range node.Sockets {
			s := Socket{
				Id:    socket.ID,
				Cores: make([]Core, 0),
			}
			//temp.Sockets = append(temp.Sockets, s)
			for _, core := range socket.Cores {
				c := Core{
					Id: core.ID,
				}
				s.Cores = append(s.Cores, c)
			}
			temp.Sockets = append(temp.Sockets, s)

		}
		NameToNode[name] = temp
	}
	klog.Infof("Read from yaml: %v", NameToNode["ns51"])
}

type Socket struct {
	Cores []Core
	Id    int
}

type Core struct {
	Id int
}

type Node struct {
	Sockets        []Socket
	Name           string
	Uuid           string
	ThreadsPerCore int
	MaxGHz         float64
	L1DCache       int
	L1ICache       int
	L2Cache        int
	L3Cache        int
}

type scorerInput struct {
	metricName string
	metrics    map[string]float64
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Agent struct {
		Port string `yaml:"port"`
	}
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
		Uuid           string  `yaml:"uuid"`
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
	"ns50": "a5822889-cf81-4db6-b5ae-8180dc90e9b6",
	"ns51": "506c38e1-cdd7-4bfc-b9c9-18a21caf16f8",
	"ns54": "69032188-ac59-4896-adc6-806346b0d1a5",
	"ns55": "c34c2ecb-6291-4bf2-b018-1435e3e9534b",
	"ns56": "521b6728-7856-4275-98bc-f05515271371",
	"ns57": "e88ee014-f124-4b43-8e80-48e7fc06f0f7",
	"ns58": "4a51cf0b-5540-4366-a064-376e1c3a3385",
	"ns59": "fef31bb5-6a32-4c42-ac04-d6f9334bcfc2",
	"ns60": "4d3cc778-9213-40c9-9735-cfd98fa26e5c",
	"ns61": "dded9f12-47e9-44a8-8d6c-5d189a4f3d32",
	"ns62": "eca73b50-ce7f-4ae8-88ac-9f53843cf239",
	"ns63": "f3160faf-27af-4822-93d4-662b3610dd3b",
	"ns64": "24ec9b01-3781-4d10-a148-3b1f999ecc22",
	"ns65": "1804b1f7-ed47-4ebd-9f0d-4e942dcb47e4",
	"ns66": "820d2305-9d1d-4098-a9be-0240a717f067",
}

var links = map[string][]float32{
	"e77467ad-636e-4e7e-8bc9-53e46ae51da1": []float32{3, 10.4},
	"c4766d29-4dc1-11ea-9d98-0242ac110002": []float32{2, 9.6},
}

var maxSpeed = map[string]float32{
	"e77467ad-636e-4e7e-8bc9-53e46ae51da1": 2.2,
	"c4766d29-4dc1-11ea-9d98-0242ac110002": 2.0,
}

// var Sockets = map[string]int{
// 	"kube-01": 1,
// 	"kube-02": 0,
// 	"kube-03": 0,
// 	"kube-04": 1,
// 	"kube-05": 0,
// 	"kube-06": 1,
// 	"kube-07": 0,
// 	"kube-08": 1,
// }

// type Core struct {
// 	ServerName string
// 	SocketId   int
// 	CoreId     int
// }

// type Socket struct {
// 	ServerName string
// 	SocketId   int
// 	Cores      []Core
// }

// type Node struct {
// 	Name           string
// 	L1dCache       int
// 	L1iCache       int
// 	L2Cache        int
// 	L3Cache        int
// 	UUid           string
// 	Sockets        []Socket
// 	ThreadsPerCore int
// }

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

// var Cores = map[string][]int{
// 	"kube-01": []int{20, 21, 22, 23},
// 	"kube-02": []int{2, 3, 4, 5, 6, 7, 8, 9},
// 	"kube-03": []int{40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55},
// 	"kube-04": []int{24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75},
// 	"kube-05": []int{0, 1, 2, 3},
// 	"kube-06": []int{12, 13, 14, 15, 16, 17, 18, 19},
// 	"kube-07": []int{4, 5, 6, 7, 8, 9, 10, 11, 24, 25, 26, 27, 28, 29, 30, 31},
// 	"kube-08": []int{20, 21, 22, 23, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
// }

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
