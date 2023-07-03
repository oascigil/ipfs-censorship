`cd` into this directory, that is `experimentCombined`

## Run experiment
Run IPFS daemon in parallel
```
../common/kubo/cmd/ipfs/ipfs daemon
```
Compile experiment
```
cd experiment
go build
cd ..
```
Run experiment
```
./experiment/main -cids 10 -clients 5 -sybils 20 -region 20 -outpath "experiment_results/sybil20Combined"
```
You may change the number of CIDS and clients. Each experiment takes about 5 minutes per CID per client. You may repeat the experiment for different number of Sybils to generate all the required data. Change `outpath` accordingly.

## Generate plots
Fig. 7
```
python plotting/attackProb.py --output "attack_success_rate.pdf"
```

Fig. 10
```
python plotting/calcAccuracy.py --output "fp-fn"
```

Fig. 11
```
python plotting/plotKLs.py --output "KLs.pdf"
```