build:
	go build -v .

deploy:
	echo "$DEPLOY_KEY" > deploy_key
	chmod 600 deploy_key
	ssh -i deploy_key -o StrictHostKeyChecking=no "root@$HOST" 'systemctl stop web-chat.service; rm -rf /web-chat'
	rsync -r -e 'ssh -i deploy_key -o StrictHostKeyChecking=no' --files-from=deploy_files . "root@$HOST:/web-chat"
	ssh -i deploy_key -o StrictHostKeyChecking=no "root@$HOST" 'systemctl start web-chat.service'
