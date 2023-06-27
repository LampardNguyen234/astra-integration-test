package common

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	"github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ITestSuite interface {
	Name() string
	RunTest()
	SetMasterKey(string)
	GetMasterKey() string
	FundAccount(recipient string, amount float64)
	ClawBack(privateKey string)
}

type TestClient struct {
	*client.CosmosClient
	Log       logger.Logger
	masterKey string
}

func NewTestClient(cosmosClient *client.CosmosClient, log logger.Logger) *TestClient {
	return &TestClient{
		CosmosClient: cosmosClient,
		Log:          log,
	}
}

func (c *TestClient) SetMasterKey(k string) {
	c.masterKey = k
}

func (c *TestClient) GetMasterKey() string {
	return c.masterKey
}

func (c *TestClient) FundAccount(recipient string, amount float64) {
	if amount == 0 {
		return
	}

	amt := common.Float64ToBigInt(amount)
	resp, err := c.CosmosClient.TxSend(msg_params.TxSendRequestParams{
		TxParams: msg_params.TxParams{
			PrivateKey: c.GetMasterKey(),
		},
		ToAddr: recipient,
		Amount: amt,
	})
	if err != nil {
		c.Log.Panicf("failed to fund account %v: %v", recipient, err)
	}

	c.TxShouldPass(resp.TxHash)

	c.BalanceCompare(recipient, sdk.NewIntFromBigInt(amt), assert.OpGTE)
}

func (c *TestClient) ClawBack(privateKey string) {
	ki, _ := account.NewKeyInfoFromPrivateKey(privateKey)
	balance, err := c.Balance(ki.CosmosAddress)
	if err != nil {
		c.Log.Errorf("failed to get balance of refunded account %v: %v", ki.CosmosAddress, err)
		return
	}
	tmp := common.BigIntToFloat64(balance.Total.BigInt())
	if tmp > 0.02 {
		_, err = c.TxSend(
			msg_params.TxSendRequestParams{
				TxParams: msg_params.TxParams{
					PrivateKey:    privateKey,
					GasLimit:      200000,
					GasAdjustment: 1,
					GasPrice:      msg_params.DefaultTxParams().GasPrice,
				},
				ToAddr: account.MustNewPrivateKeyFromString(c.masterKey).AccAddress().String(),
				Amount: common.Float64ToBigInt(tmp - 0.02),
			},
		)
		if err != nil {
			c.Log.Errorf("failed to perform refunding for %v: %v", ki.CosmosAddress, err)
			return
		}
	}
	c.Log.Debugf("succeeded to perform refunding for %v", ki.CosmosAddress)
}

func (c *TestClient) Start() {
	c.Log.Infof("==================== STARTED ====================")
}

func (c *TestClient) Finished() {
	c.Log.Infof("==================== FINISHED ====================")
	fmt.Println("")
}
