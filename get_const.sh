#!/bin/bash
username=$(whoami)
currentdir=$(pwd)

#copying the setup.conf file to the current working directory
cp /home/setup.conf .

#extracting the local IP and RPC port from setup.conf
awk -F'=' '{ print $2 }' setup.conf > orig-data-file
val1a=$(sed -n '1,1 p' orig-data-file)
val2a=$(sed -n '2,2 p' orig-data-file)
val4a=$(sed -n '4,4 p' orig-data-file)
rm -rf setup.conf
rm -rf orig-data-file

#sanitization procedure for removing whitespaces
ipaddr1="$(echo "$val1a" | sed -e 's/[[:space:]]*$//')"
rpcprt="$(echo "$val2a" | sed -e 's/[[:space:]]*$//')"
constel="$(echo "$val4a" | sed -e 's/[[:space:]]*$//')"

#Constellation Port is outputted to file to be read by Java Service
echo $constel

