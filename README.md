# statemachine
A simple state machine processor

States are functions that return the next state. The state machine provides a "state processor" that runs through the various states until done.

## Example
This example shows a state machine with two main states (Ping and Pong). The processor flips between the two for 5 times and then terminates.

    type SmProcessor struct {
	statemachine.StateMachine

        // your stateful info here
        count int
    }

    func (p *SmProcessor) Ping() statemachine.State {
	log.Println("in Ping")
	return p.Pong
    }

    func (p *SmProcessor) Pong() statemachine.State {
	log.Println("in Pong")
        if p.count < 5 {
            p.count += 1
            p.PushState(p.Wait) // add p.Wait if you know you need to wait for an event

            return p.Ping
        }
            
	return p.Last
    }

    func (p *SmProcessor) Last() statemachine.State {
	log.Println("Done!")
	return nil // return nil to terminate the state machine
    }

    func main() {
	sm := &SmProcessor{}

	sm.Init()

	sm.PushState(sm.Ping) // push initial state

	sm.Run()
	sm.Terminate() // call this to terminate early
    }
