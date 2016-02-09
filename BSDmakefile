# BSDmakefile

.PHONY: all
all:
	pkg info -s libucl >/dev/null || (echo "first: sudo pkg install libucl" && false)
	go test
