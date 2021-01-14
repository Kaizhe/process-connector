package enricher

import (
	"fmt"
	"github.com/kaizhe/proc-connector/pkg/types"
	"io/ioutil"
	"strconv"
	"strings"
)

type Enricher struct {

}

func NewEnricher() *Enricher {
	return &Enricher{}
}

func (e *Enricher) getCmdline(pid uint32) (processName string, err error) {
	cmdFile := fmt.Sprintf("/proc/%s/cmdline", strconv.FormatUint(uint64(pid), 10))

	content, err := ioutil.ReadFile(cmdFile)

	if err != nil {
		return
	}

	processName = strings.Join(strings.Split(string(content), "golang\000"), " ")

	return
}

func (e *Enricher) getContainerID(pid uint32) (containerID string, err error) {
	pidStr := strconv.FormatUint(uint64(pid), 10)
	file := fmt.Sprintf("/proc/%s/cpuset", pidStr)

	content, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	// CPU set associate with root
	c := strings.TrimSpace(string(content))
	fmt.Println("cpuset: ", c)
	if c == "/" {
		containerID = "host"
		return
	}

	// CPU set associate with container
	list := strings.Split(string(content), "/")
	containerID = list[len(list)-1]

	return
}


func (e *Enricher) getImage(containerID string) (imageName, imageSHA string, err error) {
	return
}

func (e *Enricher) Enrich(input <-chan *types.Message) error  {
	for {
		select {
		case msg := <- input:
			var err error
			pid := msg.PID
			ts := msg.Timestamp

			process, err := e.getCmdline(pid)

			if err != nil {
				fmt.Println(err)
				return err
			}

			containerID, err := e.getContainerID(pid)

			if err != nil {
				fmt.Println(err)
				return err
			}

			imageName, imageSHA, err := e.getImage(containerID)

			if err != nil {
				fmt.Println(err)
				return err
			}

			var eMsg types.EnrichedMessage
			eMsg.Timestamp = ts
			eMsg.PID = pid
			eMsg.ProcessName = process
			eMsg.ContainerID = containerID
			eMsg.ImageSHA = imageSHA
			eMsg.Image = imageName

			fmt.Println(eMsg)
		}
	}
}

