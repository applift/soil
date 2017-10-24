// +build ide test_unit

package api_test

import (
	"context"
	"github.com/akaspin/logx"
	"github.com/akaspin/soil/agent/api"
	"github.com/akaspin/soil/agent/bus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStatusNodeProcessor_Process(t *testing.T) {

	e := api.NewStatusNodeGet(logx.GetLog("test")).Processor()
	go e.(bus.Consumer).ConsumeMessage(bus.NewMessage("a", map[string]string{
		"a-k1": "a-v1",
	}))
	go e.(bus.Consumer).ConsumeMessage(bus.NewMessage("b", map[string]string{
		"b-k1": "b-v1",
	}))
	time.Sleep(time.Millisecond * 200)

	res, err := e.Process(context.Background(), nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, res, map[string]map[string]string{
		"a": {"a-k1": "a-v1"},
		"b": {"b-k1": "b-v1"},
	})
}
