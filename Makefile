configure:
	cd /tmp && rm -f -R semverbuild && git clone https://github.com/karamani/semverbuild && cp semverbuild/svbuild $$GOPATH/bin/svbuild

compile:
	goimports -w ./*.go
	go tool vet ./*.go
	golint
	go test
	go install

build:
	$(eval newver := $(shell (svbuild -$(VER_LVL) | tail -n 1)))
	@echo $(newver)

	if [ -z "$(newver)" ] ; then \
		echo "Can't make new version."; \
	fi

patch: export VER_LVL = l3
patch: compile build

minor: export VER_LVL = l2
minor: compile build

major: export VER_LVL = l1
major: compile build
