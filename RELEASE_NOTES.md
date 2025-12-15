# Saga Security Chain (SSC) Release Notes

## Version 1.0.0

### 🚀 Major Features

This release introduces a comprehensive re-architecture of the Saga Security Chain with five core modules that work together to provide a robust, scalable blockchain infrastructure for chainlet management, billing, escrow services, and liquid staking capabilities.

### 🆕 **New in This Release:**
- **x/escrow Module**: Multi-denomination escrow system with KV-based storage and predictable gas costs
- **x/chainlet Module**: Complete chainlet lifecycle management with stack-based architecture and auto-upgrade capabilities
- **x/billing Module**: Universal billing system with epoch-based billing and comprehensive history tracking
- **x/epochs Module**: Time-based event system with multiple epoch types and hook system for other modules
- **x/liquid Module**: Complete liquid staking system for tokenized delegation and reward management
- **Comprehensive Testing**: Full test suite for all module functionality
- **CLI Integration**: Complete command-line interface for all module operations
- **Tokenized Staking**: Advanced tokenization system for delegation shares and reward management
- **Multi-Denomination Support**: Flexible fee structures and token support across all modules
- **Storage Optimization**: KV-based architecture with deterministic gas costs

---

## 📦 Core Modules

The Saga Security Chain now includes five core modules that provide comprehensive blockchain infrastructure capabilities:

### 🔐 **x/escrow** - Multi-Denomination Escrow System

The escrow module provides a sophisticated fund management system with support for multiple denominations and predictable gas costs.

#### **Key Features:**
- **Multi-Denomination Support**: Each chainlet can support multiple token denominations
- **KV-Based Storage**: Prevents state bloat from dust deposits with efficient per-funder storage
- **Deterministic Gas Costs**: Predictable transaction costs regardless of number of funders
- **Pool-Based Architecture**: Independent pools per `{chainId, denom}` combination
- **Share-Based System**: Proportional share tracking for fair fund distribution

#### **CLI Commands:**
```bash
# Deposit funds into escrow
sscd tx escrow deposit <amount> <denom> --from <key> --chain-id <chain-id>

# Withdraw funds from escrow
sscd tx escrow withdraw <amount> <denom> --from <key> --chain-id <chain-id>

# Query escrow balance
sscd query escrow funder-balance <address>

# Query chainlet account information
sscd query escrow chainlet <chainlet-id> <denom>

# List all funders for a specific pool
sscd query escrow funders <chain-id> <denom>

# Query all pools for a chainlet
sscd query escrow pools <chain-id>

# Query escrow parameters
sscd query escrow params
```

#### **Storage Architecture:**
- **Chainlet Keys**: `escrow/chainlet/{chainId}`
- **Pool Keys**: `escrow/pool/{chainId}/{denom}`
- **Funder Keys**: `escrow/funder/{chainId}/{denom}/{addr}`
- **By-Funder Keys**: `escrow/byFunder/{addr}/{chainId}/{denom}`

---

### 🏗️ **x/chainlet** - Chainlet Lifecycle Management

The chainlet module manages the complete lifecycle of blockchain chainlets, from creation to upgrades and maintenance.

#### **Key Features:**
- **Chainlet Stack Management**: Create and manage reusable chainlet configurations
- **Multi-Fee Support**: Each stack can support multiple fee denominations
- **Version Management**: Sophisticated versioning system with compatibility checks
- **Auto-Upgrade Capability**: Automatic stack version upgrades
- **CCV Consumer Support**: Inter-Blockchain Communication (IBC) consumer chain support
- **Service Chainlet Support**: Special chainlets for system services
- **Genesis Validator Management**: Custom validator sets for chainlets

#### **CLI Commands:**
```bash
# Create a new chainlet stack
sscd tx chainlet create-chainlet-stack <stack-name> <description> --from <key> --chain-id <chain-id>

# Launch a new chainlet
sscd tx chainlet launch-chainlet <stack-id> <chainlet-name> <chain-id> --from <key> --chain-id <chain-id>

# Update chainlet stack
sscd tx chainlet update-chainlet-stack <stack-id> --from <key> --chain-id <chain-id>

# Update stack fees (NEW FEATURE)
sscd tx chainlet update-stack-fees <stack-id> "1000denom1,2000denom2" --from <key> --chain-id <chain-id>

# Upgrade chainlet
sscd tx chainlet upgrade-chainlet <chainlet-id> <new-version> --from <key> --chain-id <chain-id>

# Disable chainlet stack version
sscd tx chainlet disable-chainlet-stack-version <stack-id> <version> --from <key> --chain-id <chain-id>

# Query chainlet information
sscd query chainlet get-chainlet <chainlet-id>

# Query chainlet stack
sscd query chainlet get-chainlet-stack <stack-id>

# List all chainlets
sscd query chainlet list-chainlets

# Query chainlet count
sscd query chainlet get-chainlet-count

# Query chainlet parameters
sscd query chainlet params
```

