package actor

type AutoReceiveMessage interface {
	AutoReceiveMessage()
}

// An actor which receives a NotInfluenceReceiveTimeout message will not reset the ReceiveTimeout duration
type NotInfluenceReceiveTimeout interface {
	NotInfluenceReceiveTimeout()
}

// A SystemMessage message is reserved for specific lifecycle messages used by the actor system
type SystemMessage interface {
	SystemMessage()
}

// A ReceiveTimeout message is sent to an actor after the Context.ReceiveTimeout duration has expired
type ReceiveTimeout struct{}

// A Restarting message is sent to an actor when the actor is being restarted by the system due to a failure
type Restarting struct{}

// A Stopping message is sent to an actor prior to the actor being stopped
type Stopping struct{}

// A Stopped message is sent to the actor once it has been stopped. A stopped actor will receive no further messages
type Stopped struct{}

// A Started message is sent to an actor once it has been started and ready to begin receiving messages.
type Started struct{}

// Restart is message sent by the actor system to control the lifecycle of an actor
type Restart struct{}

// Stop is message sent by the actor system to control the lifecycle of an actor
//
// This will not be forwarded to the Receive method
type Stop struct{}

// ResumeMailbox is message sent by the actor system to control the lifecycle of an actor.
//
// This will not be forwarded to the Receive method
type ResumeMailbox struct{}

// SuspendMailbox is message sent by the actor system to control the lifecycle of an actor.
//
// This will not be forwarded to the Receive method
type SuspendMailbox struct{}

//TODO: make private?
type Failure struct {
	Who        *PID
	Reason     interface{}
	ChildStats *ChildRestartStats
	Message    interface{}
}

func (*Restarting) AutoReceiveMessage() {}
func (*Stopping) AutoReceiveMessage()   {}
func (*Stopped) AutoReceiveMessage()    {}
func (*PoisonPill) AutoReceiveMessage() {}

func (*Started) SystemMessage()        {}
func (*Stop) SystemMessage()           {}
func (*Watch) SystemMessage()          {}
func (*Unwatch) SystemMessage()        {}
func (*Terminated) SystemMessage()     {}
func (*Failure) SystemMessage()        {}
func (*Restart) SystemMessage()        {}
func (*ResumeMailbox) SystemMessage()  {}
func (*SuspendMailbox) SystemMessage() {}

var (
	restartingMessage     interface{} = &Restarting{}
	stoppingMessage       interface{} = &Stopping{}
	stoppedMessage        interface{} = &Stopped{}
	poisonPillMessage     interface{} = &PoisonPill{}
	receiveTimeoutMessage interface{} = &ReceiveTimeout{}
)

var (
	restartMessage        SystemMessage = &Restart{}
	startedMessage        SystemMessage = &Started{}
	stopMessage           SystemMessage = &Stop{}
	resumeMailboxMessage  SystemMessage = &ResumeMailbox{}
	suspendMailboxMessage SystemMessage = &SuspendMailbox{}
)