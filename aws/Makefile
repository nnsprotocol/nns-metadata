FILE_TEMPLATE = template.yml
FILE_PACKAGE = packaged.yml
AWS_BUCKET_NAME = sam-deployments-packages
BUCKET_PREFIX = nns-metadata
AWS_REGION = "eu-west-1"
AWS_STACK_NAME = nns-metadata

build:
	@env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/api ../cmd/api/main.go

package:
	@sam package \
		--template-file $(FILE_TEMPLATE) \
		--s3-bucket $(AWS_BUCKET_NAME) \
		--s3-prefix $(BUCKET_PREFIX) \
		--region $(AWS_REGION) \
		--output-template-file $(FILE_PACKAGE)

deploy:
	@sam deploy \
		--template-file $(FILE_PACKAGE) \
		--region $(AWS_REGION) \
		--capabilities CAPABILITY_IAM \
		--stack-name $(AWS_STACK_NAME) \
		--force-upload

pipeline: build package deploy