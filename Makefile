BINAME := stonks

build:
	go get
	go build -v -o ${BINAME}

install:
	go install -v

uninstall:
	rm -f ${HOME}/go/bin/stonks
