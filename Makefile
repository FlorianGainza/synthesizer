build-img :
	docker build -t synthesizer .

build-bin :
	docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp synthesizer go build -v -o bin/synthesizer

exec-bin :
	docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp -p 8080:8080 synthesizer ./bin/synthesizer
