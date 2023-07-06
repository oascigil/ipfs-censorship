# Figure 8
Below you can find instructions to run experiments to generate data to plot Figure 8, which executes the live censorship with the content being added prior to execution.

## Setup
For the experiment, one IPFS instance needs to provide the files to be censored, while another launches the Sybil nodes (during our experiments these are on separate server machines). The host machine needs to run the following with a random file 'initial.txt':
```
go run setup.go
```

This creates random files locally and adds them to the network. The resulting 'cidList.csv' contains the cids to be censored in the next part of the experiment. 


## Generating data
The second instance can now launch the censorship, logging the time taken for all queries to be censored and the percentage of the 20 closest peers which are Sybils (example in 'sample_results.csv'). To run specify the number of Sybils and the 'cidList.csv' file:
```
go run main.go <sybilnumber> <filename>
```
The results log the time taken to generate the Sybil keys, time for the eclipse to be successful (defined as 3 hours eclipsed in a row), the number of eclipsed queries per hour (out of 5 tries), and the percentage of closest 20 DHT nodes which are our Sybils.

Note that this experiment takes a long time as it checks hourly (up to 50 hours) if the cid is censored, for a large number of cids.

