#!/bin/bash

#Extracting raftID
sed -n '5,5 p' /home/node/start_*.sh > raftid
echo $(awk -F'=' '{ print $2 }' raftid )
rm -rf raftid