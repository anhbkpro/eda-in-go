#!/bin/bash
# gRPC Request Examples for Stores Service
# Usage: ./grpc-requests.sh <command> [args]
#
# Make executable: chmod +x scripts/grpc-requests.sh

GRPC_HOST="${GRPC_HOST:-localhost:8086}"

# ============================================================================
# STORE COMMANDS
# ============================================================================

# Create a new store
# Usage: ./grpc-requests.sh create-store "Store Name" "Location"
create_store() {
    local name="${1:-My Store}"
    local location="${2:-123 Main Street}"
    grpcurl -plaintext -d "{
        \"name\": \"$name\",
        \"location\": \"$location\"
    }" "$GRPC_HOST" storespb.StoresService/CreateStore
}

# Get a store by ID
# Usage: ./grpc-requests.sh get-store <store_id>
get_store() {
    local store_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$store_id\"
    }" "$GRPC_HOST" storespb.StoresService/GetStore
}

# Get all stores
# Usage: ./grpc-requests.sh get-stores
get_stores() {
    grpcurl -plaintext -d '{}' "$GRPC_HOST" storespb.StoresService/GetStores
}

# Get participating stores
# Usage: ./grpc-requests.sh get-participating-stores
get_participating_stores() {
    grpcurl -plaintext -d '{}' "$GRPC_HOST" storespb.StoresService/GetParticipatingStores
}

# Enable store participation
# Usage: ./grpc-requests.sh enable-participation <store_id>
enable_participation() {
    local store_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$store_id\"
    }" "$GRPC_HOST" storespb.StoresService/EnableParticipation
}

# Disable store participation
# Usage: ./grpc-requests.sh disable-participation <store_id>
disable_participation() {
    local store_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$store_id\"
    }" "$GRPC_HOST" storespb.StoresService/DisableParticipation
}

# Rebrand a store
# Usage: ./grpc-requests.sh rebrand-store <store_id> "New Name"
rebrand_store() {
    local store_id="$1"
    local name="$2"
    grpcurl -plaintext -d "{
        \"id\": \"$store_id\",
        \"name\": \"$name\"
    }" "$GRPC_HOST" storespb.StoresService/RebrandStore
}

# ============================================================================
# PRODUCT COMMANDS
# ============================================================================

# Add a product to a store
# Usage: ./grpc-requests.sh add-product <store_id> "Name" "Description" "SKU" <price>
add_product() {
    local store_id="$1"
    local name="${2:-Espresso}"
    local description="${3:-Strong Italian coffee}"
    local sku="${4:-ESP-001}"
    local price="${5:-4.50}"
    grpcurl -plaintext -d "{
        \"store_id\": \"$store_id\",
        \"name\": \"$name\",
        \"description\": \"$description\",
        \"sku\": \"$sku\",
        \"price\": $price
    }" "$GRPC_HOST" storespb.StoresService/AddProduct
}

# Get a product by ID
# Usage: ./grpc-requests.sh get-product <product_id>
get_product() {
    local product_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$product_id\"
    }" "$GRPC_HOST" storespb.StoresService/GetProduct
}

# Get catalog (all products for a store)
# Usage: ./grpc-requests.sh get-catalog <store_id>
get_catalog() {
    local store_id="$1"
    grpcurl -plaintext -d "{
        \"store_id\": \"$store_id\"
    }" "$GRPC_HOST" storespb.StoresService/GetCatalog
}

# Rebrand a product
# Usage: ./grpc-requests.sh rebrand-product <product_id> "New Name" "New Description"
rebrand_product() {
    local product_id="$1"
    local name="$2"
    local description="$3"
    grpcurl -plaintext -d "{
        \"id\": \"$product_id\",
        \"name\": \"$name\",
        \"description\": \"$description\"
    }" "$GRPC_HOST" storespb.StoresService/RebrandProduct
}

# Increase product price
# Usage: ./grpc-requests.sh increase-price <product_id> <amount>
increase_price() {
    local product_id="$1"
    local price="$2"
    grpcurl -plaintext -d "{
        \"id\": \"$product_id\",
        \"price\": $price
    }" "$GRPC_HOST" storespb.StoresService/IncreaseProductPrice
}

# Decrease product price
# Usage: ./grpc-requests.sh decrease-price <product_id> <amount>
decrease_price() {
    local product_id="$1"
    local price="$2"
    grpcurl -plaintext -d "{
        \"id\": \"$product_id\",
        \"price\": $price
    }" "$GRPC_HOST" storespb.StoresService/DecreaseProductPrice
}

# Remove a product
# Usage: ./grpc-requests.sh remove-product <product_id>
remove_product() {
    local product_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$product_id\"
    }" "$GRPC_HOST" storespb.StoresService/RemoveProduct
}

# ============================================================================
# UTILITY
# ============================================================================

# List all available services
list_services() {
    grpcurl -plaintext "$GRPC_HOST" list
}

# Describe a service
describe_service() {
    local service="${1:-storespb.StoresService}"
    grpcurl -plaintext "$GRPC_HOST" describe "$service"
}

# Show help
show_help() {
    echo "gRPC Request Examples for Stores Service"
    echo ""
    echo "Usage: $0 <command> [args]"
    echo ""
    echo "Store Commands:"
    echo "  create-store [name] [location]     Create a new store"
    echo "  get-store <store_id>               Get store by ID"
    echo "  get-stores                         Get all stores"
    echo "  get-participating-stores           Get participating stores"
    echo "  enable-participation <store_id>    Enable store participation"
    echo "  disable-participation <store_id>   Disable store participation"
    echo "  rebrand-store <store_id> <name>    Rebrand a store"
    echo ""
    echo "Product Commands:"
    echo "  add-product <store_id> [name] [desc] [sku] [price]  Add product"
    echo "  get-product <product_id>           Get product by ID"
    echo "  get-catalog <store_id>             Get all products for store"
    echo "  rebrand-product <id> <name> <desc> Rebrand a product"
    echo "  increase-price <id> <amount>       Increase product price"
    echo "  decrease-price <id> <amount>       Decrease product price"
    echo "  remove-product <product_id>        Remove a product"
    echo ""
    echo "Utility Commands:"
    echo "  list-services                      List all gRPC services"
    echo "  describe-service [service]         Describe a service"
    echo "  help                               Show this help"
}

# ============================================================================
# MAIN
# ============================================================================

case "${1:-help}" in
    create-store)           shift; create_store "$@" ;;
    get-store)              shift; get_store "$@" ;;
    get-stores)             get_stores ;;
    get-participating-stores) get_participating_stores ;;
    enable-participation)   shift; enable_participation "$@" ;;
    disable-participation)  shift; disable_participation "$@" ;;
    rebrand-store)          shift; rebrand_store "$@" ;;
    add-product)            shift; add_product "$@" ;;
    get-product)            shift; get_product "$@" ;;
    get-catalog)            shift; get_catalog "$@" ;;
    rebrand-product)        shift; rebrand_product "$@" ;;
    increase-price)         shift; increase_price "$@" ;;
    decrease-price)         shift; decrease_price "$@" ;;
    remove-product)         shift; remove_product "$@" ;;
    list-services)          list_services ;;
    describe-service)       shift; describe_service "$@" ;;
    help|--help|-h)         show_help ;;
    *)                      echo "Unknown command: $1"; show_help; exit 1 ;;
esac

