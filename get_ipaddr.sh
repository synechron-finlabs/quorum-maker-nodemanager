#!/bin/bash

#copying the setup.conf file to the current working directory
cp /home/setup.conf .

#extracting the IP address from setup.conf
awk -F'=' '{ print $2 }' setup.conf > orig-data-file

val1a=$(sed -n '1,1 p' orig-data-file)

rm -rf setup.conf
rm -rf orig-data-file

#sanitization procedure for removing whitespaces
ipaddr="$(echo "$val1a" | sed -e 's/[[:space:]]*$//')"

#Constellation Port is outputted to file to be read by Java Service
echo $ipaddr