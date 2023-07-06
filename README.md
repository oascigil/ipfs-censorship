# Artifacts for NDSS #153: Content Censorship in the InterPlanetary File System
## SSH Instructions
For the purpose of evaluating this artifact during the NDSS Artifact evaluation process, we provide access to an AWS instance on which the reviewers can run the experiments in this artifact. To control SSH access, we provide a private key file **ndss24-ae-#14.pem** to the reviewers. Store this file in your local machine. Adjust the permissions of the private key file by running the following command once.
```
chmod 400 ndss24-ae-#14.pem
```
To log in to the AWS instance, use the command
```
ssh -i "ndss24-ae-#14.pem" ubuntu@ec2-54-149-215-250.us-west-2.compute.amazonaws.com
```
## Network Setup
If you are using the AWS machine mentioned above, then you can skip this step.

For the experiments involving launching a censorship attack, the Sybils nodes must be launched in 'server mode'. This requires allowing other DHT clients and servers to connect to the Sybil peers. This requires the machine hosting the Sybil peers to have a publicly dialable IP address, i.e., it must not be behind a NAT. We use TCP ports 63800 onwards for the Sybil peers and IPFS uses ports 4001 and 5001, so make sure to allow these ports if your machine has a firewall. This is usually done using
```
sudo ufw allow 63800:63900/tcp
```
In some cases, this might happen through a different network administration interface. To check if the machine is publicly dialable and has the required ports open, you can use the tool `nc`. Run the following to listen on port 63800 on the machine on which the experiments will be run:
```
nc -vvv -l -p 62800
```
On any other machine, try to connect to the experiment machine with IP address `ip_addr` by running:
```
nc -vvv ip_addr 62800
```
If the connection was successful, then the experiment machine should be dialable with the required ports open.

## Install Required Dependencies

Running our artifacts require installing Go and Python. Below are instructions for Linux. First, start with updating your system.
```
sudo apt update
sudo apt upgrade
```
### Install Go
The following instructions for Linux are taken from [go.dev/doc/install](https://go.dev/doc/install). Visit the link for other OSes. Please use Go version 1.19.10. One of the dependencies of **kubo** is [not compatible with Go 1.20](https://github.com/quic-go/quic-go/wiki/quic-go-and-Go-versions).

1. Download the source for Go.
```
wget https://go.dev/dl/go1.19.10.linux-amd64.tar.gz
```
Remove any previous Go installation by deleting the /usr/local/go folder (if it exists), then extract the archive you just downloaded into /usr/local, creating a fresh Go tree in /usr/local/go:
```
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.19.10.linux-amd64.tar.gz
```
(You may need to run the command with `sudo`).

**Do not** untar the archive into an existing /usr/local/go tree. This is known to produce broken Go installations.

1. Add /usr/local/go/bin to the PATH environment variable.
You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):
```
export PATH=$PATH:/usr/local/go/bin
```
Note: Changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just run the shell commands directly or execute them from the profile using a command such as source $HOME/.profile.

If you ever find the error "Command 'go' not found", try running the above command again to update the PATH variable.

1. Verify that you've installed Go by opening a command prompt and typing the following command:
```
go version
```
Confirm that the command prints the installed version of Go.
### Install gcc and make
```
sudo apt install make
sudo apt install gcc
```
### Install Python
We recommend installing Python 3.10. Most machines will already have Python 3. Check if your machine has it by running `python3 --version`. Also, install `pip` and `venv` if you do not already have it. We recommend `pip` version 22.0.2 if you are using Python 3.10.
```
sudo apt install python3.10
sudo apt install -y python3-pip
sudo apt install -y python3-venv
```
Move to the **python/** directory, create a new virtual environment, and install all the required modules.
```
cd python
python3 -m venv env
source env/bin/activate
pip install -r requirements.txt
```
### Build and initialize IPFS
We run an IPFS node to obtain information such as the closest peers to the target ID to set up the attack. The **kubo** implementation of IPFS source code is already copied in this repository. We need to build the source code and initialize the IPFS node. We write these instructions assuming you start from the directory which contains this README file.
```
cd common/kubo
make build
./cmd/ipfs/ipfs init
```
If you see the following error, it can be safely ignored.
```
ERROR   provider.queue  queue/queue.go:125      Failed to enqueue cid: leveldb: closed
```
To test that your IPFS node is set up correctly, try the following:
```
./cmd/ipfs/ipfs cat /ipfs/QmQPeNsJPyVWPFDVHb77w8G42Fvo15z4bG2X8D2GhfbSXc/readme
```
## Outline

1. **common/** folder:
    a) **common/go-libp2p-kad-dht/** contains the network size computation, censorship detection and mitigation extensions for the Kademlia DHT. 
    b) **common/go-libp2p-kad-dht-Sybil/** contains the Kademlia DHT code for Sybil nodes.
    c) **common/kubo** is the IPFS kubo code.
    d) **common/Sybil_DHT_Nodes** is the Sybil DHT node application.
1. **experimentCombined/** folder:
    a) **experiment/** contains code to launch Sybil nodes to attack several target CIDs, attempt to find providers to check for the attack's success, perform attack detection, and perform mitigation.
    b) **detection_results/** and **mitigation/results** contain results from experiments that we have already run and that we reported in the paper.
1. **experimentLatency/** contains code to measure the latency of the standard `Provide` and `FindProviders` operations as well as those using our mitigation, under attack and no attack.
1. **python/** contains code to process results collected from the above experiments and generate plots as in the paper. It also contains simulations for other results that are not obtained from the above two experiments.
1. **Fig 8/** contains code to launch the censorship attack after the content has been provided and observe the success of the attack over a period of 50 hours (Fig. 8 in the paper).

To reproduce the main results of our paper, first follow the instructions in **experimentCombined/** to generate the experiment results, then move to **python/** for instructions to generate each figure in the paper. 