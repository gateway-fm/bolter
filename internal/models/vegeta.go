package models

import (
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Vegeta struct {
	Rate     vegeta.Rate
	Targeter vegeta.Targeter
	Target   vegeta.Target
	Duration time.Duration
	Attacker *vegeta.Attacker
}
