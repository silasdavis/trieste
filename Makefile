CAPNP_FILES = $(shell find . -type f -name '*.capnp' -print)
CAPNP_GO_FILES = $(patsubst %.capnp, %.capnp.go, $(CAPNP_FILES))
CAPNP_GO_FILES_REAL = $(shell find . -type f -name '*.capnp.go' -print)


.PHONY: test
test:
	go test ./...



## compile hoard.proto interface definition
%.capnp.go: %.capnp
	go mod vendor
	capnp compile -I./vendor/zombiezen.com/go/capnproto2/std -ogo $<

.PHONY: capnp
capnp: $(CAPNP_GO_FILES)

.PHONY: clean_capnp
clean_capnp:
	@rm -f $(CAPNP_GO_FILES_REAL)
