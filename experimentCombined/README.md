`cd` into this directory, that is `experimentCombined`

## Run experiment
Run IPFS daemon in parallel during the experiments. If you are running this on a remote machine, we recommend starting a `tmux` or `screen` session and running the following command on one window.
```
../common/kubo/cmd/ipfs/ipfs daemon
```
Compile experiment
```
cd experiment
go build
```
Run experiment (your working directory is still **experimentCombined/experiment**)
```
./main -cids 10 -clients 5 -sybils 20 -region 20 -outpath "../experiment_results_new"
```
The arguments are:
* `-cids`: number of CIDs to eclipse for each client (default 1)
* `-clients`: number of clients from which to query the eclipsed CID (default 1)
* `-sybils`: number of Sybils to generate in each attack (default 45)
* `-region`: size of query region in mitigation, specified as expected number of honest peers in the region (default 20 is the value used in all experiments)
* `-outpath`: path to store output of the experiment (specify a directory without a trailing "/")

The total number of CIDs that will be attacked and tested is `clients*cids`. For each CID, the following steps are done: launch attack, test attack detection, provide content, find providers and test if attack is successful, provide content again with mitigation, find providers with mitigation and test if mitigation is successful. All the data for each client is then written to a json file in the directory `outpath`/sybil`X`Combined where `X` is the number of Sybils. For each CID, the experiment takes about 5 minutes.

Repeat the experiment for different number of Sybils to generate all the required data. Use `-sybils 0` to run an experiment with no attack.