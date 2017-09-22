package api_v1_types

type NodeResponse struct {
	Id string
	Advertise string
	Drain string
	Version string
	API string
}

type NodesResponse []NodeResponse

