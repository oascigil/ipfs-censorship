import pandas as pd
from matplotlib import pyplot as plt
import matplotlib
import os.path
import numpy as np
from subprocess import PIPE, Popen
from style import *
from matplotlib.ticker import FuncFormatter
import random

from style import *



# Define a custom formatting function
def human_readable_formatter(x, pos):
    """
    Custom formatting function to convert tick values to human-readable names.
    For example, 1000 will be converted to '1K', and 1500000 will be converted to '1.5M'.
    """
    if x >= 1e6:
        return f"{x/1e6:.1f}M"
    elif x >= 1e3:
        return f"{x/1e3:.0f}K"
    else:
        return f"{x:.0f}"

def probability_dht_resolution():
    #average number of bitswap connections
    b = 300
    #popularity of the file
    x = np.linspace(0,1,1000)
    y = 1 - (1 - x)**b
    
    _, ax = plt.subplots(figsize=default_figsize)
    ax.plot(x,y, 'r')
    

def sybil_generation_time2(ec_filename = '../experiment/generate_timing.csv', rsa_filename = '../experiment/generate_timing_rsa.csv'):
    if(not os.path.exists(rsa_filename) or not os.path.exists(ec_filename)):
        print("sybil_generation_time:", ec_filename, "or",  rsa_filename, "doesn't exist - skipping")
        return
    
    df = pd.read_csv(ec_filename)

    df_rsa = pd.read_csv(rsa_filename)
    #df_rsa.rename(columns={'time[us]': 'rsa_time[us]'}, inplace=True)
    df['rsa_time[us]'] = df_rsa['time[us]']
    print(df)
    #return
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
    #ax2.set_ylabel("time[h]")




    # ax1.spines['right'].set_visible(False)
    # ax1.spines['top'].set_visible(False)
    ax2.set_xlabel('Network Size')
    # #ax2.spines['right'].set_visible(False)
    ax1.spines['top'].set_visible(False)
    ax2.spines['top'].set_visible(False)
    ax12.spines['top'].set_visible(False)
    ax22.spines['top'].set_visible(False)
    ax1.grid(True, color = "grey", linewidth = "0.3",axis = 'y')
    ax2.grid(True, color = "grey", linewidth = "0.3",axis = 'y')
    #ax3.spines['right'].set_visible(False)
    #ax3.spines['top'].set_visible(False)

    ax1.xaxis.set_major_formatter(human_readable_formatter)
    ax2.xaxis.set_major_formatter(human_readable_formatter)
    

    ax12.legend(loc='upper left')
    ax22.legend(loc='upper left')

def sybil_generation_time(ec_filename = '../experiment/generate_timing.csv', rsa_filename = '../experiment/generate_timing_rsa.csv'):
    if(not os.path.exists(rsa_filename) or not os.path.exists(ec_filename)):
        print("sybil_generation_time:", ec_filename, "or",  rsa_filename, "doesn't exist - skipping")
        return
    
    df = pd.read_csv(ec_filename)

    df_rsa = pd.read_csv(rsa_filename)
    #df_rsa.rename(columns={'time[us]': 'rsa_time[us]'}, inplace=True)
    df['rsa_time[us]'] = df_rsa['time[us]']
    print(df)
    #return
    sybils_num = 35
    cores = 4
    #convert to [s] and account for 20 Sybil (*(20/1000,000) = /50,000)
    df['time[h]'] = df['time[us]'].div((1000000/sybils_num)*(cores*3600))
    df['time[s]'] = df['time[us]'].div((1000000/sybils_num)*(cores))
    df['rsa_time[h]'] = df['rsa_time[us]'].div((1000000/sybils_num)*(cores*3600))
    df['rsa_time[s]'] = df['rsa_time[us]'].div((1000000/sybils_num)*(cores))
    #calculate $ cost assuming 0.16$/h. 
    df['cost'] = df['time[h]'].mul(0.16)
    df['rsa_cost'] = df['rsa_time[h]'].mul(0.16)
    df['tries'] = df['tries'].mul(sybils_num)
    print(df)
    
    _, ax1 = plt.subplots(figsize=default_figsize)

    grouped = df.groupby('network_size')['time[h]']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    ax1.errorbar(x, y, yerr=err, fmt='-o', c='w', alpha=0)
    ax1.set_ylabel("time[h]")
    print("y", y)
    grouped = df.groupby('network_size')['rsa_time[h]']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    
    ax1.errorbar(x, y, yerr=err, fmt='-o', c='w', alpha=0)




    ax2 = ax1.twinx()
    ax2.spines['right'].set_position(('axes', 1.0))
    grouped = df.groupby('network_size')['cost']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    ax2.errorbar(x, y, yerr=err, fmt='-o', lw=3, label='EdDSA')
    ax2.set_ylabel("cost[$]")

    grouped = df.groupby('network_size')['rsa_cost']
    x = grouped.mean().index.values
    y = grouped.mean().values
    err = grouped.std().values
    ax2.errorbar(x, y, yerr=err, fmt='-o', lw=3, label='RSA')

    ax1.spines['right'].set_visible(False)
    ax1.spines['top'].set_visible(False)
    ax1.set_xlabel('Network Size')
    #ax2.spines['right'].set_visible(False)
    ax2.spines['top'].set_visible(False)
    #ax3.spines['right'].set_visible(False)
    #ax3.spines['top'].set_visible(False)

    #ax1.yaxis.set_major_formatter(matplotlib.ticker.StrMethodFormatter('{x:,.0f}'))
    #ax1.xaxis.set_major_formatter(matplotlib.ticker.StrMethodFormatter('{x:,.0f}'))
    
    ax1.yaxis.set_major_formatter(FuncFormatter(human_readable_formatter))
    ax1.xaxis.set_major_formatter(FuncFormatter(human_readable_formatter))
    
    ax1.grid(True, color = "grey", linewidth = "0.3",axis = 'y')

    ax1.annotate('time=21s, cost=0.0012$', xy=(30000, 0), xytext=(22000, 3.5),
            arrowprops=dict(facecolor='black', arrowstyle='->'),
            #fontsize=12,
            ha='center')

    #ax.set_ylim(0, 1000)
    ax1.set_xlim(0, x[-1] + 1000)
    ax2.legend()

    
