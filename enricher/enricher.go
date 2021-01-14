package enricher

import (
	"fmt"
	"io/ioutil"
	"os"
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
	c, err := readFromProcFile(pid, "cmdline")

	if err != nil {
		return
	}

	cmdline = strings.TrimSpace(strings.Join(strings.Split(c, "golang\000"), " "))

	return
}

func (e *Enricher) getUIDs(pid uint32) (hostUID, containerUID string, err error) {
	return e.getUserNamespaceInfo(pid, "uid_map")
}

func (e *Enricher) getGIDs(pid uint32) (hostGID, containerGID string, err error) {
	return e.getUserNamespaceInfo(pid, "gid_map")
}

func (e *Enricher) getUserNamespaceInfo(pid uint32, mapFile string) (host, container string, err error) {
	c, err := readFromProcFile(pid, mapFile)

	if err != nil {
		return
	}

	list := strings.Fields(c)

	host = list[0]
	container = list[1]

	return
}

func (e *Enricher) readExe(pid uint32) (exe string, err error) {
	pidStr := strconv.FormatUint(uint64(pid), 10)
	file := fmt.Sprintf("/proc/%s/exe", pidStr)

	exe, err = os.Readlink(file)
	return
}

func (e *Enricher) getContainerID(pid uint32) (containerID string, err error) {
	c, err := readFromProcFile(pid, "cpuset")

	if err != nil {
		return
	}

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
			ignoreError(err)

			// skip host process
			if containerID == types.Host {
				continue
			}

			process, err := e.getCmdline(pid)
			ignoreError(err)

			imageName, imageSHA, err := e.getImage(containerID)
			ignoreError(err)

			hostUID, containerUID, err := e.getUIDs(pid)
			ignoreError(err)

			hostGID, containerGID, err := e.getUIDs(pid)
			ignoreError(err)

			exe, err := e.readExe(pid)
			ignoreError(err)

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
			eMsg.Exe = exe

			fmt.Printf("%+v\n", eMsg)
		}
	}
}

func ignoreError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func readFromProcFile(pid uint32, fileName string) (c string, err error) {
	pidStr := strconv.FormatUint(uint64(pid), 10)
	file := fmt.Sprintf("/proc/%s/cpuset", pidStr)

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	c = strings.TrimSpace(string(content))

	return
}
