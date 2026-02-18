package events

import (
	eventsrepo "rootstock/web-server/repo/events"
)

type getContextReq struct {
	resp chan eventsrepo.WorkflowContext
}

type shutdownReq struct {
	resp chan struct{}
}

var (
	getContextCh = make(chan getContextReq)
	shutdownCh   = make(chan shutdownReq)
)

// Initialize starts the manager goroutine that delegates to the provided repository.
func Initialize(repo eventsrepo.Repository) {
	go manage(repo)
}

// manage owns the repo reference. All access goes through channels.
func manage(repo eventsrepo.Repository) {
	for {
		select {
		case req := <-getContextCh:
			req.resp <- repo.GetContext()

		case req := <-shutdownCh:
			repo.Shutdown()
			req.resp <- struct{}{}
			return
		}
	}
}

// Shutdown gracefully shuts down the events runtime.
func Shutdown() {
	resp := make(chan struct{}, 1)
	shutdownCh <- shutdownReq{resp: resp}
	<-resp
}

// GetContext returns the workflow context for event execution.
func GetContext() eventsrepo.WorkflowContext {
	resp := make(chan eventsrepo.WorkflowContext, 1)
	getContextCh <- getContextReq{resp: resp}
	return <-resp
}
