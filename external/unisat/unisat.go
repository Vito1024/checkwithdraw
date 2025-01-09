package unisat

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"withdraw"
	"withdraw/config"

	"github.com/go-resty/resty/v2"
)

const DEFAULT_TIMEOUT = time.Second * 5
const START_CURSOR_FILE_PATH = "start_cursor.txt"
const REQUEST_WITHDRAW_LIMIT = 10

var START_CURSOR int

type Service struct {
	config     config.UnisatConfig
	httpClient *resty.Client
}

func New(config config.UnisatConfig) *Service {
	svc := &Service{
		config: config,
		httpClient: resty.New().
			SetTimeout(DEFAULT_TIMEOUT).
			SetRetryCount(100).
			SetRetryWaitTime(500 * time.Millisecond).
			AddRetryCondition(func(response *resty.Response, err error) bool {
				// network error
				if err != nil {
					return true
				}
				// http error
				return response.StatusCode() > 299
			}),
	}

	_, err := os.Stat(START_CURSOR_FILE_PATH)
	if os.IsNotExist(err) {
		START_CURSOR = 0
	} else {
		filecontent, err := os.ReadFile(START_CURSOR_FILE_PATH)
		if err != nil {
			panic(err)
		}
		START_CURSOR, err = strconv.Atoi(strings.Trim(string(filecontent), "\n"))
		if err != nil {
			panic(err)
		}
	}

	if withdraw.START_CURSOR > START_CURSOR {
		START_CURSOR = withdraw.START_CURSOR
	}

	return svc
}

func (svc *Service) FollowWithdrawTransactions(ctx context.Context) <-chan withdraw.UnisatWithdrawTransaction {
	out := make(chan withdraw.UnisatWithdrawTransaction)

	go func() {
		defer close(out)
		defer func() {
			bs, err := json.Marshal(START_CURSOR)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile(START_CURSOR_FILE_PATH, bs, 0644)
			if err != nil {
				panic(err)
			}
			fmt.Printf("stop listening new withdraw transactions from unisat, write current cursor(%d) to file", START_CURSOR)
		}()

		fmt.Printf("following unisat withdraw tx, start cursor:%d\n", START_CURSOR)

		path := "/brc20-module/withdraw-history"
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var unisatResponse struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
					Data struct {
						Height int                                  `json:"height"`
						Total  int                                  `json:"total"`
						Cursor int                                  `json:"cursor"`
						Detail []withdraw.UnisatWithdrawTransaction `json:"detail"`
					} `json:"data"`
				}
				resp, err := svc.httpClient.R().
					SetQueryParams(map[string]string{
						"cursor": strconv.Itoa(START_CURSOR),
						"size":   strconv.Itoa(REQUEST_WITHDRAW_LIMIT),
					}).
					SetResult(&unisatResponse).
					Get(svc.config.Host + path)
				if err != nil {
					panic(err)
				}
				if resp.StatusCode() != 200 {
					panic(fmt.Sprintf("resp code not 200, %+v", resp))
				}

				// tx 按时间升序
				for _, tx := range unisatResponse.Data.Detail {
					// 延后一下
					{
						blockTime := time.Unix(int64(tx.BlockTime), 0)
						for {
							if time.Since(blockTime) > time.Hour {
								break
							}
							select {
							case <-ctx.Done():
								return
							case <-time.After(time.Minute):
							}
						}
					}
					out <- tx
				}
				START_CURSOR += len(unisatResponse.Data.Detail)

				// empty
				if len(unisatResponse.Data.Detail) < REQUEST_WITHDRAW_LIMIT {
					sleepWithContext(ctx, time.Minute)
				}
			}
		}
	}()

	return out
}

func (svc *Service) GetWithdrawTransactionsFromFile(ctx context.Context) []string {
	withdraws := parseWithdrawResponse()
	fmt.Printf("withdraw transactions from unisat: %d\n", len(withdraws))

	txIds := make([]string, len(withdraws))
	for idx, withdraw := range withdraws {
		txIds[len(withdraws)-1-idx] = withdraw.TxID
	}

	return txIds
}

func sleepWithContext(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
