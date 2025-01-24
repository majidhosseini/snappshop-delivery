# Delivery Service:
## Responsibility:
1. Integration of 3PLs
2. Requesting couriers for shipments from a specific origin to a destination within a selected time frame
3. Management and updating of courier and request statuses

## TODO:
- [ ] Project directory structure
- [x] Database Design
- [ ] 3PL Proxy Pattern
- [ ] Webhook and API Design
- [ ] Database Migration
- [ ] Entity and Repositories
- [ ] Services
- [ ] Controller
- [ ] Test and Mock
- [ ] Documentation
- [ ] Request Validation and Idempotency
- [ ] Log

## Receive new request:
```mermaid
flowchart TD
    A((Receive new request)) --> B(Validation Request,Idempotency,Timeframe, ...)
    B --> C{is Valid?}
    C --> |is valid| D[Store request in DB]
    C --> |Not valid| E[Return Error]
    D --> F{is lower than 1 Hour?}
    F -->|Yes| G[Request for shipment]
    F -->|No| H[Schedule it]
```

## Fetch schedule requests:
```mermaid
flowchart TD
    A((Scheduler)) --> B(Fetch shipment ready Orders or no response requests)
    B --> DAS@{ shape: das, label: "Queue Requests" }
    DAS --> 3PL(Send to provider)
    3PL --> C{is success?}
    C --> |Yes| D[(Change state to finding)]
    C --> |No!| 3PL
```

## Proces requests:
```mermaid
flowchart TD
    A((Delivery Webhook)) --> B(Validate body)
    B --> C{is correct?}
    C --> |Yes| D(Fetch from DB)
    D --> E{Has Data?}
    E --> |Yes| U(Update state)

    C --> |No| END((Do nothing))
    E --> |No| END
```

# User Story:
## Requirements:
Request for shipment model:
| Field    | type |
| -------- | ------- |
| OrderId  | unique    |
| UserInfo | UserInfo     |
| fromLoc  | [lat,lng]    |
| toLoc    | [lat,lng] |
| delveryTimeFrame | Timeframe |

Time Frame: Two-hour time slots from 9 AM to 11 PM (for the next 4 days)


### Delivery State
    1. init
    2. isFinding
    3. found
    4. notFound
    5. delivered

```mermaid
sequenceDiagram
    participant Core
    box Green Delivery Service
    participant Delivery
    participant Database
    end
    participant 3PL

    Core->>Delivery: Success Order
    Delivery->>Database: Store Order Delivery request
    loop
        Delivery->>3PL: Shipment ready request
    end
    Note right of Delivery: Ready request
    activate 3PL
    3PL->>3PL: Finding shipment
    3PL-->>Delivery: Send "Found", "NotFound", "Delivered" state
    deactivate 3PL
    Delivery-->>Core: Send latest state of order's delivery
```

## Schedule Rules:
1. Request for shipment if time to delivery is <= 1 Hour OR is >= 15' from init state.

## Delivery Acceptance 
1. Process 10,000 Requests in 1 Hour
2. Scalability
3. Seeder for generate orders in difference time
4. Find the shipment after maximum 3 retries.
5. 95% Success (or uptime)
6. Core get 5% fault for get state.
7. 3PL has 5% fault 

