package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"main/internal/megaeth"
	"main/pkg/global"
	"main/pkg/types"
	"main/pkg/utils"
)

func inputUser(inputText string) string {
	if inputText != "" {
		fmt.Print(inputText)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func handlePanic() {
	if r := recover(); r != nil {
		log.Printf("Unexpected Error: %v", r)
		fmt.Println("Press Enter to Exit..")
		_, err := fmt.Scanln()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(1)
	}
}

func processAccounts(func_obj types.ModuleFunction, threads int) {
	if threads == 1 {
		for _, account := range global.AccountsList {
			func_obj(account)

			global.CurrentProgress++

			utils.Sleep(
				global.Config.DelayBetweenAccs.Min,
				global.Config.DelayBetweenAccs.Max,
			)
		}

	} else {
		var wg sync.WaitGroup
		sem := make(chan struct{}, threads)
	
		for _, account := range global.AccountsList {
			wg.Add(1)
			sem <- struct{}{}
	
			go func(acc types.AccountData) {
				defer wg.Done()
				func_obj(acc)
	
				global.CurrentProgress++
	
				utils.Sleep(
					global.Config.DelayBetweenAccs.Min,
					global.Config.DelayBetweenAccs.Max,
				)
				<-sem
			}(account)
		}
	
		wg.Wait()
	}

}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func main() {
    var err error

	// init log
	initLog()

	// parse config.yaml file
	utils.ParseConfig(filepath.Join("config", "config.yaml"))

	wr, err := os.OpenFile(filepath.Join("log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Panicf("Error When Opening Log File: %s\n", err)
	}

	defer func(wr *os.File) {
		err = wr.Close()
		if err != nil {
			log.Panicf("Error When Closing Log File: %s\n", err)
		}
	}(wr)
	mw := io.MultiWriter(os.Stdout, wr)
	log.SetOutput(mw)

    // handle panic
	defer handlePanic()

	// init Proxies
	err = utils.InitProxies(filepath.Join("config", "proxies.txt"))

	if err != nil {
		log.Panicf("Error initializing proxies: %s\n", err)
	}

	if len(utils.Proxies) <= 0 {
		global.Clients = append(global.Clients, utils.CreateClient(""))
	} else {
		for _, proxy := range utils.Proxies {
			global.Clients = append(global.Clients, utils.CreateClient(proxy))
		}
	}

    accountsListString, err := utils.ReadFileByRows(filepath.Join("config", "private_keys.txt"))

    if err != nil {
		log.Panicln(err.Error())
	}

	global.AccountsList, err = utils.GetAccounts(accountsListString, false)

	if err != nil {
		log.Panicln(err.Error())
	}

	log.Printf("Successfully Loaded %d Accounts\n", len(global.AccountsList))

	inputUserData := inputUser("\nThreads: ")
	threads, err := strconv.Atoi(inputUserData)

	if err != nil {
		log.Panicf("Wrong Threads Number: %s\n", inputUserData)
	}

	fmt.Printf("\n")
	global.TargetProgress = int64(len(accountsListString))

	// build CLI
	utils.Cli()

	// sleep before start
	delayMin, delayMax := global.Config.DelayBeforeStart.Min, global.Config.DelayBeforeStart.Max
	utils.Sleep(delayMin, delayMax)

	switch global.Module {
	case "Mint FUN Starts NFT":
		processAccounts(megaeth.MintFunNFT, threads)
	case "Faucet test tokens":
		processAccounts(megaeth.FaucetTokens, threads)
	case "Mint Megamafia NFT":
		processAccounts(megaeth.MintMegamafiaNFT, threads)
	case "Mint Mega Cat NFT":
		processAccounts(megaeth.MintMegaCatNFT, threads)
	case "Mint Blackhole NFT":
		processAccounts(megaeth.MintBlackholeNFT, threads)
	case "Mint Xyroph NFT":
		processAccounts(megaeth.MintXyrophNFT, threads)
	}

	log.Printf("The Work Has Been Successfully Finished\n")
	inputUser("\nPress Enter to Exit..")
}
