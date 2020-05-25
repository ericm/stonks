BINAME := stonks

build:
	go get
	go build -v -o ${BINAME}

web:
	go get
	cd stonks.icu && go build -v -o ${BINAME}

install:
	go install -v

uninstall:
	rm -f ${GOPATH}/bin/stonks
