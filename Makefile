build:
	cd ui && npm run build
	packr build

deploy: build
	GOOS=linux GOARCH=amd64 packr build
	ssh fileserver.grz "systemctl stop photogallery.service"
	scp photogallery-api fileserver.grz:/usr/local/bin/photogallery
	ssh fileserver.grz "systemctl start photogallery.service"
	ssh fileserver.grz "journalctl -fu photogallery.service"