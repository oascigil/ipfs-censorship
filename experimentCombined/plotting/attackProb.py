import argparse
import json
from glob import glob
import numpy as np
import matplotlib
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
args = argParser.parse_args()

numSybilsList = [15,20,25,30,35,40,45]

fixed_threshold = 0.94

num_attack_experiments = np.zeros(len(numSybilsList))
num_attack_eclipsed = np.zeros(len(numSybilsList))
num_no_honest_resolver = np.zeros(len(numSybilsList))
num_detected = np.zeros(len(numSybilsList))

counter = 0
plt.figure()
for numSybils in numSybilsList:
    for filename in glob("./experiment_results/sybil" + str(numSybils) + "Combined" + "/*.json"):
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

# Create a bar plot showing for each number of sybils, the percetage of attacks that were eclipsed and the percentage of attacks that were detected
# Increse the width of the bars to make them more visible
# width = 1
# plt.bar(np.array(numSybilsList) - width, num_no_honest_resolver/num_attack_experiments*100, width, label="No honest resolver found")
# plt.bar(np.array(numSybilsList), num_attack_eclipsed/num_attack_experiments*100, width, label="Content unavailable")
# plt.bar(np.array(numSybilsList) + width, num_detected/num_attack_experiments*100, width, label="Attack detected")
# Instead of bar plots, use line plots with only markers and no lines
# plt.plot(numSybilsList, num_no_honest_resolver/num_attack_experiments*100, marker="o", label="No honest resolver found")
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
#plt.title("Percentage of attacks that were effective and detected")
if args.output is not None:
    plt.savefig("./experiment_plots/"+args.output)
plt.show()