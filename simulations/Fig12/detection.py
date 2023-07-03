import random
from random import sample
import numpy as np
from binary_trie import Trie, bytes_to_bitstring, bitstring_to_bytes, int_to_bitstring, bitstring_to_int
from scipy.special import comb
from scipy.special import rel_entr

k = 20
ideal_pmf_global = 0.5 * 0.5**np.arange(0,256)

# `gen_random_key` generates a random 256-bit key
def gen_random_key(length=256):
    bitstring = "{0:b}".format(random.getrandbits(length))
    return "0"*(length-len(bitstring)) + bitstring

# `transform_netsize`takes a netsize `n` and a list of peerids `peerids` and returns a list of peerids of length `n`. If `n` is less than the length of `peerids`, then `peerids` is randomly sampled to be of length `n`. If `n` is greater than the length of `peerids`, then `peerids` is padded with random peerids to be of length `n`.
def transform_netsize(n, peerids):
    num_peers = len(peerids)
    if n <= num_peers:
        peerids = sample(peerids, n)
    for i in range(num_peers, n):
        peerids.append(gen_random_key())
    return peerids

# `estimate_pmf_best_from_netsize` takes a netsize `n` and returns the estimated ideal PMF of the prefix length counts as described in the paper.
def estimate_pmf_best_from_netsize(n):
    estimate_cdfs_best = np.zeros((k,256))
    s = np.zeros(256)
    for i in range(k):
        x = np.arange(256)
        s += comb(n,i) * (1 - 0.5**(x+1))**(n-i) * 0.5**((x+1)*i)
        estimate_cdfs_best[i] = s
    estimate_pmf_best = np.mean(np.diff(estimate_cdfs_best, axis=1, prepend=0), axis=0)
    return estimate_pmf_best

# `compute_prefix_len_counts_with_peerids` Simulates a DHT with the given peerids and then finds the k closest peers to each id in ids. Then, for each id, it computes the number of peers in each prefix length from 0 to 255. It returns a list of lists, where each sublist is the number of peers in each prefix length for the corresponding id in ids.
def compute_prefix_len_counts_with_peerids(ids, peerids):
    trie = Trie()
    for p in peerids:
        trie.add(p)
        
    prefix_len_counts = []
    
    for ID in ids:
        counts = [0]*256
        closest_ids = trie.n_closest(ID, k)
        for cid in closest_ids:
            prefix_len = [list(ID)[i] == list(cid)[i] for i in range(len(ID))].index(0)
            counts[prefix_len] += 1
        prefix_len_counts.append(counts)
    return prefix_len_counts

# `compute_kl_with_counts` takes a list of prefix length counts and computes the KL divergence between the ideal PMF and the empirical PMF of prefix length counts. It returns a list of KL divergences with the same length as the list of prefix length counts.
def compute_kl_with_counts(prefix_len_counts, ideal_dist, ll=9, ul=20):
    kls = [0]*len(prefix_len_counts)
    for i in range(len(prefix_len_counts)):
        for p in range(ll,ul):
            if prefix_len_counts[i][p] > 0:
                kls[i] += rel_entr(prefix_len_counts[i][p]/k, ideal_dist[p])
    return kls

# `eclipse_attack` simulates an eclipse attack on queried id. Given a set of prefix length counts for a query, it returns a new set of prefix length counts that is the result of an eclipse attack. It also returns the number of Sybil peers added to each prefix length.
def eclipse_attack(prefix_len_counts):
    new_counts = prefix_len_counts.copy()
    adv_counts = np.zeros((len(new_counts),256))
    for h in range(len(prefix_len_counts)):
        best_prefix = np.nonzero(new_counts[h])[0][-1]
        num_added_peers = 20
        arrangement = np.random.choice(np.arange(best_prefix,256), num_added_peers, p=ideal_pmf_global[best_prefix:]/np.sum(ideal_pmf_global[best_prefix:]))
        num_removed_peers = 0
        while num_removed_peers < num_added_peers:
            worst_prefix = np.nonzero(new_counts[h])[0][0]
            if num_added_peers - num_removed_peers < new_counts[h][worst_prefix]:
                new_counts[h][worst_prefix] -= num_added_peers - num_removed_peers
                num_removed_peers = num_added_peers
            else:
                num_removed_peers += new_counts[h][worst_prefix]
                new_counts[h][worst_prefix] = 0
        for l in arrangement:
            new_counts[h][l] += 1
            adv_counts[h][l] += 1
    return new_counts, adv_counts

# get the corresponding k-bucket for the given XOR distance in bytes
def bucket_number_for_distance(d: bytes) -> int:
    count=0
    # iterate on the bytes from left to right
    for b in d:
        # while the byte==0, add 8 (bits) to the counter
        count+=8
        if b!=0:
            # at the first non null byte, shift right until this byte==0
            while b!=0:
                b>>=1
                # for each right shift, remove 1 to counter
                count-=1
            break
    # return the length of the byte string minus the number of leading 0 bits
    return 256-(8*len(d)-count)

def xor_bitstring(bs0: str, bs1: str) -> str:
    s = ""
    if len(bs0) == len(bs1):
        for i in range(len(bs0)):
            if bs0[i]==bs1[i]:
                s+='0'
            else:
                s+='1'
    return s

# 
def estimate_size_from_buckets(peerid, random_id_per_bucket, k_best_ids_per_bucket):
    basket_avgs = np.zeros(k)
    sum_bucket_weights = 0
    all_distances = np.zeros((k,256))
    for i in range(256):
        bucket_size = 0
        for j in range(len(k_best_ids_per_bucket[i])):
            if bucket_number_for_distance(bitstring_to_bytes(xor_bitstring(peerid, k_best_ids_per_bucket[i][j]))) == i:
                bucket_size += 1
        bucket_weight = 2**(bucket_size - k)
        sum_bucket_weights += bucket_weight
        for j in range(k):
            distance = bitstring_to_int(xor_bitstring(random_id_per_bucket[i], k_best_ids_per_bucket[i][j])) / (2**256)
            basket_avgs[j] += bucket_weight * distance
            all_distances[j,i] = distance
    basket_avgs /= sum_bucket_weights
#     return all_distances
    slope = np.sum(basket_avgs)/np.sum(np.arange(1,k+1))
    return 1/slope

# `experiment_size_from_buckets_with_peerids` construct a DHT with `peerids` and returns network size estimates as computed by each peer with id in `ids`.
def experiment_size_from_buckets_with_peerids(ids, peerids):
    trie = Trie()
    for p in peerids:
        trie.add(p)
    
    estimates = []
    
    for ID in ids:
        k_best_ids_per_bucket = []
        random_id_per_bucket = []
        for j in range(256):
            rand_id = ID[:j] + gen_random_key(256-j)
            random_id_per_bucket.append(rand_id)
            closest_ids = trie.n_closest(rand_id, k)
            k_best_ids_per_bucket.append(closest_ids)
        estimates.append(estimate_size_from_buckets(ID, random_id_per_bucket, k_best_ids_per_bucket))
    return estimates