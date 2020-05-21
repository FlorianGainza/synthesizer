build-img :
	docker build -t synthesizer .

build-bin :
	docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp synthesizer go build -v -o bin/synthesizer cmd/synthesizer/synthesizer.go

exec-bin :
	docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp -p 8080:8080 synthesizer ./bin/synthesizer

mod-vendor :
	docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp synthesizer go mod vendor