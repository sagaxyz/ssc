#!/bin/bash

#from container environment vars
CHAINID=$CHAINID
MONIKER=$MONIKER
LOGLEVEL=$LOGLEVEL
KEYNAME=$KEYNAME
TRACE=$TRACE
S3BUCKETROOT=`echo $S3BUCKET|awk -F "/gentx" '{print $1}'`
S3CONFIGBUCKET=${S3BUCKETROOT}/config
S3BUCKET=$S3BUCKET/$CHAINID
QUORUM_COUNT=$QUORUM_COUNT
SLEEPTIME=${SLEEPTIME:-5} # if unset, set to 5 seconds
EXTERNAL_ADDRESS=$EXTERNAL_ADDRESS
NODE_MNEMONIC=$NODE_MNEMONIC
VALIDATOR_MNEMONIC=$VALIDATOR_MNEMONIC
KEYRING=${KEYRING:-file}
TESTKEYRING="test"
KEYALGO=${KEYALGO:-"secp256k1"}
SPCDIR="/root/.ssc"
KEYRINGDIR="/root/.ssc/keyring-test"
CONFIGDIR=${SPCDIR}/config
CONFIG_TOML=${SPCDIR}/config/config.toml
TMPDIR=/tmp/${KEYNAME}
KEYPASSWD=${KEYPASSWD:-"passw0rdK3y"}

SESSION_STAMP=SPCD_`date +%m%d%Y%H%M%S`
LOGDIR=/tmp
LOGFILE=${LOGDIR}/${SESSION_STAMP}.log
ERRFILE=${LOGDIR}/${SESSION_STAMP}.err

Logger()
{
	MSG=$1
	echo "`date` $MSG" >> $LOGFILE
	echo "`date` $MSG"
}

CheckRetcode()
{
	# ERRTYPE 1 = HARD ERROR (Exit script), ERRTYPE <> 1 = SOFT ERROR (Report and Continue)
	local RETCODE=$1
	local ERRTYPE=$2
	local MSG=$3
	if [ $RETCODE -ne 0 ];
	then
		if [ $ERRTYPE -eq 1 ];
		then
			Logger "$MSG"
			exit 1
		else
			Logger "$MSG"
		fi
	else
		Logger "Return code was $RETCODE. Success!"
	fi
}

ValidateEnvVar()
{
  local ENVVAR=$1
  Logger "Validating environment variable $ENVVAR"
  local EXITIFUNSET=${2:-1}  # exit if env var is not set. Pass 1 for true, 0 for false i.e. if 0, script will continue executing. Default: True (exit)
  local ECHOVAL=${3:-1} # echo the value of the variable in a log entry. Pass 1 = true, 0 = false. Default: True (will echo)
  if [[ -z ${!ENVVAR} ]];
  then
    Logger "Environment variable $ENVVAR is not set"
    if [ $EXITIFUNSET -eq 1 ];
    then
      Logger "Exiting in error as environment variable $ENVVAR is not set"
      exit 1
    else
      Logger "Continuing even though environment variable $ENVVAR is not set"
    fi
  fi
  if [ $ECHOVAL -eq 1 ];
  then
    Logger "Environment variable $ENVVAR is set to ${!ENVVAR}"
  fi
  Logger "Finished validating environment variable $ENVVAR"
}

