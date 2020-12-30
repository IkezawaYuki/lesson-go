package main

type State int

const (
	OffHook State = iota
	Connecting
	Connected
	OnHold
	OnHook
)

func (s State) String() string {
	switch s {
	case OffHook:
		return "OffHook"
	case Connecting:
		return "Connecting"
	case Connected:
		return "Connected"
	case OnHold:
		return "OnHold"
	case OnHook:
		return "OnHook"
	}
	return "Unknown"
}

type Trigger int

const (
	CallDialed Trigger = iota
	HungUp
	CallConnected
	PlacedOnHold
	TakenOffHold
	LeftMessage
)

func (t Trigger) String() string {
	switch t {
	case CallDialed:
		return "CallDialed"
	case HungUp:
		return "HungUp"
	case CallConnected:
		return "CallConnected"
	case PlacedOnHold:
		return "PlacedOnHold"
	case TakenOffHold:
		return "TakenOffHold"
	case LeftMessage:
		return "LeftMessage"
	}
	return "Unknown"
}

type TriggerResult struct {
	Trigger Trigger
	State   State
}

var rules = map[State][]TriggerResult{
	OffHook: {
		{CallDialed, Connecting},
	},
	Connecting: {
		{HungUp, OffHook},
		{CallConnected, Connected},
	},
	Connected: {
		{LeftMessage, OnHook},
		{HungUp, OffHook},
		{PlacedOnHold, OnHold},
	},
	OnHold: {
		{TakenOffHold, Connected},
		{HungUp, OnHook},
	},
}

func main() {

}
