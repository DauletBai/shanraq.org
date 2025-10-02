# Shanraq.org Transaction Core Implementation Guide

## Overview

This guide provides step-by-step instructions for implementing the transaction core system for Shanraq.org's fintech platform. The implementation covers all aspects of the transaction processing system.

## Prerequisites

### 1. System Requirements
- **Operating System**: Linux (Ubuntu 20.04+ recommended)
- **Memory**: Minimum 16GB RAM, 32GB recommended
- **Storage**: Minimum 500GB SSD storage
- **Network**: High-speed internet connection with load balancer

### 2. Software Dependencies
- **Node.js**: Version 18.0 or higher
- **PostgreSQL**: Version 14 or higher
- **Redis**: Version 6.0 or higher
- **Kafka**: Version 3.0 or higher (or RabbitMQ/NATS)
- **Docker**: Version 20.10 or higher
- **Kubernetes**: Version 1.24 or higher

### 3. Database Setup
```sql
-- Create databases
CREATE DATABASE shanraq_financial;
CREATE DATABASE shanraq_events;
CREATE DATABASE shanraq_audit;

-- Create users
CREATE USER shanraq_financial WITH PASSWORD 'secure_password';
CREATE USER shanraq_events WITH PASSWORD 'secure_password';
CREATE USER shanraq_audit WITH PASSWORD 'secure_password';

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE shanraq_financial TO shanraq_financial;
GRANT ALL PRIVILEGES ON DATABASE shanraq_events TO shanraq_events;
GRANT ALL PRIVILEGES ON DATABASE shanraq_audit TO shanraq_audit;
```

## Implementation Steps

### Step 1: Double-Entry Accounting System

#### 1.1 Database Schema Setup
```sql
-- Accounts table
CREATE TABLE accounts (
    account_id VARCHAR(36) PRIMARY KEY,
    account_type VARCHAR(20) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    balance DECIMAL(18,2) NOT NULL DEFAULT 0,
    debit_balance DECIMAL(18,2) NOT NULL DEFAULT 0,
    credit_balance DECIMAL(18,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Journal entries table
CREATE TABLE journal_entries (
    journal_id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Journal entry lines table
CREATE TABLE journal_entry_lines (
    line_id VARCHAR(36) PRIMARY KEY,
    journal_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    side VARCHAR(10) NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    description TEXT,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (journal_id) REFERENCES journal_entries(journal_id),
    FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

-- Transactions table
CREATE TABLE transactions (
    transaction_id VARCHAR(36) PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    committed_at BIGINT,
    rolled_back_at BIGINT
);
```

#### 1.2 Account Management
```tenge
// Initialize double-entry accounting system
atqar double_entry_system_konfig_jasau() -> aqıqat {
    // Create system accounts
    jasau system_accounts: Account[] = [
        account_jasau("system_cash", "asset", "USD", "system"),
        account_jasau("system_revenue", "revenue", "USD", "system"),
        account_jasau("system_expense", "expense", "USD", "system"),
        account_jasau("system_liability", "liability", "USD", "system")
    ];
    
    jasau i: san = 0;
    azirshe (i < system_accounts.length) {
        jasau stored: aqıqat = double_entry_store_account(system_accounts[i]);
        
        eгер (!stored) {
            korset("❌ Failed to store system account: " + system_accounts[i].account_id);
            qaytar jin;
        }
        
        i = i + 1;
    }
    
    qaytar jan;
}
```

#### 1.3 Transaction Processing
```tenge
// Process financial transaction
atqar financial_transaction_process_jasau(transaction_data: JsonObject) -> aqıqat {
    jasau transaction_id: jol = json_object_get_string(transaction_data, "transaction_id");
    jasau from_account: jol = json_object_get_string(transaction_data, "from_account");
    jasau to_account: jol = json_object_get_string(transaction_data, "to_account");
    jasau amount: jol = json_object_get_string(transaction_data, "amount");
    jasau currency: jol = json_object_get_string(transaction_data, "currency");
    
    // Create journal entries
    jasau debit_entry: JournalEntryLine = journal_entry_line_create();
    debit_entry.account_id = from_account;
    debit_entry.side = "debit";
    debit_entry.amount = decimal128_from_string(amount);
    debit_entry.description = "Transfer to " + to_account;
    
    jasau credit_entry: JournalEntryLine = journal_entry_line_create();
    credit_entry.account_id = to_account;
    credit_entry.side = "credit";
    credit_entry.amount = decimal128_from_string(amount);
    credit_entry.description = "Transfer from " + from_account;
    
    jasau journal: JournalEntry = journal_entry_jasau(
        transaction_id,
        "Financial Transfer",
        [debit_entry, credit_entry]
    );
    
    // Execute two-phase commit
    jasau commit_result: aqıqat = two_phase_commit_jasau(transaction_id, [journal]);
    
    qaytar commit_result;
}
```

