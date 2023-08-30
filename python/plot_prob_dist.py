import numpy as np
import matplotlib.pyplot as plt
import argparse
import os

from utils import *
from detection import *
from style import *

# Take a commandline argument for the output file name
argParser = argparse.ArgumentParser()
argParser.add_argument("-o", "--output", help="output filename")
args = argParser.parse_args()

# Read prefix_len_counts, hon_prefix_len_counts, and len(peerids) from a JSON file
import json
with open("crawl_data/prefix_len_counts.json", 'r') as file:
    data = json.load(file)
    prefix_len_counts = data["prefix_len_counts"]
    hon_prefix_len_counts = data["hon_prefix_len_counts"]
    num_peerids = data["peerids"]

ideal_pmf_best = estimate_pmf_best_from_netsize(num_peerids)

fig = plt.figure()
hon_mean_prefix_len_counts = np.mean(hon_prefix_len_counts, axis=0)
adv_mean_prefix_len_counts = np.mean(prefix_len_counts, axis=0)
min_pref = 0
max_pref = 32
plt.plot(np.arange(min_pref,max_pref), hon_mean_prefix_len_counts[min_pref:max_pref], 'o', markersize=8, color='blue', label="No attack")
plt.plot(np.arange(min_pref,max_pref), ideal_pmf_best[min_pref:max_pref]*k, color='cornflowerblue', linewidth=3, label="Model distribution")
plt.plot(np.arange(min_pref,max_pref), adv_mean_prefix_len_counts[min_pref:max_pref], 'x', markersize=8, color='red', label="Censorship attack (e=20)")
plt.xlabel("Common prefix length with target CID (bits)")
plt.xticks(np.arange(min_pref,max_pref,4))
plt.ylabel("#peer IDs (out of 20 closest)")
plt.grid(True, color = "grey", linewidth = "0.3")
plt.gca().spines['top'].set_visible(False)
plt.gca().spines['right'].set_visible(False)
plt.legend()
if args.output is not None:
    if not os.path.exists("plots"):
        os.makedirs("plots")
    plt.savefig(os.path.join("plots", args.output), bbox_inches='tight')
    print("Plots saved at ", os.path.join("plots", args.output))
plt.show()