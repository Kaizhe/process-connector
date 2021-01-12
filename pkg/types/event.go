package types

// CbID corresponds to cb_id in connector.h
type CbID struct {
	Idx uint32
	Val uint32
}

// CnMsg corresponds to cn_msg in connector.h
type CnMsg struct {
	ID    CbID
	Seq   uint32
	Ack   uint32
	Len   uint16
	Flags uint16
}

// ProcEventHeader corresponds to proc_event in cn_proc.h
type ProcEventHeader struct {
	What      uint32
	CPU       uint32
	Timestamp uint64
}

// ProcExitEvent corresponds to exit_proc_event in cn_proc.h
type ProcExitEvent struct {
	ProcessPid  uint32
	ProcessTgid uint32
	ExitCode    uint32
	ExitSignal  uint32
}

// ProcExecEvent corresponds to exit_proc_event in cn_proc.h
type ProcExecEvent struct {
	ProcessPid  uint32
	ProcessTgid uint32
}
