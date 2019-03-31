// SPDX-License-Identifier: Unlicense OR MIT

package pointer

import (
	"time"

	"gioui.org/ui/f32"
)

type Event struct {
	Type      Type
	Source    Source
	PointerID ID
	Priority  Priority
	Time      time.Duration
	Hit       bool
	Position  f32.Point
	Scroll    f32.Point
}

type OpHandler struct {
	Key  Key
	Area Area
	Grab bool
}

type Area func(pos f32.Point) HitResult

type Key interface{}

type Events interface {
	For(k Key) []Event
}

type HitResult uint8

const (
	HitNone HitResult = iota
	HitTransparent
	HitOpaque
)

type ID uint16
type Type uint8
type Priority uint8
type Source uint8

const (
	Cancel Type = iota
	Press
	Release
	Move
)

const (
	Mouse Source = iota
	Touch
)

const (
	Shared Priority = iota
	Foremost
	Grabbed
)

func (t Type) String() string {
	switch t {
	case Press:
		return "Press"
	case Release:
		return "Release"
	case Cancel:
		return "Cancel"
	case Move:
		return "Move"
	default:
		panic("unknown Type")
	}
}

func (p Priority) String() string {
	switch p {
	case Shared:
		return "Shared"
	case Foremost:
		return "Foremost"
	case Grabbed:
		return "Grabbed"
	default:
		panic("unknown priority")
	}
}

func (s Source) String() string {
	switch s {
	case Mouse:
		return "Mouse"
	case Touch:
		return "Touch"
	default:
		panic("unknown source")
	}
}

func (OpHandler) ImplementsOp() {}
func (Event) ImplementsEvent()  {}