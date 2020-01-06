docker:
	docker build -t stobias123/terraform-provider-orion . && docker push stobias123/terraform-provider-orion
local:
	rm example/terraform-provider-orion || true
	go build -o example/terraform-provider-orion .
tftest:
	rm example/terraform.tfstate* || true
	rm -rf example/.terraform || true
	cd test
	terraform init && terraform apply -auto-approve
clean:
	rm example/terraform.tfstate* || true
	rm -rf example/.terraform || true
	rm example/terraform-provider-orion || true
	rm example/log || true
