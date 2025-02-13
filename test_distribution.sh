#!/bin/bash

set -e  # Exit immediately if a command fails
set -o pipefail  # Exit if a command in a pipeline fails

# Configuration
CHAIN_BIN="chikkud"  # Replace with actual binary name
DENOM="egv"
BLOCK_WAIT_TIME=2  # Seconds to wait per block
TOTAL_TXS=10
NUM_OPERATORS=10
REWARD_DISTRIBUTION_BLOCK_INTERVAL=100  # Adjust based on actual chain params

echo "Starting reward distribution test..."

# Get CurrentBalance of each operator
for i in $(seq 1 $NUM_OPERATORS); do
  OPERATOR_ADDR=$($CHAIN_BIN keys show operator$i -a)
  OPERATOR_BALANCE=$($CHAIN_BIN query bank balances $OPERATOR_ADDR --output json | jq -r '.balances[] | select(.denom=="egv") | .amount')
  echo "Operator$i balance: $OPERATOR_BALANCE"egv
done

# Do some transactions to accumulate rewards
# Each of the 10 operators sends 100 transactions to the target address before sleeping
# TARGET_ADDR="chikku1ea8fzwncryh08nwxvltuvl02gnsvkxu39pcx2e"

# for j in $(seq 1 $TOTAL_TXS); do
#   for i in $(seq 1 $NUM_OPERATORS); do
#     OPERATOR_ADDR=$($CHAIN_BIN keys show operator$i -a)
#     $CHAIN_BIN tx bank send $OPERATOR_ADDR $TARGET_ADDR 1$DENOM --from operator$i -y &
#   done
#   wait  # Wait for all background transactions to complete
#   sleep $BLOCK_WAIT_TIME  # Sleep only after all operators have sent their transactions
# done