package runner

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/bigquery/mock_bigquery"
	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/harmony/mock_harmony"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestBackfillBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	testCases := []struct {
		name         string
		currentCount int64
		endBlock     int64
		setupMocks   func(*mock_harmony.MockHarmonyClient, *mock_bigquery.MockBigQueryClient)
	}{
		{
			name:         "backfill all",
			currentCount: 1,
			endBlock:     10,
			setupMocks: func(harmonyClient *mock_harmony.MockHarmonyClient, bqClient *mock_bigquery.MockBigQueryClient) {
				for i := 1; i <= 10; i++ {
					bqClient.EXPECT().GetBlock(gomock.Any(), int64(i)).Return(nil, nil)
					bqClient.EXPECT().GetTransactions(gomock.Any(), int64(i)).Return(nil, nil)

					block := &model.Block{
						Number: strconv.Itoa(i),
						Transactions: []*model.Transaction{
							{
								BlockNumber: strconv.Itoa(i),
							},
						},
					}

					harmonyClient.EXPECT().GetBlockByNumber(gomock.Any()).Return(block, nil)

					bqClient.EXPECT().InsertBlock(gomock.Any(), block).Return(nil)
					bqClient.EXPECT().InsertTransactions(gomock.Any(), block.Transactions).Return(nil)
				}
			},
		},
		{
			name:         "skips filled blocks and txns",
			currentCount: 1,
			endBlock:     10,
			setupMocks: func(harmonyClient *mock_harmony.MockHarmonyClient, bqClient *mock_bigquery.MockBigQueryClient) {
				for i := 1; i <= 10; i++ {
					var (
						returnBlock *model.Block
						returnTxns  []*model.Transaction
					)

					if i%2 == 0 {
						returnBlock = &model.Block{
							Number: strconv.Itoa(i),
						}
						returnTxns = []*model.Transaction{
							{
								BlockNumber: strconv.Itoa(i),
							},
						}
					}

					bqClient.EXPECT().GetBlock(gomock.Any(), int64(i)).Return(returnBlock, nil)
					bqClient.EXPECT().GetTransactions(gomock.Any(), int64(i)).Return(returnTxns, nil)

					block := &model.Block{
						Number: strconv.Itoa(i),
						Transactions: []*model.Transaction{
							{
								BlockNumber: strconv.Itoa(i),
							},
						},
					}

					harmonyClient.EXPECT().GetBlockByNumber(gomock.Any()).Return(block, nil)

					if i%2 != 0 {
						bqClient.EXPECT().InsertBlock(gomock.Any(), block).Return(nil)
						bqClient.EXPECT().InsertTransactions(gomock.Any(), block.Transactions).Return(nil)
					}
				}
			},
		},
		{
			name:         "skips filled blocks but fills in missing txns",
			currentCount: 1,
			endBlock:     10,
			setupMocks: func(harmonyClient *mock_harmony.MockHarmonyClient, bqClient *mock_bigquery.MockBigQueryClient) {
				for i := 1; i <= 10; i++ {
					var returnBlock *model.Block

					if i == 5 {
						returnBlock = &model.Block{
							Number: strconv.Itoa(i),
						}
					}

					bqClient.EXPECT().GetBlock(gomock.Any(), int64(i)).Return(returnBlock, nil)
					bqClient.EXPECT().GetTransactions(gomock.Any(), int64(i)).Return(nil, nil)

					block := &model.Block{
						Number: strconv.Itoa(i),
						Transactions: []*model.Transaction{
							{
								BlockNumber: strconv.Itoa(i),
							},
						},
					}

					harmonyClient.EXPECT().GetBlockByNumber(gomock.Any()).Return(block, nil)

					if i != 5 {
						bqClient.EXPECT().InsertBlock(gomock.Any(), block).Return(nil)
					}
					bqClient.EXPECT().InsertTransactions(gomock.Any(), block.Transactions).Return(nil)
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			harmonyClient := mock_harmony.NewMockHarmonyClient(ctrl)
			bqClient := mock_bigquery.NewMockBigQueryClient(ctrl)

			runner := &BackfillRunner{
				harmonyClient:  harmonyClient,
				bigQueryClient: bqClient,
				retryBlockChan: make(chan *model.RetryBlock),
				retryTxnChan:   make(chan *model.RetryTransaction),
				logger:         zap.NewNop(),
			}

			tc.setupMocks(harmonyClient, bqClient)

			counter := &counter{
				tc.currentCount,
				&sync.Mutex{},
			}

			wg := &sync.WaitGroup{}
			wg.Add(1)

			runner.backfillBlocks(context.TODO(), counter, wg, 10)
		})
	}
}