ValidateAndEchoEnvVars()
{
  Logger "Starting function ValidateAndEchoEnvVars"
  ValidateEnvVar CHAINID
  ValidateEnvVar MONIKER
  ValidateEnvVar LOGLEVEL
  ValidateEnvVar KEYNAME
  ValidateEnvVar TRACE
  ValidateEnvVar GENESIS_TIME 0 1
  ValidateEnvVar S3BUCKET
  ValidateEnvVar QUORUM_COUNT
  ValidateEnvVar SLEEPTIME 0 1
  ValidateEnvVar EXTERNAL_ADDRESS
  ValidateEnvVar VALIDATOR_ADDRESSES
  ValidateEnvVar VALIDATOR_MNEMONIC 1 0
  ValidateEnvVar AWS_ACCESS_KEY_ID 1 0
  ValidateEnvVar AWS_SECRET_ACCESS_KEY 1 0
  ValidateEnvVar AWS_DEFAULT_REGION 1 1
  ValidateEnvVar KEYRING
  ValidateEnvVar KEYALGO
  ValidateEnvVar KEYPASSWD
  ValidateEnvVar SPCDIR
  ValidateEnvVar KEYRINGDIR
  ValidateEnvVar CONFIG_TOML
  ValidateEnvVar TMPDIR
  ValidateEnvVar CONFIGDIR
  ValidateEnvVar S3CONFIGBUCKET
  ValidateEnvVar NODE_MNEMONIC 1 0
  Logger "Exiting function ValidateAndEchoEnvVars"
}

CheckIfRestarting()
{
  Logger "Starting function CheckIfRestarting"
  if [[ -d $SPCDIR && -f $SPCDIR/config/addrbook.json ]]; then
    # If SPCDIR exists we just start the sscd
    CURRENT_VALIDATOR_ADDRESS=`echo $KEYPASSWD | sscd keys show -a $KEYNAME`
    export CURRENT_VALIDATOR_ADDRESS
    ValidateEnvVar CURRENT_VALIDATOR_ADDRESS
    UpdateExternalAddress
    UpdatePersistentPeers
    sscd start $TRACE --log_level $LOGLEVEL
    exit 0
  fi
  Logger "Exiting function CheckIfRestarting"
}

UpdateExternalAddress()
{
  Logger "Starting function UpdateExternalAddress"
  sed -i "s/external_address = \".*/external_address = \"$EXTERNAL_ADDRESS:26656\"/g" $CONFIG_TOML
  local BACKUPLOC=${S3CONFIGBUCKET}/${CHAINID}_${CURRENT_VALIDATOR_ADDRESS}
  echo $EXTERNAL_ADDRESS:26656 > $TMPDIR/addr.info
  aws s3 cp ${TMPDIR}/addr.info ${BACKUPLOC}/addr.info
  Logger "Exiting function UpdateExternalAddress"
}

UpdatePersistentPeers()
{
  Logger "Starting function UpdatePersistentPeers"
  for file in `ls $TMPDIR/*.ctrl`
  do
    GENTX_FILE=`cat $file|awk -F":" '$0 ~ /gentxfile/ {print $2}'`
    VALIDATOR_ADDRESS=`cat $file|awk -F":" '$0 ~ /account/ {print $2}'`
    CheckValidatorInList $VALIDATOR_ADDRESS
    RETCODE=$?
    CheckRetcode $RETCODE 1 "Control file validator address $VALIDATOR_ADDRESS is not in the validators list"
    if [[ $VALIDATOR_ADDRESS != $CURRENT_VALIDATOR_ADDRESS ]];
    then
      local ADDRLOC=${S3CONFIGBUCKET}/${CHAINID}_${VALIDATOR_ADDRESS}/addr.info
      set +e
      aws s3 cp ${ADDRLOC} ${TMPDIR}/
      RETCODE=$?
      set -e
      if [[ ${RETCODE} -eq 0 ]];
      then
        local id=`jq .body.memo $TMPDIR/$GENTX_FILE |awk -F"@" '{gsub("\"","",$0);print $0}' | cut -d '@' -f 1`
        IPADD+=$id@
        IPADD+=$(cat ${TMPDIR}/addr.info)
        IPADD+=","
      else
        IPADD+=`jq .body.memo $TMPDIR/$GENTX_FILE |awk -F"@" '{gsub("\"","",$0);print $0}'`
        IPADD+=","
      fi
    fi
  done
  # Add persistent peers
  if [ -n $IPADD ];
  then
    Logger "Gathered persistent peers as $IPADD"
    sed -i "s/persistent_peers = \".*/persistent_peers = \"$IPADD\"/g" $CONFIG_TOML
    Logger "Successfully modified $CONFIG_TOML with persistent_peers"
  fi
  Logger "Exiting function UpdatePersistentPeers"
}

