#!/bin/bash

#Reading data sent by slave into variables
slave_enode1=$1

#Removing whitespaces
slave_enode="$(echo -e "$slave_enode1" | sed -e 's/[[:space:]]*$//')"

#Adding the slave to the network using the addPeer command
echo "var a = raft.addPeer(\"${slave_enode}\"); console.log(JSON.stringify(a));" > scriptAddPeer.js

geth --exec 'loadScript("scriptAddPeer.js")' attach ipc:/home/node/qdata/geth.ipc > addedPeer

readarray raftID < addedPeer
raft_id="$(echo -e "${raftID[0]}" | sed -e 's/[[:space:]]*$//')"
rm -rf scriptAddPeer.js
rm -rf addedPeer

#RaftID is outputted to file to be read by Java Service
echo $raft_id

