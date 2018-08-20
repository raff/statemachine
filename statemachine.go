package statemachine

import (
	"log"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type State func() State

func (s State) String() string {
	if s == nil {
		return "<TERMINATE>"
	}

	path := runtime.FuncForPC(reflect.ValueOf(s).Pointer()).Name()
	return strings.TrimPrefix(path, "github.com/raff/statemachine.(*StateMachine).")
}

func (s State) Equal(s1 State) bool {
	p := reflect.ValueOf(s).Pointer()
	p1 := reflect.ValueOf(s1).Pointer()
	return p == p1
}

type StateMachine struct {
	state       chan State
	idle        *time.Timer
	idleTimeout time.Duration
}

func (sm *StateMachine) Init() {
	sm.state = make(chan State, 10)
	sm.state <- sm.Wait
}

func (sm *StateMachine) IdleTimeout(d time.Duration, cb func()) {
	if sm.idle != nil {
		sm.idle.Stop()
	}

	if cb != nil {
		sm.idleTimeout = d
		sm.idle = time.AfterFunc(sm.idleTimeout, cb)
	}
}

func (sm *StateMachine) Terminate() {
	if sm.state != nil {
		close(sm.state)
		sm.state = nil
	}

	if sm.idle != nil {
		sm.idle.Stop()
		sm.idle = nil
	}

	if len(sm.state) > 0 {
		log.Println("State queue:")
		for s := range sm.state {
			log.Println(" ", s)
		}
	}
}

func (sm *StateMachine) PushState(s State) {
	log.Println("PUSH STATE:", s, "pending:", len(sm.state))
	sm.state <- s
}

func (sm *StateMachine) Wait() State {
	time.Sleep(time.Second)
	return sm.Wait
}

func (sm *StateMachine) Run() {
	log.Println("STATE: <START>")

	for state := range sm.state {
		log.Println("STATE:", state)

		if state == nil {
			break
		} else {
			if sm.idle != nil {
				sm.idle.Reset(sm.idleTimeout)
			}

			next := state()
			log.Println("NEXT STATE:", next, len(sm.state))
			if state.Equal(sm.Wait) && next.Equal(sm.Wait) {
				continue
			}
			sm.state <- next
		}
	}

	log.Println("STATE: <STOP>")

	if sm.idle != nil {
		sm.idle.Stop()
	}
}
