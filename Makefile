.PHONY: all

all: clean build init plan apply

start-nomad:
	systemd-run --user --unit local-development-nomad /usr/bin/nomad agent -dev

stop-nomad:
	systemctl --user stop local-development-nomad

restart-nomad:
	systemctl --user restart local-development-nomad

logs-nomad:
	journalctl --user --follow --unit local-development-nomad

clean:
	rm -rf .terraform.lock.hcl .terraform/ *.tfstate*
	rm -rf ~/.terraform.d/local/adriennecohea/nomadutility

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
