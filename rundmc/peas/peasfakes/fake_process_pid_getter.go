// Code generated by counterfeiter. DO NOT EDIT.
package peasfakes

import (
	"sync"

	"code.cloudfoundry.org/guardian/rundmc/peas"
)

type FakeProcessPidGetter struct {
	PidStub        func(pidFilePath string) (int, error)
	pidMutex       sync.RWMutex
	pidArgsForCall []struct {
		pidFilePath string
	}
	pidReturns struct {
		result1 int
		result2 error
	}
	pidReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProcessPidGetter) Pid(pidFilePath string) (int, error) {
	fake.pidMutex.Lock()
	ret, specificReturn := fake.pidReturnsOnCall[len(fake.pidArgsForCall)]
	fake.pidArgsForCall = append(fake.pidArgsForCall, struct {
		pidFilePath string
	}{pidFilePath})
	fake.recordInvocation("Pid", []interface{}{pidFilePath})
	fake.pidMutex.Unlock()
	if fake.PidStub != nil {
		return fake.PidStub(pidFilePath)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.pidReturns.result1, fake.pidReturns.result2
}

func (fake *FakeProcessPidGetter) PidCallCount() int {
	fake.pidMutex.RLock()
	defer fake.pidMutex.RUnlock()
	return len(fake.pidArgsForCall)
}

func (fake *FakeProcessPidGetter) PidArgsForCall(i int) string {
	fake.pidMutex.RLock()
	defer fake.pidMutex.RUnlock()
	return fake.pidArgsForCall[i].pidFilePath
}

func (fake *FakeProcessPidGetter) PidReturns(result1 int, result2 error) {
	fake.PidStub = nil
	fake.pidReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeProcessPidGetter) PidReturnsOnCall(i int, result1 int, result2 error) {
	fake.PidStub = nil
	if fake.pidReturnsOnCall == nil {
		fake.pidReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.pidReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeProcessPidGetter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.pidMutex.RLock()
	defer fake.pidMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeProcessPidGetter) recordInvocation(key string, args []interface{}) {
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

var _ peas.ProcessPidGetter = new(FakeProcessPidGetter)
