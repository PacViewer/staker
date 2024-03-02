# staker
Pactus bond tools for multiple stake by config

## how to run?

1. create config file `cfg.json`

```json
{
  "pactus_wallet_exec_path": "/home/user/pactus/pactus-wallet",
  "wallet_path": "/home/user/pactus/wallets/default_wallet",
  "amount": 2.2,
  "wallet_address": "",
  "validators": [
    {
      "address": "",
      "pub": ""
    }
  ]
}
```

2. run tools with `./staker -config ./cfg.json -password foobar`
- `--config` config file path default is `./cfg.json`
- `--password` is optional for wallet password
- `--server` is for custom node rpc address
- `--total` is a flag that ignore amount it config and it will stake whole of account balance
