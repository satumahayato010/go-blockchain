package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const baseURL = "https://api.bitflyer.com/v1/"

// APIClient APIでアクセスしてくる、クライアントの情報を格納するための構造体。
type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

// New APIClient構造体をイニシャライズするメソッド。
func New(key, secret string) *APIClient {
	apiClient := &APIClient{key, secret, &http.Client{}}
	return apiClient
}

// 引数に、Getなどのメソッドと、URLのエンドポイント、渡したいデータをbyteで渡す。
func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + endpoint + string(body)

	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}

func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	// Parseメソッドで、引数のURLを、goのURL構造体に解析する。
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	//　引数で入ってきた、urlPathのurlを、goのURL構造体に解析する。
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	//ResolveReferenceメソッドで、相対リンクを絶対URLに変換する。
	endpoint := baseURL.ResolveReference(apiURL).String()
	// NewRequestメソッドは、新しいRequestを作成する、引数には、リクエストである。
	// getなどのメソッド、URI、渡したいデータを渡して、Request構造体に格納して、返ってくる。
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	// クエリを渡された時に、格納できるようにセットする
	//　Queryメソッドで、クエリを追加する。 これでクエリ付きでURlが生成される。
	q := req.URL.Query()
	// 引数のqueryを取り出して、アペンドする。
	for key, value := range query {
		q.Add(key, value)
	}
	// RawQueryするときは、エンコードしないといけない。
	// クエリが格納されている変数をエンコードし、ひとつのURLを生成する。
	req.URL.RawQuery = q.Encode()
	// api.headerメソッドで、引数の情報使って、ヘッダー情報を作成してから、
	// ループでヘッダーの情報を取り出して、リクエストのヘッダー情報に付け足していく。
	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}
	// httpClient.Doで、作成したリクエストを投げることで、URLにアクセスする。
	// レスポンス構造体が返ってくる。レスポンス情報が入っている。
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// レスポンスのボディを読み込んで、その中身を返している。
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Balance ビットフライヤーのAPI叩いた時のレスポンスを格納する構造体
type Balance struct {
	CurrentCode string  `json:"currentCode"`
	Amount      float64 `json:"amount"`
	Available   float64 `json:"available"`
}

// GetBalance ビットフライヤーのAPIを叩いて、Balanceのレスポンスを取得してくる関数。
func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "me/getbalance"
	// 引数の内容のリクエストを作成して送信、返り値は、レスポンスのボディ。
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}
	var balance []Balance
	// レスポンスできたボディの中身を、Balance構造体のスライスにデコード。
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return nil, err
	}
	// レスポンスの内容を、Balance構造体に変換して、スライスで返している。要は、APIからのレスポンスが返ってきているイメージ。
	return balance, nil
}

// Ticker API叩いて返ってきたレスポンスの情報を格納する構造体。
type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:" volume_by_product"`
}

// GetMidPrice 売りと買いの中間の値を取得してきたいので、それを取得する関数。要は、自分オリジナルで
// 必要ないといえばない。
func (t *Ticker) GetMidPrice() float64 {
	return (t.BestBid + t.BestAsk) / 2
}

func (t *Ticker) DateTime() time.Time {
	dateTime, err := time.Parse(time.RFC3339, t.Timestamp)
	if err != nil {
		log.Printf("action=DateTime, err=%s", err.Error())
	}
	return dateTime
}

func (t *Ticker) TruncateDateTime(duration time.Duration) time.Time {
	return t.DateTime().Truncate(duration)
}

func (api *APIClient) GetTicker(productCode string) (*Ticker, error) {
	url := "ticker"
	// 引数の内容のリクエストを作成して送信、返り値は、レスポンスのボディ。
	resp, err := api.doRequest("GET", url, map[string]string{"product_code": productCode}, nil)
	if err != nil {
		return nil, err
	}
	var ticker Ticker
	// レスポンスできたボディの中身を、Balance構造体のスライスにデコード。
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		return nil, err
	}
	// レスポンスの内容を、Balance構造体に変換して、スライスで返している。要は、APIからのレスポンスが返ってきているイメージ。
	return &ticker, nil
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

func (api *APIClient) GetRealTimeTicker(symbol string, ch chan<- Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", symbol)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{channel}}); err != nil {
		log.Fatal("subscribe:", err)
		return
	}

OUTER:
	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Println("read:", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshaTic, err := json.Marshal(binary)
						if err != nil {
							continue OUTER
						}
						var ticker Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}

type Order struct {
	ID                     int     `json:"id"`
	ChildOrderAcceptanceID string  `json:"child_order_acceptance_id"`
	ProductCode            string  `json:"product_code"`
	ChildOrderType         string  `json:"child_order_type"`
	Side                   string  `json:"side"`
	Price                  float64 `json:"price"`
	Size                   float64 `json:"size"`
	MinuteToExpires        int     `json:"minute_to_expire"`
	TimeInForce            string  `json:"time_in_force"`
	Status                 string  `json:"status"`
	ErrorMessage           string  `json:"error_message"`
	AveragePrice           float64 `json:"average_price"`
	ChildOrderState        string  `json:"child_order_state"`
	ExpireDate             string  `json:"expire_date"`
	ChildOrderDate         string  `json:"child_order_date"`
	OutstandingSize        float64 `json:"outstanding_size"`
	CancelSize             float64 `json:"cancel_size"`
	ExecutedSize           float64 `json:"executed_size"`
	TotalCommission        float64 `json:"total_commission"`
	Count                  int     `json:"count"`
	Before                 int     `json:"before"`
	After                  int     `json:"after"`
}

type ResponseSendChildOrder struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

//ChildOrderAcceptanceIDを取得するための関数
func (api *APIClient) SendOrder(order *Order) (*ResponseSendChildOrder, error) {
	// 引数で入ってきた、Order構造体をJOSNにエンコード
	data, _ := json.Marshal(order)
	// リクエストを送る、urlを指定
	url := "me/sendchildorder"
	// POSTで、Order構造体をJSONに変換したデータを、上のURLにリクエストを投げる。
	// resp変数には、リクエストのレスポンスが入っている。
	resp, _ := api.doRequest("POST", url, map[string]string{}, data)

	var response ResponseSendChildOrder
	// 返ってきたレスポンスを、ResponseSendChildOrder構造体にデコードして、JOSNから構造体に変換
	_ = json.Unmarshal(resp, &response)
	// レスポンスから返ってきた情報が格納されている、ResponseSendChildOrder構造体を返す。
	return &response, nil
}

// Order構造体に情報を取得するための関数。
func (api *APIClient) ListOrder(query map[string]string) ([]Order, error) {
	// GETで、引数に入ってきたクエリを付与して、リクエスト送る。
	resp, err := api.doRequest("GET", "me/getchildorders", query, nil)
	if err != nil {
		return nil, err
	}
	var responseListOrder []Order
	// 返ってきたレスポンスを、上の変数の構造体に格納する
	err = json.Unmarshal(resp, &responseListOrder)
	if err != nil {
		return nil, err
	}
	// それをリターンする。
	return responseListOrder, nil
}