### Step 2: Event Sourcing System

#### 2.1 Event Store Setup
```sql
-- Events table
CREATE TABLE events (
    event_id VARCHAR(36) PRIMARY KEY,
    aggregate_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    version INTEGER NOT NULL,
    timestamp BIGINT NOT NULL,
    correlation_id VARCHAR(36),
    causation_id VARCHAR(36),
    user_id VARCHAR(36),
    source_ip VARCHAR(45),
    event_hash VARCHAR(64) NOT NULL,
    created_at BIGINT NOT NULL
);

-- Event snapshots table
CREATE TABLE event_snapshots (
    snapshot_id VARCHAR(36) PRIMARY KEY,
    aggregate_id VARCHAR(36) NOT NULL,
    version INTEGER NOT NULL,
    state JSONB NOT NULL,
    created_at BIGINT NOT NULL
);

-- Event projections table
CREATE TABLE event_projections (
    projection_id VARCHAR(36) PRIMARY KEY,
    projection_name VARCHAR(100) NOT NULL,
    event_types TEXT[] NOT NULL,
    state JSONB NOT NULL,
    last_processed_event INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
```

#### 2.2 Event Sourcing Implementation
```tenge
// Initialize event sourcing system
atqar event_sourcing_system_konfig_jasau() -> aqıqat {
    // Create event store
    jasau event_store: EventStore = event_store_jasau();
    
    // Initialize event streams
    jasau streams: jol[] = [
        "transaction_stream",
        "account_stream",
        "user_stream",
        "payment_stream",
        "audit_stream"
    ];
    
    jasau i: san = 0;
    azirshe (i < streams.length) {
        event_store_initialize_stream(event_store, streams[i]);
        i = i + 1;
    }
    
    // Create event projections
    jasau projections: EventProjection[] = [
        event_projection_jasau("transaction_projection", ["TransactionCreated", "TransactionExecuted"]),
        event_projection_jasau("account_projection", ["AccountCreated", "AccountUpdated"]),
        event_projection_jasau("user_projection", ["UserCreated", "UserUpdated"])
    ];
    
    jasau j: san = 0;
    azirshe (j < projections.length) {
        jasau stored: aqıqat = event_sourcing_store_projection(projections[j]);
        
        eгер (!stored) {
            korset("❌ Failed to store event projection: " + projections[j].projection_name);
            qaytar jin;
        }
        
        j = j + 1;
    }
    
    qaytar jan;
}
```

### Step 3: Idempotency System

#### 3.1 Idempotency Storage Setup
```sql
-- Idempotency keys table
CREATE TABLE idempotency_keys (
    key_id VARCHAR(36) PRIMARY KEY,
    idempotency_key VARCHAR(128) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    endpoint VARCHAR(200) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    request_data JSONB,
    response_data JSONB,
    status_code INTEGER,
    error_message TEXT,
    error_code INTEGER,
    created_at BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    completed_at BIGINT,
    failed_at BIGINT,
    UNIQUE(idempotency_key, user_id, endpoint)
);

-- Create indexes
CREATE INDEX idx_idempotency_key ON idempotency_keys(idempotency_key);
CREATE INDEX idx_idempotency_user ON idempotency_keys(user_id);
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);
```

