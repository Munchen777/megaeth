package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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

func processAccount(
	ctx context.Context,
	sem chan struct{},
	account types.AccountData,
	func_obj types.ModuleFunction,
	wg *sync.WaitGroup,
	) {
	defer wg.Done()
	defer func() { <-sem }()
	func_obj(ctx, account)
}

func processAccounts(func_obj types.ModuleFunction, threads int) {
	wg := &sync.WaitGroup{}
	sem := make(chan struct{}, threads)
	ctx := context.Background()
	ticker := time.NewTicker(
		time.Duration(global.Config.DelayBetweenAccs.Min) * time.Second,
	)
	defer ticker.Stop()

	for i, account := range global.AccountsList {
		if i >= 1 {
			<-ticker.C
			log.Infof("Sleep %d seconds ...", global.Config.DelayBetweenAccs.Min)
		}

		wg.Add(1)
		sem <- struct{}{}

		go processAccount(ctx, sem, account, func_obj, wg)

	}
	wg.Wait()
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

	// shuffle accounts
	if global.Config.ShuffleAccs {
		rand.NewSource(time.Now().Unix())
		rand.Shuffle(len(global.AccountsList), func(i, j int) {
			global.AccountsList[i], global.AccountsList[j] = global.AccountsList[j], global.AccountsList[i]
		})
	}

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
	case "Faucet test tokens":
		processAccounts(megaeth.FaucetTokens, threads)
	case "Mint Megamafia NFT",
		 "Mint FUN Starts NFT",
		 "Mint Mega Cat NFT",
		 "Mint Blackhole NFT",
		 "Mint Xyroph NFT",
		 "Mint Lord Lapin NFT",
		 "Mint Angry Monkeys",
		 "Mint Bloom NFT":
		minter, err := megaeth.GetInterfaceByModuleName()
		if err != nil {
			log.Panic(err)
		}
		func_obj := megaeth.StartMint(minter)
		processAccounts(func_obj, threads)
	}

	log.Printf("The Work Has Been Successfully Finished\n")
	inputUser("\nPress Enter to Exit..")
}
