package utils

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/manifoldco/promptui"

	"main/pkg/global"
)

func Cli() {
	prompt := promptui.Select{
		Label: "Select module",
		Items: []string{
			"Faucet test tokens",
			"Mint FUN Starts NFT",
			"Mint Megamafia NFT",
			"Mint Mega Cat NFT",
			"Mint Blackhole NFT",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Errorf("Error while CLI module selection: %s\n", err)
	}

	global.Module = result

	fmt.Printf("You've selected %v module\n", result)
}