#### 3.2 Idempotency Implementation
```tenge
// Initialize idempotency system
atqar idempotency_system_konfig_jasau() -> aqıqat {
    // Configure idempotency storage
    jasau storage_config: JsonObject = json_object_create();
    json_object_set_string(storage_config, "backend", "postgresql");
    json_object_set_string(storage_config, "connection_string", "postgresql://user:pass@localhost/shanraq_financial");
    json_object_set_number(storage_config, "key_expiration_seconds", 86400); // 24 hours
    
    jasau storage_initialized: aqıqat = idempotency_storage_initialize(storage_config);
    
    eгер (!storage_initialized) {
        korset("❌ Failed to initialize idempotency storage");
        qaytar jin;
    }
    
    // Start cleanup job
    jasau cleanup_job: aqıqat = idempotency_start_cleanup_job();
    
    eгер (!cleanup_job) {
        korset("❌ Failed to start idempotency cleanup job");
        qaytar jin;
    }
    
    qaytar jan;
}
```

### Step 4: Message Queue System

#### 4.1 Kafka Setup
```bash
# Install Kafka
wget https://downloads.apache.org/kafka/2.8.0/kafka_2.13-2.8.0.tgz
tar -xzf kafka_2.13-2.8.0.tgz
cd kafka_2.13-2.8.0

# Start Zookeeper
bin/zookeeper-server-start.sh config/zookeeper.properties

# Start Kafka
bin/kafka-server-start.sh config/server.properties

# Create topics
bin/kafka-topics.sh --create --topic financial.transaction.created --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic financial.transaction.executed --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic financial.account.updated --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic financial.payment.completed --bootstrap-server localhost:9092
```

#### 4.2 Message Queue Implementation
```tenge
// Initialize message queue system
atqar message_queue_system_konfig_jasau() -> aqıqat {
    // Configure Kafka
    jasau kafka_config: KafkaConfig = kafka_konfig_jasau();
    kafka_config.bootstrap_servers = ["localhost:9092"];
    kafka_config.security_protocol = "PLAINTEXT";
    kafka_config.acks = "all";
    kafka_config.retries = 3;
    kafka_config.enable_idempotence = jan;
    
    // Create connection
    jasau connection: MessageQueueConnection = message_queue_connection_jasau(kafka_config);
    
    eгер (connection == NULL) {
        korset("❌ Failed to create message queue connection");
        qaytar jin;
    }
    
    // Create event handlers
    jasau handlers: MessageHandler[] = [
        financial_event_handler_jasau("TransactionCreated"),
        financial_event_handler_jasau("TransactionExecuted"),
        financial_event_handler_jasau("AccountUpdated"),
        financial_event_handler_jasau("PaymentCompleted")
    ];
    
    jasau i: san = 0;
    azirshe (i < handlers.length) {
        jasau consumer: aqıqat = message_consume_jasau(connection, "financial.events", "shanraq_consumer", handlers[i]);
        
        eгер (!consumer) {
            korset("❌ Failed to create message consumer for: " + handlers[i].event_type);
            qaytar jin;
        }
        
        i = i + 1;
    }
    
    qaytar jan;
}
```

### Step 5: P2P Transfer System

#### 5.1 P2P Transfer Database Setup
```sql
-- P2P transfers table
CREATE TABLE p2p_transfers (
    transfer_id VARCHAR(36) PRIMARY KEY,
    from_user_id VARCHAR(36) NOT NULL,
    to_user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(18,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    completed_at BIGINT,
    failed_at BIGINT,
    rolled_back_at BIGINT,
    failure_reason TEXT,
    rollback_reason TEXT,
    block_reason TEXT
);

-- P2P transfer limits table
CREATE TABLE p2p_transfer_limits (
    user_id VARCHAR(36) PRIMARY KEY,
    daily_limit DECIMAL(18,2) NOT NULL DEFAULT 10000,
    monthly_limit DECIMAL(18,2) NOT NULL DEFAULT 100000,
    max_single_transfer DECIMAL(18,2) NOT NULL DEFAULT 5000,
    min_single_transfer DECIMAL(18,2) NOT NULL DEFAULT 1,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- P2P transfer history table
CREATE TABLE p2p_transfer_history (
    history_id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    transfer_id VARCHAR(36) NOT NULL,
    action VARCHAR(50) NOT NULL,
    details JSONB,
    created_at BIGINT NOT NULL
);
```

