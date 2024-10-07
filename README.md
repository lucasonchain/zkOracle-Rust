# Rust Implementation of Cross-Chain Oracle Using an Off-Chain Aggregation Mechanism Based on zk-SNARKs

This is a Rust port of the [zkOracle](https://github.com/0xEigenLabs/zkOracle) project, originally implemented in Go. zkOracle-Rust aims to provide the same functionality as the original project while leveraging the benefits of the Rust programming language.

Original Work: https://github.com/soberm/zkOracle
Paper: https://arxiv.org/abs/2405.08395

zkOracle-Rust is a zero-knowledge oracle system implemented in Rust. It allows for the creation and verification of zero-knowledge proofs for oracle data, enabling privacy-preserving data feeds for blockchain applications.

This project is a work in progress, porting the original Go implementation to Rust.

## Prerequisites

Before you begin, ensure you have met the following requirements:

* Rust (latest stable version) - [Install Rust](https://www.rust-lang.org/tools/install)[Golang](https://golang.org/doc/install) (version 1.19)
* [Node.js](https://nodejs.org/) (version >= 19.4.0)
* [Hardhat](https://hardhat.org/) (version >= 2.11.1)
* [Solidity](https://docs.soliditylang.org/en/latest/installing-solidity.html) (^0.8.0)

## Installation

To install zkOracle-Rust, follow these steps:

1. Change into the following directory: `cd node/cmd/compiler`
2. Install all dependencies: `go mod download`
3. Adapt the parameters of the circuits in the respective src files
4. Build the constraint system setup: `go build -o compiler`
5. Run the constraint system setup: `./compiler -b ./build`
6. Update the verifier contracts using the generated contracts

### Smart Contracts

1. Change into the contract directory: `cd contracts/`
2. Install all dependencies: `npm install`
3. Compile contracts: `hardhat compile`
4. Update the parameters in ./scripts/deploy.ts
5. Deploy contracts: `hardhat run --network <your_network> ./scripts/deploy.ts`

### Node

```
cargo run
```

## Contributing

Contributions to zkOracle-Rust are welcome. Please feel free to submit a Pull

## Licence

This project is licensed under the [MIT License](LICENSE).
