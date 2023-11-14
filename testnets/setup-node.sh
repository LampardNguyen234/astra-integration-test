#!/bin/bash

KEY_GENESIS_VALIDATOR="genesis_validator"
KEY_PARTNER="genesis_partner"
KEY_STRATEGIC="strategic"
KEY_CORE_TEAM="core_team"
CHAIN_ID="astra_11115-1"
KEYRING="test"
KEYALGO="eth_secp256k1"

# making home_dir
HOME_DIR="$1"
if [[ -z $HOME_DIR ]];
then
  export HOME_DIR=`pwd`
fi
echo "HOME_DIR=$HOME_DIR"

# installing cosmovisor
if [[ ! -f $HOME/go/bin/cosmovisor ]];
then
  echo "cosmovisor not installed. Installing..."
  go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@v1.3.0
fi

rm -rf $HOME_DIR/.testnets
mkdir -p $HOME_DIR
kill -9 $(pgrep -f "rpc.laddr=tcp://0.0.0.0") || true
tmux kill-session -t node0
tmux kill-session -t node1

BINARY_NAME=2astrad

# prepare for initializing
$BINARY_NAME config keyring-backend $KEYRING
$BINARY_NAME config chain-id $CHAIN_ID

# start a testnet
$BINARY_NAME testnet init-files --keyring-backend=$KEYRING --chain-id=$CHAIN_ID --output-dir=$HOME_DIR/.testnets --v=2

if [[ ! -d "$HOME_DIR/.testnets" ]];
then
  exit 1
fi

# if $KEY exists it should be deleted
$BINARY_NAME keys add $KEY_GENESIS_VALIDATOR --keyring-backend=$KEYRING --algo=$KEYALGO --home=$HOME_DIR/.testnets/node0/astrad
$BINARY_NAME keys add $KEY_PARTNER --keyring-backend=$KEYRING --algo=$KEYALGO --home=$HOME_DIR/.testnets/node0/astrad
$BINARY_NAME keys add $KEY_STRATEGIC --keyring-backend=$KEYRING --algo=$KEYALGO --home=$HOME_DIR/.testnets/node0/astrad
$BINARY_NAME keys add $KEY_CORE_TEAM --keyring-backend=$KEYRING --algo=$KEYALGO --home=$HOME_DIR/.testnets/node0/astrad

# change consensus params
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="30000000"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.consensus_params["block"]["time_iota_ms"]="3000"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# change staking denom to aastra
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="aastra"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["staking"]["params"]["unbonding_time"]="60s"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# update crisis variable to aastra
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="aastra"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# update fee-market
#cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["feemarket"]["params"]["min_gas_price"]="100000000000"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# udpate gov genesis
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="aastra"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["max_deposit_period"]="10s"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["gov"]["voting_params"]["voting_period"]="60s"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# update distribution
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["distribution"]["params"]["community_tax"]="0.000000000000000000"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# update mint genesis
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="aastra"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# create genesis accounts and transactions
$BINARY_NAME add-genesis-account $KEY_PARTNER 528000000000000000000000000aastra --keyring-backend=$KEYRING --home=$HOME_DIR/.testnets/node0/astrad
$BINARY_NAME add-genesis-account $KEY_STRATEGIC 432000000000000000000000000aastra --keyring-backend=$KEYRING --home=$HOME_DIR/.testnets/node0/astrad
$BINARY_NAME add-genesis-account astra13wjs7d3z8hra6rp7vjmryuulwxjrd232sceuen 120000000000000000000000000aastra --keyring-backend=$KEYRING --home=$HOME_DIR/.testnets/node0/astrad
$BINARY_NAME add-genesis-account $KEY_CORE_TEAM 119990000000000000000000000aastra --keyring-backend=$KEYRING --home=$HOME_DIR/.testnets/node0/astrad
#$BINARY_NAME add-genesis-account astra1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8y4fl3r 60000000000000000000000000aastra --home=$HOME_DIR/.testnets/node0/astrad

# update total supply
cat $HOME_DIR/.testnets/node0/astrad/config/genesis.json | jq '.app_state["bank"]["supply"][0]["amount"]="1200000000000000000000000000"' > $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json && mv $HOME_DIR/.testnets/node0/astrad/config/tmp_genesis.json $HOME_DIR/.testnets/node0/astrad/config/genesis.json

# Validate genesis
$BINARY_NAME validate-genesis --home=$HOME_DIR/.testnets/node0/astrad

