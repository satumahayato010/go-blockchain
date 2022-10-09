package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Block Blockの情報が格納されている。構造体
// フィールド四つの情報を一塊として、ブロックを作る。それをハッシュ化する。
type Block struct {
	nonce        int
	previousHash [32]byte       // 前のハッシュの情報が入っている。
	timestamp    int64          // 取引した時のタイムスタンプ
	transactions []*Transaction // 取引内容 Pool
}

// NewBlock Blockのコンストラクタ。インスタンス化する。
func NewBlock(nonce int, previousHash [32]byte, transaction []*Transaction) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transaction
	return b
}

// Print Block構造体の中身を見やすいように作成したメソッド。
func (b *Block) Print() {
	fmt.Printf("timestamp        %d\n", b.timestamp)
	fmt.Printf("nonce            %d\n", b.nonce)
	fmt.Printf("previous_hash    %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

// Hash Block構造体をJSON形式にして、sha256(ハッシュ関数）でハッシュ化する。
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

// MarshalJSON Hashメソッドでは、マーシャルした時にBlock構造体のフィールドが小文字だから、MarshalJSONメソッドで、マーシャルをカスタマイズ
// する。先頭を大文字に変換しているだけ。
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

// Blockchain Block同士をchain（つなぐ）情報が格納されている構造体。
type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

// NewBlockchain Blockchain構造体のコンストラクタ。
func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

// LastBlock ブロックが格納されている、Blockchain構造体のchainフィールドの長さ-１で最後のブロックを取得してきている。
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

// CreateBlock Block構造体をインスタンス化して、Blockchain構造体のcahinフィールドに格納する。
func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

// AddTransaction トランザクションプールの中に、トランザクションを追加するメソッド。
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}

// Transaction 取引データを格納する構造体。自分のアドレス、取引相手のアドレス、いくら送ったかの値。
type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

// NewTransaction トランザクション構造体のコンストラクタ。
func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

// Print 格納されている値を見やすく出力するためのカスタムメソッド
func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address          %s\n", t.senderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address       %s\n", t.recipientBlockchainAddress)
	fmt.Printf(" value                              %.1f\n", t.value)

}

// MarshalJSON トランザクション構造体のフィールドが小文字だから、大文字にしてマーシャルするためのカスタムマーシャル。
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float32 `json:"value"`
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})
}

func main() {
	blockChain := NewBlockchain()
	blockChain.Print()
}
