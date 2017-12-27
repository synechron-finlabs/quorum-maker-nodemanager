#!/bin/bash

#copying the setup.conf file to the current working directory
cp /home/setup.conf .

#extracting the constellation port from setup.conf
awk -F'=' '{ print $2 }' setup.conf > orig-data-file
val4a=$(sed -n '4,4 p' orig-data-file)
rm -rf setup.conf
rm -rf orig-data-file

#sanitization procedure for removing whitespaces
constel="$(echo "$val4a" | sed -e 's/[[:space:]]*$//')"

echo $constel
