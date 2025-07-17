package utils

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"main/pkg/global"
)

// Sleep func make a random delay in range [delayMin, delayMax]
func Sleep(delayMin int, delayMax int) {
	delay := rand.Intn(delayMax + 1 - delayMin) + delayMin

	log.Infof("[%d/%d] | [checkBalance] | Sleep %d seconds ...\n",
		global.CurrentProgress, global.TargetProgress, delay,
	)

	time.Sleep(time.Duration(delay) * time.Second)
}
