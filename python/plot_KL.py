import argparse
import numpy as np
import matplotlib
import matplotlib.pyplot as plt
from glob import glob
import os

from utils import *
from detection import *

from style import *

matplotlib.rcParams["figure.figsize"] = default_figsize

argParser = argparse.ArgumentParser()
argParser.add_argument("-o", "--output", help="output filename")
argParser.add_argument("-i", "--input", help="input experiment results path", default="../experimentCombined/detection_results")
args = argParser.parse_args()

fixed_threshold = 0.94

sybilNumbers = [20,30,45]
fileIterList = [os.path.join(args.input, "sybil" + str(numSybils) + "Combined", "*.json") for numSybils in sybilNumbers]

columnwidth = 300
padding = 25
text_location = 9
ylim = ((0, 10))
plt.figure()
column = 0


honest_kl_minus_thresholds = []
num_honest_eclipsed = 0
for filename in glob(os.path.join(args.input, "sybil0Combined", "*.json")):
    client, results = read_data(filename)
    for result in results:
        if result["Counts"] is None:
            continue
        if result["IsEclipsed"]:
            num_honest_eclipsed += 1
        netsize = result["Netsize"]
        model_pmf = estimate_pmf_best_from_netsize(netsize)
        honest_kl_minus_thresholds.append(compute_kl_with_counts([result["Counts"]], model_pmf, 7, 30))

num_honest = len(honest_kl_minus_thresholds)
# Count number of positive elements of honest_kl_minus_thresholds
num_fp = len(np.where(np.array(honest_kl_minus_thresholds) > fixed_threshold)[0])
print("False positive rate: " + str(num_fp / num_honest * 100) + "%")

plt.scatter(np.arange(columnwidth*column+padding, columnwidth*column+padding+num_honest), honest_kl_minus_thresholds, color="blue", s=32, facecolors="none")
report = "Censored: " + str(round(num_honest_eclipsed / num_honest * 100, 1)) + "%\nDetected: " + str(round(num_fp / num_honest * 100,1)) + "%"
plt.text(columnwidth*(column+0.5), text_location, report, ha="center", va="center", bbox=dict(facecolor='white', alpha=0.5, edgecolor='none'))

column += 1

for fileIterString in fileIterList:
    attack_kl_minus_thresholds = []
    attack_successes = []
    for filename in glob(fileIterString):
        client, results = read_data(filename)
        for result in results:
            if result["Counts"] is None:
                continue
            netsize = result["Netsize"]
            model_pmf = estimate_pmf_best_from_netsize(netsize)
            attack_kl_minus_thresholds.append(compute_kl_with_counts([result["Counts"]], model_pmf, 7, 30))
            attack_successes.append(result["IsEclipsed"])
    # Find indices of attack_successes that are True
    attack_success_indices = np.where(np.array(attack_successes) == True)[0]
    attack_failure_indices = np.where(np.array(attack_successes) == False)[0]
    # Extract elements of attack_kl_minus_thresholds at those indices
    successful_attack_kls_minus_thresholds = np.array(attack_kl_minus_thresholds)[attack_success_indices]
    # Extract all other elements of attack_kl_minus_thresholds
    failed_attack_kl_minus_thresholds = np.array(attack_kl_minus_thresholds)[attack_failure_indices]
    num_successful_attacks = len(successful_attack_kls_minus_thresholds)
    num_failed_attacks = len(failed_attack_kl_minus_thresholds)

    print("Eclipse attack success rate: " + str(num_successful_attacks / (num_successful_attacks + num_failed_attacks) * 100) + "%")
    # Count number of non-positive elements of attack_kl_minus_thresholds
    num_fn = len(np.where(np.array(attack_kl_minus_thresholds) <= fixed_threshold)[0])
    print("False negative rate: " + str(num_fn / (num_successful_attacks + num_failed_attacks) * 100) + "%")

    plt.scatter(columnwidth*column + padding + attack_success_indices, successful_attack_kls_minus_thresholds, color="red", marker="o", s=32)
    plt.scatter(columnwidth*column + padding + attack_failure_indices, failed_attack_kl_minus_thresholds, color="red", marker="o", s=32, facecolors='none')
    plt.axvline(x = columnwidth*(column), color="grey", linestyle = "-")
    report = "Censored: " + str(round(num_successful_attacks / (num_successful_attacks + num_failed_attacks) * 100, 1)) + "%\nDetected: " + str(round(100 - num_fn / (num_successful_attacks + num_failed_attacks) * 100,1)) + "%"
    plt.text(columnwidth*(column+0.5), text_location, report, ha="center", va="center", bbox=dict(facecolor='white', alpha=0.5, edgecolor='none'))
    column += 1


plt.axhline(y = fixed_threshold, color="black", linestyle = "--", label="Threshold = " + str(fixed_threshold))
plt.ylabel("KL divergence")
plt.xticks([columnwidth*(i+0.5) for i in range(column)], ["No sybils"] + [str(numSybils) + " sybils" for numSybils in sybilNumbers])
plt.tick_params(axis='x',bottom=False)
plt.xlim((0, columnwidth*(column)))
plt.ylim(ylim)
yticks = plt.yticks()
plt.yticks([fixed_threshold] + list(yticks[0]), [r"$\mathsf{thr}$="+str(fixed_threshold)] + [str(t) for t in list(yticks[0])])
# plt.legend()
if args.output is not None:
    if not os.path.exists("plots"):
        os.makedirs("plots")
    plt.savefig(os.path.join("plots", args.output), bbox_inches='tight')
    print("Plots saved at ", os.path.join("plots", args.output))
plt.show()




