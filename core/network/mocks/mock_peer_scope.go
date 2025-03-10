// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/zbliujia/go-libp2p/core/network (interfaces: PeerScope)
//
// Generated by this command:
//
//	mockgen -package mocknetwork -destination mock_peer_scope.go github.com/zbliujia/go-libp2p/core/network PeerScope
//
// Package mocknetwork is a generated GoMock package.
package mocknetwork

import (
	reflect "reflect"

	network "github.com/zbliujia/go-libp2p/core/network"
	peer "github.com/zbliujia/go-libp2p/core/peer"
	gomock "go.uber.org/mock/gomock"
)

// MockPeerScope is a mock of PeerScope interface.
type MockPeerScope struct {
	ctrl     *gomock.Controller
	recorder *MockPeerScopeMockRecorder
}

// MockPeerScopeMockRecorder is the mock recorder for MockPeerScope.
type MockPeerScopeMockRecorder struct {
	mock *MockPeerScope
}

// NewMockPeerScope creates a new mock instance.
func NewMockPeerScope(ctrl *gomock.Controller) *MockPeerScope {
	mock := &MockPeerScope{ctrl: ctrl}
	mock.recorder = &MockPeerScopeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerScope) EXPECT() *MockPeerScopeMockRecorder {
	return m.recorder
}

// BeginSpan mocks base method.
func (m *MockPeerScope) BeginSpan() (network.ResourceScopeSpan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginSpan")
	ret0, _ := ret[0].(network.ResourceScopeSpan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginSpan indicates an expected call of BeginSpan.
func (mr *MockPeerScopeMockRecorder) BeginSpan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginSpan", reflect.TypeOf((*MockPeerScope)(nil).BeginSpan))
}

// Peer mocks base method.
func (m *MockPeerScope) Peer() peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peer")
	ret0, _ := ret[0].(peer.ID)
	return ret0
}

// Peer indicates an expected call of Peer.
func (mr *MockPeerScopeMockRecorder) Peer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peer", reflect.TypeOf((*MockPeerScope)(nil).Peer))
}

// ReleaseMemory mocks base method.
func (m *MockPeerScope) ReleaseMemory(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReleaseMemory", arg0)
}

// ReleaseMemory indicates an expected call of ReleaseMemory.
func (mr *MockPeerScopeMockRecorder) ReleaseMemory(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReleaseMemory", reflect.TypeOf((*MockPeerScope)(nil).ReleaseMemory), arg0)
}

// ReserveMemory mocks base method.
func (m *MockPeerScope) ReserveMemory(arg0 int, arg1 byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReserveMemory", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReserveMemory indicates an expected call of ReserveMemory.
func (mr *MockPeerScopeMockRecorder) ReserveMemory(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReserveMemory", reflect.TypeOf((*MockPeerScope)(nil).ReserveMemory), arg0, arg1)
}

// Stat mocks base method.
func (m *MockPeerScope) Stat() network.ScopeStat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat")
	ret0, _ := ret[0].(network.ScopeStat)
	return ret0
}

// Stat indicates an expected call of Stat.
func (mr *MockPeerScopeMockRecorder) Stat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockPeerScope)(nil).Stat))
}
