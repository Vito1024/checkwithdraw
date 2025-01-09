package oklink

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"withdraw"
	"withdraw/config"

	"github.com/go-resty/resty/v2"
)

const defaultTimeout = 5 * time.Second
const intervalWithoutReachLimit = 334 * time.Millisecond

const (
	HEADER_OK_ACCESS_KEY = "Ok-Access-Key"
)

type Service struct {
	config     config.OkLinkConfig
	httpClient *resty.Client
}

func New(config config.OkLinkConfig) *Service {
	return &Service{
		config: config,
		httpClient: resty.New().
			SetTimeout(defaultTimeout).
			SetRetryCount(100). // oklink免费版，限速3/s
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
}

// GetFractalBitcoinTransactionDetail 获取交易详情
// Doc: https://www.oklink.com/docs/zh/#quickstart-guide-api-authentication
// Doc: https://www.oklink.com/docs/zh/#btc-inscription-data-get-inscription-token-transaction-details-for-specific-hash
func (svc *Service) GetFractalBitcoinBRC20TransactionDetail(ctx context.Context, txId string, options ...func(*withdraw.Option)) (withdraw.OKLinkBRC20TransactionDetail, error) {
	path := "/api/v5/explorer/inscription/transaction-detail"

	var option withdraw.Option
	for _, opt := range options {
		opt(&option)
	}
	if option.RateLimit {
		time.Sleep(intervalWithoutReachLimit)
	}

	var okxResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Limit           string                                  `json:"limit"`
			Page            string                                  `json:"page"`
			TotalPage       string                                  `json:"totalPage"`
			TransactionList []withdraw.OKLinkBRC20TransactionDetail `json:"transactionList"`
		} `json:"data"`
	}
	resp, err := svc.httpClient.R().
		SetHeader(HEADER_OK_ACCESS_KEY, svc.config.Key).
		SetQueryParams(map[string]string{
			"chainShortName": "FRACTAL",
			"protocolType":   "brc20",
			"txId":           txId,
		}).
		SetResult(&okxResponse).
		Get(svc.config.Host + path)
	if err != nil {
		panic(fmt.Sprintf("error: %v, txId: %s", err, txId))
	}
	if resp.StatusCode() != 200 {
		panic(fmt.Sprintf("resp: %s, txId: %s", resp.String(), txId))
	}
	if len(okxResponse.Data) != 1 {
		panic(fmt.Sprintf("empty data, txId: %s", txId))
	}
	if len(okxResponse.Data[0].TransactionList) == 0 {
		return withdraw.OKLinkBRC20TransactionDetail{TxId: txId}, withdraw.ErrTransactionNotBRC20Withdraw
	}

	return okxResponse.Data[0].TransactionList[0], nil
}

func (svc *Service) GetFractalBitcoinBRC20TransactionDetailBatch(ctx context.Context, txIds []string, options ...func(*withdraw.Option)) <-chan withdraw.OKLinkBRC20TransactionDetailBatch {
	out := make(chan withdraw.OKLinkBRC20TransactionDetailBatch)

	go func() {
		defer close(out)
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %v\n", r)
			}
		}()

		var option withdraw.Option
		for _, opt := range options {
			opt(&option)
		}

		mismatch := 0
		for idx, txId := range txIds {
			detail, err := svc.GetFractalBitcoinBRC20TransactionDetail(ctx, txId)
			out <- withdraw.OKLinkBRC20TransactionDetailBatch{OKLinkBRC20TransactionDetail: detail, Err: err}

			if option.Progress {
				if err == withdraw.ErrTransactionNotBRC20Withdraw {
					mismatch++
					fmt.Printf("progress: %d/%d, mismatch:%d, txId:%s\n", idx+1, len(txIds), mismatch, txId)
				} else {
					fmt.Printf("progress: %d/%d, mismatch:%d\r", idx+1, len(txIds), mismatch)
				}
			}
			if option.RateLimit {
				time.Sleep(intervalWithoutReachLimit)
			}
		}
	}()

	return out
}

