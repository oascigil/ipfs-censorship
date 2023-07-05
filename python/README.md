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
### Fig. 7
```
python plot_attack_prob.py --output "attack_success_rate.pdf" --input "../experimentCombined/experiment_results" --sybils 15 20 25 30 35 40 45
```
### Fig. 9
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