package chord

import (
	"fmt"
	"time"
)

type Tracer struct {
	tracerStart    time.Time
	tracerDuration time.Duration
	startNodeID    []byte
	resultNodeID   []byte
	keyID          []byte
	numHops        int    //hop count
	hops           string //hopped node ID
	traces         []string
}

func MakeTracer() *Tracer {
	tracer := Tracer{}
	tracer.numHops = 0
	return &tracer
}

//reset tracer
func (tracer *Tracer) resetTracer() {

	tracer.numHops = 0
	tracer.hops = ""
	tracer.traces = nil
}

func (tracer *Tracer) startTracer(startNodeID []byte, keyID []byte) {
	tracer.resetTracer()
	tracer.tracerStart = time.Now()
	tracer.startNodeID = startNodeID
	tracer.keyID = keyID
}

func (tracer *Tracer) endTracer(resultNodeID []byte) {
	tracer.tracerDuration = time.Since(tracer.tracerStart)
	tracer.resultNodeID = resultNodeID
	tracer.writeTrace()
}

func (tracer *Tracer) incrementHops(numHop int) {
	tracer.numHops += numHop
}

func (tracer *Tracer) traceAction(action string) {
	tracer.hops += fmt.Sprintf("\n%s: ", action)
}

func (tracer *Tracer) traceNode(nodeID []byte) {
	tracer.hops += fmt.Sprintf("%d, ", nodeID)
	tracer.incrementHops(1)
}

func (tracer *Tracer) writeTrace() {
	startNodeID := fmt.Sprintf("Caller Node: %d", tracer.startNodeID)
	resultNodeID := fmt.Sprintf("Result Node: %d", tracer.resultNodeID)
	keyLookup := fmt.Sprintf("Key lookup: %d", tracer.keyID)
	numHops := fmt.Sprintf("Number of hops: %d", tracer.numHops)
	hopTrace := fmt.Sprintf("Trace: %s", tracer.hops)
	latency := fmt.Sprintf("Latency: %s", tracer.tracerDuration)
	currentTrace := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", keyLookup, startNodeID, hopTrace, resultNodeID, numHops, latency)
	tracer.traces = append(tracer.traces, currentTrace)
}

func (tracer *Tracer) String() string {
	traces := ""
	for i := 0; i < len(tracer.traces); i++ {
		traces += fmt.Sprintf("%s\n", tracer.traces[i])
	}

	return traces
}

func (tracer *Tracer) Hops() string {
	return fmt.Sprintf("hops: %d", tracer.numHops)
}

func (tracer *Tracer) Latency() string {
	nano := tracer.tracerDuration / time.Nanosecond
	mili := float64(nano)
	mili = mili / 1000000
	return fmt.Sprintf("latency: %v", mili)
}
