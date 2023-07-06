import argparse
import json
from glob import glob
import numpy as np
import matplotlib
import matplotlib.pyplot as plt
import os

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
# 	NumProvsFound []int
# }

# Take a commandline argument for the output file name
argParser = argparse.ArgumentParser()
argParser.add_argument("-i", "--input", help="input experiment results path", default="../experimentCombined/experiment_results")
args = argParser.parse_args()


# Rewrite this code to generate the same table as in the paper
# Rewrite experimentLatency to generate specific folder names for the expeirment results that can be used here
# Optional: Add an option in go-libp2p-kad-dht to set the number of providers to wait for. Add this as an argument in experimentLatency.

font = {'family' : 'normal',
        'size'   : 14}

matplotlib.rc('font', **font)
matplotlib.rcParams['pdf.fonttype'] = 42
matplotlib.rcParams['ps.fonttype'] = 42
matplotlib.rcParams["figure.figsize"] = (10,10)

plt.figure()

numSybilsList = [0,45]

provideDefault0 = []
provideMit0 = []
provideMit45 = []
findProvsDefault0 = []
findProvsMit0 = []
findProvsMit45 = []
find1ProvDefault0 = []
find1ProvMit0 = []
find1ProvMit45 = []

# Read the data from the experiment results
# Discard the first value in each list because the first operation for each client involve initializing the network size estimator which takes additional time. However, this additional latency is amortized during normal DHT operation.
filename = os.path.join(args.input, "sybil0Latency", "latency.json")
with open(filename) as f:
        for line in f:
                data = json.loads(line)
                provideDefault0.extend(data["ProvideLatencyMs"][1:])
                findProvsDefault0.extend(data["FindProvsLatencyMs"][1:])
                provideMit0.extend(data["ProvideMitLatencyMs"][1:])
                findProvsMit0.extend(data["FindProvsMitLatencyMs"][1:])

filename = os.path.join(args.input, "sybil45Latency", "latency.json")
with open(filename) as f:
        for line in f:
                data = json.loads(line)
                provideMit45.extend(data["ProvideMitLatencyMs"][1:])
                findProvsMit45.extend(data["FindProvsMitLatencyMs"][1:])

filename = os.path.join(args.input, "sybil0LatencyProvider1", "latency.json")
if os.path.exists(filename):
        with open(filename) as f:
                for line in f:
                        data = json.loads(line)
                        find1ProvDefault0.extend(data["FindProvsLatencyMs"][1:])
                        find1ProvMit0.extend(data["FindProvsMitLatencyMs"][1:])

filename = os.path.join(args.input, "sybil45LatencyProvider1", "latency.json")
if os.path.exists(filename):
        with open(filename) as f:
                for line in f:
                        data = json.loads(line)
                        find1ProvMit45.extend(data["FindProvsMitLatencyMs"][1:])

# print the average latency of each operation

print("Number of sybils: 45")
print("Provide (mitigation) latency: " + str(np.mean(provideMit45)))
print("FindProviders (mitigation) latency: " + str(np.mean(findProvsMit45)))
print("Find 1 provider (mitigation) latency: " + str(np.mean(find1ProvMit45)))

print("")
print("Number of sybils: 0 (no attack)")
print("Provide (default) latency: " + str(np.mean(provideDefault0)))
print("Provide (mitigation) latency: " + str(np.mean(provideMit0)))
print("FindProviders (default) latency: " + str(np.mean(findProvsDefault0)))
print("FindProviders (mitigation) latency: " + str(np.mean(findProvsMit0)))
print("Find 1 provider (default) latency: " + str(np.mean(find1ProvDefault0)))
print("Find 1 provider (mitigation) latency: " + str(np.mean(find1ProvMit0)))