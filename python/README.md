Please setup the network and install dependencies as per the instructions in the top-level README before running any of the steps below.

This directory also contains Python scripts to plot all the figures in the paper, and run simulations where necessary to do so. Instructions to plot each figure are given below. The generated plots will be stored in the **plots** directory.

Check if the virtual environment created during setup is active (you see `(env)` before your command prompt). If not, run
```
cd python
source env/bin/activate
```
Install the required dependencies and create a folder for plots
```
python3 -m pip install -r requirements.txt
mkdir -p plots
```
### Fig. 3
To generate Fig. 3, we use data from a crawl of the IPFS DHT. The following script processes the crawl data. This step takes about 8 minutes.
```
python3 process_crawl_data.py
```
Now plot the processed data.
```
python3 plot_prob_dist.py
```
### Fig. 6
#### Getting the required software
This experiment compares the performance of multiple `libp2p` versions. In the paper, we compare [libp2p v0.17](https://github.com/libp2p/go-libp2p/releases/tag/v0.17.0) and [libp2p v0.23.4](https://github.com/libp2p/go-libp2p/releases/tag/v0.23.4).

For each version you need to:
<!-- * download [Kubo](https://github.com/ipfs/kubo) -->
* Go into the Kubo folder `$ cd ../common/kubo`
* change the required `libp2p` version in the `go.mod` file. For instance, to use `libp2p v0.23.4`, open `go.mod` in a text editor and replace `github.com/libp2p/go-libp2p *` with `github.com/libp2p/go-libp2p v0.23.4`
* recompile Kubo with `make install` or `make build`


#### Collecting Data
To collect data, you need to run `./k_closests.sh` script. The script will use the kubo version available via `ipfs` command. You can change the binary location in the script. Note that the script will remove and re-create your public/private IPFS keys. If that's not what you want (e.g., you want to keep your existing IPFS identity) you need to modify the script or backup your IPFS keys before running `./k_closests.sh`. On Ubuntu, the keys are written to `~/.ipfs`. However, you might need to modify the `IPFS_HOME` variable if that's not the case on your system. 

Once you run `./k_closests.sh`, the script will create temporary `*.log` files with the closest peers found by the script and then summarize them in a single `k_closest.dat` file. The temporary files `*.log` will be removed. 

You need to manually rename the `k_closest.dat` to the libp2p version you've used. For instance, `mv k_closest.dat k_closest_kubo_0_23.dat`

#### Plotting
To plot the data, simply run `python3 k_closests.py`. The script will automatically read `./simulation_results/k_closest_kubo_0_23.dat` and `./simulation_results/k_closest_kubo_0_17.dat` files and display the graph. We provide files from our tests in the `./simulation_results/` folder.



### Fig. 7
To generate this figure using our provided results, run
```
python3 plot_attack_prob.py --output "attack_success_rate.pdf" --input "../experimentCombined/detection_results" --sybils 15 20 25 30 35 40 45
```
If you have run the experiment in **experimentCombined/** as per the instructions, then use
```
python3 plot_attack_prob.py --output "attack_success_rate.pdf" --input "../experimentCombined/experiment_results_new" --sybils 20 30 40 45
```
Use `--sybils` to choose which values of number of Sybils you see on the X-axis. Use only those values for which you have experiment results in the specified `--input` path.
### Fig. 9
Below you can find instructions to (i) run experiments to generate data to plot Figure 9 and (ii) to generate the plots.
#### Generating data
Run two sets of experiments measuring generation time for RSA and EDDSA:
```
go run measureGenerateSybilKeys rsa > timing_rsa.csv
go run measureGenerateSybilKeys eddsa > timing_eddsa.csv
```

We provide our `./simulation_results/sample_rsa.csv` and `./simulation_results/sample_eddsa.csv` files for reference (generating RSA keys takes a while).

#### Plotting
If you have run the experiment to generate data, plot the results using:
`python3 plot_key_generation_time.py`

To plot the sample results that we have provided, use:
`python3 plot_key_generation_time.py -i ./simulation_results/timing_eddsa.csv ./simulation_results/timing_rsa.csv`

### Fig. 10
To generate this figure using our provided results, run
```
python3 plot_accuracy.py --output "fp-fn.pdf" --input "../experimentCombined/detection_results"
```
If you have run the experiment in **experimentCombined/**, then use
```
python3 plot_accuracy.py --output "fp-fn.pdf" --input "../experimentCombined/experiment_results_new"
```
### Fig. 11
To generate this figure using our provided results, run
```
python3 plot_KL.py --output "KLs.pdf" --input "../experimentCombined/detection_results"
```
If you have run the experiment in **experimentCombined/**, then use
```
python3 plot_KL.py --output "KLs.pdf" --input "../experimentCombined/experiment_results_new"
```
### Fig. 12
Run simulation
```
python3 simulateHonEclKLs.py
```
Plot KL divergence values
```
python3 plot_KL_netsize.py --output "KLs_netsize.pdf"
```
### Figs. 13-15
To generate this figure using our provided results, run
```
python3 plot_mitigation_Figs13-15.py --input "../experimentCombined/mitigation_results"
```
If you have run the experiment in **experimentCombined/**, then use
```
python3 plot_mitigation_Figs13-15.py --input "../experimentCombined/experiment_results_new"
```
(Note that the detection and mitigation results we have provided are in different folders because we ran the experiments separately. But if you have run the experiment in **experimentCombined**, then you would be using the same `--input` path for Figs. 7, 10, 11 and 13-15.)
### Table II
```
python3 plotLatency.py --input "../experimentLatency/latency_results_new"
```