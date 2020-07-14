// Copyright 2019 Conflux Foundation. All rights reserved.
// Conflux is free software and distributed under GNU General Public License.
// See http://www.gnu.org/licenses/

package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

// signedTransactionForRlp is a intermediate struct for encoding rlp data
type signedTransactionForRlp struct {
	UnsignedData *unsignedTransactionForRlp
	V            byte
	R            *big.Int
	S            *big.Int
}

// SignedTransaction represents a transaction with signature,
// it is the transaction information for sending transaction.
type SignedTransaction struct {
	UnsignedTransaction UnsignedTransaction
	V                   byte
	R                   hexutil.Bytes
	S                   hexutil.Bytes
}

// Decode decodes RLP encoded data to tx
func (tx *SignedTransaction) Decode(data []byte) error {
	txForRlp := new(signedTransactionForRlp)
	err := rlp.DecodeBytes(data, txForRlp)
	if err != nil {
		msg := fmt.Sprintf("decode data {%+x} to rlp error", data)
		return WrapError(err, msg)
	}

	*tx = *txForRlp.toSignedTransaction()
	return nil
}

//Encode encodes tx and returns its RLP encoded data
func (tx *SignedTransaction) Encode() ([]byte, error) {
	txForRlp := *tx.toStructForRlp()
	encoded, err := rlp.EncodeToBytes(txForRlp)
	if err != nil {
		msg := fmt.Sprintf("encode data {%+v} to bytes error", txForRlp)
		return nil, WrapError(err, msg)
	}

	return encoded, nil
}

func (tx *SignedTransaction) toStructForRlp() *signedTransactionForRlp {
	txForRlp := signedTransactionForRlp{
		UnsignedData: tx.UnsignedTransaction.toStructForRlp(),
		V:            tx.V,
		R:            big.NewInt(0).SetBytes(tx.R),
		S:            big.NewInt(0).SetBytes(tx.S),
	}
	return &txForRlp
}

func (tx *signedTransactionForRlp) toSignedTransaction() *SignedTransaction {
	unsigned := tx.UnsignedData.toUnsignedTransaction()
	return &SignedTransaction{
		UnsignedTransaction: *unsigned,
		V:                   tx.V,
		R:                   tx.R.Bytes(),
		S:                   tx.S.Bytes(),
	}
}
