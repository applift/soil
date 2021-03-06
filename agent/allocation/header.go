package allocation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akaspin/soil/manifest"
	"github.com/mitchellh/hashstructure"
	"strings"
)

type Header struct {
	Name      string
	PodMark   uint64
	AgentMark uint64
	Namespace string
}

func (h *Header) Mark() (res uint64) {
	res, _ = hashstructure.Hash(h, nil)
	return
}

func (h *Header) Unmarshal(src string, paths SystemPaths) (units []*Unit, blobs []*Blob, resources []*Resource, err error) {
	split := strings.Split(src, "\n")
	// extract header
	var jsonSrc string
	if _, err = fmt.Sscanf(split[0], "### POD %s %s", &h.Name, &jsonSrc); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(jsonSrc), &h); err != nil {
		return
	}
	for _, line := range split[1:] {
		if strings.HasPrefix(line, "### UNIT") {
			u := &Unit{
				UnitFile: UnitFile{
					SystemPaths: paths,
				},
				Transition: manifest.Transition{},
			}
			if _, err = fmt.Sscanf(line, "### UNIT %s %s", &u.UnitFile.Path, &jsonSrc); err != nil {
				return
			}
			if err = json.Unmarshal([]byte(jsonSrc), &u); err != nil {
				return
			}
			units = append(units, u)
		}
		if strings.HasPrefix(line, "### BLOB") {
			b := &Blob{}
			if _, err = fmt.Sscanf(line, "### BLOB %s %s", &b.Name, &jsonSrc); err != nil {
				return
			}
			if err = json.Unmarshal([]byte(jsonSrc), &b); err != nil {
				return
			}
			blobs = append(blobs, b)
		}
		if strings.HasPrefix(line, resourceHeaderPrefix) {
			resource := defaultResource()
			if err = resource.unmarshalHeader(line); err != nil {
				return
			}
			resources = append(resources, resource)
		}

	}
	return
}

func (h *Header) Marshal(name string, units []*Unit, blobs []*Blob, resources []*Resource) (res string, err error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)

	if _, err = fmt.Fprintf(buf, "### POD %s ", name); err != nil {
		return
	}
	if err = encoder.Encode(map[string]interface{}{
		"PodMark":   h.PodMark,
		"AgentMark": h.AgentMark,
		"Namespace": h.Namespace,
	}); err != nil {
		return
	}
	for _, u := range units {
		if err = u.MarshalHeader(buf, encoder); err != nil {
			return
		}
	}

	for _, b := range blobs {
		if err = b.MarshalHeader(buf, encoder); err != nil {
			return
		}
	}
	for _, resource := range resources {
		if err = resource.marshalHeader(buf, encoder); err != nil {
			return
		}
	}
	res = string(buf.Bytes())
	return
}
