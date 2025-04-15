#!/bin/bash
set -u

# expand globs to nothing
shopt -s nullglob

# Validate dependencies are installed
command -v jq &> /dev/null || fail "jq not installed"

function log() {
  local msg=$1

  echo "$(date) $msg" >&2
}

function fail() {
  if [ $# -gt 0 ]; then
    local msg=$1
    log "$msg"
  fi

  exit 1
}

function check_env_vars() {
  for name in "$@"; do
    local value="${!name:-}"
    if [ -z "$value" ]; then
      log "Variable $name is empty"
      return 1
    else
      return 0
    fi
  done
}

check_env_vars GOBIN
RETCODE=$?
if [ $RETCODE -eq 1 ];
then
    log "Please set all necessary environment variables before proceeding"
    exit 1
fi

# Constant vars
SSC_DIR="/tmp/ssc"
CONFIG_DIR=${SSC_DIR}/config
CONFIG_TOML=${CONFIG_DIR}/config.toml
SPCD=$GOBIN/spcd
SSCD=$GOBIN/sscd
TMP=/tmp
SLEEP_TIME=5

# Vars to be parameterized
RPC_NODE_SPC=${1:-"tcp://spc.testnet.sagarpc.io:26657"}
RPC_NODE_SSC=${2:-"tcp://localhost:26657"}
SSC_CHAIN_ID=${3:-pegasus_100-1}
SSC_STACK_OWNER=${4:-ashish}
KEYRING_BACKEND=${5:-test}
DENOM=${6:-utsaga}
GAS_LIMIT=${7:-500000}
FEES=${8:-5000}${DENOM}

function create_chainlet_stacks() {
    # Fetch the chainlet stacks data from current SPC and store it temporarily
    ${SPCD} q chainlet list-chainlet-stack --output json --node ${RPC_NODE_SPC} > ${TMP}/spc-migration-chainlet-stacks.txt
    local input_json_file=${TMP}/spc-migration-chainlet-stacks.txt
    local COUNT_STACKS=$(jq '.ChainletStacks | length' ${input_json_file})
    if [ $COUNT_STACKS -gt 0 ];
    then
        log "Total of $COUNT_STACKS chainlet stacks found to migrate"
        for i in $(seq 0 $(($COUNT_STACKS-1)));
        do
            log "Processing chainlet stack at array position $i"
            local creator=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].creator' ${input_json_file})
            local displayName=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].displayName' ${input_json_file})
            local description=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].description' ${input_json_file})
            local versionsSize=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].versions | length' ${input_json_file})
            
            # The code below is imperfect. Theoretically, there can be multiple versions of a chainlet stack
            # but that is not the case right now. Currently, we only have a single version of SagaOS unless there is 
            # an upgrade. The code below will set the last element in the versions array to be the stack image/version
            # to be used. This code will not work if there are multiple versions as we will need to tie each SPC chainlet
            # with the version it was running on, and ensure we have a corresponding chainlet stack in SSC. That is not the case now
            # even in Mainnet. So, if this situation changes, this code will need to be revisited and we will need to add all of the
            # relevant images and versions via an update-chainlet-stack call.

            for j in $(seq 0 $(($versionsSize-1)));
            do
                log "j is $j"
                local image=$(jq -r --argjson ivar ${i} --argjson jvar ${j} '.ChainletStacks[$ivar].versions[$jvar].image' ${input_json_file})
                local version=$(jq -r --argjson ivar ${i} --argjson jvar ${j} '.ChainletStacks[$ivar].versions[$jvar].version' ${input_json_file})
                local checksum=$(jq -r --argjson ivar ${i} --argjson jvar ${j} '.ChainletStacks[$ivar].versions[$jvar].checksum' ${input_json_file})
            done
            local setupFee=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].fees.setupFee' ${input_json_file})
            setupFee=$(echo $setupFee| tr -d '[:alpha:]')
            setupFee+=$DENOM
            local epochLength=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].fees.epochLength' ${input_json_file})
            local epochFee=$(jq -r --argjson ivar ${i} '.ChainletStacks[$ivar].fees.epochFee' ${input_json_file})
            epochFee=$(echo $epochFee| tr -d '[:alpha:]')
            epochFee+=$DENOM
            set -x
            $SSCD tx chainlet create-chainlet-stack $displayName "$description" $image $version $checksum $epochFee $epochLength $setupFee --from $SSC_STACK_OWNER --home $SSC_DIR --keyring-backend $KEYRING_BACKEND --chain-id $SSC_CHAIN_ID -y
            set +x
            sleep $SLEEP_TIME
        done
    else
        log "No chainlet stacks found to migrate"
    fi
}

function create_chainlets() {
    # Fetch the chainlets data from current SPC and store it temporarily
    ${SPCD} q chainlet list-chainlets --output json --node ${RPC_NODE_SPC} > ${TMP}/spc-migration-chainlets.txt
    local input_json_file=${TMP}/spc-migration-chainlets.txt
    local COUNT_CHAINLETS=$(jq '.Chainlets | length' ${input_json_file})
    if [ $COUNT_CHAINLETS -gt 0 ];
    then
        log "Total of $COUNT_CHAINLETS chainlets found to migrate"
        for i in $(seq 0 $(($COUNT_CHAINLETS-1)));
        do
            unset maintainers
            log "Processing chainlet at array position $i"
            local launcher=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].launcher' ${input_json_file})
            local chainletStackName=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].chainletStackName' ${input_json_file})
            local chainletStackVersion=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].chainletStackVersion' ${input_json_file})
            local chainId=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].chainId' ${input_json_file})
            local chainletName=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].chainletName' ${input_json_file})
            local chainId=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].chainId' ${input_json_file})
            local maintainersSize=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].maintainers | length' ${input_json_file})
            for j in $(seq 0 $(($maintainersSize-1)));
            do
                log "j is $j"
                maintainers+=$(jq -r --argjson ivar ${i} --argjson jvar ${j} '.Chainlets[$ivar].maintainers[$jvar]' ${input_json_file})
                if [ $j -ne $(($maintainersSize-1)) ];
                then
                    maintainers+=","
                fi
                export maintainers
            done
            log "Maintainers: $maintainers"
            local denom=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.denom' ${input_json_file})
            local gasLimit=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.gasLimit' ${input_json_file})
            local createEmptyBlocks=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.createEmptyBlocks' ${input_json_file})
            local dacEnable=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.dacEnable' ${input_json_file})
            local genAcctBalances=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.genAcctBalances' ${input_json_file})
            local fixedBaseFee=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.fixedBaseFee' ${input_json_file})
            local feeAccount=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].params.feeAccount' ${input_json_file})
            local status=$(jq -r --argjson ivar ${i} '.Chainlets[$ivar].status' ${input_json_file})
            set -x
            $SSCD tx chainlet launch-chainlet "$maintainers" $chainletStackName $chainletStackVersion $chainletName $denom '{"gasLimit":'$gasLimit',"genAcctBalances":"'"$genAcctBalances"'","fixedBaseFee":"'"$fixedBaseFee"'","feeAccount":"'"$feeAccount"'"}' --gas $GAS_LIMIT --from $($SSCD keys show -a $SSC_STACK_OWNER --home $SSC_DIR --keyring-backend $KEYRING_BACKEND) --fees $FEES --home $SSC_DIR --keyring-backend $KEYRING_BACKEND --chain-id $SSC_CHAIN_ID -y
            set +x
            sleep $SLEEP_TIME
        done
    else
        log "No chainlets found to migrate"
    fi
}


create_chainlet_stacks
create_chainlets