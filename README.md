# Kube-scheduler

1. The compilation of the scheduler in *kubernetes-v1.15.4/pkg/scheduler* needs Go v1.12
2. The binary file produced needs to be compiled with the flag `CGO_ENABLED=0`
3. There are 2 additional yaml configuration files needed (/etc/kubernetes/infra.yaml and /etc/kubernetes/monitoringDB.yaml)
4. The executable of the compilation (`make WHAT=cmd/kube-scheduler`) can be found in the path: _output/local/bin/....