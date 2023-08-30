import numpy as np
import matplotlib
import matplotlib.pyplot as plt
from glob import glob
import argparse
import os

from utils import *

from style import *

matplotlib.rcParams["figure.figsize"] = default_figsize
ylim = (0, 10)

argParser = argparse.ArgumentParser()
argParser.add_argument("-o", "--output", help="output filename")
args = argParser.parse_args()

num_query_ids = 256
fixed_threshold = 0.94

columnwidth = num_query_ids + num_query_ids/4
padding = num_query_ids/8
text_location = 9
plt.figure()

n_array = [10000, 15000, 20000, 25000]
column = 0
for n in n_array:
    honest_kls = []
    category = "estnsrand_"
    num_fp = 0
    num_fn = 0
    num_honest = 0
    num_eclipsed = 0
    for (status, color) in [("honest", "blue"), ("eclipsed", "red")]:
        kls = []
        for filename in glob("./simulation_results/{}{}_{}.json".format(category, status, int(n))):
            with open(filename, "r") as f:
                for line in f:
                    if line.startswith("{"):
                        result = json.loads(line)
                        kls.append(result["KL"])
        plt.scatter(np.arange(columnwidth*column+padding, columnwidth*column+padding+num_query_ids), kls[:num_query_ids], color=color, facecolors="none", s=32)
        if status == "honest":
            num_fp = len(np.where(np.array(kls) > fixed_threshold)[0])
            num_honest = len(kls)
        if status == "eclipsed":
            num_fn = len(np.where(np.array(kls) < fixed_threshold)[0])
            num_eclipsed = len(kls)
    report = "FN: " + str(round(num_fn / num_eclipsed * 100,1)) + "%\nFP: " + str(round(num_fp / num_honest * 100,1)) + "%"
    plt.text(columnwidth*(column+0.5), text_location, report, ha="center", va="center", bbox=dict(facecolor='white', alpha=0.5, edgecolor='none'))

    plt.axvline(x = columnwidth*(column+1), color="grey", linestyle = "-")
    column += 1

plt.axhline(y = fixed_threshold, color="black", linestyle = "--")

# Set ylabel for both axes
plt.ylabel("KL divergence")
plt.xticks([columnwidth*(i+0.5) for i in range(column)], ["N = " + human_readable_formatter(n, None) for n in n_array])
plt.tick_params(axis='x',bottom=False)
plt.xlim((0, columnwidth*column))
plt.ylim(ylim)
yticks = plt.yticks()
plt.yticks([fixed_threshold] + list(yticks[0]), [r"$\mathsf{thr}$="+str(fixed_threshold)] + [str(t) for t in list(yticks[0])])
# plt.legend()
if args.output is not None:
    if not os.path.exists("plots"):
        os.makedirs("plots")
    plt.savefig(os.path.join("plots", args.output), bbox_inches='tight')
plt.show()