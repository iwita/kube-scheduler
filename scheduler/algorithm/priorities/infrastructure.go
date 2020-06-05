package priorities

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

var Nodes = map[string]string{
	"kube-01": "e77467ad-636e-4e7e-8bc9-53e46ae51da1",
	"kube-02": "e77467ad-636e-4e7e-8bc9-53e46ae51da1",
	"kube-03": "e77467ad-636e-4e7e-8bc9-53e46ae51da1",
	"kube-04": "e77467ad-636e-4e7e-8bc9-53e46ae51da1",
	"kube-05": "c4766d29-4dc1-11ea-9d98-0242ac110002",
	"kube-06": "c4766d29-4dc1-11ea-9d98-0242ac110002",
	"kube-07": "c4766d29-4dc1-11ea-9d98-0242ac110002",
	"kube-08": "c4766d29-4dc1-11ea-9d98-0242ac110002",
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