#### 5.2 P2P Transfer Implementation
```tenge
// Initialize P2P transfer system
atqar p2p_transfer_system_konfig_jasau() -> aqıqat {
    // Create default transfer limits
    jasau default_limits: JsonObject = json_object_create();
    json_object_set_string(default_limits, "daily_limit", "10000");
    json_object_set_string(default_limits, "monthly_limit", "100000");
    json_object_set_string(default_limits, "max_single_transfer", "5000");
    json_object_set_string(default_limits, "min_single_transfer", "1");
    
    // Set up fraud detection
    jasau fraud_config: JsonObject = json_object_create();
    json_object_set_number(fraud_config, "risk_threshold", 70);
    json_object_set_number(fraud_config, "amount_multiplier", 5);
    json_object_set_number(fraud_config, "frequency_limit", 10);
    json_object_set_number(fraud_config, "time_window", 3600);
    
    jasau fraud_initialized: aqıqat = p2p_fraud_detection_initialize(fraud_config);
    
    eгер (!fraud_initialized) {
        korset("❌ Failed to initialize fraud detection");
        qaytar jin;
    }
    
    // Set up notifications
    jasau notification_config: JsonObject = json_object_create();
    json_object_set_string(notification_config, "email_enabled", "true");
    json_object_set_string(notification_config, "sms_enabled", "true");
    json_object_set_string(notification_config, "push_enabled", "true");
    
    jasau notifications_initialized: aqıqat = p2p_notifications_initialize(notification_config);
    
    eгер (!notifications_initialized) {
        korset("❌ Failed to initialize notifications");
        qaytar jin;
    }
    
    qaytar jan;
}
```

## API Implementation

### 1. Transaction API Endpoints
```tenge
// Transaction API implementation
atqar transaction_api_jasau() -> JsonObject {
    jasau api: JsonObject = json_object_create();
    
    // POST /api/v1/transactions
    jasau create_transaction: JsonObject = json_object_create();
    json_object_set_string(create_transaction, "method", "POST");
    json_object_set_string(create_transaction, "path", "/api/v1/transactions");
    json_object_set_string(create_transaction, "description", "Create a new transaction");
    json_object_set_string(create_transaction, "handler", "transaction_create_handler");
    
    // GET /api/v1/transactions/{transaction_id}
    jasau get_transaction: JsonObject = json_object_create();
    json_object_set_string(get_transaction, "method", "GET");
    json_object_set_string(get_transaction, "path", "/api/v1/transactions/{transaction_id}");
    json_object_set_string(get_transaction, "description", "Get transaction details");
    json_object_set_string(get_transaction, "handler", "transaction_get_handler");
    
    // POST /api/v1/transactions/{transaction_id}/rollback
    jasau rollback_transaction: JsonObject = json_object_create();
    json_object_set_string(rollback_transaction, "method", "POST");
    json_object_set_string(rollback_transaction, "path", "/api/v1/transactions/{transaction_id}/rollback");
    json_object_set_string(rollback_transaction, "description", "Rollback transaction");
    json_object_set_string(rollback_transaction, "handler", "transaction_rollback_handler");
    
    json_object_set_object(api, "create_transaction", create_transaction);
    json_object_set_object(api, "get_transaction", get_transaction);
    json_object_set_object(api, "rollback_transaction", rollback_transaction);
    
    qaytar api;
}
```

### 2. P2P Transfer API Endpoints
```tenge
// P2P Transfer API implementation
atqar p2p_transfer_api_jasau() -> JsonObject {
    jasau api: JsonObject = json_object_create();
    
    // POST /api/v1/p2p/transfer
    jasau create_transfer: JsonObject = json_object_create();
    json_object_set_string(create_transfer, "method", "POST");
    json_object_set_string(create_transfer, "path", "/api/v1/p2p/transfer");
    json_object_set_string(create_transfer, "description", "Create a new P2P transfer");
    json_object_set_string(create_transfer, "handler", "p2p_transfer_create_handler");
    
    // GET /api/v1/p2p/transfer/{transfer_id}
    jasau get_transfer: JsonObject = json_object_create();
    json_object_set_string(get_transfer, "method", "GET");
    json_object_set_string(get_transfer, "path", "/api/v1/p2p/transfer/{transfer_id}");
    json_object_set_string(get_transfer, "description", "Get P2P transfer status");
    json_object_set_string(get_transfer, "handler", "p2p_transfer_get_handler");
    
    // GET /api/v1/p2p/transfers
    jasau list_transfers: JsonObject = json_object_create();
    json_object_set_string(list_transfers, "method", "GET");
    json_object_set_string(list_transfers, "path", "/api/v1/p2p/transfers");
    json_object_set_string(list_transfers, "description", "List user P2P transfers");
    json_object_set_string(list_transfers, "handler", "p2p_transfer_list_handler");
    
    json_object_set_object(api, "create_transfer", create_transfer);
    json_object_set_object(api, "get_transfer", get_transfer);
    json_object_set_object(api, "list_transfers", list_transfers);
    
    qaytar api;
}
```

