package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"regexp"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var blanceRegex = regexp.MustCompile(`balance: (\d+\.\d+)`)

type Cfg struct {
	PactusWalletExecPath string      `json:"pactus_wallet_exec_path"`
	WalletPath           string      `json:"wallet_path"`
	WalletAddress        string      `json:"wallet_address"`
	Amount               float64     `json:"amount"`
	Validators           []Validator `json:"validators"`
}

type Validator struct {
	Address string `json:"address"`
	Pub     string `json:"pub"`
}

func main() {
	cfgPath := flag.String("config", "./cfg.json", "config file path")
	password := flag.String("password", "", "pactus wallet password")
	rpc := flag.String("server", "", "custom node rpc")
	total := flag.Bool("total", false, "determine that all balance of account will be staked")

	flag.Parse()

	b, err := os.ReadFile(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Cfg

	if err := json.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}

	amount := strconv.FormatFloat(cfg.Amount, 'g', -1, 64)

	if *total {
		args := make([]string, 0)
		args = append(args, "address", "balance", cfg.WalletAddress)
		out, err := exec.Command(cfg.PactusWalletExecPath, args...).Output()
		if err != nil {
			log.Fatalf("err: %s, msg: %s", err.Error(), string(out))
		}

		match := blanceRegex.FindStringSubmatch(string(out))
		if len(match) > 1 {
			amount = match[1]
		} else {
			log.Fatalf("err: can't find the address balance, msg: %s", string(out))
		}
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		log.Println("Exiting...")
		os.Exit(0)
	}()


	for _, val := range cfg.Validators {
		args := make([]string, 0)
		args = append(args, "--path", cfg.WalletPath, "tx", "bond")

		if len(*password) != 0 {
			args = append(args, "-p", *password)
		}

		args = append(args, cfg.WalletAddress, "--no-confirm", "--pub", val.Pub, val.Address, amount)

		if len(*rpc) != 0 {
			args = append(args, "--server", *rpc)
		}

		wg.Add(1)
		go runCmd(ctx, cfg.PactusWalletExecPath, val.Address, &wg, args...)
	}

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		select {
		case s := <-interrupt:
			cancel()
			log.Printf("task canceled by user, %s", s.String())
		}
	}()

	wg.Wait()
}

func runCmd(ctx context.Context, pactusWalletExecPath, validator string, wg *sync.WaitGroup, args ...string) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
		out, err := exec.CommandContext(ctx, pactusWalletExecPath, args...).Output()
		if err != nil {
			log.Printf("validator: %s err: %s, msg: %s", validator, err.Error(), string(out))
		}
		log.Println(string(out))
	}
}
