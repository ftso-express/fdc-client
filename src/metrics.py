from prometheus_client import Counter, Gauge

# Collector metrics
COLLECTOR_ITEMS_FETCHED = Counter(
    "fdc_collector_items_fetched_total",
    "Total number of items fetched by the collector",
    ["item_type"]  # e.g., "attestation_request", "bit_vote", "signing_policy"
)

# Manager metrics
MANAGER_ITEMS_PROCESSED = Counter(
    "fdc_manager_items_processed_total",
    "Total number of items processed by the manager",
    ["item_type"]  # e.g., "request", "bit_vote", "policy"
)

# Queue metrics
QUEUE_SIZE = Gauge(
    "fdc_queue_size_current",
    "Current size of the data pipe queues",
    ["queue_name"] # e.g., "requests", "bit_votes", "signing_policies"
)