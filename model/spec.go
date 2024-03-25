package model

import "time"

type ScriptSpec struct {
	Organization int
	Assets       []int
	StartTs      time.Time
	EndTs        time.Time
	Options      map[string][]string
}
