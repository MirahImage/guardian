// Code generated by counterfeiter. DO NOT EDIT.
package runruncfakes

import (
	"sync"

	"code.cloudfoundry.org/guardian/rundmc/runrunc"
	"code.cloudfoundry.org/lager"
)

type FakePidGetter struct {
	GetPidStub        func(log lager.Logger, containerHandle string) (int, error)
	getPidMutex       sync.RWMutex
	getPidArgsForCall []struct {
		log             lager.Logger
		containerHandle string
	}
	getPidReturns struct {
		result1 int
		result2 error
	}
	getPidReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePidGetter) GetPid(log lager.Logger, containerHandle string) (int, error) {
	fake.getPidMutex.Lock()
	ret, specificReturn := fake.getPidReturnsOnCall[len(fake.getPidArgsForCall)]
	fake.getPidArgsForCall = append(fake.getPidArgsForCall, struct {
		log             lager.Logger
		containerHandle string
	}{log, containerHandle})
	fake.recordInvocation("GetPid", []interface{}{log, containerHandle})
	fake.getPidMutex.Unlock()
	if fake.GetPidStub != nil {
		return fake.GetPidStub(log, containerHandle)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPidReturns.result1, fake.getPidReturns.result2
}

func (fake *FakePidGetter) GetPidCallCount() int {
	fake.getPidMutex.RLock()
	defer fake.getPidMutex.RUnlock()
	return len(fake.getPidArgsForCall)
}

func (fake *FakePidGetter) GetPidArgsForCall(i int) (lager.Logger, string) {
	fake.getPidMutex.RLock()
	defer fake.getPidMutex.RUnlock()
	return fake.getPidArgsForCall[i].log, fake.getPidArgsForCall[i].containerHandle
}

func (fake *FakePidGetter) GetPidReturns(result1 int, result2 error) {
	fake.GetPidStub = nil
	fake.getPidReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePidGetter) GetPidReturnsOnCall(i int, result1 int, result2 error) {
	fake.GetPidStub = nil
	if fake.getPidReturnsOnCall == nil {
		fake.getPidReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.getPidReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePidGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getPidMutex.RLock()
	defer fake.getPidMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePidGetter) recordInvocation(key string, args []interface{}) {
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

var _ runrunc.PidGetter = new(FakePidGetter)
