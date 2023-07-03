`cd` into this directory, that is `experimentCombined`

## Run experiment
Run IPFS daemon in parallel
```
../common/kubo/cmd/ipfs/ipfs daemon
```
Compile experiment
```
cd experiment
go build
cd ..
```
Run experiment
```
./experiment/main -clients 5 -runs 10 -sybils 0 -region 20 -outpath "experiment_results/sybil0LatencySpecial20"
```
Use the flag `-mitigation` to enable mitigation, otherwise it is disabled. You may change the number of and clients and runs. You should do at least 2 runs per client and discard the first value because in the first experiment, the client makes several queries to initializes its network size estimate leading to additional latency. For subsequent runs, this cost is absent. You may repeat the experiment for different number of Sybils (0 and 45) to generate all the required data. Change `outpath` accordingly.

## Compute average latencies
```
python plotLatency.py
```
Change the variable `numSybils` in `plotLatency.py` to see the statistics under attack and without attack