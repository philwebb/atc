// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/concourse/atc/api/jobserver"
	"github.com/concourse/atc/builds"
)

type FakeJobsDB struct {
	GetAllJobBuildsStub        func(job string) ([]builds.Build, error)
	getAllJobBuildsMutex       sync.RWMutex
	getAllJobBuildsArgsForCall []struct {
		job string
	}
	getAllJobBuildsReturns struct {
		result1 []builds.Build
		result2 error
	}
	GetCurrentBuildStub        func(job string) (builds.Build, error)
	getCurrentBuildMutex       sync.RWMutex
	getCurrentBuildArgsForCall []struct {
		job string
	}
	getCurrentBuildReturns struct {
		result1 builds.Build
		result2 error
	}
	GetJobBuildStub        func(job string, build string) (builds.Build, error)
	getJobBuildMutex       sync.RWMutex
	getJobBuildArgsForCall []struct {
		job   string
		build string
	}
	getJobBuildReturns struct {
		result1 builds.Build
		result2 error
	}
	GetJobFinishedAndNextBuildStub        func(job string) (*builds.Build, *builds.Build, error)
	getJobFinishedAndNextBuildMutex       sync.RWMutex
	getJobFinishedAndNextBuildArgsForCall []struct {
		job string
	}
	getJobFinishedAndNextBuildReturns struct {
		result1 *builds.Build
		result2 *builds.Build
		result3 error
	}
}

func (fake *FakeJobsDB) GetAllJobBuilds(job string) ([]builds.Build, error) {
	fake.getAllJobBuildsMutex.Lock()
	fake.getAllJobBuildsArgsForCall = append(fake.getAllJobBuildsArgsForCall, struct {
		job string
	}{job})
	fake.getAllJobBuildsMutex.Unlock()
	if fake.GetAllJobBuildsStub != nil {
		return fake.GetAllJobBuildsStub(job)
	} else {
		return fake.getAllJobBuildsReturns.result1, fake.getAllJobBuildsReturns.result2
	}
}

func (fake *FakeJobsDB) GetAllJobBuildsCallCount() int {
	fake.getAllJobBuildsMutex.RLock()
	defer fake.getAllJobBuildsMutex.RUnlock()
	return len(fake.getAllJobBuildsArgsForCall)
}

func (fake *FakeJobsDB) GetAllJobBuildsArgsForCall(i int) string {
	fake.getAllJobBuildsMutex.RLock()
	defer fake.getAllJobBuildsMutex.RUnlock()
	return fake.getAllJobBuildsArgsForCall[i].job
}

func (fake *FakeJobsDB) GetAllJobBuildsReturns(result1 []builds.Build, result2 error) {
	fake.GetAllJobBuildsStub = nil
	fake.getAllJobBuildsReturns = struct {
		result1 []builds.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeJobsDB) GetCurrentBuild(job string) (builds.Build, error) {
	fake.getCurrentBuildMutex.Lock()
	fake.getCurrentBuildArgsForCall = append(fake.getCurrentBuildArgsForCall, struct {
		job string
	}{job})
	fake.getCurrentBuildMutex.Unlock()
	if fake.GetCurrentBuildStub != nil {
		return fake.GetCurrentBuildStub(job)
	} else {
		return fake.getCurrentBuildReturns.result1, fake.getCurrentBuildReturns.result2
	}
}

func (fake *FakeJobsDB) GetCurrentBuildCallCount() int {
	fake.getCurrentBuildMutex.RLock()
	defer fake.getCurrentBuildMutex.RUnlock()
	return len(fake.getCurrentBuildArgsForCall)
}

func (fake *FakeJobsDB) GetCurrentBuildArgsForCall(i int) string {
	fake.getCurrentBuildMutex.RLock()
	defer fake.getCurrentBuildMutex.RUnlock()
	return fake.getCurrentBuildArgsForCall[i].job
}

func (fake *FakeJobsDB) GetCurrentBuildReturns(result1 builds.Build, result2 error) {
	fake.GetCurrentBuildStub = nil
	fake.getCurrentBuildReturns = struct {
		result1 builds.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeJobsDB) GetJobBuild(job string, build string) (builds.Build, error) {
	fake.getJobBuildMutex.Lock()
	fake.getJobBuildArgsForCall = append(fake.getJobBuildArgsForCall, struct {
		job   string
		build string
	}{job, build})
	fake.getJobBuildMutex.Unlock()
	if fake.GetJobBuildStub != nil {
		return fake.GetJobBuildStub(job, build)
	} else {
		return fake.getJobBuildReturns.result1, fake.getJobBuildReturns.result2
	}
}

func (fake *FakeJobsDB) GetJobBuildCallCount() int {
	fake.getJobBuildMutex.RLock()
	defer fake.getJobBuildMutex.RUnlock()
	return len(fake.getJobBuildArgsForCall)
}

func (fake *FakeJobsDB) GetJobBuildArgsForCall(i int) (string, string) {
	fake.getJobBuildMutex.RLock()
	defer fake.getJobBuildMutex.RUnlock()
	return fake.getJobBuildArgsForCall[i].job, fake.getJobBuildArgsForCall[i].build
}

func (fake *FakeJobsDB) GetJobBuildReturns(result1 builds.Build, result2 error) {
	fake.GetJobBuildStub = nil
	fake.getJobBuildReturns = struct {
		result1 builds.Build
		result2 error
	}{result1, result2}
}

func (fake *FakeJobsDB) GetJobFinishedAndNextBuild(job string) (*builds.Build, *builds.Build, error) {
	fake.getJobFinishedAndNextBuildMutex.Lock()
	fake.getJobFinishedAndNextBuildArgsForCall = append(fake.getJobFinishedAndNextBuildArgsForCall, struct {
		job string
	}{job})
	fake.getJobFinishedAndNextBuildMutex.Unlock()
	if fake.GetJobFinishedAndNextBuildStub != nil {
		return fake.GetJobFinishedAndNextBuildStub(job)
	} else {
		return fake.getJobFinishedAndNextBuildReturns.result1, fake.getJobFinishedAndNextBuildReturns.result2, fake.getJobFinishedAndNextBuildReturns.result3
	}
}

func (fake *FakeJobsDB) GetJobFinishedAndNextBuildCallCount() int {
	fake.getJobFinishedAndNextBuildMutex.RLock()
	defer fake.getJobFinishedAndNextBuildMutex.RUnlock()
	return len(fake.getJobFinishedAndNextBuildArgsForCall)
}

func (fake *FakeJobsDB) GetJobFinishedAndNextBuildArgsForCall(i int) string {
	fake.getJobFinishedAndNextBuildMutex.RLock()
	defer fake.getJobFinishedAndNextBuildMutex.RUnlock()
	return fake.getJobFinishedAndNextBuildArgsForCall[i].job
}

func (fake *FakeJobsDB) GetJobFinishedAndNextBuildReturns(result1 *builds.Build, result2 *builds.Build, result3 error) {
	fake.GetJobFinishedAndNextBuildStub = nil
	fake.getJobFinishedAndNextBuildReturns = struct {
		result1 *builds.Build
		result2 *builds.Build
		result3 error
	}{result1, result2, result3}
}

var _ jobserver.JobsDB = new(FakeJobsDB)