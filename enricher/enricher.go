package enricher

import (
	"encoding/json"
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

	cmdline = replaceNullCharacter(c)

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
	file := fmt.Sprintf("/host/proc/%s/exe", pidStr)

	exe, err = os.Readlink(file)
	return
}

func (e *Enricher) getContainerID(pid uint32) (containerID string, err error) {
	c, err := readFromProcFile(pid, "cpuset")

	if err != nil {
		return
	}

	// CPU set associate with root
	if c == "/" {
		containerID = types.Host
		return
	}

	// CPU set associate with container
	list := strings.Split(c, "/")
	containerID = list[len(list)-1]

	return
}

func (e *Enricher) getImage(containerID string) (containerInfo types.Container, err error) {
	return readFromDocker(containerID)
}

func (e *Enricher) getPWD(pid uint32) (pwd string, err error) {
	c, err := readFromProcFile(pid, "environ")

	if err != nil {
		return
	}

	envList := strings.Split(replaceNullCharacter(c), " ")

	for _, env := range envList {
		if strings.HasPrefix(env, "PWD=") {
			pwd = env[4:]
			return
		}
	}

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

			// skip host process or the process exited and no container info is available
			if containerID == types.Host || containerID == "" {
				continue
			}

			process, err := e.getCmdline(pid)
			ignoreError(err)

			exe, err := e.readExe(pid)
			ignoreError(err)

			//hostUID, containerUID, err := e.getUIDs(pid)
			//ignoreError(err)
			//
			//hostGID, containerGID, err := e.getUIDs(pid)
			//ignoreError(err)

			pwd, err := e.getPWD(pid)
			ignoreError(err)

			container, err := e.getImage(containerID)
			ignoreError(err)

			var eMsg types.EnrichedMessage
			eMsg.Timestamp = ts
			eMsg.PID = pid
			eMsg.ProcessName = process
			eMsg.ContainerID = containerID
			eMsg.ImageSHA = extraImageSHA(container.Image)
			eMsg.Image = container.Config.ImageName
			//eMsg.HostUID = hostUID
			//eMsg.ContainerUID = containerUID
			//eMsg.HostGID = hostGID
			//eMsg.ContainerGID = containerGID
			eMsg.Exe = exe
			eMsg.PWD = pwd

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
	file := fmt.Sprintf("/host/proc/%s/%s", pidStr, fileName)

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	c = strings.TrimSpace(string(content))

	return
}

func readFromDocker(containerID string) (c types.Container, err error) {
	file := fmt.Sprintf("/host/containers/%s/config.v2.json", containerID)

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	err = json.Unmarshal(content, &c)

	return
}

func replaceNullCharacter(c string) string {
	list := []byte(c)
	newList := []byte{}

	for _, b := range list {
		if b == 0 {
			newList = append(newList, 32)
		} else {
			newList = append(newList, b)
		}
	}

	return string(newList)
}

func extraImageSHA(input string) string {
	if strings.HasPrefix(input, "sha256:") {
		return input[7:]
	}

	return input
}
