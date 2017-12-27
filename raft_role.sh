#!/bin/bash

#Getting the raft role of the node 
echo "var a = raft; console.log(JSON.stringify(a));" > scriptRaftRole.js

geth --exec 'loadScript("scriptRaftRole.js")' attach ipc:/home/node/qdata/geth.ipc > raftRole

readarray raftroleArr < raftRole
rm -rf raftRole

raft_role="$(echo -e "${raftroleArr[0]}" | sed -e 's/[[:space:]]*$//')"
echo $raft_role > raft_role

echo $(awk -F'"' '{ print $4 }' raft_role )
rm -rf raft_role