build-redoc-cli:
	docker build -t smith-30/redoc-cli:0.10 ./docker/redoc-cli/.

gen-aws-status-redoc:
	docker run -v $(PWD)/awsstatus/apidoc:/data smith-30/redoc-cli:0.10 bundle --title aws-status-api --output /data/out/aws-status.html /data/reference/api.v1.yaml