RestartFromS3Config()
{
  set +e
  Logger "Starting function RestartFromS3Config"
  if [ -z "${CURRENT_VALIDATOR_ADDRESS}" ];
  then
    Logger "Environment variable CURRENT_VALIDATOR_ADDRESS is not set. Exiting function RestartFromS3Config"
    return 0
  fi  
  Logger "Checking if this ssc has a prior S3 config backup"
  local BACKUPLOC=${S3CONFIGBUCKET}/${CHAINID}_${CURRENT_VALIDATOR_ADDRESS}
  Logger "Checking if config bucket $BACKUPLOC exists"
  aws s3 ls ${BACKUPLOC}
  RETCODE=$?
  if [[ ${RETCODE} -eq 0 ]]; then
    Logger "Copying config files from S3 bucket ${BACKUPLOC} for to ${CONFIGDIR}"
    rm -r -f ${CONFIGDIR}/
    aws s3 cp ${BACKUPLOC}/ ${CONFIGDIR}/ --recursive
    RETCODE=$?
    if [[ ${RETCODE} -eq 0 ]]; then
      UpdateExternalAddress
      UpdatePersistentPeers
      GENTX_COUNT=`ls ${CONFIGDIR}/gentx/gentx-*.json | wc -l`
      Logger "Found $GENTX_COUNT gentx files in ${CONFIGDIR}/gentx/ directory"
      if [[ $GENTX_COUNT -eq $QUORUM_COUNT ]];
      then
        Logger "There are $GENTX_COUNT gentx files in ${CONFIGDIR}/gentx/ which equals the required quorum of $QUORUM_COUNT. Restarting sscd"
        sscd start $TRACE --log_level $LOGLEVEL
        Logger "Started sscd successfully using backed up config at ${BACKUPLOC}"
        exit 0
      else
        Logger "There are $GENTX_COUNT gentx files in ${CONFIGDIR}/gentx/ which does not equal the required quorum of $QUORUM_COUNT"
      fi
    fi
  fi
  Logger "No previous config was found, or the config is not valid. Will need to recreate genesis artifacts"
  Logger "Exiting function RestartFromS3Config"
  set -e
  return 0
}


InitTmpDir()
{
  Logger "Starting function InitTmpDir"
  rm -r -f $TMPDIR
  RETCODE=$?
  CheckRetcode $RETCODE 1 "Could not delete the temporary folder $TMPDIR from a previous run. Exiting"
  mkdir -p $TMPDIR
  RETCODE=$?
  CheckRetcode $RETCODE 1 "Could not create the temporary folder $TMPDIR. Exiting"
  Logger "Exiting function InitTmpDir"
}


# used to exit on first error (any non-zero exit code)
set -e

ValidateDependencies()
{
  Logger "Validating dependencies"
  # validate dependencies are installed
  command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }
  Logger "Finished validating all dependencies"
}

CopyFilesFromS3Bucket()
{
  Logger "Starting function CopyFilesFromS3Bucket"
  set +e
  Logger "Downloading validator control, genesis and gentx files from S3 bucket $S3BUCKET"
  aws s3 cp $S3BUCKET/ $TMPDIR/ --recursive
  RETCODE=$?
  CheckRetcode $RETCODE 0 "No files exist on S3 bucket $S3BUCKET"
  set -e
  Logger "Exiting function CopyFilesFromS3Bucket"
}