#### **Chainlet Properties:**
- **Spawn Time**: When the chainlet was created
- **Launcher**: Account that launched the chainlet
- **Maintainers**: List of accounts with maintenance privileges
- **Stack Information**: Reference to the chainlet stack configuration
- **Chain Parameters**: Custom blockchain parameters (gas limits, block settings)
- **Status Tracking**: Online/offline status monitoring
- **Auto-Upgrade**: Automatic stack version upgrades
- **Genesis Validators**: Custom validator sets
- **Tags**: Categorization and metadata
- **Service Chainlet**: System service designation
- **CCV Consumer**: IBC consumer chain support

---

### 💰 **x/billing** - Universal Billing System

The billing module provides a centralized billing system with universal epoch configuration and comprehensive billing history tracking.

#### **Key Features:**
- **Universal Epoch Configuration**: Centralized billing epoch management
- **Multi-Denomination Billing**: Support for billing from different token pools
- **Comprehensive History**: Detailed billing and payout history tracking
- **Validator Payouts**: Automated validator reward distribution
- **Epoch-Based Billing**: Automatic billing at epoch boundaries
- **Failed Billing Handling**: Graceful handling of insufficient funds

#### **CLI Commands:**
```bash
# Query billing history for an account
sscd query billing get-billing-history <address> <denom>

# Query validator payout history
sscd query billing get-validator-payout-history <validator> <denom>

# Query billing parameters
sscd query billing params
```

#### **Billing Process:**
1. **Epoch Trigger**: Billing occurs at the start of each epoch
2. **Multi-Fee Support**: Tries multiple fee denominations until one succeeds
3. **Automatic Stopping**: Chainlets are stopped if billing fails
4. **History Recording**: All billing events are recorded with timestamps
5. **Validator Payouts**: Automatic distribution to validators

#### **Billing Parameters:**
- **Validator Payout Epoch**: Epoch identifier for validator payouts
- **Billing Epoch**: Universal epoch identifier for billing cycles

---

### ⏰ **x/epochs** - Time-Based Event System

The epochs module provides a generalized timing system for other modules to execute code at regular intervals.

#### **Key Features:**
- **Multiple Epoch Types**: Support for different time intervals (minute, hour, day, week)
- **Hook System**: Other modules can register epoch hooks
- **Panic Isolation**: Failed epoch hooks don't affect other modules
- **Flexible Configuration**: Customizable epoch durations and start times
- **Genesis Initialization**: Pre-configured epoch types

#### **CLI Commands:**
```bash
# Query all epoch information
sscd query epochs epoch-infos

# Query current epoch for specific identifier
sscd query epochs current-epoch <identifier>
```

#### **Default Epoch Types:**
- **Minute**: 1-minute intervals
- **Hour**: 1-hour intervals  
- **Day**: 24-hour intervals
- **Week**: 7-day intervals

#### **Epoch Information:**
- **Identifier**: Unique epoch type name
- **Duration**: Time interval between epochs
- **Current Epoch**: Current epoch number
- **Start Time**: When the current epoch began
- **Start Height**: Block height when epoch started
- **Counting Status**: Whether epoch counting has started

---

### 💧 **x/liquid** - Liquid Staking System

The liquid module provides a sophisticated liquid staking system that enables tokenization of delegation shares, allowing users to maintain liquidity while earning staking rewards.

#### **Key Features:**
- **Tokenized Delegation**: Convert delegation shares into tradeable tokens
- **Liquid Staking**: Maintain liquidity while earning staking rewards
- **Reward Management**: Automated reward collection and distribution
- **Validator Support**: Support for multiple validators with liquid staking
- **Authorization System**: Controlled tokenization with governance oversight
- **Lock Management**: Time-based locking system for tokenized shares
- **Fee Management**: Configurable fees for tokenization operations

