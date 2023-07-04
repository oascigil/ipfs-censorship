import argparse
import json
from glob import glob
import numpy as np
import matplotlib
import matplotlib.pyplot as plt

from utils import *
from detection import *

## Data format:
# type Experiment struct {
# 	ClientPeerID string
# 	ProviderPeerID string
# 	NumSybils int
# 	SpecialProvideNumber int
# 	ProvideLatencyMs []int32
# 	FindProvsLatencyMs []int32
# 	ProvideMitLatencyMs []int32
# 	FindProvsMitLatencyMs []int32
# 	NumProvsFound []int
# 	NumProvsFoundMit []int
# }

# Take a commandline argument for the output file name
argParser = argparse.ArgumentParser()
argParser.add_argument("-o", "--output", help="output filename")
argParser.add_argument("-i", "--input", help="input experiment results path", default="../experimentCombined/experiment_results")
args = argParser.parse_args()

numSybilsList = [15,20,30,45]

font = {'family' : 'normal',
        'size'   : 14}

matplotlib.rc('font', **font)
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42
matplotlib.rcParams["figure.figsize"] = (10,10)

plt.figure()

numSybils = 45
oneProvider = False

filename = os.path.join(args.input, "sybil" + str(numSybils) + "LatencySpecial20" + ("Provider1" if oneProvider else ""), "*.json")
# The file has many lines, each line has one json object. Read the fields "ProvideLatencyMs", "FindProvsLatencyMs", "ProvideMitLatencyMs" and "FindProvsMitLatencyMs" and append each of them to a separate list. That's all.
provideLatencyMs = []
findProvsLatencyMs = []
provideMitLatencyMs = []
findProvsMitLatencyMs = []
with open(filename) as f:
    for line in f:
        data = json.loads(line)
        provideLatencyMs.extend(data["ProvideLatencyMs"][1:])
        findProvsLatencyMs.extend(data["FindProvsLatencyMs"][1:])
        provideMitLatencyMs.extend(data["ProvideMitLatencyMs"][1:])
        findProvsMitLatencyMs.extend(data["FindProvsMitLatencyMs"][1:])

# print the average latency of each operation
print("Number of sybils: " + str(numSybils))
print("One provider: " + str(oneProvider))
print("Provide latency: " + str(np.mean(provideLatencyMs)))
print("FindProvs latency: " + str(np.mean(findProvsLatencyMs)))
print("ProvideMit latency: " + str(np.mean(provideMitLatencyMs)))
print("FindProvsMit latency: " + str(np.mean(findProvsMitLatencyMs)))

# Box plot showing the latency of the four operations without showing outliers
# plt.boxplot([findProvsLatencyMs, findProvsMitLatencyMs], showfliers=False)
# plt.xticks([1,2], ["FindProvs", "FindProvsMit"])
# plt.ylabel("Latency (ms)")
# plt.title("Latency of the four operations")
# if args.output is not None:
#     plt.savefig(os.path.join("plots", args.output))
# plt.show()