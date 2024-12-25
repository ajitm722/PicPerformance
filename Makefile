# Define some variables
DOCKER_COMPOSE = sudo docker-compose
CURL = curl
HEYCMD = hey
URL = localhost:8080

# Docker-related targets
.PHONY: up down logs

up: 
	$(DOCKER_COMPOSE) up --build

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

# Make requests for image creation
.PHONY: create_image

create_image:
	@echo "Please enter the image details."
	@read -p "Image ID (integer): " id; \
	read -p "Image Format (e.g., JPEG, PNG): " format; \
	read -p "Image Resolution (e.g., 1920x1080): " resolution; \
	$(CURL) -d "{ \"id\": $${id}, \"format\": \"$${format}\", \"resolution\": \"$${resolution}\", \"img_status\": \"Processing\" }" -H "Content-Type: application/json" $(URL)/images

# Get metrics
.PHONY: get_metrics

get_metrics:
	@echo "Please enter the metric key you want to filter by (e.g. myapp_login_request_duration_seconds , myapp_request_duration_seconds , images):"
	@read -p "Metric Key: " metric_key; \
	$(CURL) localhost:8081/metrics | grep "$${metric_key}"


.PHONY: process_image

process_image:
	@echo "Please enter the image ID to update the status."
	@read -p "Image ID: " image_id; \
	$(CURL) -X PUT -d '{"img_status": "Processed"}' $(URL)/images/$${image_id}; \
	echo "Images left to process:"; \
	$(CURL) $(URL)/images

# Normal login
.PHONY: login

login:
	$(CURL) $(URL)/login

# Extreme load testing login
.PHONY: load_test_login

load_test_login:
	@echo "Please enter the values for load testing:"
	@read -p "Number of requests (e.g., 10): " num_requests; \
	read -p "Concurrency level (e.g., 1): " concurrency; \
	read -p "QPS (queries per second, e.g., 2): " qps; \
	$(HEYCMD) -n $${num_requests} -c $${concurrency} -q $${qps} http://$(URL)/login

# Clean up (docker-compose down)
.PHONY: clean

clean:
	$(DOCKER_COMPOSE) down