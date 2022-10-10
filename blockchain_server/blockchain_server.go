package main

import (
	"go-trading/block"
	"go-trading/wallet"
	"io"
	"log"
	"net/http"
	"strconv"
)

// キャッシュにブロックを格納する。
var cache map[string]*block.Blockchain = make(map[string]*block.Blockchain)

// BlockchainServer ブロックチェーンサーバーの情報を格納する構造体。
type BlockchainServer struct {
	port uint16
}

// NewBlockchainServer ブロックチェーンサーバーの構造体のコンストラクタ。
func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

// Port サーバーのポート番号を返すメソッド。
func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

// GetBlockchain ブロックチェーン構造体（情報）を取得するメソッド。
func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minerWallet.BlockchainAddress(), bcs.Port())
		cache["blockchain"] = bc
	}
	return bc
}

// GetChain ブロックチェーン構造体の情報を取得するメソッド。
func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

// Run サーバーを走らせるメソッド。
func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.port)), nil))
}
