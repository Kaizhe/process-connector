package types

const (
	Host = "host"
)

type Message struct {
	PID       uint32
	Timestamp uint64
}

func NewMessage(pid uint32, ts uint64) *Message {
	return &Message{
		PID:       pid,
		Timestamp: ts,
	}
}

func (m *Message) IsEmpty() bool {
	return m.Timestamp == 0
}

type EnrichedMessage struct {
	Message
	ProcessName  string
	Image        string
	ImageSHA     string
	ContainerID  string
	HostUID      string
	ContainerUID string
	HostGID      string
	ContainerGID string
	Exe          string
}
