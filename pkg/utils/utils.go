package utils

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	log "github.com/sirupsen/logrus"

	"main/pkg/global"
)

// Sleep func make a random delay in range [delayMin, delayMax]
func Sleep(delayMin int, delayMax int) {
	delay := rand.Intn(delayMax + 1 - delayMin) + delayMin

	log.Infof("[%d/%d] | [Sleep] | Sleep %d seconds ...\n",
		global.CurrentProgress, global.TargetProgress, delay,
	)

	time.Sleep(time.Duration(delay) * time.Second)
}

// LoadABI loads ABI from defined filepath
func LoadABI(filepath string) (abi.ABI, error) {
	abiFile, err := os.ReadFile(filepath)

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | [LoadABI] | Problem with reading %v file with ABI\n",
			global.CurrentProgress, global.TargetProgress, filepath,
		)
		return abi.ABI{}, errors.New(msg)
	}

	contractABI, err := abi.JSON(bytes.NewReader(abiFile))

	if err != nil {
		msg := fmt.Sprintf("[%d/%d] | [LoadABI] | Problem with parsing ABI\n",
			global.CurrentProgress, global.TargetProgress,
		)
		return abi.ABI{}, errors.New(msg)
	}

	return contractABI, nil
}
