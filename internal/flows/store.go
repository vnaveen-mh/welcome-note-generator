package flows

import "sync"

type FlowStore struct {
	mu    sync.Mutex
	flows map[string]any
}

var flowStore FlowStore

func init() {
	if flowStore.flows == nil {
		flowStore.flows = map[string]any{}
	}
}

func GetFlow(name string) (any, bool) {
	flowStore.mu.Lock()
	defer flowStore.mu.Unlock()

	val, ok := flowStore.flows[name]
	return val, ok
}

func SetFlow(name string, val any) {
	flowStore.mu.Lock()
	defer flowStore.mu.Unlock()

	flowStore.flows[name] = val
}
