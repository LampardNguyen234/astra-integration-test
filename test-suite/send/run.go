package send

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	sdkCommon "github.com/LampardNguyen234/astra-go-sdk/common"
	repoCommon "github.com/LampardNguyen234/astra-integration-test/common"
	. "github.com/LampardNguyen234/astra-integration-test/framework"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func (s *SendSuite) RunTest() {
	s.Start()

	tc := s.registerTests()
	tc.Run()
	tc.Report()

	s.Finished()
}

func (s *SendSuite) registerTests() ITestNode {
	RegisterFailHandler(ginkgo.Fail)

	var from, to, operator *account.KeyInfo
	var txParams msg_params.TxParams

	defaultBeforeEach := BeforeEach(func() {
		from = repoCommon.RandKeyInfo()
		to = repoCommon.RandKeyInfo()
		operator = repoCommon.RandKeyInfo()
		txParams = msg_params.TxParams{
			PrivateKey: from.PrivateKey,
			GasLimit:   200000,
		}
		_ = operator

		s.FundAccount(from.CosmosAddress, 1)
	})
	defaultAfterEach := AfterEach(func() {
		s.Refund(from.PrivateKey)
		s.Refund(operator.PrivateKey)
	})

	root := Describe(fmt.Sprintf("%v", s.Name()),
		Context(
			"sufficient balance",
			defaultBeforeEach,
			defaultAfterEach,

			It("should be able to send by self", func() {
				oldRecipientBalance, err := s.CosmosClient.Balance(to.CosmosAddress)
				Expect(err).To(BeNil())

				sentAmt := sdk.NewIntWithDecimal(1, 17)
				resp, err := s.CosmosClient.BuildAndSendTx(
					txParams,
					bankTypes.NewMsgSend(
						account.MustParseCosmosAddress(from.CosmosAddress),
						account.MustParseCosmosAddress(to.CosmosAddress),
						sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sentAmt)),
					),
				)
				Expect(err).To(BeNil())
				s.TxShouldPass(resp.TxHash)

				newRecipientBalance, err := s.CosmosClient.Balance(to.CosmosAddress)
				Expect(err).To(BeNil())
				Expect(newRecipientBalance.Total.Sub(oldRecipientBalance.Total)).To(Equal(
					sentAmt,
				))
			}),
		),

		Context("insufficient balance",
			defaultBeforeEach,
			defaultAfterEach,

			It("should not be able to send", func() {
				sentAmt := sdk.NewIntWithDecimal(2, 20)
				resp, err := s.CosmosClient.BuildAndSendTx(
					txParams,
					bankTypes.NewMsgSend(
						account.MustParseCosmosAddress(from.CosmosAddress),
						account.MustParseCosmosAddress(to.CosmosAddress),
						sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sentAmt)),
					),
				)
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("insufficient funds"))
				} else {
					s.TxShouldFailWithError(resp.TxHash, "insufficient funds")
				}
			}),
		),

		Context("invalid fee",
			defaultBeforeEach,
			defaultAfterEach,

			It("should not be able to send with insufficient gas limit", func() {
				sentAmt := sdk.NewIntWithDecimal(1, 17)
				resp, err := s.CosmosClient.BuildAndSendTx(
					msg_params.TxParams{
						PrivateKey: from.PrivateKey,
						GasLimit:   10000,
					},
					bankTypes.NewMsgSend(
						account.MustParseCosmosAddress(from.CosmosAddress),
						account.MustParseCosmosAddress(to.CosmosAddress),
						sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sentAmt)),
					),
				)
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("out of gas"))
				} else {
					s.TxShouldFailWithError(resp.TxHash, "out of gas")
				}
			}),

			It("should not be able to send with gasPrice=0", func() {
				fp, err := s.FeemarketParams()
				Expect(err).To(BeNil())
				if fp.MinGasPrice.IsZero() {
					return
				}

				sentAmt := sdk.NewIntWithDecimal(1, 17)
				resp, err := s.CosmosClient.BuildAndSendTx(
					msg_params.TxParams{
						PrivateKey: from.PrivateKey,
						GasLimit:   1000000,
						GasPrice:   "0aastra",
					},
					bankTypes.NewMsgSend(
						account.MustParseCosmosAddress(from.CosmosAddress),
						account.MustParseCosmosAddress(to.CosmosAddress),
						sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sentAmt)),
					),
				)
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("insufficient fee"))
				} else {
					s.TxShouldFailWithError(resp.TxHash, "insufficient fee")
				}
			}),
		),
	)

	return root
}
