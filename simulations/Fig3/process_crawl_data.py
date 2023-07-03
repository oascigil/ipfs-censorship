import csv
import time

from detection import *

class NebulaPeer:
    def __init__(self, nebula_id, peer_id, neighbors_ids):
        self.nebula_id = nebula_id
        self.peer_id = peer_id
        self.neighbors_ids = neighbors_ids
        
        self.key = multihash_to_kad_id(peer_id)
        
        self.alive = len(neighbors_ids)>0
        
        self.buckets = [[] for _ in range(257)]
        self.neighbors = {}
        
    def distance(self, p):
        return xor_distance(self.key, p.key)
        
    def addNeighbor(self, peer):
        if bytes_to_bitstring(peer.key) not in self.neighbors:
            self.neighbors[bytes_to_bitstring(peer.key)]=(peer)
            self.buckets[bucket_number_for_distance(self.distance(peer))].append(peer)
        
    def __str__(self):
        return "nebula_id: "+str(self.nebula_id)+", peer_id: "+str(self.peer_id)+", neighbors: "+str(self.neighbors_ids)


# Gathering data for plots

## Global variables used to collect data from the crawls before plotting it

peerids_distribution_data = []
peerids_distribution_chunk88_data = [] # 1011000

## Routing table distribution in the k-buckets

levels=7

def peerids_distribution(peers):
    # getting the list of bitstrings for the given peerids
    peerids = [bytes_to_bitstring(peers[p].key) for p in peers]
    
    trie = Trie()
    for p in peerids:
        trie.add(p)
        
    density=[0] * 2**levels

    for i in range(2**levels):
        path="{0:b}".format(i)           # getting binary number
        path='0'*(levels-len(path))+path # zero padding for constant size bitstring

        density[i]=trie.find_trie(path).size
        
    peerids_distribution_data.append(density)
    
    levels88 = 7
    density88 = [0] * 2**levels88
    for i in range(2**levels88):
        path="{0:b}".format(i)           # getting binary number
        path='0'*(levels88-len(path))+path # zero padding for constant size bitstring
        path=int_to_bitstring(88,7)+path

        for p in peerids:
            if p[:len(path)] == path:
                density88[i]+=1
                
        
    peerids_distribution_chunk88_data.append(density88)


with open("data/all-peerids.csv", 'r') as file:
    nebula_peerids = {line[0]:line[1] for line in csv.reader(file)}

filename = "data/nebula-peers-2crawls.csv"

with open(filename, 'r') as file:
    all_crawls = [line for line in csv.reader(file)]

crawl_id = all_crawls[0][0]
peers = {}
startTime = time.time()

for i in range(len(all_crawls)):
    entry = all_crawls[i]
        
    # add entry to peers
    peers[entry[1]] = NebulaPeer(entry[1], entry[2], entry[3:])

    if i == len(all_crawls)-1 or all_crawls[i+1][0] != crawl_id:
        # do all computations for the crawl
        print("Crawl", crawl_id, ":", len(peers),"alive peers crawled, time elapsed:", time.time() - startTime)
        
        # stale peers count
        stale_count = 0
        # define neighbor relationships
        for p in peers.copy():
            for n in peers[p].neighbors_ids:
                if n not in peers:
                    stale_count += 1
                    peers[n] = NebulaPeer(n, nebula_peerids[n], [])
                peers[p].addNeighbor(peers[n])

        
        peerids_distribution(peers)
            
        print("Crawl", crawl_id, "finished, stale peers:", stale_count,", total time elapsed:", time.time() - startTime)
        
        # reset variables for next crawl
        if i < len(all_crawls) - 1:
            peers = {}
            crawl_id = all_crawls[i+1][0]
            startTime = time.time()

peerids = [bytes_to_bitstring(peers[p].key) for p in peers]

# Write the peerids to a file
with open("data/peerids_bitstrings.csv", 'w') as file:
    writer = csv.writer(file)
    for p in peerids:
        writer.writerow([p])

# Read the peerids from a file
with open("data/peerids_bitstrings.csv", 'r') as file:
    peerids = [line[0] for line in csv.reader(file)]

## Attack scenarios
query_ids = []
num_query_ids = 1024
for h in range(num_query_ids):
    query_ids.append(gen_random_key())

prefix_len_counts = compute_prefix_len_counts_with_peerids(query_ids, peerids)
new_counts, adv_counts = eclipse_attack(prefix_len_counts)

## Honest scenarios

hon_query_ids = []
num_query_ids = 1024
for h in range(num_query_ids):
    hon_query_ids.append(gen_random_key())

hon_prefix_len_counts = compute_prefix_len_counts_with_peerids(hon_query_ids, peerids)

# Write prefix_len_counts, hon_prefix_len_counts, and len(peerids) as a JSON file
import json
with open("data/prefix_len_counts.json", 'w') as file:
    json.dump({"prefix_len_counts": prefix_len_counts, "hon_prefix_len_counts": hon_prefix_len_counts, "peerids": len(peerids)}, file)