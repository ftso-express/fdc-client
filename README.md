<p align="left">
  <a href="https://flare.network/" target="blank"><img src="https://flare.network/wp-content/uploads/Artboard-1-1.svg" width="400" height="300" alt="Flare Logo" /></a>
</p>

# Flare Data Connector Client

Flare Data Connector client supports tha attestation process. It does the following tasks:

-   Queries Flare C-Chain indexer for signing policies, attestation requests, and bitVotes.
-   Assigns the attestation requests to the correct voting rounds and begins their verification process.
-   Provides bitVote for each round.
-   Computes consensus bitVote for each round.
-   For each round, provides Merkle root of Merkle tree build on hashes of confirmed attestations.

The client has no direct interactions with the Flare blockchain/node. The data is read through C-Chain indexer and submitted through Flare System Client.

## Protocol

Link to specs (TODO)

## Tests

To run all tests locally and generate a coverage report:

```
$ ./gencover.sh
```

## Configurations

The configurations are set in `userConfig.toml` file in `configs` folder.

```toml
# options are: "coston", "songbird", "coston2", "flare"
chain = <chainName>

protocolId = <protocolId>
```

The client needs access to C-chain indexer

```toml
[db]
host = "localhost"
port = 3306
database = "flare_ftso_indexer"
username = "root"
password = "root"
log_queries = false
```

FSP client access data from FDC client through the rest server.

```toml
[rest_server]
# Addr optionally specifies the TCP address for the server to listen on, in the form "host:port". If empty, ":http" (port 80) is used. The service names are defined in RFC 6335 and assigned by IANA. See net.Dial for details of the address format.
addr = ":8080"
api_key_name = "X-API-KEY"
api_keys = ["12345", "123456"]
title = "FDC protocol data provider API"
fsp_sub_router_title = "FDC protocol data provider for FSP client"
fsp_sub_router_path = "/fsp"
version = "0.0.0"
swagger_path = "/api-doc"

```

For each supported attestation type, the ABI of the attestation response struct should be provided.
The ABI in json file should be saved in json file in `configs/abis` folder.
It is recommended for a file to be named `<attestationType>.json`.
In `userConfig.toml`, a path to the json file is specified.

For each supported source of an attestation type, an url and an API key of a verifier server should be specified.
In addition, LUT limit of the pair must be provided as a string representing a non-negative number smaller than $2^{64}$.

```toml
# Verifiers for <attestationType>
[verifiers.<attestationType>]
abi_path = "configs/abis/<attestationType>.json"

## <source1>
[verifiers.<attestationType>.Sources.<source1>]
url = "http://url/of/the/verifier1"
api_key = "api-key1"
lut_limit = "123124124"

## <source1>
[verifiers.<attestationType>.Sources.<source2>]
url = "http://url/of/the/verifier2"
api_key = "api-key2"
lut_limit = "123124124"

```

System configs for a pair of chain and protocol ID should be specified in
`configs/systemConfigs/<protrocolId>/<chain>.toml`

The client needs data from three contracts `Submit` for bitVotes, `Relay` for signing policies, and `FDC` for attestation requests.
The addresses must be specified in the systemConfig file.

```toml
[addresses]
submit_contract = "0x2cA6571Daa15ce734Bbd0Bf27D5C9D16787fc33f"
relay_contract = "0x32D46A1260BB2D8C9d5Ab1C9bBd7FF7D7CfaabCC"
fdc_contract = "0xCf6798810Bc8C0B803121405Fee2A5a9cc0CA5E5"
```
