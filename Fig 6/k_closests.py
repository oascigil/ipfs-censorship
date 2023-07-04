import matplotlib.pyplot as plt
import numpy as np

from style import *

log_files = {'libp2p v0.23': './k_closest_kubo_0_23.dat', 
             'libp2p v0.17': './k_closest_kubo_0_17.dat'}

colors = ['#0173b2', '#de8f05', '#029e73', '#d55e00', '#cc78bc', '#ca9161', '#fbafe4', '#949494', '#ece133', '#56b4e9']

data = {}
data['perfect routing'] = ''

num_experiments = 100
max_size = 0
for version, log_file in log_files.items():
    if version not in data:
        data[version] = []
    size = 0
    with open(log_file) as f:
        for line in f:
            count = line.split()[0]
            data[version].append((100*int(count))/num_experiments)
            size += 1
    if size > max_size:
        max_size = size

data['perfect routing'] = [100]*20 + [0]*(max_size - 20)

fig, ax = plt.subplots(figsize=default_figsize)

x = np.arange(max_size)
width = 0.25  # the width of the bars
multiplier = 0


#expand the lists in data to the same size
for version, l in data.items():
    offset = width * multiplier
    y = l + [0]*(max_size - len(l))
    ax.bar(x + 1 + offset, sorted(list(y), reverse=True), width, align='center', label=version, color=colors[multiplier])
    multiplier += 1


ax.set_ylabel("Dicovered by ratio of queries [%]")
ax.set_xlabel("Node Index")
ax.set_xlim(left=0.5, right=31)
ax.spines['top'].set_visible(False)
ax.spines['right'].set_visible(False)
ax.legend()
#plt.xticks(range(len(data)), list(data.keys()))
plt.show()