#### change config.toml values
# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME_DIR/.testnets/node0/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:40001|tcp://0.0.0.0:40000|g' $HOME_DIR/.testnets/node0/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:50002|tcp://0.0.0.0:50000|g' $HOME_DIR/.testnets/node0/astrad/config/config.toml
sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME_DIR/.testnets/node0/astrad/config/app.toml
# validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26558|g' $HOME_DIR/.testnets/node1/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:40001|tcp://0.0.0.0:40001|g' $HOME_DIR/.testnets/node1/astrad/config/config.toml
sed -i -E 's|tcp://0.0.0.0:50002|tcp://0.0.0.0:50002|g' $HOME_DIR/.testnets/node1/astrad/config/config.toml
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME_DIR/.testnets/node1/astrad/config/config.toml
sed -i -E 's|address = "tcp://0.0.0.0:1317"|address = "tcp://0.0.0.0:2317"|g' $HOME_DIR/.testnets/node1/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9090"|address = "0.0.0.0:9190"|g' $HOME_DIR/.testnets/node1/astrad/config/app.toml
sed -i -E 's|address = "0.0.0.0:9091"|address = "0.0.0.0:9191"|g' $HOME_DIR/.testnets/node1/astrad/config/app.toml
sed -i -E 's|address = "127.0.0.1:8545"|address = "127.0.0.1:18545"|g' $HOME_DIR/.testnets/node1/astrad/config/app.toml
sed -i -E 's|enabled-unsafe-cors = false|enabled-unsafe-cors = true|g' $HOME_DIR/.testnets/node1/astrad/config/app.toml

sed -i -E 's|timeout_commit = "5s"|timeout_commit = "1s"|g' $HOME_DIR/.testnets/node0/astrad/config/config.toml
sed -i -E 's|timeout_commit = "5s"|timeout_commit = "1s"|g' $HOME_DIR/.testnets/node1/astrad/config/config.toml

# copy validator1 genesis file to validator2
cp $HOME_DIR/.testnets/node0/astrad/config/genesis.json $HOME_DIR/.testnets/node1/astrad/config/genesis.json

# create cosmosvisor folders
mkdir -p $HOME_DIR/.testnets/node0/astrad/cosmovisor
mkdir -p $HOME_DIR/.testnets/node0/astrad/data
mkdir -p $HOME_DIR/.testnets/node0/astrad/cosmovisor/genesis
mkdir -p $HOME_DIR/.testnets/node0/astrad/cosmovisor/genesis/bin
mkdir -p $HOME_DIR/.testnets/node0/astrad/cosmovisor/upgrades
cp ~/go/bin/$BINARY_NAME $HOME_DIR/.testnets/node0/astrad/cosmovisor/genesis/bin/astrad

mkdir -p $HOME_DIR/.testnets/node1/astrad/cosmovisor
mkdir -p $HOME_DIR/.testnets/node1/astrad/data
mkdir -p $HOME_DIR/.testnets/node1/astrad/cosmovisor/genesis
mkdir -p $HOME_DIR/.testnets/node1/astrad/cosmovisor/genesis/bin
mkdir -p $HOME_DIR/.testnets/node1/astrad/cosmovisor/upgrades
cp ~/go/bin/$BINARY_NAME $HOME_DIR/.testnets/node1/astrad/cosmovisor/genesis/bin/astrad

echo "node0:" $($BINARY_NAME keys unsafe-export-eth-key node0 --keyring-backend=test --home=$HOME_DIR/.testnets/node0/astrad) $($BINARY_NAME keys show node0 -a --keyring-backend=test --home=$HOME_DIR/.testnets/node0/astrad)
echo "node1:" $($BINARY_NAME keys unsafe-export-eth-key node1 --keyring-backend=test --home=$HOME_DIR/.testnets/node1/astrad) $($BINARY_NAME keys show node1 -a --keyring-backend=test --home=$HOME_DIR/.testnets/node1/astrad)
echo "strategic:" $($BINARY_NAME keys unsafe-export-eth-key $KEY_STRATEGIC --keyring-backend=$KEYRING --home=$HOME_DIR/.testnets/node0/astrad) $($BINARY_NAME keys show $KEY_STRATEGIC -a --keyring-backend=$KEYRING --home=$HOME_DIR/.testnets/node0/astrad)

# run with cosmovior
tmux new-session -s node0 -d
tmux send-keys -t node0 "DAEMON_ALLOW_DOWNLOAD_BINARIES=true DAEMON_NAME=astrad DAEMON_HOME=$HOME_DIR/.testnets/node0/astrad cosmovisor run start --home=$HOME_DIR/.testnets/node0/astrad --rpc.laddr=tcp://0.0.0.0:26657" Enter
tmux new-session -s node1 -d
tmux send-keys -t node1 "DAEMON_ALLOW_DOWNLOAD_BINARIES=true DAEMON_NAME=astrad DAEMON_HOME=$HOME_DIR/.testnets/node1/astrad cosmovisor run start --home=$HOME_DIR/.testnets/node1/astrad --rpc.laddr=tcp://0.0.0.0:26557" Enter
tmux ls

sleep 10
$BINARY_NAME tx bank send $KEY_PARTNER astra1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8y4fl3r 60000000000000000000000000aastra --home=$HOME_DIR/.testnets/node0/astrad --keyring-backend=$KEYRING --chain-id=$CHAIN_ID -y