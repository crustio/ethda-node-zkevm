// Code generated by mockery v2.39.0. DO NOT EDIT.

package sequencesender

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	coretypes "github.com/ethereum/go-ethereum/core/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/0xPolygonHermez/zkevm-node/etherman/types"
)

// EthermanMock is an autogenerated mock type for the etherman type
type EthermanMock struct {
	mock.Mock
}

// BuildPostZkBlobTxData provides a mock function with given fields: sender, batchNumber, hashes
func (_m *EthermanMock) BuildPostZkBlobTxData(sender common.Address, batchNumber int64, hashes []common.Hash) (*common.Address, []byte, error) {
	ret := _m.Called(sender, batchNumber, hashes)

	if len(ret) == 0 {
		panic("no return value specified for BuildPostZkBlobTxData")
	}

	var r0 *common.Address
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func(common.Address, int64, []common.Hash) (*common.Address, []byte, error)); ok {
		return rf(sender, batchNumber, hashes)
	}
	if rf, ok := ret.Get(0).(func(common.Address, int64, []common.Hash) *common.Address); ok {
		r0 = rf(sender, batchNumber, hashes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Address, int64, []common.Hash) []byte); ok {
		r1 = rf(sender, batchNumber, hashes)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(common.Address, int64, []common.Hash) error); ok {
		r2 = rf(sender, batchNumber, hashes)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// BuildSequenceBatchesTxData provides a mock function with given fields: sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase
func (_m *EthermanMock) BuildSequenceBatchesTxData(sender common.Address, sequences []types.Sequence, maxSequenceTimestamp uint64, initSequenceBatchNumber uint64, l2Coinbase common.Address) (*common.Address, []byte, error) {
	ret := _m.Called(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)

	if len(ret) == 0 {
		panic("no return value specified for BuildSequenceBatchesTxData")
	}

	var r0 *common.Address
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) (*common.Address, []byte, error)); ok {
		return rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	}
	if rf, ok := ret.Get(0).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) *common.Address); ok {
		r0 = rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) []byte); ok {
		r1 = rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) error); ok {
		r2 = rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// EstimateGasSequenceBatches provides a mock function with given fields: sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase
func (_m *EthermanMock) EstimateGasSequenceBatches(sender common.Address, sequences []types.Sequence, maxSequenceTimestamp uint64, initSequenceBatchNumber uint64, l2Coinbase common.Address) (*coretypes.Transaction, error) {
	ret := _m.Called(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)

	if len(ret) == 0 {
		panic("no return value specified for EstimateGasSequenceBatches")
	}

	var r0 *coretypes.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) (*coretypes.Transaction, error)); ok {
		return rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	}
	if rf, ok := ret.Get(0).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) *coretypes.Transaction); ok {
		r0 = rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Address, []types.Sequence, uint64, uint64, common.Address) error); ok {
		r1 = rf(sender, sequences, maxSequenceTimestamp, initSequenceBatchNumber, l2Coinbase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBatchNumber provides a mock function with given fields:
func (_m *EthermanMock) GetLatestBatchNumber() (uint64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBatchNumber")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestBlockHeader provides a mock function with given fields: ctx
func (_m *EthermanMock) GetLatestBlockHeader(ctx context.Context) (*coretypes.Header, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlockHeader")
	}

	var r0 *coretypes.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*coretypes.Header, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *coretypes.Header); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEthermanMock creates a new instance of EthermanMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEthermanMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *EthermanMock {
	mock := &EthermanMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
