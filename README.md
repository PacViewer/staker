# staker
Pactus bond tools for multiple stake by config, you can stake many validators in same time with specific configuration.

- Note: this approach is one to many (1 account address and many validators address)

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

- `pactus_wallet_exec_path` is patus-wallet file address
- `wallet_path` is address wallet file for example default path is `/home/{user}/pactus/wallets/default_wallet`
- `amount` a certain amount that is shared between the validators
- `wallet_address` sender address wallet have coin
- `validators` is list of validators for stake **note: if your address is first time for validator need set public key in `pub`**

2. run tools with `./staker -config ./cfg.json -password foobar`
- `-config` config file path default is `./cfg.json`
- `-password` is optional for wallet password
- `-server` is for custom node rpc address