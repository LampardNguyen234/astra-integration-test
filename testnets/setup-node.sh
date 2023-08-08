#!/bin/bash
rm -rf .testnets/
kill -9 $(pgrep -f "rpc.laddr=tcp://127.0.0.1") || true
mkdir -p .testnets

# start a testnet
2astrad testnet init-files --keyring-backend=test --chain-id=astra_11115-1

# change staking denom to ubaby
cat .testnets/node0/astrad/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aastra"' > .testnets/node0/astrad/config/tmp_genesis.json && mv .testnets/node0/astrad/config/tmp_genesis.json .testnets/node0/astrad/config/genesis.json

# update crisis variable to ubaby
cat .testnets/node0/astrad/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aastra"' > .testnets/node0/astrad/config/tmp_genesis.json && mv .testnets/node0/astrad/config/tmp_genesis.json .testnets/node0/astrad/config/genesis.json

# udpate gov genesis
cat .testnets/node0/astrad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aastra"' > .testnets/node0/astrad/config/tmp_genesis.json && mv .testnets/node0/astrad/config/tmp_genesis.json .testnets/node0/astrad/config/genesis.json
cat .testnets/node0/astrad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["max_deposit_period"]="10s"' > .testnets/node0/astrad/config/tmp_genesis.json && mv .testnets/node0/astrad/config/tmp_genesis.json .testnets/node0/astrad/config/genesis.json
cat .testnets/node0/astrad/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' > .testnets/node0/astrad/config/tmp_genesis.json && mv .testnets/node0/astrad/config/tmp_genesis.json .testnets/node0/astrad/config/genesis.json

# update mint genesis
cat .testnets/node0/astrad/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aastra"' > .testnets/node0/astrad/config/tmp_genesis.json && mv .testnets/node0/astrad/config/tmp_genesis.json .testnets/node0/astrad/config/genesis.json

#### change config.toml values
# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' .testnets/node0/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:40003|tcp://0.0.0.0:40000|g' .testnets/node0/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:50003|tcp://0.0.0.0:50000|g' .testnets/node0/astrad/config/config.toml
# validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26558|g' .testnets/node1/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:40003|tcp://0.0.0.0:40001|g' .testnets/node1/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:50003|tcp://0.0.0.0:50001|g' .testnets/node1/astrad/config/config.toml
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' .testnets/node1/astrad/config/config.toml
sed -i -E 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:2317"|g' .testnets/node1/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9190"|g' .testnets/node1/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9191"|g' .testnets/node1/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:8545"|address = "0.0.0.0:18545"|g' .testnets/node1/astrad/config/app.toml
# validator3
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26458|g' .testnets/node2/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:40003|tcp://0.0.0.0:40002|g' .testnets/node2/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:50003|tcp://0.0.0.0:50002|g' .testnets/node2/astrad/config/config.toml
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' .testnets/node2/astrad/config/config.toml
sed -i -E 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:3317"|g' .testnets/node2/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9290"|g' .testnets/node2/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9291"|g' .testnets/node2/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:8545"|address = "0.0.0.0:28545"|g' .testnets/node2/astrad/config/app.toml

# validator 4
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26358|g' .testnets/node3/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:40003|tcp://0.0.0.0:40003|g' .testnets/node3/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:50003|tcp://0.0.0.0:50003|g' .testnets/node3/astrad/config/config.toml
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' .testnets/node3/astrad/config/config.toml
sed -i -E 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:4317"|g' .testnets/node3/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9390"|g' .testnets/node3/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9391"|g' .testnets/node3/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:8545"|address = "0.0.0.0:38545"|g' .testnets/node3/astrad/config/app.toml

sed -i -E 's|timeout_commit = "5s"|timeout_commit = "1s"|g' .testnets/node0/astrad/config/config.toml
sed -i -E 's|timeout_commit = "5s"|timeout_commit = "1s"|g' .testnets/node1/astrad/config/config.toml
sed -i -E 's|timeout_commit = "5s"|timeout_commit = "1s"|g' .testnets/node2/astrad/config/config.toml
sed -i -E 's|timeout_commit = "5s"|timeout_commit = "1s"|g' .testnets/node3/astrad/config/config.toml