## Testing Implementation

### 1. Unit Testing
```tenge
// Unit tests for double-entry accounting
atqar test_double_entry_accounting() -> aqıqat {
    // Test account creation
    jasau account: Account = account_jasau("test_account", "asset", "USD", "test_user");
    
    eгер (account == NULL) {
        korset("❌ Account creation test failed");
        qaytar jin;
    }
    
    // Test journal entry validation
    jasau debit_entry: JournalEntryLine = journal_entry_line_create();
    debit_entry.account_id = "test_account";
    debit_entry.side = "debit";
    debit_entry.amount = decimal128_jasau("100", 18, 2);
    
    jasau credit_entry: JournalEntryLine = journal_entry_line_create();
    credit_entry.account_id = "test_account";
    credit_entry.side = "credit";
    credit_entry.amount = decimal128_jasau("100", 18, 2);
    
    jasau journal: JournalEntry = journal_entry_jasau("test_transaction", "Test Entry", [debit_entry, credit_entry]);
    
    eгер (journal == NULL) {
        korset("❌ Journal entry creation test failed");
        qaytar jin;
    }
    
    qaytar jan;
}
```

### 2. Integration Testing
```tenge
// Integration tests for transaction processing
atqar test_transaction_processing() -> aqıqat {
    // Test complete transaction flow
    jasau transaction_data: JsonObject = json_object_create();
    json_object_set_string(transaction_data, "transaction_id", "test_transaction_001");
    json_object_set_string(transaction_data, "from_account", "test_account_1");
    json_object_set_string(transaction_data, "to_account", "test_account_2");
    json_object_set_string(transaction_data, "amount", "100.00");
    json_object_set_string(transaction_data, "currency", "USD");
    
    jasau result: aqıqat = financial_transaction_process_jasau(transaction_data);
    
    eгер (!result) {
        korset("❌ Transaction processing test failed");
        qaytar jin;
    }
    
    qaytar jan;
}
```

### 3. Performance Testing
```tenge
// Performance tests for high-volume transactions
atqar test_performance_high_volume() -> aqıqat {
    jasau start_time: san = current_timestamp();
    jasau transaction_count: san = 1000;
    jasau success_count: san = 0;
    
    jasau i: san = 0;
    azirshe (i < transaction_count) {
        jasau transaction_data: JsonObject = json_object_create();
        json_object_set_string(transaction_data, "transaction_id", "perf_test_" + i.toString());
        json_object_set_string(transaction_data, "from_account", "perf_account_1");
        json_object_set_string(transaction_data, "to_account", "perf_account_2");
        json_object_set_string(transaction_data, "amount", "10.00");
        json_object_set_string(transaction_data, "currency", "USD");
        
        jasau result: aqıqat = financial_transaction_process_jasau(transaction_data);
        
        eгер (result) {
            success_count = success_count + 1;
        }
        
        i = i + 1;
    }
    
    jasau end_time: san = current_timestamp();
    jasau duration: san = end_time - start_time;
    jasau tps: san = transaction_count / duration;
    
    korset("✅ Performance test completed: " + success_count + "/" + transaction_count + " transactions in " + duration + "ms (" + tps + " TPS)");
    
    qaytar success_count == transaction_count;
}
```

## Monitoring Implementation

