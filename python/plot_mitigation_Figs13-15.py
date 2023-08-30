import pandas as pd
from matplotlib import pyplot as plt
import os.path
import numpy as np
import argparse
from subprocess import PIPE, Popen
from style import *

from style import *

def count_sybils(row):
    if isinstance(row['SybilIDs'], list):
        return len(row['SybilIDs'])
    else:
        return 0

def plots_with_varying_sybil_size(path = '../mitigation_experiment/logs/'):
    
    if(not os.path.exists(path)):
        print("Path: ", path, " doesn't exist - skipping")
        return

    dir_list = [d for d in os.listdir(path) if os.path.isdir(os.path.join(path, d))]


    print(dir_list)

    df = pd.DataFrame()
    for log in dir_list:
        full_path = os.path.join(path,log) + '/'
        print(full_path)
        # Read the most recently added log files in the given path
        command = "ls " + full_path
        process = Popen(command, stdout=PIPE, stderr=None, shell=True)
        output = process.communicate()[0]
        str_output = output.decode("utf-8")
        files = str_output.split('\n')
        files.remove("")
        print("Files: ", files)

        for file in files:
            with open(full_path + file) as f:
                print('Reading...', full_path + file)
                f.readlines(1)
                df_json = pd.read_json(f, lines=True)
                df = pd.concat([df, df_json], axis=0)


    # Update and add new colmns
    df = df.rename(columns={"SpecialProvideNumber": "Region Size"})
    
    df['mitigation']= df[list(df.filter(regex='IsMitigated'))].sum(axis=1)
    print("Number of measurements for mitigation: ", len(list(df.filter(regex='IsMitigated'))))
    df['Percent Mitigation Success'] = df['mitigation']*100.0/len(list(df.filter(regex='IsMitigated')))
    df['Contacted Sybils'] = df.apply(lambda row: len(set(row['ContactedPeers']).intersection(set(row['SybilIDs']))) if isinstance(row['SybilIDs'], list) else 0, axis=1)
    df['Intersection Size'] = df[list(df.filter(regex='NumIntersection'))].mean(axis=1)
    df['Contacted Peers'] = df[list(df.filter(regex='NumContacted'))].mean(axis=1)
    df['Updated Peers'] = df[list(df.filter(regex='NumUpdated'))].mean(axis=1)
    df['Updated Sybils'] = df[list(df.filter(regex='NumSybilsUpdated'))].mean(axis=1)
    df['Number of Lookups'] = df[list(df.filter(regex='NumLookups'))].mean(axis=1)
    df['Percent Reachable w/o Mitigation'] = 100.0 - df[list(df.filter(regex='PercentEclipsed'))].mean(axis=1)
    #df['Number of Sybils'] = df['SybilIDs'].apply(lambda x: len(x))
    df['Number of Sybils'] = df.apply(lambda row: count_sybils(row), axis=1)


    print(df)

    agg_funcs = {'Percent Mitigation Success': ['mean', 'std'],
             'Intersection Size': ['mean', 'std'],
             'Contacted Peers': ['mean', 'std'],
             'Updated Peers': ['mean', 'std'],
             'Updated Sybils': ['mean', 'std'],
             'Number of Lookups': ['mean', 'std'],
             'Contacted Sybils': ['mean', 'std'],
             'Updated Sybils': ['mean', 'std'],  
             'Percent Reachable w/o Mitigation':['mean','std']}

    df = df.groupby('Number of Sybils')[['Percent Mitigation Success', 'Intersection Size', 'Contacted Peers', 'Updated Peers', 'Updated Sybils', 'Number of Lookups', 'Contacted Sybils', 'Updated Sybils', 'Percent Reachable w/o Mitigation']].agg(agg_funcs)

    #df = df.groupby('Number of Sybils')[['Percent Mitigation Success', 'Intersection Size', 'Contacted Peers', 'Updated Peers', 'Updated Sybils', 'Number of Lookups']].mean()
    df = df.reset_index()

    #df.to_csv("test.csv")
    #exit(0)
    
    ####################  Contacted, updated, intersection plot ##################
    _, ax1 = plt.subplots(figsize=default_figsize)
    
    # define hatch pattern for Sybil portion of bars
    sybil_hatch = "////"
    
    #x_values = df['Number of Sybils']
    x_values = list(range(0, 10*df['Number of Sybils'].nunique(), 10))

    # Define the y-axis values for each bar
    y1_mean = df['Intersection Size']['mean']
    y1_std = df['Intersection Size']['std']

    y2_mean = df['Contacted Peers']['mean'] - df['Contacted Sybils']['mean']
    y2_std = df['Contacted Peers']['std']
    y2Sybil_mean = df['Contacted Sybils']['mean']
    y2Sybil_std = df['Contacted Sybils']['std']

    y3_mean = df['Updated Peers']['mean'] - df['Updated Sybils']['mean']
    y3_std = df['Updated Peers']['std']
    y3Sybil_mean = df['Updated Sybils']['mean']
    y3Sybil_std = df['Updated Sybils']['std']
    

    # Set the width of each bar
    bar_width = 2.0

    # Set the positions of the bars on the x-axis
    bar1_positions = [x - bar_width for x in x_values]
    bar2_positions = x_values
    bar3_positions = [x + bar_width for x in x_values]

    # Create the three bars for each datapoint
    # Plot the first bar
    bar1 = ax1.bar(bar1_positions, y1_mean, width=bar_width, yerr=y1_std, capsize=5, color='g', align='center')

    # Plot the second bar
    bar2_sybil = ax1.bar(bar2_positions, y2Sybil_mean, width=bar_width, yerr=y2Sybil_std, capsize=5, color='r',align='center')
    bar2 = ax1.bar(bar2_positions, y2_mean, width=bar_width, bottom=y2Sybil_mean, yerr=y2_std, capsize=5, color='b', align='center')

    # Plot the third bar
    bar3 = ax1.bar(bar3_positions, y3_mean, width=bar_width, bottom=y3Sybil_mean,  yerr=y3_std, capsize=5, color='darkorange', align='center')
    bar3_sybil = ax1.bar(bar3_positions, y3Sybil_mean, width=bar_width, yerr=y3Sybil_std, capsize=5, color='r',  align='center')

    # Add the value of each bar on top of the bar
    for bar in [bar1, bar2, bar3]:
        bar_indx = 0
        for rect in bar:
            height = rect.get_height()
            rotation = 0
            if bar == bar1 or bar == bar3:
                # For the red and blue bars, rotate the text by 90 degrees
                rotation = 90
            if bar == bar1:
                ax1.text(rect.get_x() + rect.get_width()/2. - 1.3, height + 0.2, '%.2f' % height, ha='center', va='bottom', rotation=rotation)
            elif bar == bar3:
                rect_sybil = bar3_sybil[bar_indx]
                height += rect_sybil.get_height()
                ax1.text(rect.get_x() + rect.get_width()/2. + 1.3, height + 0.2, '%.2f' % height, ha='center', va='bottom', rotation=rotation)
            else: # bar2
                rect_sybil = bar2_sybil[bar_indx]
                height += rect_sybil.get_height()
                ax1.text(rect.get_x() + rect.get_width()/2., height + 0.1, '%.2f' % height, ha='center', va='bottom', rotation=rotation)
            bar_indx += 1

    # Set the x-axis label and tick labels
    ax1.set_xlabel('Number of Sybils')
    ax1.set_xticks(x_values)
    ax1.set_xticklabels(df['Number of Sybils'])

    # Set the y-axis label
    ax1.set_ylabel('Number of Peers')
    ax1.spines['top'].set_visible(False)
    ax1.spines['right'].set_visible(False)

    # Add a legend
    #ax1.legend(['Intersection Size', 'Sybils', 'Contacted Peers', 'Updated Peers'])
    #plt.legend(loc="upper left")
    ax1.legend(['Non-Sybil intersection', 'Sybils (subset)', 'Contacted peers', 'Updated peers' ], bbox_to_anchor=(0, 1.08), loc='upper left', ncol=1)    

    #plt.figure(figsize=default_figsize)
    plt.savefig('./plots/Overhead_vs_sybils.pdf', bbox_inches='tight')
    
    ####################  Mitigation success plot ##################
    _, ax2 = plt.subplots(figsize=default_figsize)
    bar_width = 6
    # Set the positions of the bars on the x-axis
    # filter the data values for NumOfSybils = 0 
    df_filtered = df[df['Number of Sybils'] != 0]
    #df_filtered.loc[df['Number of Sybils'] == 20, 'Mitigation Success'] = 0.32
    x = range(0, 20*df_filtered['Number of Sybils'].nunique(), 20)
    bar1_positions = [x - bar_width/2 for x in x]
    bar2_positions = [x + bar_width/2 for x in x]
    bar_plot1 = ax2.bar(bar2_positions, df_filtered['Percent Mitigation Success']['mean'], width=bar_width, color='b', label='Successfully mitigated')
    bar_plot2 = ax2.bar(bar1_positions, df_filtered['Percent Reachable w/o Mitigation']['mean'], width=bar_width, color='g', label='CID discoverable w/o mitigation')

    # add the mean value to the top of each bar
    for bar_plot in [bar_plot1, bar_plot2]:
        for p in bar_plot:
            if bar_plot == bar_plot2:
                ax2.annotate(f"{p.get_height():.2f}", (p.get_x() + p.get_width() / 2. - 0.5, p.get_height()), ha='center', va='bottom')

            else:
                ax2.annotate(f"{p.get_height():.2f}", (p.get_x() + p.get_width() / 2, p.get_height()), ha='center', va='bottom')

    # set the x-tick positions and labels below each bar
    #ax2.set_xticks(df_filtered['Number of Sybils'])
    ax2.set_xticks(x)

    ax2.set_xticklabels(df_filtered['Number of Sybils'])

    # set the x and y axis labels
    ax2.set_xlabel('Number of Sybils')
    ax2.set_ylabel('Percentage of Queries')
    ax2.spines['top'].set_visible(False)
    ax2.spines['right'].set_visible(False)
    plt.legend(loc="upper left")

    #ax2.get_legend().remove()
    ax2.plot()
    #plt.figure(figsize=default_figsize)
    plt.savefig('./plots/Mitigation_vs_Sybils.pdf', bbox_inches='tight')
    

    ####################  Number of lookups plot ##################
    x_values = df['Number of Sybils']
    y_means = df['Number of Lookups']['mean']
    y_stds = df['Number of Lookups']['std']
    # Set the positions of the bars on the x-axis
    x = range(0, 10*df['Number of Sybils'].nunique(), 10)
    _, ax3 = plt.subplots(figsize=default_figsize)
    #ax3 = df.plot.bar(x = 'Number of Sybils', y = 'Number of Lookups', rot=0, ax=ax)
    bar_plot = ax3.bar(x, y_means, width=bar_width, color='b', align='center', yerr=y_stds, capsize=5)
    #ax3.bar_label(ax3.containers[0], labels=[f"{val:.2f}" for val in y_means], padding=2)

    for p in bar_plot:
        ax3.annotate(f"{p.get_height():.2f}", (p.get_x() + p.get_width() / 2, p.get_height()), ha='center', va='bottom')
    ax3.set_xlabel("Number of Sybils")
    ax3.set_ylabel("Number of DHT Lookups")
    ax3.set_xticks(x)
    ax3.set_xticklabels(x_values)
    #ax3.get_legend().remove()
    ax3.plot()
    ax3.spines['top'].set_visible(False)
    ax3.spines['right'].set_visible(False)
    # NOTE: must call savefig before show; otherwise, it doesn't save the fig
    #plt.figure(figsize=default_figsize)
    plt.savefig('./plots/Lookups_vs_Sybils.pdf', bbox_inches='tight')
    
#plt.show()
    
def main():
    if not os.path.exists("plots"): 
        os.makedirs("plots")
    argParser = argparse.ArgumentParser()
    argParser.add_argument("-i", "--input", help="input experiment results path", default="../experimentCombined/mitigation_results/")
    args = argParser.parse_args()
    plots_with_varying_sybil_size(args.input)
    print("Plots saved: ./plots/Lookups_vs_Sybils.pdf ./plots/Mitigation_vs_Sybils.pdf ./plots/Overhead_vs_sybils.pdf")

if __name__ == '__main__':
    main()
