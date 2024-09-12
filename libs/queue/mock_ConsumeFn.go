// Code generated by mockery v2.45.0. DO NOT EDIT.

package queue

import mock "github.com/stretchr/testify/mock"

// MockConsumeFn is an autogenerated mock type for the ConsumeFn type
type MockConsumeFn struct {
	mock.Mock
}

type MockConsumeFn_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConsumeFn) EXPECT() *MockConsumeFn_Expecter {
	return &MockConsumeFn_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: msg
func (_m *MockConsumeFn) Execute(msg Message) {
	_m.Called(msg)
}

// MockConsumeFn_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockConsumeFn_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - msg Message
func (_e *MockConsumeFn_Expecter) Execute(msg interface{}) *MockConsumeFn_Execute_Call {
	return &MockConsumeFn_Execute_Call{Call: _e.mock.On("Execute", msg)}
}

func (_c *MockConsumeFn_Execute_Call) Run(run func(msg Message)) *MockConsumeFn_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Message))
	})
	return _c
}

func (_c *MockConsumeFn_Execute_Call) Return() *MockConsumeFn_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockConsumeFn_Execute_Call) RunAndReturn(run func(Message)) *MockConsumeFn_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockConsumeFn creates a new instance of MockConsumeFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConsumeFn(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConsumeFn {
	mock := &MockConsumeFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
