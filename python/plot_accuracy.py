import argparse
import os
from glob import glob
import numpy as np
import matplotlib
import matplotlib.pyplot as plt

from utils import *
from detection import *

from style import *

## Data format:
# {
#     "TargetCID":"QmbyCHYYXJMv4rG3J7DUDcho1MxBsGXaxFt8uZXbxqKnRR",
#     "SybilIDs":["","",...],
#     "PercentEclipsed":100,
#     "ProviderPeerID":"QmPwWekdyW2MzWUEruXtg6oKG2P2ggVxSaoKmgeyvRp33H",
#     "IsEclipsed":true,
#     "RegionSize":20,
#     "IsMitigated":true,
#     "UpdatedPeers":["","",...],
#     "ProviderDiscoveredPeers":["","",...],
#     "ContactedPeers":["","",...],
#     "RespondedPeers":["","",...],
#     "NumContacted":100,
#     "NumUpdated":20,
#     "NumSybilsUpdated":20,
#     "NumIntersection":0,
#     "NumLookups":10
#     "Counts":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,8,4,1,2,3,2,0,...],
#     "Netsize":8056.427538159879,
#     "L":7,
#     "Threshold":1.7476952544539643,
#     "KL":10.685429590665159,
#     "Detection":true
# }

argParser = argparse.ArgumentParser()
argParser.add_argument("-o", "--output", help="output filename")
argParser.add_argument("-i", "--input", help="input experiment results path", default="../experimentCombined/detection_results")
argParser.add_argument("-s", "--sybils", help="list of sybil counts to plot", nargs='+', type=int, default=[20,45])
argParser.add_argument("-t", "--threshold", help="threshold to plot", type=float, default=0.94)
args = argParser.parse_args()

numSybilsList = args.sybils
chosen_threshold = args.threshold

matplotlib.rcParams["figure.figsize"] = default_figsize
#colorsList = ['black', 'red', 'blue', 'green', 'yellow']
colorsList = ['#0173b2', '#029e73', '#de8f05', '#d55e00', '#cc78bc', '#ca9161', '#fbafe4', '#949494', '#ece133', '#56b4e9']
markerList = ['o', 's', 'v', 'D', 'X']
stylesList = ['-', '--', '-.', ':', '-']

args = argParser.parse_args()

ideal_thresholds = np.linspace(0.8, 2, 101)
print(ideal_thresholds)

attack_detections_per_threshold_exact_dist = [0] * len(ideal_thresholds)
honest_detections_per_threshold_exact_dist = [0] * len(ideal_thresholds)

num_attack_experiments = 0
num_honest_experiments = 0
num_atack_eclipsed = 0
num_honest_eclipsed = 0

num_attack_percent_but_not_eclipsed = 0
num_attack_percent = 0
num_attack_percent_but_eclipsed = 0

plt.figure()
for filename in glob(os.path.join(args.input, "sybil0Combined", "*.json")):
    client, results = read_data(filename)
    for result in results:
        if result["Counts"] is None:
            continue
        num_honest_experiments += 1
        if result["IsEclipsed"]:
            num_honest_eclipsed += 1

        netsize = result["Netsize"]

        for i in range(len(ideal_thresholds)):
            model_pmf = estimate_pmf_best_from_netsize(netsize)
            if compute_kl_with_counts([result["Counts"]], model_pmf, 7, 30) > ideal_thresholds[i]:
                honest_detections_per_threshold_exact_dist[i] += 1

fp_per_threshold_exact_dist = [det/float(num_honest_experiments)*100 for det in honest_detections_per_threshold_exact_dist]
print("False positives: " + str(fp_per_threshold_exact_dist))

for (numSybils, color, marker, style) in zip(numSybilsList, colorsList, markerList, stylesList):
    for filename in glob(os.path.join(args.input, "sybil" + str(numSybils) + "Combined", "*.json")):
        client, results = read_data(filename)
        for result in results:
            if result["Counts"] is None:
                continue
            num_attack_experiments += 1
            if result["IsEclipsed"]:
                num_atack_eclipsed += 1
            percentEclipsed = result["PercentEclipsed"]
            if percentEclipsed == 100:
                num_attack_percent += 1
            if percentEclipsed < 100 and result["IsEclipsed"]:
                num_attack_percent_but_eclipsed += 1
            if percentEclipsed == 100 and not result["IsEclipsed"]:
                num_attack_percent_but_not_eclipsed += 1

            kl = result["KL"]
            netsize = result["Netsize"]
            
            for i in range(len(ideal_thresholds)):
                model_pmf = estimate_pmf_best_from_netsize(netsize)
                if compute_kl_with_counts([result["Counts"]], model_pmf, 7, 30) > ideal_thresholds[i]:
                    attack_detections_per_threshold_exact_dist[i] += 1


    fn_per_threshold_exact_dist = [(num_attack_experiments - det)/float(num_attack_experiments)*100 for det in attack_detections_per_threshold_exact_dist]
    print("False negatives: " + str(fn_per_threshold_exact_dist))

    plt.plot(fp_per_threshold_exact_dist, fn_per_threshold_exact_dist, color=color, lw=3, marker=marker, ls=style, label = str(numSybils) + " Sybils")

    # chosen_threshold_index = 34
    # Find the index of the threshold closest to the chosen threshold
    chosen_threshold_index = min(range(len(ideal_thresholds)), key=lambda i: abs(ideal_thresholds[i]-chosen_threshold))
    plt.scatter(fp_per_threshold_exact_dist[chosen_threshold_index], fn_per_threshold_exact_dist[chosen_threshold_index], s=200, facecolors='none', edgecolors='r', linewidth=2)
    print(ideal_thresholds[chosen_threshold_index])
    print(fp_per_threshold_exact_dist[chosen_threshold_index])
    print(fn_per_threshold_exact_dist[chosen_threshold_index])

# Plot fn vs fp as percentage
plt.xlabel("False Positive Rate [%]")
plt.ylabel("False Negative Rate [%]")
plt.ylim((-0.1,10.1))
plt.xlim((-0.1,10.1))
plt.legend()

plt.gca().spines['top'].set_visible(False)
plt.gca().spines['right'].set_visible(False)
plt.grid(True, color = "grey", linewidth = "0.2")
plt.annotate('Chosen detection threshold thr=0.94', xy=(4.1, 1.5), xytext=(4.5, 4),
            arrowprops=dict(facecolor='black', arrowstyle='->'),
            #fontsize=12,
            ha='center')
# set figure size
if args.output is not None:
    if not os.path.exists("plots"):
        os.makedirs("plots")
    plt.savefig(os.path.join("plots", args.output))
plt.show()