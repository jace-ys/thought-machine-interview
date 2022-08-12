// Code generated by counterfeiter. DO NOT EDIT.
package domainfakes

import (
	"context"
	"net"
	"sync"

	"github.com/jace-ys/thought-machine-interview/internal/domain"
)

type FakeMonitoringService struct {
	GetServerStub        func(context.Context, net.IP) (*domain.Server, error)
	getServerMutex       sync.RWMutex
	getServerArgsForCall []struct {
		arg1 context.Context
		arg2 net.IP
	}
	getServerReturns struct {
		result1 *domain.Server
		result2 error
	}
	getServerReturnsOnCall map[int]struct {
		result1 *domain.Server
		result2 error
	}
	ListServersStub        func(context.Context) ([]net.IP, error)
	listServersMutex       sync.RWMutex
	listServersArgsForCall []struct {
		arg1 context.Context
	}
	listServersReturns struct {
		result1 []net.IP
		result2 error
	}
	listServersReturnsOnCall map[int]struct {
		result1 []net.IP
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMonitoringService) GetServer(arg1 context.Context, arg2 net.IP) (*domain.Server, error) {
	fake.getServerMutex.Lock()
	ret, specificReturn := fake.getServerReturnsOnCall[len(fake.getServerArgsForCall)]
	fake.getServerArgsForCall = append(fake.getServerArgsForCall, struct {
		arg1 context.Context
		arg2 net.IP
	}{arg1, arg2})
	stub := fake.GetServerStub
	fakeReturns := fake.getServerReturns
	fake.recordInvocation("GetServer", []interface{}{arg1, arg2})
	fake.getServerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeMonitoringService) GetServerCallCount() int {
	fake.getServerMutex.RLock()
	defer fake.getServerMutex.RUnlock()
	return len(fake.getServerArgsForCall)
}

func (fake *FakeMonitoringService) GetServerCalls(stub func(context.Context, net.IP) (*domain.Server, error)) {
	fake.getServerMutex.Lock()
	defer fake.getServerMutex.Unlock()
	fake.GetServerStub = stub
}

func (fake *FakeMonitoringService) GetServerArgsForCall(i int) (context.Context, net.IP) {
	fake.getServerMutex.RLock()
	defer fake.getServerMutex.RUnlock()
	argsForCall := fake.getServerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeMonitoringService) GetServerReturns(result1 *domain.Server, result2 error) {
	fake.getServerMutex.Lock()
	defer fake.getServerMutex.Unlock()
	fake.GetServerStub = nil
	fake.getServerReturns = struct {
		result1 *domain.Server
		result2 error
	}{result1, result2}
}

func (fake *FakeMonitoringService) GetServerReturnsOnCall(i int, result1 *domain.Server, result2 error) {
	fake.getServerMutex.Lock()
	defer fake.getServerMutex.Unlock()
	fake.GetServerStub = nil
	if fake.getServerReturnsOnCall == nil {
		fake.getServerReturnsOnCall = make(map[int]struct {
			result1 *domain.Server
			result2 error
		})
	}
	fake.getServerReturnsOnCall[i] = struct {
		result1 *domain.Server
		result2 error
	}{result1, result2}
}

func (fake *FakeMonitoringService) ListServers(arg1 context.Context) ([]net.IP, error) {
	fake.listServersMutex.Lock()
	ret, specificReturn := fake.listServersReturnsOnCall[len(fake.listServersArgsForCall)]
	fake.listServersArgsForCall = append(fake.listServersArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.ListServersStub
	fakeReturns := fake.listServersReturns
	fake.recordInvocation("ListServers", []interface{}{arg1})
	fake.listServersMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeMonitoringService) ListServersCallCount() int {
	fake.listServersMutex.RLock()
	defer fake.listServersMutex.RUnlock()
	return len(fake.listServersArgsForCall)
}

func (fake *FakeMonitoringService) ListServersCalls(stub func(context.Context) ([]net.IP, error)) {
	fake.listServersMutex.Lock()
	defer fake.listServersMutex.Unlock()
	fake.ListServersStub = stub
}

func (fake *FakeMonitoringService) ListServersArgsForCall(i int) context.Context {
	fake.listServersMutex.RLock()
	defer fake.listServersMutex.RUnlock()
	argsForCall := fake.listServersArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeMonitoringService) ListServersReturns(result1 []net.IP, result2 error) {
	fake.listServersMutex.Lock()
	defer fake.listServersMutex.Unlock()
	fake.ListServersStub = nil
	fake.listServersReturns = struct {
		result1 []net.IP
		result2 error
	}{result1, result2}
}

func (fake *FakeMonitoringService) ListServersReturnsOnCall(i int, result1 []net.IP, result2 error) {
	fake.listServersMutex.Lock()
	defer fake.listServersMutex.Unlock()
	fake.ListServersStub = nil
	if fake.listServersReturnsOnCall == nil {
		fake.listServersReturnsOnCall = make(map[int]struct {
			result1 []net.IP
			result2 error
		})
	}
	fake.listServersReturnsOnCall[i] = struct {
		result1 []net.IP
		result2 error
	}{result1, result2}
}

func (fake *FakeMonitoringService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getServerMutex.RLock()
	defer fake.getServerMutex.RUnlock()
	fake.listServersMutex.RLock()
	defer fake.listServersMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMonitoringService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ domain.MonitoringService = new(FakeMonitoringService)
