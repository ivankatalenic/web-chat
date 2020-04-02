build:
	go build -v .

check:
	golint ./...
	go test -v -count=1 ./...

deploy:
	ssh -i deploy_key -o StrictHostKeyChecking=no "root@$(HOST)" 'systemctl stop web-chat.service; rm -rf /web-chat'
	rsync -r -e 'ssh -i deploy_key -o StrictHostKeyChecking=no' --files-from=deploy_files . "root@$(HOST):/web-chat"
	ssh -i deploy_key -o StrictHostKeyChecking=no "root@$(HOST)" 'systemctl start web-chat.service'