# copy validator1 genesis file to validator2-3-4
cp .testnets/node0/astrad/config/genesis.json .testnets/node1/astrad/config/genesis.json
cp .testnets/node0/astrad/config/genesis.json .testnets/node2/astrad/config/genesis.json
cp .testnets/node0/astrad/config/genesis.json .testnets/node3/astrad/config/genesis.json

# create cosmosvisor folders
mkdir -p .testnets/node0/astrad/cosmovisor
mkdir -p .testnets/node0/astrad/data
mkdir -p .testnets/node0/astrad/cosmovisor/genesis
mkdir -p .testnets/node0/astrad/cosmovisor/genesis/bin
mkdir -p .testnets/node0/astrad/cosmovisor/upgrades
cp ~/go/bin/2astrad .testnets/node0/astrad/cosmovisor/genesis/bin/astrad

mkdir -p .testnets/node1/astrad/cosmovisor
mkdir -p .testnets/node1/astrad/data
mkdir -p .testnets/node1/astrad/cosmovisor/genesis
mkdir -p .testnets/node1/astrad/cosmovisor/genesis/bin
mkdir -p .testnets/node1/astrad/cosmovisor/upgrades
cp ~/go/bin/2astrad .testnets/node1/astrad/cosmovisor/genesis/bin/astrad

mkdir -p .testnets/node2/astrad/cosmovisor
mkdir -p .testnets/node2/astrad/data
mkdir -p .testnets/node2/astrad/cosmovisor/genesis
mkdir -p .testnets/node2/astrad/cosmovisor/genesis/bin
mkdir -p .testnets/node2/astrad/cosmovisor/upgrades
cp ~/go/bin/2astrad .testnets/node2/astrad/cosmovisor/genesis/bin/astrad

mkdir -p .testnets/node3/astrad/cosmovisor
mkdir -p .testnets/node3/astrad/data
mkdir -p .testnets/node3/astrad/cosmovisor/genesis
mkdir -p .testnets/node3/astrad/cosmovisor/genesis/bin
mkdir -p .testnets/node3/astrad/cosmovisor/upgrades
cp ~/go/bin/2astrad .testnets/node3/astrad/cosmovisor/genesis/bin/astrad

#
#tmux new -s node0 -d 2astrad start --home=.testnets/node0/astrad --rpc.laddr=tcp://127.0.0.1:26657 2>&1 | tee .testnets/node0/log.log
#tmux new -s node1 -d 2astrad start --home=.testnets/node1/astrad --rpc.laddr=tcp://127.0.0.1:26557
#tmux new -s node2 -d 2astrad start --home=.testnets/node2/astrad --rpc.laddr=tcp://127.0.0.1:26457
#tmux new -s node3 -d 2astrad start --home=.testnets/node3/astrad --rpc.laddr=tcp://127.0.0.1:26357
#tmux set remain-on-exit on
#tmux ls

echo "node0:" $(2astrad keys unsafe-export-eth-key node0 --keyring-backend=test --home=.testnets/node0/astrad) $(2astrad keys show node0 -a --keyring-backend=test --home=.testnets/node0/astrad)
echo "node1:" $(2astrad keys unsafe-export-eth-key node1 --keyring-backend=test --home=.testnets/node1/astrad) $(2astrad keys show node1 -a --keyring-backend=test --home=.testnets/node1/astrad)
echo "node2:" $(2astrad keys unsafe-export-eth-key node2 --keyring-backend=test --home=.testnets/node2/astrad) $(2astrad keys show node2 -a --keyring-backend=test --home=.testnets/node2/astrad)
echo "node3:" $(2astrad keys unsafe-export-eth-key node3 --keyring-backend=test --home=.testnets/node3/astrad) $(2astrad keys show node3 -a --keyring-backend=test --home=.testnets/node3/astrad)