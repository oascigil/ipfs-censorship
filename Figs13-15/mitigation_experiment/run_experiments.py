import subprocess 

num_sybils = [20, 30, 40, 45]
region_sizes = [20]

for regSize in region_sizes:
    for numSybils in num_sybils:
        cmd = './main -sybils '
        cmd += str(numSybils) 
        cmd += ' -special '
        cmd += str(regSize)
        cmd += ' -outpath ./logs/'
        cmd += 'region' + str(regSize) + '_sybils' + str(numSybils) + '/'
        cmd += ' > out' + '_region' + str(regSize) + '_sybils' + str(numSybils) 
        print(cmd)
        subprocess.run(cmd, shell=True, check=True)
