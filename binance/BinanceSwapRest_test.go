package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	. "github.com/strengthening/goghostex"
)

const (
	SWAP_API_KEY       = ""
	SWAP_API_SECRETKEY = ""
	SWAP_PROXY_URL     = "socks5://127.0.0.1:1090"
)

func TestSwap_MarketAPI(t *testing.T) {

	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(SWAP_PROXY_URL)
				},
			},
		},
		ApiKey:        SWAP_API_KEY,
		ApiSecretKey:  SWAP_API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	bn := New(config)
	// ticker unit test
	if ticker, resp, err := bn.Swap.GetTicker(
		Pair{BTC, USDT},
	); err != nil {
		t.Error(err)
		return
	} else {
		standard, err := json.Marshal(ticker)
		if err != nil {
			t.Error(err)
			return
		}

		t.Log("Ticker standard struct: ")
		t.Log(string(standard))
		t.Log("Ticker remote api response: ")
		t.Log(string(resp))
	}

	// depth unit test
	if depth, resp, err := bn.Swap.GetDepth(
		Pair{BTC, USDT},
		20,
	); err != nil {
		t.Error(err)
		return
	} else {
		standard, err := json.Marshal(depth)
		if err != nil {
			t.Error(err)
			return
		}

		t.Log("Depth standard struct:")
		t.Log(string(standard))
		t.Log("Depth remote api response: ")
		t.Log(string(resp))

		// make sure the later request get bigger sequence
		depth1, _, _ := bn.Swap.GetDepth(
			Pair{BTC, USDT},
			20,
		)

		if depth1.Sequence < depth.Sequence {
			t.Error("later request get smaller sequence!!")
			return
		}

		if err := depth.Verify(); err != nil {
			t.Error(err)
			return
		}

		if err := depth1.Verify(); err != nil {
			t.Error(err)
			return
		}
	}

	if highest, lowest, err := bn.Swap.GetLimit(Pair{BTC, USDT}); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("highest:", highest)
		t.Log("lowest:", lowest)
	}

	if klines, resp, err := bn.Swap.GetKline(
		Pair{BTC, USDT},
		KLINE_PERIOD_1DAY,
		20,
		1271149752000,
	); err != nil {
		t.Error(err)
		return
	} else {
		klineRaw, _ := json.Marshal(klines)
		t.Log(string(klineRaw))
		t.Log(string(resp))
	}

	if openAmount, timestamp, _, err := bn.Swap.GetOpenAmount(Pair{BTC, USDT}); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(openAmount)
		t.Log(timestamp)
	}

	if fees, _, err := bn.Swap.GetFundingFees(Pair{BTC, USDT}); err != nil {
		t.Error(err)
		return
	} else {
		t.Log(fees)
		//t.Log(string(resp))
	}
}

func TestFuture_TradeAPI(t *testing.T) {

	config := &APIConfig{
		Endpoint: ENDPOINT,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(PROXY_URL)
				},
			},
		},
		ApiKey:        SWAP_API_KEY,
		ApiSecretKey:  SWAP_API_SECRETKEY,
		ApiPassphrase: "",
		Location:      time.Now().Location(),
	}

	bn := New(config)
	if account, raw, err := bn.Swap.GetAccount(); err != nil {
		t.Error(err)
		return
	} else {

		rawAccount, _ := json.Marshal(account)
		t.Log(string(rawAccount))
		t.Log(string(raw))

		if account.BalanceAvail < 1 {
			t.Error("There have no enough asset to trade. ")
			return
		}
	}

	pair := Pair{LTC, USDT}
	ticker, _, err := bn.Swap.GetTicker(pair)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(ticker.Buy)

	order := SwapOrder{
		Cid:       UUID(),
		Price:     ticker.Sell * 1.03,
		Amount:    0.0123333333333,
		PlaceType: NORMAL,
		Type:      OPEN_SHORT,
		LeverRate: 20,
		Pair:      pair,
		Exchange:  BINANCE,
	}
	//preCid := order.Cid

	if resp, err := bn.Swap.PlaceOrder(&order); err != nil {
		t.Error(err)
		return
	} else {
		stdOrder, _ := json.Marshal(order)
		t.Log(string(resp))
		t.Log(string(stdOrder))
	}

	if resp, err := bn.Swap.GetOrder(&order); err != nil {
		t.Error(err)
		return
	} else {
		stdOrder, _ := json.Marshal(order)
		t.Log(string(resp))
		t.Log(string(stdOrder))
	}

	if resp, err := bn.Swap.CancelOrder(&order); err != nil {
		t.Error(err)
		return
	} else {
		stdOrder, _ := json.Marshal(order)
		t.Log(string(resp))
		t.Log(string(stdOrder))
	}

	//if orders, resp, err := bn.Swap.GetOrders(Pair{LTC, USDT}); err != nil {
	//	t.Error(err)
	//	return
	//} else {
	//	raw, _ := json.Marshal(orders)
	//	t.Log(string(raw))
	//	t.Log(string(resp))
	//}

	//if orders, resp, err := bn.Swap.GetUnFinishOrders(Pair{LTC, USDT}); err != nil {
	//	t.Error(err)
	//	return
	//} else {
	//
	//	raw, _ := json.Marshal(orders)
	//	t.Log(string(raw))
	//	t.Log(string(resp))
	//}
	//
	//if s, resp, err := bn.Swap.GetPosition(Pair{LTC, USDT}, OPEN_LONG); err != nil {
	//	t.Error(err)
	//	return
	//} else {
	//	t.Log(string(resp))
	//	t.Log(s)
	//}
}