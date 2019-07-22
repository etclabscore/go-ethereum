#!/bin/bash

set -ev

JWT=$1
BRANCH=$2

go get github.com/whiteblock/genesis
WORKING_DIR=(pwd)
cd $GOPATH/src/github.com/whiteblock/genesis
go build
cd $WORKING_DIR

wb login $JWT
echo "logged into whiteblock"

echo "building gethc node..."
wb build --git-repo https://github.com/etclabscore/go-ethereum.git --git-repo-branch $BRANCH -n 1 -m 0 -c 0 -y -p0=8570:8545

echo "building parity node..."
wb build append -b parity -n 1 -m 0 -c 0 -y -p1=8571:8545

echo "pulling latest testing script from github..."
git clone https://github.com/ChainSafe/Attalus.git
cd ./Attalus
yarn install

PK=(whiteblock get account info | jq '.[].privateKey'| sed -e 's/"//g')

echo "running script"
./attalus --pk $PK && wb done



