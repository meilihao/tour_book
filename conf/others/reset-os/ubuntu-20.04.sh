#!/usr/bin/env bash
echo "tmpfs /home/chen/tmpfs tmpfs size=256m 0 0" | sudo tee -a /etc/fstab
sudo ln -s /home/chen/opt/sources.list.d/ubuntu_2004/ /etc/apt/sources.list.d
sudo apt-key add /home/chen/opt/sources.list.d/ubuntu_2004/keys/*.gpg
sudo apt update
sudo apt install wget git locate vim tilix openvpn qemu-system-x86 bcompare code
sudo apt install clang-12 lldb-12 lld-12 gcc
sudo apt install mariadb-server
sudo apt install maven openjdk-8-jdk openjdk-11-jdk
# sudo apt install unixbench
