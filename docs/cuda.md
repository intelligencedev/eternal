# CUDA Dependencies

The commands listed below are a general reference on how to setup CUDA in an Ubuntu Linux based system. Eternal has been tested to work with ubuntu via WSL2 on a Windows 11 based host with a CUDA capable GPU.

```
$ sudo apt install cmake build-essential libtool autoconf unzip wget
$ wget https://www.nvidia.com/content/DriverDownloads/confirmation.php?url=/XFree86/Linux-x86_64/535.98/NVIDIA-Linux-x86_64-535.98.run&lang=us&type=geforcem
$ wget https://developer.download.nvidia.com/compute/cuda/repos/wsl-ubuntu/x86_64/cuda-wsl-ubuntu.pin
$ sudo mv cuda-wsl-ubuntu.pin /etc/apt/preferences.d/cuda-repository-pin-600
$ wget https://developer.download.nvidia.com/compute/cuda/12.3.1/local_installers/cuda-repo-wsl-ubuntu-12-3-local_12.3.1-1_amd64.deb
$ sudo dpkg -i cuda-repo-wsl-ubuntu-12-3-local_12.3.1-1_amd64.deb
$ sudo cp /var/cuda-repo-wsl-ubuntu-12-3-local/cuda-96064797-keyring.gpg /usr/share/keyrings/
$ sudo apt-get update
$ sudo apt-get -y install cuda-toolkit-12-3
$ sudo apt install build-essential checkinstall zlib1g-dev libssl-dev
```