CreateKeyringDirectory()
{
  set -e
  Logger "Starting function CreateKeyringDirectory"
  if [[ ! -d "$KEYRINGDIR" ]]; then
  ### if $KEYRINGDIR exists then just run sscd###
    Logger "Creating $SPCDIR related data"
    # Set client config

    Logger "Creating $TESTKEYRING keyring-backend temporarily"
    sscd config keyring-backend test
    
    sscd config chain-id $CHAINID
    Logger "Created $CHAINID config"
    # if $KEYNAME exists it should be deleted
    echo "$VALIDATOR_MNEMONIC" | sscd keys add $KEYNAME --recover --keyring-backend test --algo $KEYALGO
    RETCODE=$?
    CheckRetcode $RETCODE 1 "Unable to add a key from the provided mnemonic. Exiting"

    # Export the private key
    echo $KEYPASSWD | sscd keys export $KEYNAME > $KEYNAME.priv
    RETCODE=$?
    CheckRetcode $RETCODE 1 "Unable to export the private key $KEYNAME. Exiting"
    # Delete the key created in test keyring
    sscd keys delete $KEYNAME --keyring-backend $TESTKEYRING -y
    RETCODE=$?
    CheckRetcode $RETCODE 1 "Unable to delete key $KEYNAME key from $TESTKEYRING. Exiting"
    # Set client config
    if [[ "$KEYRING" == "file" ]]; then
      Logger "Using $KEYRING as keyring backend"
      Logger "Setting additional options to use this keyring backend"
      sscd config keyring-backend $KEYRING
      RETCODE=$?
      CheckRetcode $RETCODE 1 "Unable to set keyring-backend to $KEYRING. Exiting"
      echo $KEYPASSWD | sscd keys show -a me 2>> $ERRFILE
      RETCODE=$?
      if [ $RETCODE -ne 0 ];
      then
        # Add a dummy key to initialize file keyring-backend
        Logger "Adding a dummy key to keyring $KEYRING"
        (echo $KEYPASSWD;echo $KEYPASSWD) | sscd keys add me
        RETCODE=$?
        CheckRetcode $RETCODE 1 "Unable to initialize keyring-backend to $KEYRING. Adding a dummy key failed. Exiting"
        echo $KEYPASSWD | sscd keys show me
        RETCODE=$?
        CheckRetcode $RETCODE 1 "Unable to retrieve key from $KEYRING. Exiting"
      else
        Logger "Dummy key already exists in keyring $KEYRING. Skipping recreate"
      fi
      
      # Now import the private key created originally using test keyring
      echo $KEYPASSWD | sscd keys show -a $KEYNAME 2>> $ERRFILE
      RETCODE=$?
      if [ $RETCODE -ne 0 ];
      then
        (echo $KEYPASSWD; echo $KEYPASSWD) | sscd keys import $KEYNAME $KEYNAME.priv
        RETCODE=$?
        CheckRetcode $RETCODE 1 "Unable to import $KEYNAME into $KEYRING. Exiting"
        Logger "Created $KEYNAME key"
        Logger "Deleting dummy key created earlier"
        (echo $KEYPASSWD; echo $KEYPASSWD) | sscd keys delete me -y >> $LOGFILE 2>> $ERRFILE
        RETCODE=$?
        CheckRetcode $RETCODE 1 "Unable to delete dummy key. Exiting"
      else
        Logger "Key $KEYNAME already exists in keyring $KEYRING. Skipping re-import"
      fi
      if [ -s $KEYNAME.priv ];
      then
        Logger "Deleting key file $KEYNAME.priv created earlier"
        rm -f $KEYNAME.priv >> $LOGFILE 2>> $ERRFILE
        RETCODE=$?
        CheckRetcode $RETCODE 1 "Unable to delete key file $KEYNAME.priv. Exiting"
      fi
    else
      Logger "Using $KEYRING as keyring backend"
      sscd config keyring-backend $KEYRING
    fi

    # Set moniker and chain-id for ssc (Moniker can be anything, chain-id must be an integer)
    if [ ! -s $CONFIGDIR/genesis.json ];
    then
      echo "$NODE_MNEMONIC" | sscd init $MONIKER --chain-id $CHAINID --recover
      Logger "Init executed with moniker $MONIKER and chain-id $CHAINID"
    fi

    # Change parameter token denominations to tsaga
    echo "$(cat $HOME/.ssc/config/genesis.json)" "$(cat $HOME/defaults.genesis.json)" | jq --slurp 'reduce .[] as $item ({}; . * $item)' > $HOME/.ssc/config/tmp_genesis.json && mv $HOME/.ssc/config/tmp_genesis.json $HOME/.ssc/config/genesis.json

    # Update min gas price in app.toml
    sed -i 's/minimum-gas-prices = "0stake"/minimum-gas-prices = "0tsaga"/g' $HOME/.ssc/config/app.toml
    # disable produce empty block
    if [[ "$OSTYPE" == "darwin"* ]]; then
      sed -i '' 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.ssc/config/config.toml
      sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $HOME/.ssc/config/config.toml
      sed -i '' 's/addr_book_strict = true/addr_book_strict = false/g' $HOME/.ssc/config/config.toml
      sed -i '' 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' $HOME/.ssc/config/config.toml
    else
      sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $HOME/.ssc/config/config.toml
      sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' $HOME/.ssc/config/config.toml
      sed -i 's/addr_book_strict = true/addr_book_strict = false/g' $HOME/.ssc/config/config.toml
      sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' $HOME/.ssc/config/config.toml
    fi

    if [[ $1 == "pending" ]]; then
      if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.ssc/config/config.toml
        sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.ssc/config/config.toml
      else
        sed -i 's/create_empty_blocks_interval = "0s"/create_empty_blocks_interval = "30s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' $HOME/.ssc/config/config.toml
        sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' $HOME/.ssc/config/config.toml
      fi
    fi

    Logger "All configuration updates completed"
    # Allocate genesis accounts (cosmos formatted addresses)
    CV_ADDRESS=`echo $KEYPASSWD | sscd keys show -a $KEYNAME`
    EXIST_CNT=`cat ${CONFIGDIR}/genesis.json | jq '.app_state.auth.accounts[].base_account.address'|grep -c $CV_ADDRESS`
    if [ $EXIST_CNT -eq 0 ];
    then
      echo $KEYPASSWD | sscd add-genesis-account $KEYNAME 100000000000000000000000000tsaga,100000000stake --keyring-backend $KEYRING
      Logger "Added genesis account $KEYNAME"
    else
      Logger "Validator $CV_ADDRESS with key $KEYNAME already exists in the genesis file. Skipping"
    fi
    
    # Sign genesis transaction
    # sscd gentx $KEYNAME 1000000000000000000000tsaga --keyring-backend $KEYRING --chain-id $CHAINID
    if [ ! -d ${CONFIGDIR}/gentx ];
    then
      echo $KEYPASSWD | sscd gentx $KEYNAME 100000000stake --keyring-backend $KEYRING --chain-id $CHAINID --ip $EXTERNAL_ADDRESS
    else
      Logger "Gentx record already exists"
    fi
    
    Logger "Signed gentx using keyname $KEYNAME and keyring $KEYRING config"
    ## In case you want to create multiple validators at genesis
    ## 1. Back to `sscd keys add` step, init more keys
    ## 2. Back to `sscd add-genesis-account` step, add balance for those
    ## 3. Clone this ~/.ssc home directory into some others, let's say `~/.clonedsscd`
    ## 4. Run `gentx` in each of those folders
    ## 5. Copy the `gentx-*` folders under `~/.clonedsscd/config/gentx/` folders into the original `~/.ssc/config/gentx`
    AddTeamAddresses

    if [[ $1 == "pending" ]]; then
      Logger "pending mode is on, please wait for the first block committed."
    fi
  fi
  Logger "Setting current validator address"
  CURRENT_VALIDATOR_ADDRESS=`echo $KEYPASSWD | sscd keys show -a $KEYNAME`
  export CURRENT_VALIDATOR_ADDRESS
  ValidateEnvVar CURRENT_VALIDATOR_ADDRESS
  Logger "Setting S3 validator control file"
  S3CTRLFILE=Validator-${CURRENT_VALIDATOR_ADDRESS}.ctrl
  export S3CTRLFILE
  ValidateEnvVar S3CTRLFILE
  Logger "Exiting function CreateKeyringDirectory"
  return 0
}