def plots_with_varying_sybil_number(path = '../experiment/logs/'):
    if(not os.path.exists(path)):
        print("Path: ", path, " doesn't exist - skipping")
        return

    # Read the most recently added log files in the given path
    command = "ls -t " + path + " | head -4" 
    process = Popen(command, stdout=PIPE, stderr=None, shell=True)
    output = process.communicate()[0]
    str_output = output.decode("utf-8")
    files = str_output.split('\n')
    print("Files: ", files)

    mitigationResultsFile = None
    overheadResultsFile = None
    discoveryResultsFile = None
    for filename in files:
        if filename.__contains__('mitigation'):
            mitigationResultsFile = filename
        elif filename.__contains__('overhead'):
            overheadResultsFile = filename
        elif filename.__contains__('discovery'):
            discoveryResultsFile = filename
        else:
            continue

    if mitigationResultsFile is None or overheadResultsFile is None:
        print ("Mitigation result file not found in: ", path)
        return
    
    # Merge dataframes for mitigation, overhead, and discovery
    df_mitigation = pd.read_csv(path + mitigationResultsFile)
    df_overhead = pd.read_csv(path + overheadResultsFile)
    df_discovery = pd.read_csv(path + discoveryResultsFile)
    df = pd.merge(df_mitigation, df_overhead)
    df = pd.merge(df, df_discovery)
    print(df)

    # Update and add new colmns
    df = df.rename(columns={"numOfSybils": "Number of Sybils"})

    df['mitigation']= df[list(df.filter(regex='mitigationSuccess'))].sum(axis=1)
    df['Percent Mitigation Success'] = df['mitigation']*100.0/len(list(df.filter(regex='contactedPeers')))
    
    df['Intersection Size'] = df[list(df.filter(regex='intersectionSize'))].mean(axis=1)
    df['Contacted Peers'] = df[list(df.filter(regex='contactedPeers'))].mean(axis=1)
    df['Updated Peers'] = df[list(df.filter(regex='updatedPeers'))].mean(axis=1)
    df['Updated Sybils'] = df[list(df.filter(regex='numSybilsUpdated'))].mean(axis=1)
    df['Percent Discovered'] = df[list(df.filter(regex='Discovery'))].mean(axis=1)

    df['Number of Lookups'] = df[list(df.filter(regex='numLookups'))].mean(axis=1)
    print(df)

    _, ax1 = plt.subplots(figsize=default_figsize)
    df.plot(x = 'Number of Sybils', y = ['Intersection Size', 'Contacted Peers', 'Updated Peers', 'Updated Sybils'], kind="line", rot=0, ax=ax1)
    ax1.set_xlabel("Number of Sybils")
    ax1.set_ylabel("Number of Peers")
    ax1.plot()
    # NOTE: must call savefig before show; otherwise, it doesn't save the fig
    plt.savefig('./plots/Overhead_vs_sybils.pdf')

    _, ax = plt.subplots(figsize=default_figsize)
    ax2 = df.plot.bar(x = 'Number of Sybils', y = ['Percent Mitigation Success','Percent Discovered' ], rot=0, ax=ax)
    ax2.set_xlabel("Number of Sybils")
    ax2.set_ylabel("Percentage")
    ax2.plot()
    plt.savefig('./plots/Mitigation_vs_Sybils.pdf')

    _, ax = plt.subplots(figsize=default_figsize)
    ax3 = df.plot.bar(x = 'Number of Sybils', y = 'Number of Lookups', rot=0, ax=ax)
    ax3.set_xlabel("Number of Sybils")
    ax3.set_ylabel("Number of DHT Lookups")
    ax3.get_legend().remove()
    ax3.plot()
    # NOTE: must call savefig before show; otherwise, it doesn't save the fig
    plt.savefig('./plots/Lookups_vs_Sybils.pdf')
    
    #ax4 = df.plot.bar(x = 'Number of Sybils', y = 'Percent Discovered', rot=0)
    #ax4.set_xlabel("Number of Sybils")
    #ax4.set_ylabel("Percentage")
    #ax4.plot()
    # NOTE: must call savefig before show; otherwise, it doesn't save the fig
    #plt.savefig('./plots/Discovery_vs_Sybils.pdf')
    
