#!/bin/bash
# gRPC Request Examples for Stores, Customers, Baskets, Payments, and Ordering Services
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
# CUSTOMER COMMANDS
# ============================================================================

# Register a new customer
# Usage: ./grpc-requests.sh register-customer "Name" "+1234567890"
register_customer() {
    local name="${1:-John Doe}"
    local sms_number="${2:-+1234567890}"
    grpcurl -plaintext -d "{
        \"name\": \"$name\",
        \"sms_number\": \"$sms_number\"
    }" "$GRPC_HOST" customerspb.CustomersService/RegisterCustomer
}

# Get a customer by ID
# Usage: ./grpc-requests.sh get-customer <customer_id>
get_customer() {
    local customer_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$customer_id\"
    }" "$GRPC_HOST" customerspb.CustomersService/GetCustomer
}

# Authorize a customer
# Usage: ./grpc-requests.sh authorize-customer <customer_id>
authorize_customer() {
    local customer_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$customer_id\"
    }" "$GRPC_HOST" customerspb.CustomersService/AuthorizeCustomer
}

# Enable a customer
# Usage: ./grpc-requests.sh enable-customer <customer_id>
enable_customer() {
    local customer_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$customer_id\"
    }" "$GRPC_HOST" customerspb.CustomersService/EnableCustomer
}

# Disable a customer
# Usage: ./grpc-requests.sh disable-customer <customer_id>
disable_customer() {
    local customer_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$customer_id\"
    }" "$GRPC_HOST" customerspb.CustomersService/DisableCustomer
}

# ============================================================================
# BASKET COMMANDS
# ============================================================================

# Start a new basket for a customer
# Usage: ./grpc-requests.sh start-basket <customer_id>
start_basket() {
    local customer_id="$1"
    grpcurl -plaintext -d "{
        \"customer_id\": \"$customer_id\"
    }" "$GRPC_HOST" basketspb.BasketService/StartBasket
}

# Get a basket by ID
# Usage: ./grpc-requests.sh get-basket <basket_id>
get_basket() {
    local basket_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$basket_id\"
    }" "$GRPC_HOST" basketspb.BasketService/GetBasket
}

# Cancel a basket
# Usage: ./grpc-requests.sh cancel-basket <basket_id>
cancel_basket() {
    local basket_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$basket_id\"
    }" "$GRPC_HOST" basketspb.BasketService/CancelBasket
}

# Checkout a basket
# Usage: ./grpc-requests.sh checkout-basket <basket_id> <payment_id>
checkout_basket() {
    local basket_id="$1"
    local payment_id="$2"
    grpcurl -plaintext -d "{
        \"id\": \"$basket_id\",
        \"payment_id\": \"$payment_id\"
    }" "$GRPC_HOST" basketspb.BasketService/CheckoutBasket
}

# Add an item to a basket
# Usage: ./grpc-requests.sh add-item <basket_id> <product_id> <quantity>
add_item() {
    local basket_id="$1"
    local product_id="$2"
    local quantity="${3:-1}"
    grpcurl -plaintext -d "{
        \"id\": \"$basket_id\",
        \"product_id\": \"$product_id\",
        \"quantity\": $quantity
    }" "$GRPC_HOST" basketspb.BasketService/AddItem
}

# Remove an item from a basket
# Usage: ./grpc-requests.sh remove-item <basket_id> <product_id> <quantity>
remove_item() {
    local basket_id="$1"
    local product_id="$2"
    local quantity="${3:-1}"
    grpcurl -plaintext -d "{
        \"id\": \"$basket_id\",
        \"product_id\": \"$product_id\",
        \"quantity\": $quantity
    }" "$GRPC_HOST" basketspb.BasketService/RemoveItem
}

# ============================================================================
# PAYMENT COMMANDS
# ============================================================================

# Authorize a payment
# Usage: ./grpc-requests.sh authorize-payment <customer_id> <amount>
authorize_payment() {
    local customer_id="$1"
    local amount="$2"
    grpcurl -plaintext -d "{
        \"customer_id\": \"$customer_id\",
        \"amount\": $amount
    }" "$GRPC_HOST" paymentspb.PaymentsService/AuthorizePayment
}

# Confirm a payment
# Usage: ./grpc-requests.sh confirm-payment <payment_id>
confirm_payment() {
    local payment_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$payment_id\"
    }" "$GRPC_HOST" paymentspb.PaymentsService/ConfirmPayment
}

# Create an invoice
# Usage: ./grpc-requests.sh create-invoice <order_id> <payment_id> <amount>
create_invoice() {
    local order_id="$1"
    local payment_id="$2"
    local amount="$3"
    grpcurl -plaintext -d "{
        \"order_id\": \"$order_id\",
        \"payment_id\": \"$payment_id\",
        \"amount\": $amount
    }" "$GRPC_HOST" paymentspb.PaymentsService/CreateInvoice
}

# Adjust an invoice
# Usage: ./grpc-requests.sh adjust-invoice <invoice_id> <amount>
adjust_invoice() {
    local invoice_id="$1"
    local amount="$2"
    grpcurl -plaintext -d "{
        \"id\": \"$invoice_id\",
        \"amount\": $amount
    }" "$GRPC_HOST" paymentspb.PaymentsService/AdjustInvoice
}

# Pay an invoice
# Usage: ./grpc-requests.sh pay-invoice <invoice_id>
pay_invoice() {
    local invoice_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$invoice_id\"
    }" "$GRPC_HOST" paymentspb.PaymentsService/PayInvoice
}

