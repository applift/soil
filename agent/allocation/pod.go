package allocation

import (
	"fmt"
	"github.com/akaspin/soil/manifest"
	"github.com/mitchellh/hashstructure"
	"strings"
)

const (
	podUnitTemplate = `
[Unit]
Description=${pod.name}
Before=${pod.units}
[Service]
${system.pod_exec}
[Install]
WantedBy=${pod.target}
`
	dirSystemDLocal   = "/etc/systemd/system"
	dirSystemDRuntime = "/run/systemd/system"
)

// Pod represents pod allocated on agent
type Pod struct {
	Header
	UnitFile
	Units     []*Unit
	Blobs     []*Blob
	Resources []*Resource
}

func NewPod(systemPaths SystemPaths) (p *Pod) {
	p = &Pod{
		UnitFile: UnitFile{
			SystemPaths: systemPaths,
		},
		Header: Header{},
	}
	return
}

func (p *Pod) FromManifest(m *manifest.Pod, env map[string]string) (err error) {
	agentMark, _ := hashstructure.Hash(env, nil)
	p.Header = Header{
		Name:      m.Name,
		PodMark:   m.Mark(),
		AgentMark: agentMark,
		Namespace: m.Namespace,
	}
	p.UnitFile = NewUnitFile(fmt.Sprintf("pod-%s-%s.service", m.Namespace, m.Name), p.SystemPaths, m.Runtime)
	baseEnv := map[string]string{
		"pod.name":      m.Name,
		"pod.namespace": m.Namespace,
	}
	baseSourceEnv := map[string]string{
		"pod.target": m.Target,
	}

	// Blobs
	fileHashes := map[string]string{}
	for _, b := range m.Blobs {
		ab := &Blob{
			Name:        manifest.Interpolate(b.Name, baseEnv),
			Permissions: b.Permissions,
			Leave:       b.Leave,
			Source:      manifest.Interpolate(b.Source, baseEnv, baseSourceEnv, env),
		}
		p.Blobs = append(p.Blobs, ab)
		fileHash, _ := hashstructure.Hash(ab.Source, nil)
		fileHashes[fmt.Sprintf("blob.%s", strings.Replace(strings.Trim(ab.Name, "/"), "/", "-", -1))] = fmt.Sprintf("%d", fileHash)
	}

	// Units
	var unitNames []string
	for _, u := range m.Units {
		unitName := manifest.Interpolate(u.Name, baseEnv)
		pu := &Unit{
			Transition: u.Transition,
			UnitFile:   NewUnitFile(unitName, p.SystemPaths, m.Runtime),
		}
		pu.Source = manifest.Interpolate(u.Source, baseEnv, baseSourceEnv, fileHashes, env)
		p.Units = append(p.Units, pu)
		unitNames = append(unitNames, unitName)
	}

	// Resources
	for _, resource := range m.Resources {
		p.Resources = append(p.Resources, newResource(p.Name, resource, env))
	}

	p.Source, err = p.Header.Marshal(p.Name, p.Units, p.Blobs, p.Resources)
	p.Source += manifest.Interpolate(podUnitTemplate, baseEnv, baseSourceEnv, map[string]string{
		"pod.units": strings.Join(unitNames, " "),
	}, env)
	return
}

func (p *Pod) FromFilesystem(path string) (err error) {
	p.UnitFile.Path = path
	if err = p.UnitFile.Read(); err != nil {
		return
	}
	if p.Units, p.Blobs, p.Resources, err = p.Header.Unmarshal(p.UnitFile.Source, p.SystemPaths); err != nil {
		return
	}

	for _, u := range p.Units {
		if err = u.UnitFile.Read(); err != nil {
			return
		}
	}
	for _, b := range p.Blobs {
		if err = b.Read(); err != nil {
			return
		}
	}
	return
}

func (p *Pod) GetPodUnit() (res *Unit) {
	res = &Unit{
		UnitFile: p.UnitFile,
		Transition: manifest.Transition{
			Create:    "start",
			Update:    "restart",
			Destroy:   "stop",
			Permanent: true,
		},
	}
	return
}
