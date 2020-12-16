.PHONY: all

all: clean build init plan apply

start-nomad:
	systemd-run --user /usr/bin/nomad agent -dev -acl-enabled

clean:
	rm -rf .terraform.lock.hcl .terraform/ *.tfstate*
	rm -rf ~/.terraform.d/local/adriennecohea/nomadutility dist/

build:
	go build
	mkdir -p ~/.terraform.d/local/adriennecohea/nomadutility/0.0.3/linux_amd64
	cp terraform-provider-nomadutility ~/.terraform.d/local/adriennecohea/nomadutility/0.0.3/linux_amd64

init:
	terraform init

plan:
	terraform plan

apply:
	terraform apply
