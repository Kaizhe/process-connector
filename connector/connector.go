package connector

// #include <linux/connector.h>
// #include <linux/cn_proc.h>
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/kaizhe/proc-connector/pkg/types"
	"golang.org/x/sys/unix"
	"os"
)

type ProcessConnector struct {

}


func NewProcessConnector() *ProcessConnector {
	return &ProcessConnector{}
}

func (pc *ProcessConnector) Listen() error {
	sock, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_DGRAM, unix.NETLINK_CONNECTOR)

	if err != nil {
		return err
	}

	addr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK, Groups: C.CN_IDX_PROC, Pid: uint32(os.Getpid())}
	err = unix.Bind(sock, addr)

	if err != nil {
		return err
	}

	defer func() {
		send(sock, C.PROC_CN_MCAST_IGNORE)
		unix.Close(sock)
	}()

	err = send(sock, C.PROC_CN_MCAST_LISTEN)

	if err != nil {
		return nil
	}

	for {
		p := make([]byte, 1024)

		nlmessages, err := recvData(p, sock)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		for _, m := range nlmessages {
			parseNetlinkMessage(m)
		}
	}

	return nil
}

func parseNetlinkMessage(m syscall.NetlinkMessage) {
	if m.Header.Type == unix.NLMSG_DONE {
		buf := bytes.NewBuffer(m.Data)
		msg := &types.CnMsg{}
		hdr := &types.ProcEventHeader{}
		binary.Read(buf, binary.LittleEndian, msg)
		binary.Read(buf, binary.LittleEndian, hdr)

		if hdr.What == C.PROC_EVENT_EXIT {
			event := &types.ProcExitEvent{}
			binary.Read(buf, binary.LittleEndian, event)
			pid := int(event.ProcessTgid)
			fmt.Printf("%d just exited.\n", pid)
		} else if hdr.What == C.PROC_EVENT_EXEC {
			event := &types.ProcExecEvent{}
			binary.Read(buf, binary.LittleEndian, event)
			pid := int(event.ProcessTgid)
			fmt.Printf("%d just started.\n", pid)
		}
	}
}

func send(sock int, msg uint32) error {
	cnMsg := types.CnMsg{}
	destAddr := &unix.SockaddrNetlink{Family: unix.AF_NETLINK, Groups: C.CN_IDX_PROC, Pid: 0} // the kernel
	header := unix.NlMsghdr{
		Len:   unix.NLMSG_HDRLEN + uint32(binary.Size(cnMsg)+binary.Size(msg)),
		Type:  uint16(unix.NLMSG_DONE),
		Flags: 0,
		Seq:   1,
		Pid:   uint32(os.Getpid()),
	}
	cnMsg.ID = types.CbID{Idx: C.CN_IDX_PROC, Val: C.CN_VAL_PROC}
	cnMsg.Len = uint16(binary.Size(msg))
	cnMsg.Ack = 0
	cnMsg.Seq = 1
	buf := bytes.NewBuffer(make([]byte, 0, header.Len))
	binary.Write(buf, binary.LittleEndian, header)
	binary.Write(buf, binary.LittleEndian, cnMsg)
	binary.Write(buf, binary.LittleEndian, msg)

	return unix.Sendto(sock, buf.Bytes(), 0, destAddr)
}

func recvData(p []byte, sock int) ([]syscall.NetlinkMessage, error) {
	nr, from, err := unix.Recvfrom(sock, p, 0)

	if sockaddrNl, ok := from.(*unix.SockaddrNetlink); !ok || sockaddrNl.Pid != 0 {
		return nil, fmt.Errorf("sender was not kernel")
	}

	if err != nil {
		return nil, err
	}

	if nr < unix.NLMSG_HDRLEN {
		return nil, fmt.Errorf("received %d bytes", nr)
	}

	nlmessages, err := syscall.ParseNetlinkMessage(p[:nr])

	if err != nil {
		return nil, err
	}

	return nlmessages, nil
}