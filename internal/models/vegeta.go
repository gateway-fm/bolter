package models

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

type Vegeta struct {
	Rate     vegeta.Rate
	Targeter vegeta.Targeter
	Target   vegeta.Target
	Duration time.Duration
	Attacker *vegeta.Attacker
}