# Cancel an invoice
# Usage: ./grpc-requests.sh cancel-invoice <invoice_id>
cancel_invoice() {
    local invoice_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$invoice_id\"
    }" "$GRPC_HOST" paymentspb.PaymentsService/CancelInvoice
}

# ============================================================================
# ORDERING COMMANDS
# ============================================================================

# Create an order
# Usage: ./grpc-requests.sh create-order <customer_id> <payment_id> <store_id> <product_id> <quantity> [price]
create_order() {
    local customer_id="$1"
    local payment_id="$2"
    local store_id="$3"
    local product_id="$4"
    local quantity="${5:-1}"
    local price="${6:-10.00}"
    local store_name="${7:-Store}"
    local product_name="${8:-Product}"
    grpcurl -plaintext -d "{
        \"items\": [{
            \"store_id\": \"$store_id\",
            \"product_id\": \"$product_id\",
            \"store_name\": \"$store_name\",
            \"product_name\": \"$product_name\",
            \"price\": $price,
            \"quantity\": $quantity
        }],
        \"customer_id\": \"$customer_id\",
        \"payment_id\": \"$payment_id\"
    }" "$GRPC_HOST" orderingpb.OrderingService/CreateOrder
}

# Get an order by ID
# Usage: ./grpc-requests.sh get-order <order_id>
get_order() {
    local order_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$order_id\"
    }" "$GRPC_HOST" orderingpb.OrderingService/GetOrder
}

# Cancel an order
# Usage: ./grpc-requests.sh cancel-order <order_id>
cancel_order() {
    local order_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$order_id\"
    }" "$GRPC_HOST" orderingpb.OrderingService/CancelOrder
}

# Ready an order
# Usage: ./grpc-requests.sh ready-order <order_id>
ready_order() {
    local order_id="$1"
    grpcurl -plaintext -d "{
        \"id\": \"$order_id\"
    }" "$GRPC_HOST" orderingpb.OrderingService/ReadyOrder
}

# Complete an order
# Usage: ./grpc-requests.sh complete-order <order_id> <invoice_id>
complete_order() {
    local order_id="$1"
    local invoice_id="$2"
    grpcurl -plaintext -d "{
        \"id\": \"$order_id\",
        \"invoice_id\": \"$invoice_id\"
    }" "$GRPC_HOST" orderingpb.OrderingService/CompleteOrder
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
    echo "gRPC Request Examples for Stores, Customers, Baskets, Payments, and Ordering Services"
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
    echo "Customer Commands:"
    echo "  register-customer [name] [sms]     Register a new customer"
    echo "  get-customer <customer_id>         Get customer by ID"
    echo "  authorize-customer <customer_id>   Authorize a customer"
    echo "  enable-customer <customer_id>      Enable a customer"
    echo "  disable-customer <customer_id>     Disable a customer"
    echo ""
    echo "Basket Commands:"
    echo "  start-basket <customer_id>         Start a new basket"
    echo "  get-basket <basket_id>             Get basket by ID"
    echo "  cancel-basket <basket_id>          Cancel a basket"
    echo "  checkout-basket <basket_id> <payment_id>  Checkout a basket"
    echo "  add-item <basket_id> <product_id> [quantity]  Add item to basket"
    echo "  remove-item <basket_id> <product_id> [quantity]  Remove item from basket"
    echo ""
    echo "Payment Commands:"
    echo "  authorize-payment <customer_id> <amount>  Authorize a payment"
    echo "  confirm-payment <payment_id>       Confirm a payment"
    echo "  create-invoice <order_id> <payment_id> <amount>  Create an invoice"
    echo "  adjust-invoice <invoice_id> <amount>  Adjust invoice amount"
    echo "  pay-invoice <invoice_id>            Pay an invoice"
    echo "  cancel-invoice <invoice_id>         Cancel an invoice"
    echo ""
    echo "Ordering Commands:"
    echo "  create-order <customer_id> <payment_id> <store_id> <product_id> [qty] [price]  Create an order"
    echo "  get-order <order_id>                Get order by ID"
    echo "  cancel-order <order_id>             Cancel an order"
    echo "  ready-order <order_id>              Mark order as ready"
    echo "  complete-order <order_id> <invoice_id>  Complete an order"
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
    register-customer)      shift; register_customer "$@" ;;
    get-customer)           shift; get_customer "$@" ;;
    authorize-customer)     shift; authorize_customer "$@" ;;
    enable-customer)        shift; enable_customer "$@" ;;
    disable-customer)       shift; disable_customer "$@" ;;
    start-basket)           shift; start_basket "$@" ;;
    get-basket)             shift; get_basket "$@" ;;
    cancel-basket)          shift; cancel_basket "$@" ;;
    checkout-basket)        shift; checkout_basket "$@" ;;
    add-item)               shift; add_item "$@" ;;
    remove-item)            shift; remove_item "$@" ;;
    authorize-payment)      shift; authorize_payment "$@" ;;
    confirm-payment)        shift; confirm_payment "$@" ;;
    create-invoice)         shift; create_invoice "$@" ;;
    adjust-invoice)         shift; adjust_invoice "$@" ;;
    pay-invoice)            shift; pay_invoice "$@" ;;
    cancel-invoice)         shift; cancel_invoice "$@" ;;
    create-order)           shift; create_order "$@" ;;
    get-order)              shift; get_order "$@" ;;
    cancel-order)           shift; cancel_order "$@" ;;
    ready-order)            shift; ready_order "$@" ;;
    complete-order)         shift; complete_order "$@" ;;
    list-services)          list_services ;;
    describe-service)       shift; describe_service "$@" ;;
    help|--help|-h)         show_help ;;
    *)                      echo "Unknown command: $1"; show_help; exit 1 ;;
esac

