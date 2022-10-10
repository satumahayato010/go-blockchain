package utils

import (
	"fmt"
	"math/big"
)

// Signature シグネチャの情報を格納した構造体。
type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
