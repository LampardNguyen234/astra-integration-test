#!/bin/bash
astrad tendermint unsafe-reset-all --home=.testnets/node0/astrad
astrad tendermint unsafe-reset-all --home=.testnets/node1/astrad
astrad tendermint unsafe-reset-all --home=.testnets/node2/astrad
astrad tendermint unsafe-reset-all --home=.testnets/node3/astrad

rm -rf .testnets/node0/astrad/cosmovisor/current
rm -rf .testnets/node1/astrad/cosmovisor/current
rm -rf .testnets/node2/astrad/cosmovisor/current
rm -rf .testnets/node3/astrad/cosmovisor/current

cp ~/go/bin/2astrad .testnets/node0/astrad/cosmovisor/genesis/bin/astrad
cp ~/go/bin/2astrad .testnets/node1/astrad/cosmovisor/genesis/bin/astrad
cp ~/go/bin/2astrad .testnets/node2/astrad/cosmovisor/genesis/bin/astrad
cp ~/go/bin/2astrad .testnets/node3/astrad/cosmovisor/genesis/bin/astrad

rm -rf .testnets/node0/astrad/cosmovisor/upgrades/v3.0.0
rm -rf .testnets/node1/astrad/cosmovisor/upgrades/v3.0.0
rm -rf .testnets/node2/astrad/cosmovisor/upgrades/v3.0.0
rm -rf .testnets/node3/astrad/cosmovisor/upgrades/v3.0.0