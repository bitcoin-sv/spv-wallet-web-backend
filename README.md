# SPV WALLET WEB BACKEND

_Go application used as a web backend with storage for custodial wallet's xpriv keys, utilizing a non-custodial spv-wallet service._

---

The `spv-wallet-web-backend` is an HTTP server that serves as the referential backend server for a custodial web wallet for Bitcoin SV (BSV). It integrates and utilizes the `spv-wallet` component to handle various BSV-related operations, including the creation of transactions and listing incoming and outgoing transactions.

For in-depth information and guidance, please refer to the [SPV Wallet Documentation](https://docs.bsvblockchain.org/network-topology/applications/spv-wallet).

## How to use it

### Endpoints documentation

For endpoints documentation you can visit swagger which is exposed on port 8180 by default.

```
http://localhost:8180/swagger/index.html
```

## Contribution

To easy development process we use https://taskfile.dev/

### Installation of taskfile

Task offers many installation methods. Check out the available methods on [the installation documentation page](https://taskfile.dev/installation/).

### Tasks

#### List

To see the list of available tasks just call

```bash
task
```

#### Installing and set up tools

To install missing tools that are need just call

```bash
task install
```

##### Git hooks

Also consider installing git hooks with the command:

```bash
task git-hooks:install
```

So it will ensure that tests are passing and linter is not complaining before git push.

#### Start application

To start application locally with performed all the code checks first use the command

```bash
task start
```

if you just want to run app without all the checks

```bash
task run
```

#### Run checks & tests

To verify the code overall run verification task

```bash
task verify
```

If you just want a static check of a code run

```bash
task verify:check
```

If you want to run tests then run

```bash
task test
```

## License

[Open BSV License](https://github.com/bitcoin-sv/bitcoin-sv/blob/master/LICENSE).
