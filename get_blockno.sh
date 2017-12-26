#!/bin/bash
#Getting the blocknumber of the node 
echo "var a = eth.blockNumber; console.log(a);" > scriptBlockNumber.js

geth --exec 'loadScript("scriptBlockNumber.js")' attach ipc:/home/node/qdata/geth.ipc > blockNo

readarray blockNoArr < blockNo
rm -rf blockNo
block_no="$(echo -e "${blockNoArr[0]}" | sed -e 's/[[:space:]]*$//')"

#RaftID is outputted to file to be read by Java Service
echo $block_no