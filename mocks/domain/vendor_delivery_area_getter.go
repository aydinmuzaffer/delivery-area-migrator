// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/aydinmuzaffer/migration-tool-service/src/migrationtool/domain"
	mock "github.com/stretchr/testify/mock"
)

// VendorDeliveryAreaGetter is an autogenerated mock type for the VendorDeliveryAreaGetter type
type VendorDeliveryAreaGetter struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, ids
func (_m *VendorDeliveryAreaGetter) Get(ctx context.Context, ids *[]int64) (*[]domain.VendorDeliveryAreaModel, error) {
	ret := _m.Called(ctx, ids)

	var r0 *[]domain.VendorDeliveryAreaModel
	if rf, ok := ret.Get(0).(func(context.Context, *[]int64) *[]domain.VendorDeliveryAreaModel); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.VendorDeliveryAreaModel)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *[]int64) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetIds provides a mock function with given fields: ctx
func (_m *VendorDeliveryAreaGetter) GetIds(ctx context.Context) (*[]int64, error) {
	ret := _m.Called(ctx)

	var r0 *[]int64
	if rf, ok := ret.Get(0).(func(context.Context) *[]int64); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]int64)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewVendorDeliveryAreaGetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewVendorDeliveryAreaGetter creates a new instance of VendorDeliveryAreaGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVendorDeliveryAreaGetter(t mockConstructorTestingTNewVendorDeliveryAreaGetter) *VendorDeliveryAreaGetter {
	mock := &VendorDeliveryAreaGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
