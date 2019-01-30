// Copyright (c) 2018-2019 The MATRIX Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php

package matrixstate

import (
	"github.com/matrix/go-matrix/common"
	"github.com/matrix/go-matrix/log"
	"github.com/pkg/errors"
)

var (
	ErrOptNotExist  = errors.New("operator not exist in manager")
	ErrStateDBNil   = errors.New("state db is nil")
	ErrParamReflect = errors.New("param reflect failed")
	ErrDataEmpty    = errors.New("data is empty")
	ErrParamNil     = errors.New("param is nil")
	ErrAccountNil   = errors.New("account is empty account")
	ErrDataSize     = errors.New("data size err")
	ErrFindManager  = errors.New("find manger err")
)

type StateDB interface {
	GetMatrixData(hash common.Hash) (val []byte)
	SetMatrixData(hash common.Hash, val []byte)
}

const (
	matrixStatePrefix = "ms_"
)

func checkStateDB(st StateDB) error {
	if st == nil {
		log.Error(logInfo, "stateDB err", ErrStateDBNil)
		return ErrStateDBNil
	}
	return nil
}