// GetFractalBitcoinAddressTokenTransactionList 获取地址token交易列表
func (svc *Service) GetFractalBitcoinAddressTokenTransactionList(ctx context.Context, address string, protocolType string, inscriptionID string) (any, error) {
	path := "/api/v5/explorer/inscription/address-token-transaction-list"

	var okxResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Limit           string `json:"limit"`
			Page            string `json:"page"`
			TotalPage       string `json:"totalPage"`
			ChainFullName   string `json:"chainFullName"`
			ChainShortName  string `json:"chainShortName"`
			TotalTransfer   string `json:"totalTransfer"`
			TransactionList []struct {
				TxId               string `json:"txId"`
				BlockHash          string `json:"blockHash"`
				Height             string `json:"height"`
				TransactionTime    string `json:"transactionTime"`
				From               string `json:"from"`
				To                 string `json:"to"`
				Amount             string `json:"amount"`
				Symbol             string `json:"symbol"`
				Action             string `json:"action"`
				TokenInscriptionId string `json:"tokenInscriptionId"`
				ProtocolType       string `json:"protocolType"`
				State              string `json:"state"`
				InscriptionId      string `json:"inscriptionId"`
				InscriptionNumber  string `json:"inscriptionNumber"`
				OutputIndex        string `json:"outputIndex"`
			} `json:"transactionList"`
		} `json:"data"`
	}
	resp, err := svc.httpClient.R().
		SetHeader(HEADER_OK_ACCESS_KEY, svc.config.Key).
		SetQueryParams(map[string]string{
			"chainShortName":     "FRACTAL",
			"address":            address,
			"protocolType":       protocolType,
			"tokenInscriptionId": inscriptionID,
		}).
		SetResult(&okxResponse).
		Get(svc.config.Host + path)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode() != 200 {
		panic(resp.String())
	}
	if len(okxResponse.Data) != 1 {
		panic("empty data")
	}

	return okxResponse.Data[0], nil
}

func (svc *Service) GetFractalBitcoinInscriptionHolderList(ctx context.Context, address string, protocolType string, symbol string) (any, error) {
	path := "/api/v5/explorer/inscription/token-position-list"

	var okxResponse struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Page         string `json:"page"`
			Limit        string `json:"limit"`
			TotalPage    string `json:"totalPage"`
			PositionList []struct {
				HolderAddress string `json:"holderAddress"`
				Amount        string `json:"amount"`
				Rank          string `json:"rank"`
			} `json:"positionList"`
		} `json:"data"`
	}
	resp, err := svc.httpClient.R().
		SetHeader(HEADER_OK_ACCESS_KEY, svc.config.Key).
		SetQueryParams(map[string]string{
			"chainShortName": "FRACTAL",
			"address":        address,
			"protocolType":   protocolType,
			"symbol":         symbol,
		}).
		SetResult(&okxResponse).
		Get(svc.config.Host + path)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode() != 200 {
		panic(resp.String())
	}

	return okxResponse.Data, nil
}

type AddressBrc20InscriptionListOption struct {
	Symbol string
}

func RequestAddressBrc20InscriptionListWithSymbol(symbol string) func(*AddressBrc20InscriptionListOption) {
	return func(option *AddressBrc20InscriptionListOption) {
		option.Symbol = symbol
	}
}

type Inscription struct {
	InscriptionId      string `json:"inscriptionId"`
	TokenInscriptionId string `json:"tokenInscriptionId"`
	InscriptionNumber  string `json:"inscriptionNumber"`
	Symbol             string `json:"symbol"`
	State              string `json:"state"`
	ProtocolType       string `json:"protocolType"`
	Action             string `json:"action"`
}

func (svc *Service) GetFractalBitcoinAddressBrc20InscriptionList(ctx context.Context, address string, options ...func(*AddressBrc20InscriptionListOption)) (any, error) {
	path := "/api/v5/explorer/inscription/address-inscription-list"

	var inscriptions []Inscription

	var option AddressBrc20InscriptionListOption
	for _, opt := range options {
		opt(&option)
	}
	totalPage := 0
	for page := 1; ; page++ {
		var okxResponse struct {
			Code string `json:"code"`
			Msg  string `json:"msg"`
			Data []struct {
				Page            string        `json:"page"`
				Limit           string        `json:"limit"`
				TotalPage       string        `json:"totalPage"`
				InscriptionList []Inscription `json:"inscriptionList"`
			} `json:"data"`
		}
		resp, err := svc.httpClient.R().
			SetHeader(HEADER_OK_ACCESS_KEY, svc.config.Key).
			SetQueryParams(map[string]string{
				"chainShortName": "FRACTAL",
				"address":        address,
				"protocolType":   "brc20",
				"limit":          "100",
			}).
			SetResult(&okxResponse).
			Get(svc.config.Host + path)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode() != 200 {
			panic(resp.String())
		}
		if len(okxResponse.Data) == 0 {
			break
		}
		if len(okxResponse.Data[0].InscriptionList) == 0 {
			break
		}

		for _, inscription := range okxResponse.Data[0].InscriptionList {
			if option.Symbol != "" && inscription.Symbol != option.Symbol {
				continue
			}
			inscriptions = append(inscriptions, inscription)
		}

		if totalPage == 0 {
			totalPage, err = strconv.Atoi(okxResponse.Data[0].TotalPage)
			if err != nil {
				panic(err)
			}
		}
		if page >= totalPage {
			break
		}
	}

	return inscriptions, nil
}
