.PHONY: help run_dev run_prod clean put_mapping put_data drop_index trigger_github

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run_dev: clean
	./run_dev.sh > temp.log

run_prod:
	./run_prod.sh

clean:
	git checkout last_timestamp.txt
	rm -f temp.log

put_mapping:
	@echo "putting a new mapping"
	@curl -u ${ES_USERNAME}:${ES_PASSWORD} -XPUT ${LB_IP}/elastic/slack -H 'Content-Type: application/x-ndjson' --data-binary '@elasticsearch_utils/mapping.json'
	@echo "done puttng a new mapping"

put_data:
	@echo "putting data"
	@curl -u ${ES_USERNAME}:${ES_PASSWORD} -XPUT ${LB_IP}/elastic/slack/_bulk -H 'Content-Type: application/x-ndjson' --data-binary '@temp.log'
	@echo "done putting data"

drop_index:
	@echo "dropping index"
	@curl -u ${ES_USERNAME}:${ES_PASSWORD} -XDELETE ${LB_IP}/elastic/slack
	@echo "done dropping index"

trigger_github:
	curl -H "Accept: application/vnd.github.everest-preview+json" \
    -H "Authorization: token ${GITHUB_API_KEY}" \
    --request POST \
    --data '{"event_type": "menual-trigger"}' \
    https://api.github.com/repos/dl4ab/DFAB-Archiver-slackbot/dispatches
