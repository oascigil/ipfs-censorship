#!/bin/bash

IPFS_BINARY=ipfs
#IPFS_HOME=~/snap/ipfs/common
IPFS_HOME=~/.ipfs
#remove all the previous log files
rm *.log

for i in {1..100}
do
	#recreate IPFS keys
	rm -rf $IPFS_HOME
	$IPFS_BINARY init
	$IPFS_BINARY daemon &> /dev/null &
  	pid=$!
	echo IPFS daemon launched
 	sleep 10
	echo starting DHT query
	$IPFS_BINARY dht query bafzaajaiaejcanycz76eb3gzx5abopm5vbjvhpdp76zlcexgdcylbaua7ovyd6yl > ${i}.log
	echo DHT query done
	kill $pid
	sleep 1
done

cat *.log | sort | uniq -c > ./simulation_results/k_closests.dat
rm *.log
