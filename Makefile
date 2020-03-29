build:
	go build -v .

deploy:
	echo $(DEPLOY_KEY) > deploy_key
	chmod 600 deploy_key
	rsync $ARGS --delete --force -r -e 'ssh -i deploy_key -o StrictHostKeyChecking=no' --files-from=deploy_files root@$(HOST):/web-chat

restart-service:
	echo $(DEPLOY_KEY) > deploy_key
	chmod 600 deploy_key
	ssh -i deploy_key -o StrictHostKeyChecking=no root@$(HOST) systemctl restart web-chat.service
