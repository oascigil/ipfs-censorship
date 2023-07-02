# Generating Figures 13-15

Below you can find instructions to (i) run experiments to generate data to plot Figures 13-15 and (ii) to generate the plots. 

# Experiment Data

You can find the data used to generate the Figs. 13-15 in **mitigation_experiment/logs/**

## Running experiments

1. Build and run ipfs daemon
```
foo@bar:~$ cd kubo
foo@bar:~/kubo/$ go build .
foo@bar:~/kubo/$ cd cmd/ipfs
foo@bar:~/kubo/cmd/ipfs/$ ipfs daemon & 
```
2.  Build mitigation code and run the experiments
```
foo@bar:~$ cd mitigation_experiments
foo@bar:~/mitigation_experiments/$ go build .
foo@bar:~/mitigation_experiments/$ python3 run_experiments.py
```
## Plotting Figures
```
foo@bar:~$ cd plot_figs13-15
foo@bar:~/plot_figs13-15/$ python3 graphs.py
```
The generated plots can be found in **plot_figs13-15/plots/** folder.
