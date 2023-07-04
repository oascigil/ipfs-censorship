# Artifacts for NDSS #153: Content Censorship in the InterPlanetary File System

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
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.10.linux-amd64.tar.gz
```
(You may need to run the command as root or through sudo).

**Do not** untar the archive into an existing /usr/local/go tree. This is known to produce broken Go installations.

1. Add /usr/local/go/bin to the PATH environment variable.
You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):
```
export PATH=$PATH:/usr/local/go/bin
```
Note: Changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just run the shell commands directly or execute them from the profile using a command such as source $HOME/.profile.

1. Verify that you've installed Go by opening a command prompt and typing the following command:
```
go version
```
Confirm that the command prints the installed version of Go.

### Install Python

### Install gcc and make
```
sudo apt install make
sudo apt install gcc
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
a)  **common/go-lib-kad-dht/** contains the network size computation, censorship detection and mitigation extensions for kademlia DHT. 
b) **common/go-libp2p-kad-dht-Sybil/** contains the Kademlia DHT code for Sybil nodes
c) **common/kubo** is the ipfs kubo code
d) **common/Sybil_DHT_Nodes** is the Sybil DHT node application.

1. **Figs13-15/** folder contains the experiment code, logs, and plot generation code for Figures 13, 14, and 15 of the paper. 

1. 