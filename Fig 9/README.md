# Generating Figure 9

Below you can find instructions to (i) run experiments to generate data to plot Figure 9 and (ii) to generate the plots.

## Generating data
Run two sets of experiments measuring generation time for RSA and EDDSA:
`go run measureGenerateSybilKeys rsa > timing_rsa.csv`
`go run measureGenerateSybilKeys eddsa > timing_eddsa.csv`

We provide our `sample_rsa.csv` and `sample_eddsa.csv` files for reference (generating RSA keys takes a while).

## Plotting
Plot the results using:
`python3 plot.py`

Note that the script will look for `timing_rsa.csv` and `timing_eddsa.csv` files, so you have to generate them in advance