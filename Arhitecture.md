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

### Collector

### Attestation

An attestation is created for each request emitted by the FDC contract.

#### Manager

Manager has an infinite loop that constantly listens to the `Requests`, `BitVotes` and `SigningPolicies` channels.
SigningPolicy is received once per reward epoch. It is stored in signing policy storage.
Requests are delivered as they are emitted on chain. A request is assigned to the round according to the timestamp of its emission, it is put into a queue to be sent to the verifier.

Manager holds rounds in a cyclic storage.

BitVotes are received after the end of each choose phase.
Once they are received the consensus bitVote is computed.
After the consensus bitVote is computed, the requests that were chosen but not confirmed are resent to the verifiers.

#### Round

The round holds an attestation for each request.
The attestations are sorted by the order of the emissions of their underlying requests.

The bitVote for the round is build when queried.

The Merkle Tree can only be build after the consensus bitVote is computed and all the chosen requests are confirmed.
