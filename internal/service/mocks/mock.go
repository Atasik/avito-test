// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"
	domain "segmenter/internal/domain"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockHistory is a mock of History interface.
type MockHistory struct {
	ctrl     *gomock.Controller
	recorder *MockHistoryMockRecorder
}

// MockHistoryMockRecorder is the mock recorder for MockHistory.
type MockHistoryMockRecorder struct {
	mock *MockHistory
}

// NewMockHistory creates a new mock instance.
func NewMockHistory(ctrl *gomock.Controller) *MockHistory {
	mock := &MockHistory{ctrl: ctrl}
	mock.recorder = &MockHistoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHistory) EXPECT() *MockHistoryMockRecorder {
	return m.recorder
}

// CreateReport mocks base method.
func (m *MockHistory) CreateReport(period time.Time, userID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReport", period, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReport indicates an expected call of CreateReport.
func (mr *MockHistoryMockRecorder) CreateReport(period, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReport", reflect.TypeOf((*MockHistory)(nil).CreateReport), period, userID)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// DeleteExpiredSegments mocks base method.
func (m *MockUser) DeleteExpiredSegments() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredSegments")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExpiredSegments indicates an expected call of DeleteExpiredSegments.
func (mr *MockUserMockRecorder) DeleteExpiredSegments() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredSegments", reflect.TypeOf((*MockUser)(nil).DeleteExpiredSegments))
}

// GetSegments mocks base method.
func (m *MockUser) GetSegments(userID int) ([]domain.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSegments", userID)
	ret0, _ := ret[0].([]domain.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSegments indicates an expected call of GetSegments.
func (mr *MockUserMockRecorder) GetSegments(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSegments", reflect.TypeOf((*MockUser)(nil).GetSegments), userID)
}

// UpsertUserSegments mocks base method.
func (m *MockUser) UpsertUserSegments(userID int, segmentsToAdd, segmentToDelete []domain.Segment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertUserSegments", userID, segmentsToAdd, segmentToDelete)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertUserSegments indicates an expected call of UpsertUserSegments.
func (mr *MockUserMockRecorder) UpsertUserSegments(userID, segmentsToAdd, segmentToDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUserSegments", reflect.TypeOf((*MockUser)(nil).UpsertUserSegments), userID, segmentsToAdd, segmentToDelete)
}

// MockSegment is a mock of Segment interface.
type MockSegment struct {
	ctrl     *gomock.Controller
	recorder *MockSegmentMockRecorder
}

// MockSegmentMockRecorder is the mock recorder for MockSegment.
type MockSegmentMockRecorder struct {
	mock *MockSegment
}

// NewMockSegment creates a new mock instance.
func NewMockSegment(ctrl *gomock.Controller) *MockSegment {
	mock := &MockSegment{ctrl: ctrl}
	mock.recorder = &MockSegmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSegment) EXPECT() *MockSegmentMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSegment) Create(seg domain.Segment) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", seg)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSegmentMockRecorder) Create(seg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSegment)(nil).Create), seg)
}

// Delete mocks base method.
func (m *MockSegment) Delete(seg domain.Segment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", seg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSegmentMockRecorder) Delete(seg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSegment)(nil).Delete), seg)
}
