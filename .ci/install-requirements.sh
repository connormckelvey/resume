#! /usr/bin/env bash

set -ex;

# Install google-chrome in ubuntu
sudo wget http://www.linuxidc.com/files/repo/google-chrome.list -P /etc/apt/sources.list.d/
wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
sudo apt-get update
sudo apt-get install -y google-chrome-stable

# Install aspell
sudo apt-get install -y aspell
