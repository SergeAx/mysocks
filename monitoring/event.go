package monitoring

import "time"

type eventmetric struct {
	ev  metric
	sim *counter
}

func NewEventMetric(name string, tags map[string]string, fields ...string) eventmetric {
	return eventmetric{
		ev:  NewMetric(name, tags, append(fields, "duration")...),
		sim: NewCounter(name+"_sim", tags),
	}
}

func (e *eventmetric) Start() time.Time {
	e.sim.Increase()
	return time.Now()
}

func (e *eventmetric) Stop(t time.Time, values ...interface{}) {
	e.sim.Decrease()
	e.ev.record(append(values, time.Since(t)), t)
}
