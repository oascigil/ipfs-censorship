import subprocess 

num_sybils = [0, 20, 30, 40, 45]
region_sizes = [20]

for regSize in region_sizes:
    for numSybils in num_sybils:
        cmd = './main -sybils '
        cmd += str(numSybils) 
        cmd += ' -region '
        cmd += str(regSize)
        cmd += ' -cids 5'
        cmd += ' -clients 5'
        cmd += ' -outpath ./experiment_results_new/'
        print(cmd)
        subprocess.run(cmd, shell=True, check=True)
