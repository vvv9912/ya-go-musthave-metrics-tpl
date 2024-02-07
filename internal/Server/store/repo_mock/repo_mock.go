// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repo_mock is a generated GoMock package.
package repo_mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// GetAllCounter mocks base method.
func (m *MockStorager) GetAllCounter(ctx context.Context) (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCounter", ctx)
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCounter indicates an expected call of GetAllCounter.
func (mr *MockStoragerMockRecorder) GetAllCounter(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCounter", reflect.TypeOf((*MockStorager)(nil).GetAllCounter), ctx)
}

// GetAllGauge mocks base method.
func (m *MockStorager) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllGauge", ctx)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllGauge indicates an expected call of GetAllGauge.
func (mr *MockStoragerMockRecorder) GetAllGauge(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllGauge", reflect.TypeOf((*MockStorager)(nil).GetAllGauge), ctx)
}

// GetCounter mocks base method.
func (m *MockStorager) GetCounter(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockStoragerMockRecorder) GetCounter(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockStorager)(nil).GetCounter), ctx, key)
}

// GetGauge mocks base method.
func (m *MockStorager) GetGauge(ctx context.Context, key string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", ctx, key)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockStoragerMockRecorder) GetGauge(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockStorager)(nil).GetGauge), ctx, key)
}

// Ping mocks base method.
func (m *MockStorager) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoragerMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorager)(nil).Ping), ctx)
}

// UpdateCounter mocks base method.
func (m *MockStorager) UpdateCounter(ctx context.Context, key string, val int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockStoragerMockRecorder) UpdateCounter(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockStorager)(nil).UpdateCounter), ctx, key, val)
}

// UpdateGauge mocks base method.
func (m *MockStorager) UpdateGauge(ctx context.Context, key string, val float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGauge", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockStoragerMockRecorder) UpdateGauge(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockStorager)(nil).UpdateGauge), ctx, key, val)
}

// UpdateMetricsBatch mocks base method.
func (m *MockStorager) UpdateMetricsBatch(ctx context.Context, metrics []model.Metrics) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricsBatch", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricsBatch indicates an expected call of UpdateMetricsBatch.
func (mr *MockStoragerMockRecorder) UpdateMetricsBatch(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricsBatch", reflect.TypeOf((*MockStorager)(nil).UpdateMetricsBatch), ctx, metrics)
}