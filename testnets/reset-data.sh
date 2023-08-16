#!/bin/bash
# making home_dir
HOME_DIR="$1"
if [[ -z $HOME_DIR ]];
then
  HOME_DIR=`pwd`
fi
echo "HOME_DIR=$HOME_DIR"

astrad tendermint unsafe-reset-all --home=$HOME_DIR/.testnets/node0/astrad
astrad tendermint unsafe-reset-all --home=$HOME_DIR/.testnets/node1/astrad

rm -rf $HOME_DIR/.testnets/node0/astrad/cosmovisor/current
rm -rf $HOME_DIR/.testnets/node1/astrad/cosmovisor/current

cp ~/go/bin/2astrad $HOME_DIR/.testnets/node0/astrad/cosmovisor/genesis/bin/astrad
cp ~/go/bin/2astrad $HOME_DIR/.testnets/node1/astrad/cosmovisor/genesis/bin/astrad

rm -rf $HOME_DIR/.testnets/node0/astrad/cosmovisor/upgrades/
mkdir -p $HOME_DIR/.testnets/node0/astrad/cosmovisor/upgrades/

rm -rf $HOME_DIR/.testnets/node1/astrad/cosmovisor/upgrades/
mkdir -p $HOME_DIR/.testnets/node1/astrad/cosmovisor/upgrades/