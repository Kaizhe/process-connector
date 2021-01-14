package enricher

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/kaizhe/proc-connector/pkg/types"
)

type Enricher struct {
}

func NewEnricher() *Enricher {
	return &Enricher{}
}

func (e *Enricher) getCmdline(pid uint32) (cmdline string, err error) {
	cmdFile := fmt.Sprintf("/proc/%s/cmdline", strconv.FormatUint(uint64(pid), 10))

	content, err := ioutil.ReadFile(cmdFile)

	if err != nil {
		return
	}

	cmdline = strings.TrimSpace(strings.Join(strings.Split(string(content), "golang\000"), " "))

	return
}

func (e *Enricher) getUIDs(pid uint32) (hostUID, containerUID string, err error) {
	return e.getUserNamespaceInfo(pid, "uid_map")
}

func (e *Enricher) getGIDs(pid uint32) (hostGID, containerGID string, err error) {
	return e.getUserNamespaceInfo(pid, "gid_map")
}

func (e *Enricher) getUserNamespaceInfo(pid uint32, mapFile string) (host, container string, err error) {
	pidStr := strconv.FormatUint(uint64(pid), 10)
	file := fmt.Sprintf("/proc/%s/%s", pidStr, mapFile)

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	list := strings.Fields(strings.TrimSpace(string(content)))

	host = list[0]
	container = list[1]

	return
}

func (e *Enricher) getContainerID(pid uint32) (containerID string, err error) {
	pidStr := strconv.FormatUint(uint64(pid), 10)
	file := fmt.Sprintf("/proc/%s/cpuset", pidStr)

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	c := strings.TrimSpace(string(content))

	// CPU set associate with root
	fmt.Println("cpuset: ", c)
	if c == "/" {
		containerID = types.Host
		return
	}

	// CPU set associate with container
	list := strings.Split(c, "/")
	containerID = list[len(list)-1]

	return
}

func (e *Enricher) getImage(containerID string) (imageName, imageSHA string, err error) {
	return
}

func (e *Enricher) Enrich(input <-chan *types.Message) {
	for {
		select {
		case msg := <-input:
			var err error
			pid := msg.PID
			ts := msg.Timestamp

			containerID, err := e.getContainerID(pid)

			if err != nil {
				fmt.Println(err)
			}

			// skip host process
			if containerID == types.Host {
				continue
			}

			process, err := e.getCmdline(pid)

			if err != nil {
				fmt.Println(err)
			}

			imageName, imageSHA, err := e.getImage(containerID)

			if err != nil {
				fmt.Println(err)
			}

			hostUID, containerUID, err := e.getUIDs(pid)

			if err != nil {
				return
			}

			hostGID, containerGID, err := e.getUIDs(pid)

			if err != nil {
				return
			}

			var eMsg types.EnrichedMessage
			eMsg.Timestamp = ts
			eMsg.PID = pid
			eMsg.ProcessName = process
			eMsg.ContainerID = containerID
			eMsg.ImageSHA = imageSHA
			eMsg.Image = imageName
			eMsg.HostUID = hostUID
			eMsg.ContainerUID = containerUID
			eMsg.HostGID = hostGID
			eMsg.ContainerGID = containerGID

			fmt.Println(eMsg)
		}
	}
}
