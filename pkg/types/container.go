package types

type Container struct {
	ContainerID string          `json:"ID"`
	Created     string          `json:"Created"`
	Image       string          `json:"Image"`
	Config      ContainerConfig `json:"Config"`
}

type ContainerConfig struct {
	ImageName string `json:"Image"`
}
