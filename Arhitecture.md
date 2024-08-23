# Organization of this repo

This repository is organized into packages.
Each package contains logic for certain part of the FDC client

## Server

Implements

-   REST endpoints for Flare Systems Client
-   User facing REST endpoints for queriying merkle trees (attestation proofs)
-   exposes OpenAPIv3 schema
-   Swagger server

## Client

### Attestation

An attestation is created for each request emitted by the FDC contract. The
package includes all the functions needed to handle the requests, as well as
for well as a sub-package `bitvotes` for calculating a consensus for the
joint response.

### Collector

Collector package includes a functionality that is continuously making requests
to the C-chain database to obtain information about requests, bit votes and
signing policies posted on the chain. It submits the results to the
`Requests`, `BitVotes` and `SigningPolicies` channels.

#### Manager

Manager has an infinite loop that constantly listens to the `Requests`, `BitVotes` and `SigningPolicies` channels.

SigningPolicy is received once per reward epoch. It is stored in signing policy storage.
Requests are delivered as they are emitted on chain. 

Manager holds rounds in a cyclic storage. A request is assigned to the round according to the timestamp of its emission, it is put into a queue to be sent to the verifier.

BitVotes are received after the end of each choose phase.
Once they are received, the consensus bitVote is computed.
After the consensus bitVote is computed, the requests that were chosen but not confirmed are resent to the verifiers.

#### Round

The round holds all the information gathered and computed in one round.

It holds an attestation for each request.
The attestations are sorted by the order of the emissions of their underlying requests.

The bitVote for the round is build when queried and bitVotes consensus is saved when available.
The Merkle Tree can only be build after the consensus bitVote is computed and all the chosen requests are confirmed.

#### Config, Shared, Timing, and Utils

The packages `config`, `shared`, `timing`, and `utils`, respectively, include functions 
and structs for parsing the configuration of the client, creating shared data, calculating
the timing of various epochs, and for general work with slices and maps.

# Logic

FDC protocol works cyclically in epochs. Logic of the FDC client is aligned with these epochs and works in the following way:

-   Before the start of each epoch, the Signing Policy containing information about
    all the data providers is updated on chain. FDC client obtains the signing policy
    through the _collector_ fetching info from the C-chain indexer's database. The
    _manager_ saves the policy for future use.
-   During the epoch, users submit attestation requests on chain, which are then
    obtained by the _collector_ fetching info from the C-chain indexer's database. The
    _manager_ creates a struct _round_ and saves the attestation requests in it. It puts
    the attestations in queues and tries to confirm the attestation requests for which
    it has a verifier.
-   After the attestation collection period, the FDC client gets a request on the _server_
    from the Flare System Client to provide confirmations of the attestation requests.
    The FDC client provides a bit vector indicating which attestations can be confirmed.
    It returns the bit vector to the Flare System Client which puts it on chain.
-   The FDC client obtains all the bit votes submitted by data providers by the
    _collector_ fetching info from the C-chain indexer's database. The
    _manager_ saves the bit vectors in a _round_. Furthermore, it calculate the consensus
    bit vote, i.e. the indicator which attestations can be confirmed by the majority of
    data providers. If there are some attestations in the consensus that a provider
    did not confirm, it tries to do it again.
-   The FDC client gets a request on the _server_ from the Flare System Client to provide
    a commit of the confirmed attestations. The Flare System Client publishes the commit
    on chain.
-   The FDC client gets a request on the _server_ from the Flare System Client to provide
    a reveal of the confirmed attestations. The Flare System Client signs and publishes
    the reveal on chain.
