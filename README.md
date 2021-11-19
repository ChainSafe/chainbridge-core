# Chainbridge Core
<a href="https://discord.gg/ykXsJKfhgq">
  <img alt="discord" src="https://img.shields.io/discord/593655374469660673?label=Discord&logo=discord&style=flat" />
</a>

Chainbridge-core is the project that was born from the existing version of [Chainbridge](https://github.com/ChainSafe/chainbridge). It was built to improve the maintainability and modularity of the current solution. The fundamental distinction is that chainbridge-core is more of a framework rather than a stand-alone application.

*Project still in deep beta*
- Chat with us on [discord](https://discord.gg/ykXsJKfhgq).

### Table of Contents

1. [Installation](#installation)
2. [Modules](#modules)
3. [Usage](#usage)
<<<<<<< HEAD
4. [Metrics](#metrics)
4. [EVM-CLI](#evm-cli)
5. [Celo-CLI](#celo-cli)
6. [Substrate](#substrate)
7. [Local Setup](#local-setup)
=======
4. [EVM-CLI](#evm-cli)
5. [Celo-CLI](#celo-cli)
>>>>>>> main

## Installation
Refer to [installation](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/installation.md) guide for assistance in installing.

## Modules

The chainbridge-core-example currently supports two modules:
1. [EVM-CLI](#evm-cli)
2. [Celo-CLI](#celo-cli)
<<<<<<< HEAD
3. [Substrate](#substrate)
=======
>>>>>>> main

## Usage
Since chainbridge-core is the modular framework it will require writing some code to get it running. Here you can find some examples

[Example](https://github.com/ChainSafe/chainbridge-core-example)

<<<<<<< HEAD
&nbsp;

## Metrics

Metrics, in `chainbridge-core` are handled via injecting metrics provider compatible with [Metrics](./relayer/relayer.go) interface into `Relayer` instance. It currently needs only one method, `TrackDepositMessage(m *message.Message)`
which is called inside relayer router when a `Deposit` event appears and should contain all necessary metrics data for extraction in custom implementations.

### OpenTelementry

`chainbridge-core` provides already implemented metrics [provider](./opentelemetry/opentelemetry.go) via [OpenTelemetry](https://opentelemetry.io/) that is vendor-agnostic. It sends metrics to a separate [collector](https://opentelemetry.io/docs/collector/) that then sends metrics via configured [exporters](https://opentelemetry.io/docs/collector/configuration/#exporters) to one or many metrics back-ends.


## `EVM-CLI`
=======

## EVM-CLI
>>>>>>> main
This module provides instruction for communicating with EVM-compatible chains.

```bash
Usage:
   evm-cli [command]

Available Commands:
  accounts    Account instructions
  admin       Admin-related instructions
  bridge      Bridge-related instructions
  deploy      Deploy smart contracts
  erc20       ERC20-related instructions
  erc721      ERC721-related instructions
  utils       Utils-related instructions

Flags:
  -h, --help   help for evm-cli
```

<<<<<<< HEAD
&nbsp;

### `Accounts`
=======
### Accounts
>>>>>>> main
Account instructions, allowing us to generate keypairs or import existing keypairs for use.

```bash
Usage:
   evm-cli accounts [command]

Available Commands:
  generate    Generate bridge keystore (Secp256k1)
  import      Import bridge keystore
<<<<<<< HEAD
  transfer    Transfer base currency
=======
>>>>>>> main

Flags:
  -h, --help   help for accounts
```

<<<<<<< HEAD
#### `generate`
=======
#### generate
>>>>>>> main
The generate subcommand is used to generate the bridge keystore. If no options are specified, a Secp256k1 key will be made.

```bash
Usage:
   evm-cli accounts generate [flags]

Flags:
  -h, --help   help for generate
```

<<<<<<< HEAD
#### `import`
=======
#### import
>>>>>>> main
The import subcommand is used to import a keystore for the bridge.

```bash
Usage:
   evm-cli accounts import [flags]

Flags:
  -h, --help              help for import
      --password string   password to encrypt with
```

<<<<<<< HEAD
#### `transfer`
The generate subcommand is used to transfer the base currency.

```bash
Usage:
   evm-cli accounts transfer [flags]

Flags:
      --amount string      transfer amount
      --decimals uint      base token decimals (default 18)
  -h, --help               help for transfer
      --recipient string   recipient address
```

&nbsp;

### `Admin`
=======
### Admin
>>>>>>> main
Admin-related instructions.

```bash
Usage:
   evm-cli admin [command]

Available Commands:
  add-admin      Add a new admin
  add-relayer    Add a new relayer
  is-relayer     Check if an address is registered as a relayer
  pause          Pause deposits and proposals
  remove-admin   Remove an existing admin
  remove-relayer Remove a relayer
  set-fee        Set a new fee for deposits
  set-threshold  Set a new relayer vote threshold
  unpause        Unpause deposits and proposals
<<<<<<< HEAD
  withdraw       Withdraw tokens from a handler contract
=======
  withdraw       Withdraw tokens from the handler contract
>>>>>>> main

Flags:
  -h, --help   help for admin
```

<<<<<<< HEAD
#### `add-admin`
=======
#### add-admin
>>>>>>> main
Add a new admin.

```bash
Usage:
   evm-cli admin add-admin [flags]

Flags:
      --admin string    address to add
      --bridge string   bridge contract address
  -h, --help            help for add-admin
```

<<<<<<< HEAD
#### `add-relayer`
=======
#### add-relayer
>>>>>>> main
Add a new relayer.

```bash
Usage:
   evm-cli admin add-relayer [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for add-relayer
      --relayer string   address to add
```

<<<<<<< HEAD
#### `is-relayer`
=======
#### is-relayer
>>>>>>> main
Check if an address is registered as a relayer.

```bash
Usage:
   evm-cli admin is-relayer [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for is-relayer
      --relayer string   address to check
```

<<<<<<< HEAD
#### `pause`
=======
#### pause
>>>>>>> main
Pause deposits and proposals,

```bash
Usage:
   evm-cli admin pause [flags]

Flags:
      --bridge string   bridge contract address
  -h, --help            help for pause

```

<<<<<<< HEAD
#### `remove-admin`
=======
#### remove-admin
>>>>>>> main
Remove an existing admin.

```bash
Usage:
   evm-cli admin remove-admin [flags]

Flags:
      --admin string    address to remove
      --bridge string   bridge contract address
  -h, --help            help for remove-admin
```

<<<<<<< HEAD
#### `remove-relayer`
=======
#### remove-relayer
>>>>>>> main
Remove a relayer.

```bash
Usage:
   evm-cli admin remove-relayer [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for remove-relayer
      --relayer string   address to remove
```

<<<<<<< HEAD
#### `set-fee`
=======
#### set-fee
>>>>>>> main
Set a new fee for deposits.

```bash
Usage:
   evm-cli admin set-fee [flags]

Flags:
      --bridge string   bridge contract address
      --fee string      New fee (in ether)
  -h, --help            help for set-fee
```

<<<<<<< HEAD
#### `set-threshold`
=======
#### set-threshold
>>>>>>> main
Set a new relayer vote threshold.

```bash
Usage:
   evm-cli admin set-threshold [flags]

Flags:
      --bridge string    bridge contract address
  -h, --help             help for set-threshold
      --threshold uint   new relayer threshold
```

<<<<<<< HEAD
#### `unpause`
=======
#### unpause
>>>>>>> main
Unpause deposits and proposals.

```bash
Usage:
   evm-cli admin unpause [flags]

Flags:
      --bridge string   bridge contract address
  -h, --help            help for unpause
```

<<<<<<< HEAD
#### `withdraw`
Withdraw tokens from a handler contract.
=======
#### withdraw
Withdraw tokens from the handler contract.
>>>>>>> main

```bash
Usage:
   evm-cli admin withdraw [flags]

Flags:
      --amount string      token amount to withdraw. Should be set or ID or amount if both set error will occur
      --bridge string      bridge contract address
      --decimals uint      ERC20 token decimals
      --handler string     handler contract address
  -h, --help               help for withdraw
      --id string          token ID to withdraw. Should be set or ID or amount if both set error will occur
      --recipient string   address to withdraw to
      --token string       ERC20 or ERC721 token contract address
```

<<<<<<< HEAD
&nbsp;

### `Bridge`
=======
### Bridge
>>>>>>> main
Bridge-related instructions.

```bash
Usage:
   evm-cli bridge [command]

Available Commands:
  cancel-proposal           Cancel an expired proposal
  query-proposal            Query an inbound proposal
  query-resource            Query the contract address
  register-generic-resource Register a generic resource ID
  register-resource         Register a resource ID
  set-burn                  Set a token contract as mintable/burnable

Flags:
  -h, --help   help for bridge
```

<<<<<<< HEAD
#### `cancel-proposal`
=======
#### cancel-proposal
>>>>>>> main
Cancel an expired proposal.

```bash
Usage:
   evm-cli bridge cancel-proposal [flags]

Flags:
      --bridge string       bridge contract address
<<<<<<< HEAD
      --domainId uint       domain ID of proposal to cancel
=======
      --chainId uint        chain ID of proposal to cancel
>>>>>>> main
      --dataHash string     hash of proposal metadata
      --depositNonce uint   deposit nonce of proposal to cancel
  -h, --help                help for cancel-proposal
```

<<<<<<< HEAD
#### `query-proposal`
=======
#### query-proposal
>>>>>>> main
Query an inbound proposal.

```bash
Usage:
   evm-cli bridge query-proposal [flags]

Flags:
      --bridge string       bridge contract address
<<<<<<< HEAD
      --domainId uint       source domain ID of proposal
=======
      --chainId uint        source chain ID of proposal
>>>>>>> main
      --dataHash string     hash of proposal metadata
      --depositNonce uint   deposit nonce of proposal
  -h, --help                help for query-proposal
```

<<<<<<< HEAD
#### `query-resource`
=======
#### query-resource
>>>>>>> main
Query the contract address with the provided resource ID for a specific handler contract.

```bash
Usage:
   evm-cli bridge query-resource [flags]

Flags:
      --handler string      handler contract address
  -h, --help                help for query-resource
      --resourceId string   resource ID to query
```

<<<<<<< HEAD
#### `register-generic-resource`
=======
#### register-generic-resource
>>>>>>> main
Register a resource ID with a contract address for a generic handler.

```bash
Usage:
   evm-cli bridge register-generic-resource [flags]

Flags:
      --bridge string       bridge contract address
      --deposit string      deposit function signature (default "0x00000000")
<<<<<<< HEAD
      --depositerOffset int   depositer address position offset in the metadata, in bytes
=======
>>>>>>> main
      --execute string      execute proposal function signature (default "0x00000000")
      --handler string      handler contract address
      --hash                treat signature inputs as function prototype strings, hash and take the first 4 bytes
  -h, --help                help for register-generic-resource
      --resourceId string   resource ID to query
      --target string       contract address to be registered
```

<<<<<<< HEAD
#### `register-resource`
=======
#### register-resource
>>>>>>> main
Register a resource ID

```bash
Usage:
   evm-cli bridge register-resource [flags]

Flags:
      --bridge string       bridge contract address
      --handler string      handler contract address
  -h, --help                help for register-resource
      --resourceId string   resource ID to be registered
      --target string       contract address to be registered
```

<<<<<<< HEAD
#### `set-burn`
=======
#### set-burn
>>>>>>> main
Set a token contract as mintable/burnable

```bash
Usage:
   evm-cli bridge set-burn [flags]

Flags:
      --bridge string          bridge contract address
      --handler string         ERC20 handler contract address
  -h, --help                   help for set-burn
      --tokenContract string   token contract to be registered
```

<<<<<<< HEAD
&nbsp;

### `Deploy`
=======
### Deploy
>>>>>>> main
Deploy smart contracts.

Used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags

```bash
Usage:
   evm-cli deploy [flags]

Flags:
      --all                     deploy all
      --bridge                  deploy bridge
      --bridgeAddress string    bridge contract address. Should be provided if handlers are deployed separately
<<<<<<< HEAD
      --domainId string          domain ID for the instance (default "1")
=======
      --chainId string          chain ID for the instance (default "1")
>>>>>>> main
      --erc20                   deploy ERC20
      --erc20Handler            deploy ERC20 handler
      --erc20Name string        ERC20 contract name
      --erc20Symbol string      ERC20 contract symbol
      --erc721                  deploy ERC721
<<<<<<< HEAD
      --erc721Handler           deploy ERC721 handler
      --erc721Name string       ERC721 contract name
      --erc721Symbol string     ERC721 contract symbol
      --erc721BaseURI           ERC721 base URI
      --genericHandler          deploy generic handler
=======
>>>>>>> main
      --fee string              fee to be taken when making a deposit (in ETH, decimas are allowed) (default "0")
  -h, --help                    help for deploy
      --relayerThreshold uint   number of votes required for a proposal to pass (default 1)
      --relayers strings        list of initial relayers
```

<<<<<<< HEAD
&nbsp;

### `ERC20`
=======
### ERC20
>>>>>>> main
ERC20-related instructions.

```bash
Usage:
   evm-cli erc20 [command]

Available Commands:
  add-minter  Add a minter to an Erc20 mintable contract
  allowance   Get the allowance of a spender for an address
  approve     Approve tokens in an ERC20 contract for transfer
  balance     Query balance of an account in an ERC20 contract
  deposit     Initiate a transfer of ERC20 tokens
  mint        Mint tokens on an ERC20 mintable contract

Flags:
  -h, --help   help for erc20
```

<<<<<<< HEAD
#### `add-minter`
=======
#### add-minter
>>>>>>> main
Add a minter to an Erc20 mintable contract.

```bash
Usage:
   evm-cli erc20 add-minter [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for add-minter
<<<<<<< HEAD
      --minter string         handler contract address

```

#### `allowance`
=======
      --minter string         address of minter

```

#### allowance
>>>>>>> main
Get the allowance of a spender for an address.

```bash
Usage:
   evm-cli erc20 allowance [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for allowance
      --owner string          address of token owner
      --spender string        address of spender
```

<<<<<<< HEAD
#### `approve`
=======
#### approve
>>>>>>> main
Approve tokens in an ERC20 contract for transfer.

```bash
Usage:
   evm-cli erc20 approve [flags]

Flags:
      --amount string         amount to grant allowance
      --decimals uint         ERC20 token decimals (default 18)
      --erc20address string   ERC20 contract address
  -h, --help                  help for approve
      --recipient string      address of recipient
```

<<<<<<< HEAD
#### `balance`
=======
#### balance
>>>>>>> main
Query balance of an account in an ERC20 contract.

```bash
Usage:
   evm-cli erc20 balance [flags]

Flags:
      --accountAddress string   address to receive balance of
      --erc20Address string     ERC20 contract address
  -h, --help                    help for balance
```

<<<<<<< HEAD
#### `deposit`
=======
#### deposit
>>>>>>> main
Initiate a transfer of ERC20 tokens.

```bash
Usage:
   evm-cli erc20 deposit [flags]

Flags:
      --amount string       amount to deposit
      --bridge string       address of bridge contract
      --decimals uint       ERC20 token decimals
<<<<<<< HEAD
      --domainId string     destination domain ID
=======
      --destId string       destination chain ID
>>>>>>> main
  -h, --help                help for deposit
      --recipient string    address of recipient
      --resourceId string   resource ID for transfer
```

<<<<<<< HEAD
#### `mint`
=======
#### mint
>>>>>>> main
Mint tokens on an ERC20 mintable contract.

```bash
Usage:
   evm-cli erc20 mint [flags]

Flags:
      --amount string         amount to mint fee (in ETH)
      --decimal uint          ERC20 token decimals (default 18)
      --dstAddress string     Where tokens should be minted. Defaults to TX sender
      --erc20Address string   ERC20 contract address
  -h, --help                  help for mint
```

<<<<<<< HEAD
&nbsp;

### `ERC721`
ERC721-related instructions.

```bash
Usage:
   evm-cli erc721 [command]

Available Commands:
  add-minter  Add a minter to an Erc721 mintable contract
  approve     Approve token in an ERC721 contract for transfer
  deposit     Initiate a transfer of ERC721 token
  mint        Mint token on an ERC721 mintable contract
  owner       Get token owner from an ERC721 mintable contract
```

#### `add-minter`
=======
### ERC721
ERC721-related instructions.

#### add-minter
>>>>>>> main
Add a minter to an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 add-minter [flags]

Flags:
      --erc721Address string   ERC721 contract address
  -h, --help                   help for add-minter
      --minter string          address of minter
```

<<<<<<< HEAD
#### `approve`
Approve token in an ERC721 contract for transfer.

```bash
Usage:
   evm-cli erc721 approve [flags]

Flags:
      --contract-address string   ERC721 contract address
      --recipient string          address of recipient
      --tokenId string            ERC721 token ID
  -h, --help                      help for add-minter
```

#### `deposit`
Deposit ERC721 token.

```bash
Usage:
   evm-cli erc721 deposit [flags]

Flags:
      --contract-address string   ERC721 contract address
      --recipient string          address of recipient
      --bridge string             address of bridge contract
      --destId string             destination domain ID
      --resourceId string         resource ID for transfer
      --tokenId string            ERC721 token ID
  -h, --help                      help for add-minter
```

#### `mint`
Mint token on an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 mint [flags]

Flags:
      --contract-address string      ERC721 contract address
      --destination-address string   address of token recipient
      --tokenId string               ERC721 token ID
      --metadata string              token metadata
  -h, --help                         help for add-minter
```

#### `owner`
Get token owner from an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 owner [flags]

Flags:
      --contract-address string      ERC721 contract address
      --tokenId string               ERC721 token ID
  -h, --help                         help for add-minter
```

=======
>>>>>>> main
### Utils
Utils-related instructions.
*Useful for debugging*

```bash
Usage:
   evm-cli utils [command]

Available Commands:
  hashList    List tx hashes
  simulate    Simulate transaction invocation

Flags:
  -h, --help   help for utils
```

<<<<<<< HEAD
#### `hashlist`
=======
#### hashlist
>>>>>>> main
List tx hashes.

```bash
Usage:
   evm-cli utils hashList [flags]

Flags:
      --blockNumber string   block number
  -h, --help                 help for hashList
```

<<<<<<< HEAD
#### `simulate`
=======
#### simulate
>>>>>>> main
Replay a failed transaction by simulating invocation; not state-altering

```bash
Usage:
   evm-cli utils simulate [flags]

Flags:
      --blockNumber string   block number
      --fromAddress string   address of sender
  -h, --help                 help for simulate
      --txHash string        transaction hash
```

<<<<<<< HEAD
### `Centrifuge`
Centrifuge-related instructions.

#### `deploy`

This command can be used to deploy Centrifuge asset store contract that represents bridged Centrifuge assets.

```bash
Usage:
   evm-cli centrifuge deploy
```

#### `getHash`
Checks _assetsStored map on Centrifuge asset store contract to find if asset hash exists.

```bash
Usage:
   evm-cli centrifuge getHash [flags]

Flags:
      --address string   Centrifuge asset store contract address
      --hash string      A hash to lookup
  -h, --help             help for getHash
```

&nbsp;

## `Celo-CLI`
=======
## Celo-CLI
>>>>>>> main
Though Celo is an EVM-compatible chain, it deviates in its implementation of the original Ethereum specifications, and therefore is deserving of its own separate module.

See: [differences between EVM and Celo](#differences-between-evm-and-celo).

```bash
Usage:
   celo-cli [command]

Available Commands:
  bridge      Bridge-related instructions
  deploy      Deploy smart contracts
  erc20       erc20-related instructions

Flags:
  -h, --help   help for celo-cli
```

<<<<<<< HEAD
&nbsp;

=======
>>>>>>> main
### Differences Between EVM and Celo

The differences alluded to above in how Celo constructs transactions versus those found within Ethereum can be viewed below by taking a look at the Message structs in both implementations.

[Ethereum Message Struct](https://github.com/ethereum/go-ethereum/blob/ac7baeab57405c64592b1646a91e0a2bb33d8d6c/core/types/transaction.go#L586-L598)

Here you will find fields relating to the most recent London hardfork (EIP-1559), most notably `gasFeeCap` and `gasTipCap`.

```go
Message {
   from:       from,
   to:         to,
   nonce:      nonce,
   amount:     amount,
   gasLimit:   gasLimit,
   gasPrice:   gasPrice,
   gasFeeCap:  gasFeeCap,
   gasTipCap:  gasTipCap,
   data:       data,
   accessList: accessList,
   isFake:     isFake,
}
```

[Celo Message Struct](https://github.com/ChainSafe/chainbridge-celo-module/blob/b6d7ad422a5356500d2d5cf0b98e00da86dbb42e/transaction/tx.go#L422-L435)

In Celo's struct you will notice that there are additional fields added for `feeCurrency`, `gatewayFeeRecipient` and `gatewayFee`. You may also notice the `ethCompatible` field, a boolean value we added in order to quickly determine whether the message is Ethereum compatible or not, ie, that `feeCurrency`, `gatewayFeeRecipient` and `gatewayFee` are omitted.

```go
Message {
   from:                from,
   to:                  to,
   nonce:               nonce,
   amount:              amount,
   gasLimit:            gasLimit,
   gasPrice:            gasPrice,
   feeCurrency:         feeCurrency,         // Celo-specific
   gatewayFeeRecipient: gatewayFeeRecipient, // Celo-specific
   gatewayFee:          gatewayFee,          // Celo-specific
   data:                data,
   ethCompatible:       ethCompatible,       // Bool to check presence of: feeCurrency, gatewayFeeRecipient, gatewayFee
   checkNonce:          checkNonce,
}
```

<<<<<<< HEAD
&nbsp;

### `Bridge`
=======
### Bridge
>>>>>>> main
Bridge-related instructions.

```bash
Usage:
   celo-cli bridge [command]

Available Commands:
  register-resource Register a resource ID
  set-burn          Set a token contract as mintable/burnable

Flags:
  -h, --help   help for bridge
```

<<<<<<< HEAD
#### `register-resource`
=======
#### register-resource
>>>>>>> main
Register a resource ID with a contract address for a handler

```bash
Usage:
   celo-cli bridge register-resource [flags]

Flags:
      --bridge string       bridge contract address
      --handler string      handler contract address
  -h, --help                help for register-resource
      --resourceId string   resource ID to be registered
      --target string       contract address to be registered
```

<<<<<<< HEAD
#### `set-burn`
=======
#### set-burn
>>>>>>> main
Set a token contract as mintable/burnable in a handler

```bash
Usage:
   celo-cli bridge set-burn [flags]

Flags:
      --bridge string          bridge contract address
      --handler string         ERC20 handler contract address
  -h, --help                   help for set-burn
      --tokenContract string   token contract to be registered
```

<<<<<<< HEAD
&nbsp;

### `Deploy`
=======
### Deploy
>>>>>>> main
Deploy smart contracts.

This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags.

```bash
Usage:
   celo-cli deploy [flags]

Flags:
      --all                     deploy all
      --bridge                  deploy bridge
      --bridgeAddress string    bridge contract address. Should be provided if handlers are deployed separately
<<<<<<< HEAD
      --domainId string         domain ID for the instance (default "1")
=======
      --chainId string          chain ID for the instance (default "1")
>>>>>>> main
      --erc20                   deploy ERC20
      --erc20Handler            deploy ERC20 handler
      --erc20Name string        ERC20 contract name
      --erc20Symbol string      ERC20 contract symbol
      --erc721                  deploy ERC721
<<<<<<< HEAD
      --erc721Handler           deploy ERC721 handler
      --erc721Name string       ERC721 contract name
      --erc721Symbol string     ERC721 contract symbol
      --erc721BaseURI           ERC721 base URI
      --genericHandler          deploy generic handler
=======
>>>>>>> main
      --fee string              fee to be taken when making a deposit (in ETH, decimas are allowed) (default "0")
  -h, --help                    help for deploy
      --relayerThreshold uint   number of votes required for a proposal to pass (default 1)
      --relayers strings        list of initial relayers
```

<<<<<<< HEAD
&nbsp;

### `ERC20`
erc20-related instructions.
=======
### ERC20
erc20-related instructions
>>>>>>> main

```bash
Usage:
   celo-cli erc20 [command]

Available Commands:
  add-minter  Add a minter to an Erc20 mintable contract
  allowance   Set a token contract as mintable/burnable
  approve     Approve tokens in an ERC20 contract for transfer
  balance     Query balance of an account in an ERC20 contract
  deposit     Initiate a transfer of ERC20 tokens
  mint        Mint tokens on an ERC20 mintable contract

Flags:
  -h, --help   help for erc20
```

<<<<<<< HEAD
#### `add-minter`
=======
#### add-minter
>>>>>>> main
Add a minter to an Erc20 mintable contract.

```bash
Usage:
   celo-cli erc20 add-minter [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for add-minter
      --minter string         address of minter
```

<<<<<<< HEAD
#### `allowance`
=======
#### allowance
>>>>>>> main
Set a token contract as mintable/burnable in a handler.

```bash
Usage:
   celo-cli erc20 allowance [flags]

Flags:
      --erc20Address string   ERC20 contract address
  -h, --help                  help for allowance
      --owner string          address of token owner
      --spender string        address of spender
```

<<<<<<< HEAD
#### `approve`
=======
#### approve
>>>>>>> main
Approve tokens in an ERC20 contract for transfer.

```bash
Usage:
   celo-cli erc20 approve [flags]

Flags:
      --amount string         amount to grant allowance
      --decimals uint         ERC20 token decimals (default 18)
      --erc20address string   ERC20 contract address
  -h, --help                  help for approve
      --recipient string      address of recipient
```

<<<<<<< HEAD
#### `balance`
=======
#### balance
>>>>>>> main
Query balance of an account in an ERC20 contract.

```bash
Usage:
   celo-cli erc20 balance [flags]

Flags:
      --accountAddress string   address to receive balance of
      --erc20Address string     ERC20 contract address
  -h, --help                    help for balance
```

<<<<<<< HEAD
#### `deposit`
=======
#### deposit
>>>>>>> main
Initiate a transfer of ERC20 tokens.

```bash
Usage:
   celo-cli erc20 deposit [flags]

Flags:
      --amount string       amount to deposit
      --bridge string       address of bridge contract
      --decimals uint       ERC20 token decimals
<<<<<<< HEAD
      --domainId string     destination domain ID
=======
      --destId string       destination chain ID
>>>>>>> main
  -h, --help                help for deposit
      --recipient string    address of recipient
      --resourceId string   resource ID for transfer
```

<<<<<<< HEAD
#### `mint`
=======
#### mint
>>>>>>> main
Mint tokens on an ERC20 mintable contract.

```bash
Usage:
   celo-cli erc20 mint [flags]

Flags:
      --amount string         amount to mint fee (in ETH)
      --decimal uint          ERC20 token decimals (default 18)
      --dstAddress string     Where tokens should be minted. Defaults to TX sender
      --erc20Address string   ERC20 contract address
  -h, --help                  help for mint
```

<<<<<<< HEAD
&nbsp;

### Centrifuge
Centrifuge-related instructions.

#### `deploy`

This command can be used to deploy Centrifuge asset store contract that represents bridged Centrifuge assets.

```bash
Usage:
   evm-cli centrifuge deploy
```

#### `getHash`
Checks _assetsStored map on Centrifuge asset store contract to find if asset hash exists.

```bash
Usage:
   evm-cli centrifuge getHash [flags]

Flags:
      --address string   Centrifuge asset store contract address
      --hash string      A hash to lookup
  -h, --help             help for getHash
```

&nbsp;

## Substrate
This module provides instruction for communicating with Substrate-compatible chains.

Currently there is no CLI for this, though more information can be found about this module within its repository, listed below.

[Substrate Module Repository](https://github.com/ChainSafe/chainbridge-substrate-module)

&nbsp;

## Local Setup

This section allows developers with a way to quickly and with minimal effort stand-up a local development environment in order to fully test out functionality of the chainbridge.

### `local`

Locally deploy bridge and ERC20 handler contracts with preconfigured accounts and ERC20 handler.

```bash
Usage:
   local-setup [flags]

Flags:
  -h, --help   help for local-setup
```

This can be easily run by building the [chainbridge-core-example](https://github.com/ChainSafe/chainbridge-core-example) app, or by issuing a `Makefile` instruction directly from the root of the [chainbridge-core](https://github.com/ChainSafe/chainbridge-core) itself.
```bash
make local-setup
```
##### ^ this command will run a shell script that contains instructions for running two EVM chains via [Docker](https://www.docker.com/) (`docker-compose`). Note: this will likely take a few minutes to run.
&nbsp;

You can also review our [Local Setup Guide](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/guides/local-setup-guide.md) for a more detailed example of setting up a local development environment manually.

&nbsp;

=======
>>>>>>> main
# ChainSafe Security Policy

## Reporting a Security Bug

We take all security issues seriously, if you believe you have found a security issue within a ChainSafe
project please notify us immediately. If an issue is confirmed, we will take all necessary precautions
to ensure a statement and patch release is made in a timely manner.

Please email us a description of the flaw and any related information (e.g. reproduction steps, version) to
[security at chainsafe dot io](mailto:security@chainsafe.io).

## License

_GNU Lesser General Public License v3.0_
