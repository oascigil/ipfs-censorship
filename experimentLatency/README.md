This folder contains the experiment to measure the latency of the `Provide` and `FindProviders` operations under the default and mitigation modes. Follow the steps below to run the experiment and collect the results. Then follow the instructions in **./python/** to read the results and generate data for Table II in the paper.

## Run experiment
In a terminal, make sure your working directory is this directory, that is `experimentCombined`.

Run IPFS daemon in parallel in a different terminal window. If you are running this on a remote machine, you can start a `tmux` or `screen` session and running the following command on one window.
```
../common/kubo/cmd/ipfs/ipfs daemon
```
Compile experiment
```
cd experiment
go build
```
Run two different experiments for no attack (0 sybils) and attack (45 sybils).
```
./main -clients 1 -runs 5 -sybils 45 -region 20 -outpath "../latency_results_new/sybil45Latency"
./main -clients 1 -runs 5 -sybils 0 -region 20 -outpath "../latency_results_new/sybil0Latency"
```
Do not change the output paths because the python script in the next step looks for exactly these filenames.
You should do at least 2 runs per client and discard the first value because in the first experiment, the client makes several queries to initializes its network size estimate leading to additional latency. For subsequent runs, this cost is absent. In all cases, the run aborts if the latency exceeds 5 minutes.

To measure the latency of finding only one provider, you need to change `dht.bucketSize` to `1` in the function `FindProviders()` in **./common/go-libp2p-kad-dht/routing.go**. Then run `go build` again, and run the following two experiments:
```
./main -clients 1 -runs 5 -sybils 45 -region 20 -outpath "../latency_results_new/sybil45LatencyProvider1"
./main -clients 1 -runs 5 -sybils 0 -region 20 -outpath "../latency_results_new/sybil0LatencyProvider1"
```
Remember to change `1` back to `dht.bucketSize` so that it does not affect the other experiments.
## Compute average latencies
Head over to the **./python/** directory and see instructions there to calculate average latencies for each call and generate data for Table II in the paper.
```
cd ../../python
python3 plotLatency.py --input "../latency_results_new"
```