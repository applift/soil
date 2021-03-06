// +build ide test_systemd

package provision_test

import (
	"context"
	"github.com/akaspin/logx"
	"github.com/akaspin/soil/agent/allocation"
	"github.com/akaspin/soil/agent/bus"
	"github.com/akaspin/soil/agent/provision"
	"github.com/akaspin/soil/fixture"
	"github.com/akaspin/soil/lib"
	"github.com/akaspin/soil/manifest"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEvaluator_Allocate(t *testing.T) {
	sd := fixture.NewSystemd("/run/systemd/system", "pod-private")
	sd.Cleanup()
	defer sd.Cleanup()

	ctx := context.Background()

	var state allocation.Recovery
	assert.NoError(t, state.FromFilesystem(allocation.DefaultSystemPaths(), allocation.DefaultDbusDiscoveryFunc))

	evaluator := provision.NewEvaluator(ctx, logx.GetLog("test"), provision.EvaluatorConfig{
		SystemPaths:    allocation.DefaultSystemPaths(),
		Recovery:       state,
		StatusConsumer: &bus.BlackholePipe{},
	})
	assert.NoError(t, evaluator.Open())

	waitConfig := fixture.DefaultWaitConfig()

	time.Sleep(time.Millisecond * 500)
	t.Run("0 create pod-1", func(t *testing.T) {
		var buffers lib.StaticBuffers
		var registry manifest.Registry
		assert.NoError(t, buffers.ReadFiles("testdata/evaluator_test_Allocate_0.hcl"))
		assert.NoError(t, registry.Unmarshal("private", buffers.GetReaders()...))

		evaluator.Allocate(registry[0], map[string]string{
			"system.pod_exec": "ExecStart=/usr/bin/sleep inf",
		})

		fixture.WaitNoError(t, waitConfig, sd.UnitStatesFn(
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]string{
				"pod-private-pod-1.service": "active",
				"unit-1.service":            "active",
			}))
		sd.AssertUnitHashes(t,
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]uint64{
				"/run/systemd/system/pod-private-pod-1.service": 0xc43253a8821be2b,
				"/run/systemd/system/unit-1.service":            0xbca69ea672e79d81,
			},
		)
	})
	t.Run("1 update pod-1", func(t *testing.T) {
		var buffers lib.StaticBuffers
		var registry manifest.Registry
		assert.NoError(t, buffers.ReadFiles("testdata/evaluator_test_Allocate_1.hcl"))
		assert.NoError(t, registry.Unmarshal("private", buffers.GetReaders()...))
		evaluator.Allocate(registry[0], map[string]string{
			"system.pod_exec": "ExecStart=/usr/bin/sleep inf",
		})

		fixture.WaitNoError(t, waitConfig, sd.UnitStatesFn(
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]string{
				"pod-private-pod-1.service": "active",
				"unit-1.service":            "active",
			}))
		sd.AssertUnitHashes(t,
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]uint64{
				"/run/systemd/system/pod-private-pod-1.service": 0x28525a605380724b,
				"/run/systemd/system/unit-1.service":            0x448529ac4d4389a0,
			},
		)
	})
	t.Run("2 destroy non-existent", func(t *testing.T) {
		evaluator.Deallocate("pod-2")

		fixture.WaitNoError(t, waitConfig, sd.UnitStatesFn(
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]string{
				"pod-private-pod-1.service": "active",
				"unit-1.service":            "active",
			}))
		sd.AssertUnitHashes(t,
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]uint64{
				"/run/systemd/system/pod-private-pod-1.service": 0x28525a605380724b,
				"/run/systemd/system/unit-1.service":            0x448529ac4d4389a0,
			},
		)
	})
	t.Run("3 destroy pod-1", func(t *testing.T) {
		evaluator.Deallocate("pod-1")

		fixture.WaitNoError(t, waitConfig, sd.UnitStatesFn(
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]string{}))
		sd.AssertUnitHashes(t,
			[]string{"pod-private-pod-1.service", "unit-1.service"},
			map[string]uint64{},
		)
	})

	assert.NoError(t, evaluator.Close())
	assert.NoError(t, evaluator.Wait())
}

func TestEvaluator_Report(t *testing.T) {

	sd := fixture.NewSystemd("/run/systemd/system", "pod-private")
	sd.Cleanup()
	sd.DeployPod("test-1", 1)
	defer sd.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stat := bus.NewTestingConsumer(ctx)

	var state allocation.Recovery
	assert.NoError(t, state.FromFilesystem(allocation.DefaultSystemPaths(), allocation.DefaultDbusDiscoveryFunc))

	evaluator := provision.NewEvaluator(ctx, logx.GetLog("test"), provision.EvaluatorConfig{
		SystemPaths:    allocation.DefaultSystemPaths(),
		Recovery:       state,
		StatusConsumer: stat,
	})
	assert.NoError(t, evaluator.Open())

	t.Run(`ensure recovered`, func(t *testing.T) {
		fixture.WaitNoError10(t, stat.ExpectMessagesFn(
			bus.NewMessage("", map[string]map[string]string{
				"test-1": {
					"present": "true",
					"state":   "dirty",
				},
			}),
		))
	})
	t.Run("deallocate test-1", func(t *testing.T) {
		evaluator.Deallocate("test-1")
		fixture.WaitNoError10(t, stat.ExpectMessagesFn(
			// reset
			bus.NewMessage("", map[string]map[string]string{
				"test-1": {
					"present": "true",
					"state":   "dirty",
				},
			}),
			bus.NewMessage("test-1", map[string]string{
				"present": "true",
				"state":   "destroy",
			}),
			bus.NewMessage("test-1", nil),
		))
	})
	t.Run("1 create pod-1", func(t *testing.T) {
		var buffers lib.StaticBuffers
		var registry manifest.Registry
		assert.NoError(t, buffers.ReadFiles("testdata/evaluator_test_Report_1.hcl"))
		assert.NoError(t, registry.Unmarshal("private", buffers.GetReaders()...))

		evaluator.Allocate(registry[0], map[string]string{
			"system.pod_exec": "ExecStart=/usr/bin/sleep inf",
		})

		fixture.WaitNoError10(t, stat.ExpectMessagesFn(
			// reset
			bus.NewMessage("", map[string]map[string]string{
				"test-1": {
					"present": "true",
					"state":   "dirty",
				},
			}),
			bus.NewMessage("test-1", map[string]string{
				"present": "true",
				"state":   "destroy",
			}),
			bus.NewMessage("test-1", nil),
			bus.NewMessage("pod-1", map[string]string{
				"present": "true",
				"state":   "create",
			}),
			bus.NewMessage("pod-1", map[string]string{
				"present": "true",
				"state":   "done",
			}),
		))

	})

	assert.NoError(t, evaluator.Close())
	assert.NoError(t, evaluator.Wait())
}
