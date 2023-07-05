This directory contains simulations to calculate the ideal probability distribution of prefix length counts (Fig. 3) and to simulate detection for different network sizes (Fig. 12). This directory also contains Python scripts to plot all the figures in the paper. Instructions to plot each figure are given below.
<!-- The generated plots will be stored in the **plots** directory. -->

### Fig. 3
To generate Fig. 3, we use data from a crawl of the IPFS DHT. The following script processes the crawl data. This step takes about 8 minutes.
```
python process_crawl_data.py
```
Now plot the processed data
```
python plot_prob_dist.py
```
### Fig. 6
#### Getting the required software
This experiment compares the performance of multiple `libp2p` versions. In the paper, we compare [libp2p v0.17](https://github.com/libp2p/go-libp2p/releases/tag/v0.17.0) and [libp2p v0.23.4](https://github.com/libp2p/go-libp2p/releases/tag/v0.23.4).

For each version you need to:
* download [Kubo](https://github.com/ipfs/kubo)
* Go into the Kubo folder `$ cd kubo`
* change the required `libp2p` version in the `go.mod` file. For instance, to use `libp2p v0.23.4`, open `go.mod` in a text editor and replace `github.com/libp2p/go-libp2p *` with `github.com/libp2p/go-libp2p v0.23.4`
* recompile Kubo with `make install`


#### Collecting Data
To collect data, you need to run `./k_closests.sh` script. The script will use the kubo version available via `ipfs` command. You can change the binary location in the script. Note that the script will remove and re-create your public/private IPFS keys. If that's not what you want (e.g., you want to keep your existing IPFS identity) you need to modify the script or backup your IPFS keys before running `./k_closests.sh`. On Ubuntu, the keys are written to `~/.ipfs`. However, you might need to modify the `IPFS_HOME` variable if that's not the case on your system. 

Once you run `./k_closests.sh`, the script will create temporary `*.log` files with the closest peers found by the script and then summarize them in a single `k_closest.dat` file. The temporary files `*.log` will be removed. 

You need to manually rename the `k_closest.dat` to the libp2p version you've used. For instance, `mv k_closest.dat k_closest_kubo_0_23.dat`

#### Plotting
To plot the data, simply run `python3 k_closest.py`. The script will automatically read `./simulation_results/k_closest_kubo_0_23.dat` and `./simulation_results/k_closest_kubo_0_17.dat` files and display the graph. We provide files from our tests in the `./simulation_results/` folder.



### Fig. 7
```
python plot_attack_prob.py --output "attack_success_rate.pdf" --input "../experimentCombined/experiment_results" --sybils 15 20 25 30 35 40 45
```
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
Plot the results using:
`python3 plot_key_generation_time.py`

Note that the script will look for `./simulation_results/timing_rsa.csv` and `./simulation_results/timing_eddsa.csv` files, so you have to generate them in advance (or use the one we provided).

### Fig. 10
```
python plot_accuracy.py --output "fp-fn.pdf" --input "../experimentCombined/experiment_results"
```
### Fig. 11
```
python plot_KL.py --output "KLs.pdf" --input "../experimentCombined/experiment_results"
```
### Fig. 12
Run simulation
```
python simulateHonEclKLs.py
```
Plot KL divergence values
```
python plot_KL_netsize.py --output "KLs_netsize.pdf"
```
### Fig. 13