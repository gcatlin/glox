all: build

build:
	@vgo build

escapes:
	@vgo build -a -gcflags "-m -m" 2>&1 | grep -v inlin | grep -v "does not escape"

run: build
	@glox

