# Delivery Service:

## Receive new request:
```mermaid
flowchart TD
    A((Receive new request)) --> B(Validation Request)
    B --> C{is Valid?}
    C --> |is valid| D[Store request in DB]
    C --> |Not valid| E[Return Error]
    D --> F{is lower than 1 Hour?}
    F -->|Yes| G[Request for delivery]
    F -->|No| H[Schedule it]
```

## Fetch schedule requests:
```mermaid
flowchart TD
    A((Fetch from DB))
```

## Proces requests:
```mermaid
flowchart TD
    A((Request for shipment))
```





# User Story:
## Requirements:

Order Entity:

| Field    | type |
| -------- | ------- |
| OrderId  | unique    |
| UserInfo | UserInfo     |
| fromLoc  | [lat,lng]    |
| toLoc    | [lat, lng] |
| delveryDuration | [from, to] |



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

