// Copyright (c) 2018-2019 The MATRIX Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php
// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package tests

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/matrix/go-matrix/common"
	"github.com/matrix/go-matrix/common/math"
)

var _ = (*stEnvMarshaling)(nil)

func (s stEnv) MarshalJSON() ([]byte, error) {
	type stEnv struct {
		Coinbase   common.UnprefixedAddress `json:"currentCoinbase"   gencodec:"required"`
		Difficulty *math.HexOrDecimal256    `json:"currentDifficulty" gencodec:"required"`
		GasLimit   math.HexOrDecimal64      `json:"currentGasLimit"   gencodec:"required"`
		Number     math.HexOrDecimal64      `json:"currentNumber"     gencodec:"required"`
		Timestamp  math.HexOrDecimal64      `json:"currentTimestamp"  gencodec:"required"`
	}
	var enc stEnv
	enc.Coinbase = common.UnprefixedAddress(s.Coinbase)
	enc.Difficulty = (*math.HexOrDecimal256)(s.Difficulty)
	enc.GasLimit = math.HexOrDecimal64(s.GasLimit)
	enc.Number = math.HexOrDecimal64(s.Number)
	enc.Timestamp = math.HexOrDecimal64(s.Timestamp)
	return json.Marshal(&enc)
}

func (s *stEnv) UnmarshalJSON(input []byte) error {
	type stEnv struct {
		Coinbase   *common.UnprefixedAddress `json:"currentCoinbase"   gencodec:"required"`
		Difficulty *math.HexOrDecimal256     `json:"currentDifficulty" gencodec:"required"`
		GasLimit   *math.HexOrDecimal64      `json:"currentGasLimit"   gencodec:"required"`
		Number     *math.HexOrDecimal64      `json:"currentNumber"     gencodec:"required"`
		Timestamp  *math.HexOrDecimal64      `json:"currentTimestamp"  gencodec:"required"`
	}
	var dec stEnv
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Coinbase == nil {
		return errors.New("missing required field 'currentCoinbase' for stEnv")
	}
	s.Coinbase = common.Address(*dec.Coinbase)
	if dec.Difficulty == nil {
		return errors.New("missing required field 'currentDifficulty' for stEnv")
	}
	s.Difficulty = (*big.Int)(dec.Difficulty)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'currentGasLimit' for stEnv")
	}
	s.GasLimit = uint64(*dec.GasLimit)
	if dec.Number == nil {
		return errors.New("missing required field 'currentNumber' for stEnv")
	}
	s.Number = uint64(*dec.Number)
	if dec.Timestamp == nil {
		return errors.New("missing required field 'currentTimestamp' for stEnv")
	}
	s.Timestamp = uint64(*dec.Timestamp)
	return nil
}