def plots_with_varying_region_size(path = '../experimentSpecialProvideNumber/logs/'):
    if(not os.path.exists(path)):
        print("Path: ", path, " doesn't exist - skipping")
        return

    # Read the most recently added log files in the given path
    command = "ls -t " + path + " | head -4" 
    process = Popen(command, stdout=PIPE, stderr=None, shell=True)
    output = process.communicate()[0]
    str_output = output.decode("utf-8")
    files = str_output.split('\n')
    print("Files: ", files)
    
    mitigationResultsFile = None
    overheadResultsFile = None
    for filename in files:
        if filename.__contains__('mitigation'):
            mitigationResultsFile = filename
        elif filename.__contains__('overhead'):
            overheadResultsFile = filename
        elif filename.__contains__('discovery'):
            discoveryResultsFile = filename
        else:
            continue

    if mitigationResultsFile is None or overheadResultsFile is None:
        print ("Mitigation result file not found in: ", path)
        return

    # Merge mitigation and overhead results into a single df
    df_mitigation = pd.read_csv(path + mitigationResultsFile)
    df_overhead = pd.read_csv(path + overheadResultsFile)
    df = pd.merge(df_mitigation, df_overhead)
    print(df)
    
    # Update and add new colmns
    df = df.rename(columns={"SpecialProvideNumber": "Percent Region Size"})
    df['Percent Region Size'] = df['Percent Region Size'] *100.0 / 20.0
    
    df['mitigation']= df[list(df.filter(regex='mitigationSuccess'))].sum(axis=1)
    print("Number of measurements for mitigation: ", len(list(df.filter(regex='mitigationSuccess'))))
    print(df['mitigation'])
    df['Percent Mitigation Success'] = df['mitigation']*100.0/len(list(df.filter(regex='mitigationSuccess')))
    df['Intersection Size'] = df[list(df.filter(regex='intersectionSize'))].mean(axis=1)
    df['Contacted Peers'] = df[list(df.filter(regex='contactedPeers'))].mean(axis=1)
    df['Updated Peers'] = df[list(df.filter(regex='updatedPeers'))].mean(axis=1)
    df['Updated Sybils'] = df[list(df.filter(regex='numSybilsUpdated'))].mean(axis=1)
    df['Number of Lookups'] = df[list(df.filter(regex='numLookups'))].mean(axis=1)
    print(df)
    
    # Number of Sybils is fixed to 20 so no need to plot updated Sybils
    _, ax = plt.subplots(figsize=default_figsize)
    ax1 = df.plot(x = 'Percent Region Size', y = ['Intersection Size', 'Contacted Peers', 'Updated Peers'], kind="line", rot=0, ax=ax)
    ax1.set_xlabel("Percent Region Size")
    ax1.set_ylabel("Number of Peers")
    ax1.plot()
    # NOTE: must call savefig before show; otherwise, it doesn't save the fig
    plt.savefig('./plots/Overhead_vs_regionSize.pdf')
    
    _, ax = plt.subplots(figsize=default_figsize)
    ax2 = df.plot.bar(x = 'Percent Region Size', y = 'Percent Mitigation Success', rot=0, ax=ax)
    ax2.set_xlabel("Percent Region Size")
    ax2.set_ylabel("Mitigation Success Percentage")
    ax2.get_legend().remove()
    ax2.plot()
    plt.savefig('./plots/Mitigation_vs_regionSize.pdf')
    
    _, ax = plt.subplots(figsize=default_figsize)
    ax3 = df.plot.bar(x = 'Percent Region Size', y = 'Number of Lookups', rot=0, ax=ax)
    ax3.set_xlabel("Percent Region Size")
    ax3.set_ylabel("Number of DHT Lookups")
    ax3.get_legend().remove()
    ax3.plot()
    # NOTE: must call savefig before show; otherwise, it doesn't save the fig
    plt.savefig('./plots/Lookups_vs_regionSize.pdf')

def main():
    sybil_generation_time2()
    #probability_dht_resolution()
    #mitigation_success_by_sybil_number()
    #plots_with_varying_sybil_number()
    #plots_with_varying_region_size()
    plt.show()

if __name__ == '__main__':
    main()