#### **CLI Commands:**
```bash
# Tokenize delegation shares
sscd tx liquid tokenize-share <validator-address> <amount> <owner-address> --from <key> --chain-id <chain-id>

# Redeem tokenized shares back to delegation
sscd tx liquid redeem-tokens <amount> <denom> <owner-address> --from <key> --chain-id <chain-id>

# Enable tokenization for a validator
sscd tx liquid enable-tokenize-shares <validator-address> --from <key> --chain-id <chain-id>

# Disable tokenization for a validator
sscd tx liquid disable-tokenize-shares <validator-address> --from <key> --chain-id <chain-id>

# Transfer tokenize share record ownership
sscd tx liquid transfer-tokenize-share-record <record-id> <from-address> <to-address> --from <key> --chain-id <chain-id>

# Withdraw rewards for a specific tokenize share record
sscd tx liquid withdraw-tokenize-share-rewards <record-id> <owner-address> --from <key> --chain-id <chain-id>

# Withdraw all tokenize share record rewards
sscd tx liquid withdraw-all-tokenize-share-rewards <owner-address> --from <key> --chain-id <chain-id>

# Query liquid staking parameters
sscd query liquid params

# Query all liquid validators
sscd query liquid liquid-validators

# Query specific liquid validator
sscd query liquid liquid-validator <validator-address>

# Query total liquid staked tokens
sscd query liquid total-liquid-staked

# Query tokenize share records by owner
sscd query liquid tokenize-share-records-owned <owner-address>

# Query tokenize share record by ID
sscd query liquid tokenize-share-record-by-id <record-id>

# Query tokenize share record by denom
sscd query liquid tokenize-share-record-by-denom <denom>

# Query tokenize share record rewards
sscd query liquid tokenize-share-record-rewards <record-id> <owner-address>

# Query tokenize share lock information
sscd query liquid tokenize-share-lock-info <owner-address>

# Query all tokenize share records
sscd query liquid all-tokenize-share-records

# Query last tokenize share record ID
sscd query liquid last-tokenize-share-record-id

# Query total tokenized share assets
sscd query liquid total-tokenize-share-assets
```

#### **Tokenization System:**
- **Share Tokenization**: Convert delegation shares into tradeable tokens
- **Denom Generation**: Automatic generation of unique token denominations
- **Record Management**: Comprehensive tracking of tokenized share records
- **Ownership Transfer**: Transfer ownership of tokenized share records
- **Authorization Control**: Governance-controlled tokenization permissions

#### **Liquid Staking Features:**
- **Validator Support**: Support for multiple validators with liquid staking
- **Liquid Shares**: Track liquid shares per validator
- **Staking Caps**: Configurable caps on liquid staking per validator
- **Reward Distribution**: Automated reward collection and distribution
- **Lock Management**: Time-based locking system for tokenized shares

#### **Reward Management:**
- **Automatic Collection**: Automated collection of staking rewards
- **Reward Withdrawal**: Manual withdrawal of accumulated rewards
- **Bulk Operations**: Withdraw rewards for all tokenized share records
- **Fee Distribution**: Distribution of fees to the community pool

#### **Storage Architecture:**
- **Tokenize Share Records**: `liquid/tokenize-share-record/{id}`
- **Liquid Validators**: `liquid/liquid-validator/{validator-address}`
- **Total Liquid Staked**: `liquid/total-liquid-staked`
- **Tokenize Share Locks**: `liquid/tokenize-share-lock/{owner}`
- **Authorization Queue**: `liquid/tokenize-share-auth-queue/{timestamp}`

#### **Default Parameters:**
- **Global Liquid Staking Cap**: Configurable global cap on liquid staking
- **Validator Liquid Staking Cap**: Per-validator caps on liquid staking
- **Tokenization Fees**: Configurable fees for tokenization operations
- **Lock Duration**: Time-based locking for tokenized shares

#### **Governance Integration:**
- **Parameter Updates**: Governance-controlled parameter updates
- **Authorization Management**: Governance control over tokenization permissions
- **Fee Management**: Governance-controlled fee structures
- **Cap Management**: Governance-controlled staking caps

#### **Integration Features:**
- **Staking Integration**: Seamless integration with the staking module
- **Distribution Integration**: Integration with the distribution module for rewards
- **Bank Integration**: Integration with the bank module for token operations
- **Event System**: Comprehensive event emission for all operations

