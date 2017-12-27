#!/bin/bash

#Extracting NetID
sed -n '4,4 p' /home/node/start_*.sh > net
netid=$(awk -F'=' '{ print $2 }' net )
rm -rf net
echo $netid
