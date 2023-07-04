import argparse
from glob import glob
import os
import numpy as np
import matplotlib.pyplot as plt

from utils import *
from detection import *
from style import *

## Data format:
# {"ClientPeerID":"QmeoU1o2xhntZbM3LPwyqbdoYKTYpyJowe1gomZ8F4YMEF"}
# {
#     "TargetCID":"QmbyCHYYXJMv4rG3J7DUDcho1MxBsGXaxFt8uZXbxqKnRR",
#     "SybilIDs":["","",...],
#     "PercentEclipsed":100,
#     "ProviderPeerID":"QmPwWekdyW2MzWUEruXtg6oKG2P2ggVxSaoKmgeyvRp33H",
#     "IsEclipsed":true,
#     "NumIntersection":0,
#     "Counts":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,8,4,1,2,3,2,0,...],
#     "Netsize":8056.427538159879,
#     "L":7,
#     "Threshold":1.7476952544539643,
#     "KL":10.685429590665159,
#     "Detection":true
# }

# Take a commandline argument for the output file name
argParser = argparse.ArgumentParser()
argParser.add_argument("-o", "--output", help="output filename")
argParser.add_argument("-i", "--input", help="input experiment results path", default="../experimentCombined/experiment_results")
argParser.add_argument("-s", "--sybils", help="list of sybil counts to plot", nargs='+', type=int, default=[15,20,25,30,35,40,45])
args = argParser.parse_args()

numSybilsList = args.sybils

fixed_threshold = 0.94

num_attack_experiments = np.zeros(len(numSybilsList))
num_attack_eclipsed = np.zeros(len(numSybilsList))
num_no_honest_resolver = np.zeros(len(numSybilsList))
num_detected = np.zeros(len(numSybilsList))

counter = 0
plt.figure()
for numSybils in numSybilsList:
    for filename in glob(os.path.join(args.input, "sybil" + str(numSybils) + "Combined", "*.json")):
        client, results = read_data(filename)
        for result in results:
            if result["Counts"] is None:
                continue
            num_attack_experiments[counter] += 1
            if result["IsEclipsed"]:
                num_attack_eclipsed[counter] += 1
            if result["PercentEclipsed"] > 99:
                num_no_honest_resolver[counter] += 1
            kl = result["KL"]
            netsize = result["Netsize"]
            model_pmf = estimate_pmf_best_from_netsize(netsize)
            if compute_kl_with_counts([result["Counts"]], model_pmf, 7, 30)[0] > fixed_threshold:
                num_detected[counter] += 1
    counter += 1

plt.plot(numSybilsList, num_attack_eclipsed/num_attack_experiments*100, marker="o", lw = 3, label="Attack effective (content censored)")
plt.plot(numSybilsList, num_detected/num_attack_experiments*100, ls='--', marker="o", lw = 3, label="Attack detected")

plt.ylabel("Percentage of attacks")
plt.legend()
plt.xlabel("Number of Sybils placed near target CID")
plt.grid(True, color = "grey", linewidth = "0.1")
plt.xticks(numSybilsList)
plt.gca().spines['top'].set_visible(False)
plt.gca().spines['right'].set_visible(False)
plt.tight_layout()
if args.output is not None:
    plt.savefig(os.path.join("plots", args.output))
plt.show()