CreateS3ControlFile()
{
  set -e
  Logger "Starting function CreateS3ControlFile"
  if [[ ! -s $TMPDIR/$S3CTRLFILE ]]; then
  # Create the S3 control file
    echo "account:`echo $KEYPASSWD | sscd keys show -a $KEYNAME`" > $TMPDIR/$S3CTRLFILE #Overwrite any existing files
    GENTX_FILE=`ls ~/.ssc/config/gentx/gentx-* | awk -F"/" '{print $NF}'`
    echo "gentxfile:$GENTX_FILE" >> $TMPDIR/$S3CTRLFILE
    # Copy the gentx, genesis and control file to aws s3
    Logger "Copying files to S3"
    aws s3 cp ~/.ssc/config/gentx/ $S3BUCKET/ --recursive --exclude "*" --include "gentx-*"
    aws s3 cp $TMPDIR/ $S3BUCKET/ --recursive --exclude "*" --include "${S3CTRLFILE}"
  fi
  Logger "Exiting function CreateS3ControlFile"
  return 0
}


CheckQuorum()
{
  Logger "Starting function CheckQuorum"
  Logger "Checking if we have a quorum of gentx records. Current quorum is $QUORUM_COUNT"
  while true
  do
    GENTX_COUNT=`ls $TMPDIR/gentx-*.json | wc -l`
    CTRL_COUNT=`ls $TMPDIR/Validator*.ctrl | wc -l`
    Logger "Currently found $GENTX_COUNT gentx and $CTRL_COUNT validator control files"
    if [[ $GENTX_COUNT -lt $QUORUM_COUNT || $CTRL_COUNT -lt $QUORUM_COUNT ]];
    then
      set +e
      Logger "Downloading gentx and control files from S3 bucket $S3BUCKET"
      aws s3 sync $S3BUCKET/ $TMPDIR/
      set -e
      Logger "Sleeping for $SLEEPTIME seconds"
      sleep $SLEEPTIME
    else
      Logger "Found $GENTX_COUNT gentx files and $CTRL_COUNT control files which satisfies quorum of $QUORUM_COUNT"
      break
    fi
  done
  Logger "Exiting function CheckQuorum"
  return 0
}

