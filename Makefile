run: 
	docker run pch-client

build: 
	docker build . -f Dockerfile -t pch-client
