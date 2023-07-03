This directory contains simulations to calculate the ideal probability distribution of prefix length counts (Fig. 3) and to simulate detection for different network sizes (Fig. 12)

## Fig. 3
Process data
```
cd Fig3
python process_crawl_data.py
```

Plot data
```
python plot_prob_dist.py
```

## Fig. 12
Run simulation
```
cd Fig12
python simulateHonEclKLs.py
```
Plot KL divergence values
```
python plotKLsN.py
```