CheckValidatorInList()
{
  local address_to_check=$1
  Logger "Checking if validator address $address_to_check is allowed"
  for ADDRESS in `echo $VALIDATOR_ADDRESSES | tr " " " "`;
  do
    if [[ $ADDRESS = $address_to_check ]];
    then
      Logger "Success - validator address matches our expected address. Continuing"
      return 0
    fi;
  done  
  Logger "Failure - validator address $address_to_check does not match any of our expected addresses $VALIDATOR_ADDRESSES"
  return 1
}

AddValidatorAccountsInfo()
{
  Logger "Starting function AddValidatorAccountsInfo"
  for file in `ls $TMPDIR/*.ctrl`
  do
    Logger "Working with control file $file"
    VALIDATOR_ADDRESS=`cat $file|awk -F":" '$0 ~ /account/ {print $2}'`
    GENTX_FILE=`cat $file|awk -F":" '$0 ~ /gentxfile/ {print $2}'`
    Logger " Evaluating validator $VALIDATOR_ADDRESS and gentx $GENTX_FILE"

    CheckValidatorInList $VALIDATOR_ADDRESS
    RETCODE=$?
    CheckRetcode $RETCODE 1 "Control file validator address $VALIDATOR_ADDRESS is not in the validators list"

    # Validate validator address in gentx file
    if [[ `jq '.body.messages | length' $TMPDIR/$GENTX_FILE` -ne 1 ]];
    then
      Logger "$TMPDIR/$GENTX_FILE expected to contain one element array as .body.message but got `jq '.body.messages | length' $TMPDIR/$GENTX_FILE`"
      exit 1
    fi

    if [[ `jq '.body.messages[0].delegator_address' $TMPDIR/$GENTX_FILE` != \"$VALIDATOR_ADDRESS\" ]];
    then
      Logger "Gentx file address `jq '.body.messages[0].delegator_address' $TMPDIR/$GENTX_FILE` differs from control file address $VALIDATOR_ADDRESS"
      exit 1
    fi

    if [ $VALIDATOR_ADDRESS != $CURRENT_VALIDATOR_ADDRESS ];
    then
      # Add the other validators' address to genesis
      set +e
      EXIST_CNT=`cat ${CONFIGDIR}/genesis.json | jq '.app_state.auth.accounts[].base_account.address'|grep -c $VALIDATOR_ADDRESS`
      set -e
      if [ $EXIST_CNT -eq 0 ];
      then
        Logger "Adding other validators' address $VALIDATOR_ADDRESS to genesis file"
        sscd add-genesis-account $VALIDATOR_ADDRESS 100000000000000000000000000tsaga,100000000stake --keyring-backend $KEYRING
      else
        Logger "Validator $VALIDATOR_ADDRESS already exists in genesis file. Skipping."
      fi
      # Copy the gentx file to the config/gentx directory
      Logger "Copying other validators' gentx $GENTX_FILE to ~/.ssc/config/gentx/"
      cp -f $TMPDIR/$GENTX_FILE ~/.ssc/config/gentx/
      local ADDRLOC=${S3CONFIGBUCKET}/${CHAINID}_${VALIDATOR_ADDRESS}/addr.info
      set +e
      aws s3 cp ${ADDRLOC} ${TMPDIR}/
      RETCODE=$?
      set -e
      if [[ ${RETCODE} -eq 0 ]];
      then
        IPADD+=$(cat ${TMPDIR}/addr.info)
        IPADD+=","
      else
        IPADD+=`jq .body.memo $TMPDIR/$GENTX_FILE |awk -F"@" '{gsub("\"","",$0);print $0}'`
        IPADD+=","
      fi
    else
      Logger "This validator is $VALIDATOR_ADDRESS and current is $CURRENT_VALIDATOR_ADDRESS"
    fi
  done
  # Add external_address
  Logger "Configuring config.toml with the external address $EXTERNAL_ADDRESS"
  sed -i 's@external_address = ""@external_address = "'"$EXTERNAL_ADDRESS"':26656"@g' $CONFIG_TOML
  # Add persistent peers
  if [ -n $IPADD ];
  then
    Logger "Gathered persistent peers as $IPADD"
    sed -i "s/persistent_peers = \"\"/persistent_peers = \"$IPADD\"/g" $CONFIG_TOML
    Logger "Successfully modified $CONFIG_TOML with external_address and persistent_peers"
  fi
  Logger "Exiting function AddValidatorAccountsInfo"
  return 0
}

UploadGenesisToS3Bucket()
{
  Logger "Starting function UploadGenesisToS3Bucket"
  # Now upload the genesis file to S3
  GENESIS_FILE=`ls $SPCDIR/config/genesis.json | awk -F"/" '{print $NF}' | awk -F"." -v VAL=$CURRENT_VALIDATOR_ADDRESS '{print $1"-"VAL".json"}'`
  cp $SPCDIR/config/genesis.json $SPCDIR/config/$GENESIS_FILE
  aws s3 cp ~/.ssc/config/ $S3BUCKET/ --recursive --exclude "*" --include "${GENESIS_FILE}"
  Logger "Successfully uploaded genesis file $GENESIS_FILE to S3"
  sleep $SLEEPTIME # Sleep before checking and getting all the genesis files from S3
  while true
  do
    GEN_COUNT=`ls $TMPDIR/genesis*.json | wc -l`
    Logger "Currently found $GEN_COUNT genesis files in S3"
    if [[ $GEN_COUNT -lt QUORUM_COUNT ]];
    then
      set +e
      Logger "Downloading genesis files from S3 bucket $S3BUCKET"
      aws s3 sync $S3BUCKET/ $TMPDIR/
      set -e
      Logger "Sleeping for $SLEEPTIME seconds"
      sleep $SLEEPTIME
    else
      Logger "Found $GEN_COUNT genesis files which satisfies quorum of $QUORUM_COUNT"
      GENESIS_FILE_TO_USE=`ls $TMPDIR/genesis*.json | sort | awk 'NR==1'`
      ls -ltr $TMPDIR/genesis*.json
      ls $TMPDIR/genesis*.json | sort
      ls $TMPDIR/genesis*.json | sort | awk 'NR==1'
      Logger "Determined the genesis file to be used will be $GENESIS_FILE_TO_USE"
      cp -f $GENESIS_FILE_TO_USE $SPCDIR/config/genesis.json
      Logger "Successfully copied genesis.json to $SPCDIR/config"
      break
    fi
  done 
  Logger "Exiting function UploadGenesisToS3Bucket"
  return 0
}

AddTeamAddresses() 
{
  Logger "Adding team accounts"
  echo $KEYPASSWD | sscd add-genesis-account saga1rdssl22ysxyendrkh2exw9zm7hvj8d2ju346g3 1000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # ashish
  echo $KEYPASSWD | sscd add-genesis-account saga16p4cejpaqpuha65hqyj85k5lx4umw7qzku37eg 1000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # konstantin
  echo $KEYPASSWD | sscd add-genesis-account saga17gk4chqd0lrkyamrxdmu62czmu0dpnemmxlymn 1000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # bogdan
  echo $KEYPASSWD | sscd add-genesis-account saga1sz83y27774xwrahwmv5afutv86grc286hcf7w5 1000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # roberto
  echo $KEYPASSWD | sscd add-genesis-account saga1gme3rzzddpf4hkdngpruz5e4739lqsyyakgu0j 1000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # roman
  echo $KEYPASSWD | sscd add-genesis-account saga1rcs5sw5yy9r04xsultcqv6tj73408qnawmlxqw 1000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # brian
  echo $KEYPASSWD | sscd add-genesis-account saga120nzke36a2s44w0f0dhndknndhu5ytyxmsmgrs 1000000000000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # jin
  echo $KEYPASSWD | sscd add-genesis-account saga129nuu6nclj6ta5r8pgvuvqc7qara706gwlpyh4 1000000000000000000000tsaga --keyring-backend $KEYRING 2>> $ERRFILE # relayer
}

CreateConfigS3Backup()
{
  set -e
  Logger "Starting function CreateConfigS3Backup"
  local BACKUPLOC=${S3CONFIGBUCKET}/${CHAINID}_${CURRENT_VALIDATOR_ADDRESS}
  if [[ -d ${CONFIGDIR} ]]; then
    Logger "Copying config files to S3 bucket ${BACKUPLOC} for backup"
    aws s3 cp ${CONFIGDIR} ${BACKUPLOC}/ --recursive
  else
    Logger "${CONFIGDIR} does not exist. This is a serious error and should never be the case."
  fi
  Logger "Exiting function CreateConfigS3Backup"
  return 0
}

### MAIN ###
ValidateAndEchoEnvVars
InitTmpDir
ValidateDependencies
CopyFilesFromS3Bucket
CheckIfRestarting
if CreateKeyringDirectory;
then
  RestartFromS3Config # if this is executed, the script will successfully exit with a return code of 0
  CreateS3ControlFile  
  if CheckQuorum;
  then
    AddValidatorAccountsInfo
    # Collect genesis tx
    cp $CONFIG_TOML $CONFIG_TOML.bak
    sscd collect-gentxs
    cp $CONFIG_TOML.bak $CONFIG_TOML
    # Run this to ensure everything worked and that the genesis file is setup correctly
    sscd validate-genesis
    if UploadGenesisToS3Bucket;
    then
      # Create the config backup
      CreateConfigS3Backup
      # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
      sscd start $TRACE --log_level $LOGLEVEL
    fi
  fi
fi
