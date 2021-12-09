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
4. [Metrics](#metrics)
4. [EVM-CLI](#evm-cli)
5. [Celo-CLI](#celo-cli)
6. [Substrate](#substrate)
7. [Local Setup](#local-setup)

## `Installation`
Refer to [installation](https://github.com/ChainSafe/chainbridge-docs/blob/develop/docs/installation.md) guide for assistance in installing.

## `Modules`

The chainbridge-core-example currently supports two modules:
1. [EVM-CLI](#evm-cli)
2. [Celo-CLI](#celo-cli)
3. [Substrate](#substrate)

## Usage
Since chainbridge-core is the modular framework it will require writing some code to get it running. 

You can find some examples [here](https://github.com/ChainSafe/chainbridge-core-example).

&nbsp;

## `Metrics`

Metrics, in `chainbridge-core` are handled via injecting metrics provider compatible with [Metrics](./relayer/relayer.go) interface into `Relayer` instance. It currently needs only one method, `TrackDepositMessage(m *message.Message)`
which is called inside relayer router when a `Deposit` event appears and should contain all necessary metrics data for extraction in custom implementations.

### `OpenTelementry`

`chainbridge-core` provides already implemented metrics [provider](./opentelemetry/opentelemetry.go) via [OpenTelemetry](https://opentelemetry.io/) that is vendor-agnostic. It sends metrics to a separate [collector](https://opentelemetry.io/docs/collector/) that then sends metrics via configured [exporters](https://opentelemetry.io/docs/collector/configuration/#exporters) to one or many metrics back-ends.


## `EVM-CLI`
This module provides instruction for communicating with EVM-compatible chains.

```bash
Usage:
   evm-cli [flags]
   evm-cli [command]

Available Commands:
  accounts    Set of commands for managing accounts
  admin       Set of commands for executing various admin actions
  bridge      Set of commands for interacting with a bridge
  centrifuge  Set of commands for interacting with a cetrifuge asset store contract
  deploy      Deploy smart contracts
  erc20       Set of commands for interacting with ERC20 contract
  erc721      Set of commands for interacting with an ERC721 contract
  utils       Set of utility commands

Flags:
  -h, --help   help for evm-cli
      --gas-limit uint              Gas limit used in transactions (default 6721975)
      --gas-price uint              Used as upperLimitGasPrice for transactions if not 0. Transactions gasPrice is defined by estimating it on network for pre London fork networks and by estimating BaseFee and MaxTipFeePerGas in post London networks
      --json-wallet string          Encrypted JSON wallet
      --json-wallet-password string Password for encrypted JSON wallet
      --network uint                ID of the Network
      --private-key string          Private key to use
      --url string                  URL of the node to receive RPC calls (default "ws://localhost:8545")
```

&nbsp;

### `Accounts`
Set of commands for managing accounts.

```bash
Usage:
   evm-cli accounts [command]

Available Commands:
  generate    Generate bridge keystore (Secp256k1)
  import      Import a bridge keystore
  transfer    Transfer base currency

Flags:
  -h, --help   help for accounts
```

#### `generate`
The generate subcommand is used to generate the bridge keystore. If no options are specified, a Secp256k1 key will be made.

```bash
Usage:
   evm-cli accounts generate [flags]

Flags:
  -h, --help   help for generate
```

#### `import`
The import subcommand is used to import a keystore for the bridge.

```bash
Usage:
   evm-cli accounts import [flags]

Flags:
  -h, --help               help for import
      --private-key string Private key to use
      --password string    Password to encrypt with
```

#### `transfer`
The transfer subcommand is used to transfer the base currency.

```bash
Usage:
   evm-cli accounts transfer [flags]

Flags:
      --amount string      Transfer amount
      --decimals uint      Base token decimals (default 0)
  -h, --help               help for transfer
      --recipient string   Recipient address
```

&nbsp;

### `Admin`
Set of commands for executing various admin actions.

```bash
Usage:
   evm-cli admin [command]

Available Commands:
  add-admin         Add a new admin
  add-relayer       Add a new relayer
  get-threshold     Get the relayer vote threshold
  is-relayer        Check if an address is registered as a relayer
  pause             Pause deposits and proposals
  remove-admin      Remove an existing admin
  remove-relayer    Remove an existing relayer
  set-deposit-nonce Set the deposit nonce
  set-fee           Set a new fee for deposits
  set-threshold     Set a new relayer vote threshold
  unpause           Unpause deposits and proposals
  withdraw          Withdraw tokens from a handler contract

Flags:
  -h, --help   help for admin
```

#### `add-admin`
The add-admin subcommand sets an address as a bridge admin.

```bash
Usage:
   evm-cli admin add-admin [flags]

Flags:
      --admin string    Address to add
      --bridge string   Bridge contract address
  -h, --help            help for add-admin
```

#### `add-relayer`
The add-relayer subcommand sets an address as a bridge relayer.

```bash
Usage:
   evm-cli admin add-relayer [flags]

Flags:
      --bridge string    Bridge contract address
  -h, --help             help for add-relayer
      --relayer string   Address to add
```

#### `get-threshold`
The get-threshold subcommand returns the relayer vote threshold.

```bash
Usage:
   evm-cli admin get-threshold [flags]

Flags:
      --bridge string    Bridge contract address
  -h, --help             help for get-threshold
```

#### `is-relayer`
The is-relayer subcommand checks if an address is registered as a relayer.

```bash
Usage:
   evm-cli admin is-relayer [flags]

Flags:
      --bridge string    Bridge contract address
  -h, --help             help for is-relayer
      --relayer string   Address to check
```

#### `pause`
The pause subcommand pauses deposits and proposals.

```bash
Usage:
   evm-cli admin pause [flags]

Flags:
      --bridge string   Bridge contract address
  -h, --help            help for pause

```

#### `remove-admin`
The remove-admin subcommand removes an existing admin.

```bash
Usage:
   evm-cli admin remove-admin [flags]

Flags:
      --admin string    Address to remove
      --bridge string   Bridge contract address
  -h, --help            help for remove-admin
```

#### `remove-relayer`
The remove-relayer subcommand removes an existing relayer.

```bash
Usage:
   evm-cli admin remove-relayer [flags]

Flags:
      --bridge string    Bridge contract address
  -h, --help             help for remove-relayer
      --relayer string   Address to remove
```

#### `set-deposit-nonce`
The set-deposit-nonce subcommand sets the deposit nonce. This nonce cannot be less than what is currently stored in the contract.

```bash
Usage:
   evm-cli admin set-deposit-nonce [flags]

Flags:
      --domain string   Domain ID of chain
      --deposit-nonce   Deposit nonce to set (does not decrement)
      --bridge string   Bridge contract address
  -h, --help            help for set-fee
```

#### `set-fee`
The set-fee subcommand sets a new fee for deposits.

```bash
Usage:
   evm-cli admin set-fee [flags]

Flags:
      --bridge string   Bridge contract address
      --fee string      New fee (in ether)
  -h, --help            help for set-fee
```

#### `set-threshold`
The set-threshold subcommand sets a new relayer vote threshold.

```bash
Usage:
   evm-cli admin set-threshold [flags]

Flags:
      --bridge string    Bridge contract address
  -h, --help             help for set-threshold
      --threshold uint   New relayer threshold
```

#### `unpause`
The unpause subcommand unpauses deposits and proposals.

```bash
Usage:
   evm-cli admin unpause [flags]

Flags:
      --bridge string   Bridge contract address
  -h, --help            help for unpause
```

#### `withdraw`
The withdraw subcommand withdrawals tokens from a handler contract.

```bash
Usage:
   evm-cli admin withdraw [flags]

Flags:
      --amount string         Token amount to withdraw, use only if ERC20 token is withdrawn. If both amount and token-id are set an error will occur
      --token string          Token ID to withdraw, use only if ERC721 token is withdrawn. If both amount and token-id are set an error will occur
      --bridge string         Bridge contract address
      --decimals uint         ERC20 token decimals
      --handler string        Handler contract address
  -h, --help                  help for withdraw
      --recipient string      Address to withdraw to
      --token-contract string ERC20 or ERC721 token contract address
```

&nbsp;

### `Bridge`
Set of commands for interacting with a bridge.

```bash
Usage:
   evm-cli bridge [command]

Available Commands:
  cancel-proposal           Cancel an expired proposal
  query-proposal            Query an inbound proposal
  query-resource            Query the resource ID for a handler contract
  register-generic-resource Register a generic resource ID
  register-resource         Register a resource ID
  set-burn                  Set a token contract as mintable/burnable

Flags:
  -h, --help   help for bridge
```

#### `cancel-proposal`
The cancel-proposal subcommand cancels an expired proposal.

```bash
Usage:
   evm-cli bridge cancel-proposal [flags]

Flags:
      --bridge string        Bridge contract address
      --data-hash string     Hash of proposal metadata
      --deposit-nonce uint   Deposit nonce of proposal to cancel
      --domain uint          Domain ID of proposal to cancel
  -h, --help                 help for cancel-proposal
```

#### `query-proposal`
The query-proposal subcommand queries an inbound proposal.

```bash
Usage:
   evm-cli bridge query-proposal [flags]

Flags:
      --bridge string        Bridge contract address
      --data-hash string     Hash of proposal metadata
      --deposit-nonce uint   Deposit nonce of proposal
      --domain uint          Source domain ID of proposal
  -h, --help                 help for query-proposal
```

#### `query-resource`
The query-resource subcommand queries the contract address with the provided resource ID for a specific handler contract.

```bash
Usage:
   evm-cli bridge query-resource [flags]

Flags:
      --handler string    Handler contract address
  -h, --help              help for query-resource
      --resource string   Resource ID to query
```

#### `register-generic-resource`
The register-generic-resource subcommand registers a resource ID with a contract address for a generic handler.

```bash
Usage:
   evm-cli bridge register-generic-resource [flags]

Flags:
      --bridge string     Bridge contract address
      --deposit string    Deposit function signature (default "0x00000000")
      --execute string    Execute proposal function signature (default "0x00000000")
      --handler string    Handler contract address
      --hash              Treat signature inputs as function prototype strings, hash and take the first 4 bytes
  -h, --help              help for register-generic-resource
      --resource string   Resource ID to query
      --target string     Contract address or hash storage to be registered
```

#### `register-resource`
The register-resource subcommand registers a resource ID with a contract address for a handler.

```bash
Usage:
   evm-cli bridge register-resource [flags]

Flags:
      --bridge string     Bridge contract address
      --handler string    Handler contract address
  -h, --help              help for register-resource
      --resource string   Resource ID to be registered
      --target string     Contract address to be registered
```

#### `set-burn`
The set-burn subcommand sets a token contract as mintable/burnable in a handler.

```bash
Usage:
   evm-cli bridge set-burn [flags]

Flags:
      --bridge string         Bridge contract address
      --handler string        ERC20 handler contract address
  -h, --help                  help for set-burn
      --token-contract string Token contract to be registered
```

&nbsp;

### `Deploy`
The deploy subcommand deploys all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags.

```bash
Usage:
   evm-cli deploy [flags]

Flags:
      --all                      Deploy all
      --bridge                   Deploy bridge
      --bridge-address string    Bridge contract address. Should be provided if handlers are deployed separately
      --domain uint8             Domain ID for the instance (default 1)
      --erc20                    Deploy ERC20
      --erc20-handler            Deploy ERC20 handler
      --erc20-name string        ERC20 contract name
      --erc20-symbol string      ERC20 contract symbol
      --erc721                   Deploy ERC721
      --erc721-base-uri string   ERC721 base URI
      --erc721-handler           Deploy ERC721 handler
      --erc721-name string       ERC721 contract name
      --erc721-symbol string     ERC721 contract symbol
      --fee uint                 Fee to be taken when making a deposit (in ETH, decimals are allowed)
      --generic-handler          Deploy generic handler
  -h, --help                     help for deploy
      --relayer-threshold uint   Number of votes required for a proposal to pass (default 1)
      --relayers strings         List of initial relayers
```

&nbsp;

### `ERC20`
Set of commands for interacting with an ERC20 contract.

```bash
Usage:
   evm-cli erc20 [command]

Available Commands:
  add-minter    Add a new ERC20 minter
  approve       Approve an ERC20 tokens
  balance       Query an ERC20 token balance
  deposit       Deposit an ERC20 token
  get-allowance Get the allowance of a spender for an address
  mint          Mint an ERC20 token

Flags:
  -h, --help   help for erc20
```

#### `add-minter`
The add-minter subcommand adds a minter to an ERC20 mintable contract.

```bash
Usage:
   evm-cli erc20 add-minter [flags]

Flags:
      --contract string   ERC20 contract address
  -h, --help              help for add-minter
      --minter string     Minter address
```

#### `get-allowance`
The get-allowance subcommand returns the allowance of a spender for an address.

```bash
Usage:
   evm-cli erc20 get-allowance [flags]

Flags:
      --contract string   ERC20 contract address
  -h, --help              help for get-allowance
      --owner string      Address of token owner
      --spender string    Address of spender
```

#### `approve`
The approve subcommand approves tokens in an ERC20 contract for transfer.

```bash
Usage:
   evm-cli erc20 approve [flags]

Flags:
      --amount string      Amount to grant allowance
      --contract string    ERC20 contract address
      --decimals uint      ERC20 token decimals
  -h, --help               help for approve
      --recipient string   Recipient address
```

#### `balance`
The balance subcommand queries the balance of an account in an ERC20 contract.

```bash
Usage:
   evm-cli erc20 balance [flags]

Flags:
      --address string    Address to receive balance of
      --contract string   ERC20 contract address
  -h, --help              help for balance
```

#### `deposit`
The deposit subcommand creates a new ERC20 token deposit on the bridge contract.

```bash
Usage:
   evm-cli erc20 deposit [flags]

Flags:
      --amount string      Amount to deposit
      --bridge string      Address of bridge contract
      --decimals uint      ERC20 token decimals
      --domain uint        Destination domain ID
  -h, --help               help for deposit
      --recipient string   Address of recipient
      --resource string    Resource ID for transfer
```

#### `mint`
The mint subcommand mints a token on an ERC20 mintable contract.

```bash
Usage:
   evm-cli erc20 mint [flags]

Flags:
      --amount string        Deposit amount
      --contract string      ERC20 contract address
      --decimals uint        ERC20 token decimals
      --destination string   Where tokens should be minted. Defaults to TX sender
  -h, --help                 help for mint
```

&nbsp;

### `ERC721`
ERC721-related instructions.

```bash
Usage:
   evm-cli erc721 [command]

Available Commands:
  add-minter  Add a new ERC721 minter
  approve     Approve an ERC721 token
  deposit     Deposit an ERC721 token
  mint        Mint an ERC721 token
  owner       Get an ERC721 token owner
```

#### `add-minter`
The add-minter subcommand adds a new minter address to an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 add-minter [flags]

Flags:
      --contract string   ERC721 contract address
  -h, --help              help for add-minter
      --minter string     Minter address
```

#### `approve`
The approve subcommand approves a token in an ERC721 contract for transfer.

```bash
Usage:
   evm-cli erc721 approve [flags]

Flags:
      --contract string    ERC721 contract address
  -h, --help               help for approve
      --recipient string   Recipient address
      --token string       ERC721 token ID
```

#### `deposit`
The deposit subcommand creates a new ERC721 token deposit on the bridge contract.

```bash
Usage:
   evm-cli erc721 deposit [flags]

Flags:
      --bridge string        Bridge contract address
      --destination string   Destination domain ID
  -h, --help                 help for deposit
      --metadata string      ERC721 token metadata
      --recipient string     Recipient address
      --resource string      Resource ID for transfer
      --token string         ERC721 token ID
```

#### `mint`
The mint subcommand mints a token on an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 mint [flags]

Flags:
      --contract string    ERC721 contract address
  -h, --help               help for mint
      --metadata string    ERC721 token metadata
      --recipient string   Recipient address
      --token string       ERC721 token ID
```

#### `owner`
The owner subcommand gets a token owner from an ERC721 mintable contract.

```bash
Usage:
   evm-cli erc721 owner [flags]

Flags:
      --contract string   ERC721 contract address
  -h, --help              help for owner
      --token string      ERC721 token ID
```

### `Utils`
Set of utility commands.

```bash
Usage:
   evm-cli utils [command]

Available Commands:
  hash-list   List tx hashes within N number of blocks
  simulate    Simulate the invocation of a transaction

Flags:
  -h, --help   help for utils
```

#### `hash-list`
The hash-list subcommand accepts a starting block to query, loops over N number of blocks past it, then prints this list of blocks to review hashes contained within.

```bash
Usage:
   evm-cli utils hash-list [flags]

Flags:
      --block-number string   Block number to start at
      --blocks string         Number of blocks past the provided block-number to review
  -h, --help                  help for hash-list
```

#### `simulate`
The simulate subcommand simulates a transaction result by simulating invocation; not state-altering.

```bash
Usage:
   evm-cli utils simulate [flags]

Flags:
      --block-number string   Block number
      --from string           Address of sender
  -h, --help                  help for simulate
      --tx-hash string        Transaction hash
```

### `Centrifuge`
Set of commands for interacting with a cetrifuge asset store contract.

```bash
Usage:
   evm-cli centrifuge [command]

Available Commands:
  deploy      Deploy a centrifuge asset store contract
  get-hash    Returns the status of whether a given hash exists in an asset store
```

#### `deploy`
The deploy subcommand deploys a Centrifuge asset store contract that represents bridged Centrifuge assets.

```bash
Usage:
   evm-cli centrifuge deploy

Flags:
  -h, --help   help for deploy
```

#### `get-hash`
The get-hash subcommand checks the _assetsStored map on a Centrifuge asset store contract to determine whether the asset hash exists or not.

```bash
Usage:
   evm-cli centrifuge get-hash [flags]

Flags:
      --address string   Centrifuge asset store contract address
      --hash string      A hash to lookup
  -h, --help             help for get-hash
```
&nbsp;

## `Celo-CLI`
This module provides instruction for communicating with Celo-compatible chains.

More information can be found about this module within its repository, listed below.

[Celo Module Repository](https://github.com/ChainSafe/chainbridge-celo-module)

&nbsp;

## `Substrate`
This module provides instruction for communicating with Substrate-compatible chains.

Currently there is no CLI for this, though more information can be found about this module within its repository, listed below.

[Substrate Module Repository](https://github.com/ChainSafe/chainbridge-substrate-module)

&nbsp;

## `Local Setup`

This section allows developers with a way to quickly and with minimal effort stand-up a local development environment in order to fully test out functionality of the chainbridge.

### `local`
The local-setup command deploys a bridge, ERC20, ERC721 and generic handler contracts with preconfigured accounts and appropriate handlers.

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

# ChainSafe Security Policy

## Reporting a Security Bug

We take all security issues seriously, if you believe you have found a security issue within a ChainSafe
project please notify us immediately. If an issue is confirmed, we will take all necessary precautions
to ensure a statement and patch release is made in a timely manner.

Please email us a description of the flaw and any related information (e.g. reproduction steps, version) to
[security at chainsafe dot io](mailto:security@chainsafe.io).

## License

_GNU Lesser General Public License v3.0_
