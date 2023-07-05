import numpy as np
import json
import os

from detection import *

k = 20
num_query_ids = 256
num_test_ids = 256
use_netsize_estimates = True

# Generate random query IDs which will be eclipsed and not eclipsed (honest)
hon_query_ids = []
for h in range(num_query_ids):
    hon_query_ids.append(gen_random_key())

ecl_query_ids = []
for h in range(num_query_ids):
    ecl_query_ids.append(gen_random_key())

# Generate random peerids from where the query ID will be queried
test_ids = []
for h in range(num_test_ids):
    test_ids.append(gen_random_key())

n_array = np.arange(10000,30000,5000)
peerids = []
for n in n_array:
    honest_json_string = ""
    eclipsed_json_string = ""
    print("n = " + str(n))
    peerids = transform_netsize(n, peerids)
    hon_prefix_len_counts = compute_prefix_len_counts_with_peerids(hon_query_ids, peerids)
    ecl_prefix_len_counts = compute_prefix_len_counts_with_peerids(ecl_query_ids, peerids)
    new_ecl_prefix_len_counts, _ = eclipse_attack(ecl_prefix_len_counts)

    if use_netsize_estimates:
        size_estimates = experiment_size_from_buckets_with_peerids(test_ids, peerids)
    else:
        size_estimates = [n]
    
    for i in range(num_test_ids):
        L = 7
        model_pmf = estimate_pmf_best_from_netsize(size_estimates[i])
        j = i
        hon_kls = compute_kl_with_counts([hon_prefix_len_counts[j]], model_pmf, L,30)[0]
        ecl_kls = compute_kl_with_counts([new_ecl_prefix_len_counts[j]], model_pmf, L,30)[0]
        honest_sample = {
            "Counts": hon_prefix_len_counts[j],
            "Netsize": int(size_estimates[i]),
            "KL": hon_kls
        }
        honest_json_string += json.dumps(honest_sample) + "\n"
        eclipsed_sample = {
            "Counts": new_ecl_prefix_len_counts[j],
            "Netsize": int(size_estimates[i]),
            "KL": ecl_kls
        }
        eclipsed_json_string += json.dumps(eclipsed_sample) + "\n"

    if not os.path.exists("simulation_results"):
        os.makedirs("simulation_results")
    with open("simulation_results/estnsrand_honest_"+str(n)+".json", "w") as outfile:
        outfile.write(honest_json_string)
    with open("simulation_results/estnsrand_eclipsed_"+str(n)+".json", "w") as outfile:
        outfile.write(eclipsed_json_string)