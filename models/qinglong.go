package models

type ql struct {
	Command []string
	Admin   bool
	Handle  func(sender *Sender) interface{}
}
