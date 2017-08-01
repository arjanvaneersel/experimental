package model

// A mapping of license identifiers to descriptions
var Licenses map[string]string

func init() {
	Licenses = make(map[string]string)
}