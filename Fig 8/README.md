# Figure 8
Below you can find instructions to run experiments to generate data to plot Figure 8, which executes the live censorship with the content being added prior to execution.

## Setup
For the experiment, one IPFS instance needs to provide the files to be censored, while another launches the Sybil nodes (during our experiments these are on separate server machines). The host machine needs to run:
```
go run setup.go
```

This creates random files locally and adds them to the network. The resulting 'cidlist.csv' contains the cids to be censored in the next part of the experiment. 



part 2 attack figure
step 1 setup explanation - provide content on different ipfs instance
step 2 launch attack for those cids





Below you can find instructions to (i) run experiments to generate data to plot Figure 9 and (ii) to generate the plots.

## Generating data
Run two sets of experiments measuring generation time for RSA and EDDSA:
```
go run measureGenerateSybilKeys rsa > timing_rsa.csv
go run measureGenerateSybilKeys eddsa > timing_eddsa.csv
```

We provide our `sample_rsa.csv` and `sample_eddsa.csv` files for reference (generating RSA keys takes a while).

## Plotting
Plot the results using:
`python3 plot.py`

Note that the script will look for `timing_rsa.csv` and `timing_eddsa.csv` files, so you have to generate them in advance

