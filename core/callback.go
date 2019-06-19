package core

// Callback defines a pre/post Set callback function.
type Callback func(flag Flag, value string) error
