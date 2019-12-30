docker:
	docker build -t stobias123/terraform-provider-orion . && docker push stobias123/terraform-provider-orion
local:
	rm test/terraform-provider-orion || true
	go build -o test/terraform-provider-orion .
tftest:
	rm test/terraform.tfstate* || true
	rm -rf test/.terraform || true
	cd test
	terraform init && terraform apply -auto-approve