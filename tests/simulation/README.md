# Simulation

## Simulate Flare environment

### Simulate chain and Flare System

Using the repository `flare-smart-contracts-v2`, branch `e2e-fdc-temp`
found [here](https://gitlab.com/flarenetwork/flare-smart-contracts-v2/-/tree/e2e-fdc-temp?ref_type=heads)
one can deploy all the Fast Updates contracts
together with the whole Flare system and voter repository. Navigate to the
repository and run

```bash
yarn install
yarn compile
yarn sim-node # in first terminal
yarn sim-run # in second terminal
```

This will start a Flare system on a local Hardhat node, register 4
data providers and start a simulation of FTSO v2 feed providers.

Extract from the logs of the above simulation the value `firstVotingEpochStartSec` that
is needed to configure FDC client. Put this value in the config
file `configs/systemConfigs/200/hardhat_test.toml`

```toml
[timing]
t0 = 1722567915 # to be replaced by the firstVotingEpochStartSec
```

### Run indexer and database

Using the repository `flare-system-c-chain-indexer` found [here](https://github.com/flare-foundation/flare-system-c-chain-indexer)
run an indexer of the hardhat chain used in the simulation. Deploy a database
(for example by running `docker compose up indexer-db`
in `docker/local` of this repository). Then run (in the `flare-system-c-chain-indexer` repository)

```bash
go run main.go --config c-chain-indexer-for-simulation-config.toml
```

where `c-chain-indexer-for-simulation-config.toml` is the config file found in this folder.

### Simulate participants submitting requests, other data providers, and a verification server

Run

```bash
go run tests/simulation/simulation_mock.go --config tests/configs/simulationConfig.toml
```

to simulate:

-   a user submitting a request on chain to be verified every round,
-   other data providers providing inputs,
-   a simple verification server that returns a hardcoded confirmation,
-   a system client querying the FDC server to provide data to be submitted on the blockchain.

## FDC client in simulated environment

Finally one can run a FDC client to participate in this simulation:

```bash
go run main/main.go --config tests/configs/simulationConfig.toml
```

TODO: Currently the simulation is rather static in the sense that the every round, at the
same time a request is sent, that is always confirmed by all the providers. Make the
simulation more dynamic.
