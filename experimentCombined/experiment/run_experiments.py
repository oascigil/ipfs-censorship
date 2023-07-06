import subprocess 

num_sybils = [45, 40, 30, 20, 0]
region_sizes = [20]

for regSize in region_sizes:
    for numSybils in num_sybils:
        cmd = './main -sybils '
        cmd += str(numSybils) 
        cmd += ' -region '
        cmd += str(regSize)
        cmd += ' -cids 10'
        cmd += ' -clients 1'
        cmd += ' -outpath ../experiment_results_new/'
        print(cmd)
        subprocess.run(cmd, shell=True, check=True)

# Commands that are being executed:
# ./main -sybils 0 -region 20 -cids 10 -clients 1 -outpath ../experiment_results_new/
# ./main -sybils 20 -region 20 -cids 10 -clients 1 -outpath ../experiment_results_new/
# ./main -sybils 30 -region 20 -cids 10 -clients 1 -outpath ../experiment_results_new/
# ./main -sybils 40 -region 20 -cids 10 -clients 1 -outpath ../experiment_results_new/
# ./main -sybils 45 -region 20 -cids 10 -clients 1 -outpath ../experiment_results_new/
