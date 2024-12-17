# op-enclave

`op-enclave` is a relatively small modification to the [op-stack](https://github.com/ethereum-optimism/optimism/)
that proves state transitions in a AWS Nitro Enclave, and submits the resulting state roots to the L1 chain.
This removes the need for the 7-day challenge period, and allows for immediate withdrawals.

## Directory Structure

<pre>
├── <a href="./bindings">bindings</a>: Go bindings for various contracts, generated by `make bindings`
├── <a href="./contracts">contracts</a>: Solidity contracts
├── <a href="./op-batcher">op-batcher</a>: Batcher modification that submits batches immediately after withdrawals are detected
├── <a href="./op-da">op-da</a>: Data availability service for writing to S3 / file system
├── <a href="./op-enclave">op-enclave</a>: Stateless transition function, for running in a AWS Nitro TEE
├── <a href="./op-proposer">op-proposer</a>: L2-Output Submitter, communicates with op-enclave and submits proposals to L1
├── <a href="./op-withdrawer">op-withdrawer</a>: Withdrawal utility for submitting withdrawals to L1
├── <a href="./register-signer">register-signer</a>: Registers a enclave signer key from a Nitro attestation with the SystemConfigGlobal contract
├── <a href="./testnet">testnet</a>: Dockerized testnet for running the op-enclave stack
</pre>

## Running a testnet

1. Deploy the Nitro certificate manager using `make deploy-cert-manager`:
```bash
IMPL_SALT=0 DEPLOY_PRIVATE_KEY=<privatekey> RPC_URL=https://sepolia.base.org make deploy-cert-manager
```

2. Deploy the system contracts using `make deploy`:
```bash
IMPL_SALT=0 DEPLOY_PRIVATE_KEY=<privatekey> DEPLOY_CONFIG_PATH=deploy-config/example.json RPC_URL=https://sepolia.base.org make deploy
```

3. Generate a testnet genesis block and deploy the proxy contracts for a new chain using `make testnet`:
```bash
DEPLOY_PRIVATE_KEY=<privatekey> L1_URL=https://sepolia.base.org make testnet
```

4. Copy `testnet/.env.example` to `testnet/.env` and fill in the environment variables,
in particular the `# per deploy` section at the top.

5. Run the testnet:
```bash
docker-compose -f testnet/Dockerfile up
```