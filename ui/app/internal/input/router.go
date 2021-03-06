// SPDX-License-Identifier: Unlicense OR MIT

package input

import (
	"encoding/binary"
	"time"

	"gioui.org/ui"
	"gioui.org/ui/input"
	"gioui.org/ui/internal/opconst"
	"gioui.org/ui/internal/ops"
	"gioui.org/ui/key"
	"gioui.org/ui/pointer"
	"gioui.org/ui/system"
)

// Router is a Queue implementation that routes events from
// all available input sources to registered handlers.
type Router struct {
	pqueue pointerQueue
	kqueue keyQueue

	handlers handlerEvents

	reader ops.Reader

	// InvalidateOp summary.
	wakeup     bool
	wakeupTime time.Time

	// ProfileOp summary.
	profHandlers []input.Key
}

type handlerEvents struct {
	handlers map[input.Key][]input.Event
	updated  bool
}

func (q *Router) Next(k input.Key) (input.Event, bool) {
	return q.handlers.Next(k)
}

func (q *Router) Frame(ops *ui.Ops) {
	q.handlers.Clear()
	q.wakeup = false
	q.profHandlers = q.profHandlers[:0]
	q.reader.Reset(ops)
	q.collect()

	q.pqueue.Frame(ops, &q.handlers)
	q.kqueue.Frame(ops, &q.handlers)
	if q.handlers.Updated() {
		q.wakeup = true
		q.wakeupTime = time.Time{}
	}
}

func (q *Router) Add(e input.Event) bool {
	switch e := e.(type) {
	case pointer.Event:
		q.pqueue.Push(e, &q.handlers)
	case key.EditEvent, key.Event, key.FocusEvent:
		q.kqueue.Push(e, &q.handlers)
	}
	return q.handlers.Updated()
}

func (q *Router) TextInputState() TextInputState {
	return q.kqueue.InputState()
}

func (q *Router) collect() {
	for encOp, ok := q.reader.Decode(); ok; encOp, ok = q.reader.Decode() {
		switch opconst.OpType(encOp.Data[0]) {
		case opconst.TypeInvalidate:
			op := decodeInvalidateOp(encOp.Data)
			if !q.wakeup || op.At.Before(q.wakeupTime) {
				q.wakeup = true
				q.wakeupTime = op.At
			}
		case opconst.TypeProfile:
			op := decodeProfileOp(encOp.Data, encOp.Refs)
			q.profHandlers = append(q.profHandlers, op.Key)
		}
	}
}

func (q *Router) AddProfile(e system.ProfileEvent) {
	for _, h := range q.profHandlers {
		q.handlers.Add(h, e)
	}
}

func (q *Router) Profiling() bool {
	return len(q.profHandlers) > 0
}

func (q *Router) WakeupTime() (time.Time, bool) {
	return q.wakeupTime, q.wakeup
}

func (h *handlerEvents) init() {
	if h.handlers == nil {
		h.handlers = make(map[input.Key][]input.Event)
	}
}

func (h *handlerEvents) Set(k input.Key, evts []input.Event) {
	h.init()
	h.handlers[k] = evts
	h.updated = true
}

func (h *handlerEvents) Add(k input.Key, e input.Event) {
	h.init()
	h.handlers[k] = append(h.handlers[k], e)
	h.updated = true
}

func (h *handlerEvents) Updated() bool {
	u := h.updated
	h.updated = false
	return u
}

func (h *handlerEvents) Next(k input.Key) (input.Event, bool) {
	events := h.handlers[k]
	if len(events) == 0 {
		return nil, false
	}
	e := events[0]
	h.handlers[k] = events[1:]
	return e, true
}

func (h *handlerEvents) Clear() {
	for k := range h.handlers {
		delete(h.handlers, k)
	}
}

func decodeProfileOp(d []byte, refs []interface{}) system.ProfileOp {
	if opconst.OpType(d[0]) != opconst.TypeProfile {
		panic("invalid op")
	}
	return system.ProfileOp{
		Key: refs[0].(input.Key),
	}
}

func decodeInvalidateOp(d []byte) ui.InvalidateOp {
	bo := binary.LittleEndian
	if opconst.OpType(d[0]) != opconst.TypeInvalidate {
		panic("invalid op")
	}
	var o ui.InvalidateOp
	if nanos := bo.Uint64(d[1:]); nanos > 0 {
		o.At = time.Unix(0, int64(nanos))
	}
	return o
}
