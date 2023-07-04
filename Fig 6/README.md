# Figure 6
This repo provides steps to re-create Figure 6

## Getting the required software
This experiment compares the performance of multiple `libp2p` versions. In the paper, we compare [libp2p v0.17](https://github.com/libp2p/go-libp2p/releases/tag/v0.17.0) and [libp2p v0.23.4](https://github.com/libp2p/go-libp2p/releases/tag/v0.23.4).

For each version you need to:
* download [Kubo](https://github.com/ipfs/kubo)
* Go into the Kubo folder `$ cd kubo`
* change the required `libp2p` version in the `go.mod` file. For instance, to use `libp2p v0.23.4`, open `go.mod` in a text editor and replace `github.com/libp2p/go-libp2p *` with `github.com/libp2p/go-libp2p v0.23.4`
* recompile Kubo with `make install`


## Collecting Data
To collect data, you need to run `./k_closests.sh` script. Note that the script will remove and re-create your public/private IPFS keys. If that's not what you want (e.g., you want to keep your existing IPFS identity) you need to modify the script or backup your IPFS keys before running `./k_closests.sh`. On Ubuntu, the keys are written to `~/.ipfs`. However, you might need to modify the `IPFS_HOME` variable if that's not the case on your system. 

Once you run `./k_closests.sh`, the script will create temporary `*.log` files with the closest peers found by the script and then summarize them in a single `k_closest.dat` file. The temporary files `*.log` will be removed. 

You need to manually rename the `k_closest.dat` to the libp2p version you've used. For instance, `mv k_closest.dat k_closest_kubo_0_23.dat`

## Plotting
To plot the data, simply run `python3 k_closest.py`. The script will automatically read `./k_closest_kubo_0_23.dat` and `./k_closest_kubo_0_17.dat` files and display the graph. 



