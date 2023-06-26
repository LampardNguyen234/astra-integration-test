package send

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	repoCommon "github.com/LampardNguyen234/astra-integration-test/common"
)

func (s *SendSuite) RunTest() {
	s.Log.Debugf("START")
	defaultTxParams := msg_params.TxParams{
		GasLimit:      1000000,
		GasAdjustment: 1,
		GasPrice:      "100000000000aastra",
	}

	defaultTestAmt := 0.15
	defaultPrefunded := 0.3

	tcs := []sendTestCase{
		//{
		//	name: "should be able to send where amt + fee >= balance",
		//	txParams: msg_params.TxParams{
		//		PrivateKey:    repoCommon.RandKeyInfo().PrivateKey,
		//		GasLimit:      defaultTxParams.GasLimit,
		//		GasAdjustment: defaultTxParams.GasAdjustment,
		//		GasPrice:      defaultTxParams.GasPrice,
		//	},
		//	recipient: repoCommon.RandKeyInfo(),
		//	amt:       defaultTestAmt,
		//	prefunded: defaultPrefunded,
		//	expErr:    nil,
		//},
		//{
		//	name: "should not be able to send where amt + fee < balance",
		//	txParams: msg_params.TxParams{
		//		PrivateKey:    repoCommon.RandKeyInfo().PrivateKey,
		//		GasLimit:      defaultTxParams.GasLimit,
		//		GasAdjustment: defaultTxParams.GasAdjustment,
		//		GasPrice:      defaultTxParams.GasPrice,
		//	},
		//	recipient: repoCommon.RandKeyInfo(),
		//	amt:       defaultPrefunded,
		//	prefunded: defaultPrefunded,
		//	expErr:    fmt.Errorf("insufficient funds"),
		//},
		{
			name: "should not be able to send with insufficient gas limit",
			txParams: msg_params.TxParams{
				PrivateKey:    repoCommon.RandKeyInfo().PrivateKey,
				GasLimit:      10000,
				GasAdjustment: defaultTxParams.GasAdjustment,
				GasPrice:      "100000000000aastra",
			},
			recipient: repoCommon.RandKeyInfo(),
			amt:       defaultTestAmt,
			prefunded: defaultPrefunded,
			expErr:    fmt.Errorf("out of gas"),
		},
		{
			name: "should not be able to send with gasPrice = 0",
			txParams: msg_params.TxParams{
				PrivateKey:    repoCommon.RandKeyInfo().PrivateKey,
				GasLimit:      1000000,
				GasAdjustment: defaultTxParams.GasAdjustment,
				GasPrice:      "0aastra",
			},
			recipient: repoCommon.RandKeyInfo(),
			amt:       defaultTestAmt,
			prefunded: defaultPrefunded,
			expErr:    fmt.Errorf("insufficient fee"),
		},
	}

	for _, tc := range tcs {
		s.runTestCase(tc)
	}

	s.Log.Debugf("FINISHED")
}
