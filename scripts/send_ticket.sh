#!/bin/bash

# ONLY WORKS WHEN RUN IN PROJECT ROOT DIRECTORY

echo "Sending ticket scripts/ticket.txt"

curl -i  --data-binary "@scripts/ticket.txt" http://localhost:9098/ticket


echo ""
echo "Tickets in database:"

docker exec -t -u postgres postgres_dev psql -c "SELECT * FROM tickets"

echo ""
echo "Products in database:"

docker exec -t -u postgres postgres_dev psql -c "SELECT * FROM products"