### 1. Metrics Collection
```tenge
// Metrics collection for transaction core
atqar transaction_core_metrics_jasau() -> JsonObject {
    jasau metrics: JsonObject = json_object_create();
    
    // Transaction metrics
    jasau transaction_metrics: JsonObject = json_object_create();
    json_object_set_number(transaction_metrics, "total_transactions", get_total_transaction_count());
    json_object_set_number(transaction_metrics, "successful_transactions", get_successful_transaction_count());
    json_object_set_number(transaction_metrics, "failed_transactions", get_failed_transaction_count());
    json_object_set_number(transaction_metrics, "rolled_back_transactions", get_rolled_back_transaction_count());
    
    // Performance metrics
    jasau performance_metrics: JsonObject = json_object_create();
    json_object_set_number(performance_metrics, "average_processing_time", get_average_processing_time());
    json_object_set_number(performance_metrics, "transactions_per_second", get_transactions_per_second());
    json_object_set_number(performance_metrics, "queue_length", get_queue_length());
    
    // System metrics
    jasau system_metrics: JsonObject = json_object_create();
    json_object_set_number(system_metrics, "cpu_usage", get_cpu_usage());
    json_object_set_number(system_metrics, "memory_usage", get_memory_usage());
    json_object_set_number(system_metrics, "disk_usage", get_disk_usage());
    
    json_object_set_object(metrics, "transactions", transaction_metrics);
    json_object_set_object(metrics, "performance", performance_metrics);
    json_object_set_object(metrics, "system", system_metrics);
    
    qaytar metrics;
}
```

### 2. Health Checks
```tenge
// Health checks for transaction core
atqar transaction_core_health_check() -> JsonObject {
    jasau health: JsonObject = json_object_create();
    
    // Check database connectivity
    jasau database_ok: aqıqat = database_health_check();
    
    // Check message queue connectivity
    jasau message_queue_ok: aqıqat = message_queue_health_check();
    
    // Check event store connectivity
    jasau event_store_ok: aqıqat = event_store_health_check();
    
    // Check idempotency system
    jasau idempotency_ok: aqıqat = idempotency_health_check();
    
    // Determine overall health
    jasau overall_health: jol = "healthy";
    
    eгер (!database_ok || !message_queue_ok || !event_store_ok || !idempotency_ok) {
        overall_health = "unhealthy";
    }
    
    json_object_set_string(health, "status", overall_health);
    json_object_set_boolean(health, "database_ok", database_ok);
    json_object_set_boolean(health, "message_queue_ok", message_queue_ok);
    json_object_set_boolean(health, "event_store_ok", event_store_ok);
    json_object_set_boolean(health, "idempotency_ok", idempotency_ok);
    json_object_set_number(health, "checked_at", current_timestamp());
    
    qaytar health;
}
```

## Deployment Configuration

### 1. Docker Configuration
```dockerfile
# Dockerfile for transaction core
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

EXPOSE 3000

CMD ["node", "index.js"]
```

### 2. Kubernetes Configuration
```yaml
# transaction-core-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: transaction-core
spec:
  replicas: 3
  selector:
    matchLabels:
      app: transaction-core
  template:
    metadata:
      labels:
        app: transaction-core
    spec:
      containers:
      - name: transaction-core
        image: shanraq/transaction-core:latest
        ports:
        - containerPort: 3000
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: url
        - name: KAFKA_BROKERS
          value: "kafka:9092"
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
```

### 3. Environment Configuration
```bash
# Environment variables
export DATABASE_URL="postgresql://user:pass@localhost/shanraq_financial"
export REDIS_URL="redis://localhost:6379"
export KAFKA_BROKERS="localhost:9092"
export NODE_ENV="production"
export LOG_LEVEL="info"
export METRICS_ENABLED="true"
export HEALTH_CHECK_ENABLED="true"
```

## Conclusion

This implementation guide provides comprehensive instructions for setting up the transaction core system for Shanraq.org's fintech platform. The implementation covers:

- **Database Setup**: Complete database schema and configuration
- **System Integration**: All core components and their interactions
- **API Implementation**: RESTful API endpoints for all operations
- **Testing Strategy**: Unit, integration, and performance testing
- **Monitoring**: Metrics collection and health checks
- **Deployment**: Docker and Kubernetes configuration

Following this guide will result in a robust, scalable, and secure transaction processing system capable of handling high-volume financial operations while maintaining data integrity and compliance with regulatory requirements.
