build: build-app
run:
	make build-app && ./astra-test
build-app:
	go build -o astra-test
clean:
	rm -f astra-test && go clean --modcache && go mod tidy
build-genesis:
	bash setup-node.sh
run-all-nodes:
	make build-genesis
	tmux new -s node0 -d astrad start --home=.testnets/node0/astrad --rpc.laddr=tcp://127.0.0.1:26657
	tmux new -s node1 -d astrad start --home=.testnets/node1/astrad --rpc.laddr=tcp://127.0.0.1:26557
	tmux new -s node2 -d astrad start --home=.testnets/node2/astrad --rpc.laddr=tcp://127.0.0.1:26457
	tmux new -s node3 -d astrad start --home=.testnets/node3/astrad --rpc.laddr=tcp://127.0.0.1:26357
run-node0:
	astrad start --home=.testnets/node0/astrad --rpc.laddr=tcp://127.0.0.1:26657 &>> .testnets/node0/log.log
run-node1:
	astrad start --home=.testnets/node1/astrad --rpc.laddr=tcp://127.0.0.1:26557 &>> .testnets/node1/log.log
run-node2:
	astrad start --home=.testnets/node2/astrad --rpc.laddr=tcp://127.0.0.1:26457 &>> .testnets/node2/log.log
run-node3:
	astrad start --home=.testnets/node3/astrad --rpc.laddr=tcp://127.0.0.1:26357 &>> .testnets/node3/log.log