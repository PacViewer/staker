package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
)

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

		out, err := exec.Command(cfg.PactusWalletExecPath, args...).Output()
		if err != nil {
			log.Fatalf("err: %s, msg: %s", err.Error(), string(out))
		}
		log.Println(string(out))
	}

}
