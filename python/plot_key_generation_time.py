import pandas as pd
from matplotlib import pyplot as plt
import os
import os.path
import argparse

from style import *
from utils import human_readable_formatter

def sybil_generation_time(ec_filename, rsa_filename):
    if(not os.path.exists(rsa_filename) or not os.path.exists(ec_filename)):
        print("sybil_generation_time:", ec_filename, "or",  rsa_filename, "doesn't exist (or both!) - aborting")
        return
    
    df = pd.read_csv(ec_filename)

    df_rsa = pd.read_csv(rsa_filename)
    df['rsa_time[us]'] = df_rsa['time[us]']
    print(df)
    sybils_num = 35
    cores = 4
    cloud_cost_per_h = 0.16
    #convert to [s] and account for 20 Sybil (*(20/1000,000) = /50,000)
    df['time[h]'] = df['time[us]'].div((1000000/sybils_num)*(cores*3600))
    df['time[s]'] = df['time[us]'].div((1000000/sybils_num)*(cores))
    df['rsa_time[h]'] = df['rsa_time[us]'].div((1000000/sybils_num)*(cores*3600))
    df['rsa_time[s]'] = df['rsa_time[us]'].div((1000000/sybils_num)*(cores))
    #calculate $ cost assuming 0.16$/h. 
    df['cost'] = df['time[h]'].mul(cloud_cost_per_h)
    df['rsa_cost'] = df['rsa_time[h]'].mul(cloud_cost_per_h)
    df['tries'] = df['tries'].mul(sybils_num)
    print(df)
    
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 5), sharex=True)

    ax12 = ax1.twinx()
    ax12.spines['right'].set_position(('axes', 1.0))
    grouped = df.groupby('network_size')['cost']
    x = grouped.mean().index.values
    y = grouped.mean().values
    
    err = grouped.std().values
    ax12.errorbar(x, y, yerr=err, fmt='-o', lw=3, label='EdDSA')
    #ax22.set_ylabel("cost[$]")

    ax22 = ax2.twinx()
    grouped = df.groupby('network_size')['rsa_cost']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    ax22.errorbar(x, y, yerr=err, fmt='-o', lw=3, c='#029e73', label='RSA')
    #ax22.set_ylabel("cost[$]")
    fig.text(0.00, 0.5, 'time[h]', va='center', rotation='vertical', weight='normal')
    fig.text(0.98, 0.5, 'cost[$]', va='center', rotation='vertical', weight='normal')

    grouped = df.groupby('network_size')['time[h]']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    ax1.errorbar(x, y, yerr=err, fmt='-o', c='b', alpha=0)
    #ax1.set_ylabel("time[h]")
    # print("y", y)
    
    grouped = df.groupby('network_size')['rsa_time[h]']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    ax2.errorbar(x, y, yerr=err, fmt='-o', c='b', alpha=0)

    ax2.set_xlabel('Network Size')
    ax1.spines['top'].set_visible(False)
    ax2.spines['top'].set_visible(False)
    ax12.spines['top'].set_visible(False)
    ax22.spines['top'].set_visible(False)
    ax1.grid(True, color = "grey", linewidth = "0.3",axis = 'y')
    ax2.grid(True, color = "grey", linewidth = "0.3",axis = 'y')

    ax1.xaxis.set_major_formatter(human_readable_formatter)
    ax2.xaxis.set_major_formatter(human_readable_formatter)
    

    ax12.legend(loc='upper left')
    ax22.legend(loc='upper left')

    if not os.path.exists("plots"):
        os.makedirs("plots")
    plt.savefig('./plots/sybil_generation_time.pdf', bbox_inches='tight')
    print("Plots saved at ./plots/sybil_generation_time.pdf")
    plt.show()

def main():
    argParser = argparse.ArgumentParser()
    argParser.add_argument("-i", "--input", help="input simulation results", nargs=2, default=['./simulation_results/timing_eddsa.csv', './simulation_results/timing_rsa.csv'])
    args = argParser.parse_args()

    sybil_generation_time(args.input[0], args.input[1])

if __name__ == '__main__':
    main()