---

## 🔧 Technical Improvements

### **Storage Optimization**
- **KV-Based Architecture**: Efficient storage patterns prevent state bloat
- **Compact Key Design**: Single-byte prefixes for optimal storage
- **Deterministic Gas Costs**: Predictable transaction costs

### **Multi-Denomination Support**
- **Flexible Fee Structures**: Support for multiple token types
- **Pool Isolation**: Independent pools per denomination
- **Cross-Denomination Operations**: Seamless multi-token support

### **Enhanced CLI Experience**
- **Comprehensive Commands**: Full CRUD operations for all modules
- **Rich Query Interface**: Detailed information retrieval
- **Parameter Management**: Easy configuration updates

### **Robust Error Handling**
- **Graceful Degradation**: Failed operations don't crash the system
- **Detailed Logging**: Comprehensive event tracking
- **Recovery Mechanisms**: Automatic retry and fallback systems

---

## 🚀 Getting Started

### **Prerequisites**
- Go 1.21+
- Cosmos SDK v0.50
- CometBFT v0.38

### **Installation**
```bash
git clone https://github.com/sagaxyz/ssc.git
cd ssc
make build
```

### **Quick Start**
```bash
# Initialize the chain
./build/sscd init testchain --chain-id ssc

# Create a chainlet stack
./build/sscd tx chainlet create-chainlet-stack "MyStack" "Test stack" --from alice --chain-id ssc

# Launch a chainlet
./build/sscd tx chainlet launch-chainlet stack1 "MyChainlet" "chainlet-1" --from alice --chain-id ssc

# Deposit funds
./build/sscd tx escrow deposit 1000stake stake --from alice --chain-id ssc

# Tokenize delegation shares
./build/sscd tx liquid tokenize-share <validator-address> 1000000stake $(sscd keys show alice -a) --from alice --chain-id ssc

# Query chainlet status
./build/sscd query chainlet get-chainlet chainlet-1

# Query liquid validators
./build/sscd query liquid liquid-validators
```

---

## 📊 Performance Characteristics

### **Gas Efficiency**
- **Predictable Costs**: Gas costs don't scale with number of funders
- **Optimized Storage**: Efficient key-value storage patterns
- **Batch Operations**: Support for bulk operations

### **Scalability**
- **Multi-Denomination**: Support for unlimited token types
- **Pool Isolation**: Independent scaling per denomination
- **Efficient Queries**: Fast data retrieval with indexed storage

### **Reliability**
- **Fault Tolerance**: Graceful handling of failed operations
- **State Consistency**: Atomic operations ensure data integrity
- **Recovery Mechanisms**: Automatic retry and fallback systems

---

## 🔄 Migration Notes

### **From Previous Versions**
- **Escrow Re-architecture**: Complete rewrite of escrow storage system
- **Universal Billing**: Centralized epoch configuration
- **Multi-Fee Support**: Enhanced fee management capabilities
- **Liquid Module**: New liquid staking system for tokenized delegation and reward management
- **Script Updates**: New `escrow.sh` replaces `escrow-chainlet-restart.sh`

### **Breaking Changes**
- **Storage Format**: New KV-based storage patterns
- **CLI Commands**: Updated command structure and parameters
- **Configuration**: New parameter structures for all modules

---

## 🛠️ Development Tools

### **Testing**
```bash
# Run all tests
go test ./...

# Run specific module tests
go test ./x/escrow/...
go test ./x/chainlet/...
go test ./x/billing/...
go test ./x/epochs/...
go test ./x/liquid/...

# Run liquid module keeper tests
go test ./x/liquid/keeper/... -v
```

### **Integration Testing**
```bash
# Run environment setup
./scripts/ci/prepare-env.sh

# Run happypath tests
./scripts/happypath.sh

# Run escrow tests
./scripts/escrow.sh
```

---

## 📚 Documentation

- **Module Documentation**: Each module includes comprehensive README files
- **API Reference**: Full gRPC and REST API documentation
- **CLI Reference**: Complete command-line interface documentation
- **Architecture Guide**: Detailed system architecture documentation

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 🆘 Support

- **Documentation**: [Validator/Node docs](https://nodedocs.saga.xyz/)
- **Issues**: [GitHub Issues](https://github.com/sagaxyz/ssc/issues)

---

*For more information, visit [Saga Protocol](https://saga.xyz)*
