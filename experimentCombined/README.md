Please setup the network and install dependencies as per the instructions in the top-level README before running any of the steps below.

This folder contains the main experiment that runs the attack on several CIDs and DHT clients, and collects results regarding the attack, detection and mitigation. The collected results will then be used to plot Figs. 7, 10, 11, 13, 14, 15 in the paper. Follow the steps below to run the experiment and collect the results. Then follow the instructions in **./python/** to read the results and generate plots.
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
The following command is used to run the experiment (your working directory is still **experimentCombined/experiment**). First run a small version of the experiment with the following command to test that the program runs without errors.
```
./main -cids 1 -clients 1 -sybils 20 -region 20 -outpath "../experiment_results_new"
```
In general, the arguments for the program are:
* `-cids`: number of CIDs to eclipse for each client (default 1)
* `-clients`: number of clients from which to query the eclipsed CID (default 1)
* `-sybils`: number of Sybils to generate in each attack (default 45)
* `-region`: size of query region in mitigation, specified as expected number of honest peers in the region (default 20 is the value used in all experiments)
* `-outpath`: path to store results of the experiment (specify a directory without a trailing "/")

The total number of CIDs that will be attacked and tested is `clients*cids`. For each CID, the following steps are done: launch attack, test attack detection, provide content, find providers and test if attack is successful, provide content again with mitigation, find providers with mitigation and test if mitigation is successful. All the results for each client is then written to a json file in the directory `outpath`/sybil`X`Combined where `X` is the number of Sybils. For each CID, the experiment takes 4-5 minutes.

Repeat the experiment for different number of Sybils to generate all the required results. Use `-sybils 0` to run an experiment with no attack. As a shortcut, simply run `python3 run_experiments.py` to repeat the experiment for 0, 20, 30, 40, 45 Sybils with 10 CIDs and 1 client. This set of experiments is expected to take 4 hours. The results will be written to **experiment_results_new/**.