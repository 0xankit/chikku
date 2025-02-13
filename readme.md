<h1 align="center">Chikku</h1>

<!-- show image in middle -->
<p align="center" style="text-align: center; background-color: #f0f0f0;">
  <img src="public/chikku.png" alt="Chikku" width="200"/>
</p>

**Chikku** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

Warning: This is a development version of the blockchain. Do not use in production.

- EGV Token and Dynamic Rewards: Implement a system for distributing "egv" tokens to 10 pre-defined addresses. Design a reward mechanism that distributes additional egv tokens based on transaction volume, calculated as:

  ```go
    individual_reward = (individual_address_transactions / total_network_transactions) * (inflation_rate * total_supply).
  ```

- Transaction-Weighted Governance(Bonus): Develop a governance system that allows for modification of the inflation rate via proposals and weighted voting. Voting weight should be determined by network contribution:

  ```go
  Voting weight = (individual_address_transactions / total_network_transactions)
  ```

## Architecture

<p align="center" style="text-align: center; background-color: #f0f0f0;">
  <img src="public/diag.png" alt="Chikku" />
</p>

1. x/egvmod: EGV Token and Dynamic Rewards

    - Predefined addresses that receive EGV tokens:

      - [X] Register the predefined addresses (10 operators with some egv tokens) from the config as module params file.

    - Params:

      - InflationRate: The rate at which new tokens are minted.
      - distrubutionInterval: The interval at which rewards are distributed.
      - operators: The list of predefined operators.

    - Structs:

      - TotalTrxsCount: The total number of transactions in the network since genesis at given block height.

        ```go
        total_network_trxs = sum(operator_cumulative_trxs_count)

        (block_height => [(operator_address, cumulative_trxs_count), ...])
        ```

      - OperatorTrxsCount: The number of transactions for each operator, which would get updated after each transaction and reset after each distribution has happened.

        ``` go
        [(operator_address, trxs_count), ...]
        ```

        Also storing for future reference with key as block height.

        ```go
        (block_height => [(operator_address, individual_reward), ...])
        ```

      - lastDistributionBlock: The block height at which the last distribution happened.


    - Expected behavior:

      - [x] while Gensis intialization, register the predefined operators with some egv tokens. (10 operators)

        - [x] Initialize the `OperatorTrxsCount` for each operator to 0.
        - [ ] Initialize the `distrubutionInterval` to 100 blocks through module params. so that it can be changed through governance.
        - [x] Initialize the `inflationRate` to 0.1 (default 10%) through module params. so that it can be changed through governance.
        - [x] add export genesis to the module.

      - [x] Increment the `OperatorTrxsCount` for each operator after each transaction. Reset the count after each distribution. Currently it is using store to track the `OperatorTrxsCount` but can be updated to use memStore.

        - [x] use AnteHandler to update the `OperatorTrxsCount` after each transaction.
        - [ ] Can use MsgTypeURL to check the type of transaction and update the `OperatorTrxsCount`.

      - [x] Calculate & Mint new tokens and distribute them to the operators.

        - [x] Distribute the rewards to the operators at `distributionInterval` blocks, at `EndBlocker`.

          ```go
          if current_height% distributionInterval == 0:
              for operator in operators:
                  individual_reward = (operator_trxs_count / total_network_trxs) * (inflation_rate * total_supply)
                  operator_trxs_count = 0
                  mint_tokens(operator, individual_reward)
          ```

    - [x] Permissions:
      - x/bank
      - x/mint
      - x/auth

2. x/egov: Transaction-Weighted Governance
    - [ ] Permissions:
      - [ ] x/gov
    - [ ] Expected behavior:
      - [ ] Create a `Proposal` with new InflationRate, distributionInterval, Or operators.
      - [ ] Store each `Proposal` with `snapshot_trxs_count` and `total_network_trxs`.
      - [ ] A params change proposal to modify the `InflationRate`, `distributionInterval`, and `operators`.
        - [ ] Only the governance module can modify the `InflationRate`.
        - [ ] Only the governance module can modify the `distributionInterval`.
        - [ ] Only the governance module can modify the `operators`.
      - [ ] The voting weight should be determined by the `operator_trxs_count` of the operator.

        ```go
        vote_weight = individual_address_transactions / total_network_transactions
        ```

      - [ ] The proposal should be accepted if the voting weight is greater than 50%
        - Custom Tally that calculates the voting weight.

### Assumptions

1. inflation_rate = 0.1 (default 10%) and is for a given `distributionInterval`.
2. total_network_transactions = total number of transactions in the network since genesis.

### Refrences
1. [Optimistic Execution](https://docs.cosmos.network/main/build/rfc/rfc-005-optimistic-execution)
2. [Cosmos-sdk](https://docs.cosmos.network/)