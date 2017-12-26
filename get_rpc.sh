#!/bin/bash

#copying the setup.conf file to the current working directory
cp /home/setup.conf .

#extracting the RPC port from setup.conf
awk -F'=' '{ print $2 }' setup.conf > orig-data-file

val2a=$(sed -n '2,2 p' orig-data-file)

rm -rf setup.conf
rm -rf orig-data-file

#sanitization procedure for removing whitespaces
rpcprt="$(echo "$val2a" | sed -e 's/[[:space:]]*$//')"

#Constellation Port is outputted to file to be read by Java Service
echo $rpcprt