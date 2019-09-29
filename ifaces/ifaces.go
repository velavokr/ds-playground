package ifaces

// Node env ///////////////////////////////////////

type NodeEnv interface {
	Net(handler NetHandler) Net
	Timer(handler TimerHandler) Timer
	Storage() Storage
	PKI() PKI
}

// Network events /////////////////////////////////

// NodeName is an abstraction for a node address
type NodeName = string

type Group struct {
	// Nodes are all nodes in the current group, including current node
	Nodes []NodeName
	// Self is guaranteed to be the current node's position in Nodes
	Self int
}

// Net is the ReceiveMessage counterpart.
type Net interface {
	SendMessage(dst NodeName, message []byte)
}

// NetHandler is the interface to be implemented in the assignments.
// Usually the assignments will provide helper objects.
type NetHandler interface {
	// ReceiveMessage is called by the link when a message arrives.
	// Usually the network messages are to be designed by the students.
	ReceiveMessage(src NodeName, message []byte)
}

// Time events ////////////////////////////////////

type TimerId = int

type Timer interface {
	NextTick(ctx interface{}) TimerId
	After(ticks uint32, ctx interface{}) TimerId
	CancelTimer(id TimerId)
}

type TimerHandler interface {
	HandleTimer(ctx interface{}, id TimerId)
}

// Crash-Restore //////////////////////////////////

type Storage interface {
	OpenTable(name string) DiskTable
}

type DiskTable interface {
	StoreValue(key []byte, val []byte)
	LoadValue(key []byte) []byte
	DeleteKey(key []byte)
	LoadKeys() [][]byte
}

// Signature //////////////////////////////////////

type Signature = string

type PKI interface {
	SignMessage(dst NodeName, msg []byte) Signature
	CheckSignature(mac Signature, src NodeName, dst NodeName, msg []byte) bool
}
