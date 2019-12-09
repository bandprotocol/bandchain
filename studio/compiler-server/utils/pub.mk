APP ?= powerful-wave-45685
APP_START=node/bin/node .
SLUG_JSON=downloads/slug.json
GET_SLUG_URL=${shell node utils/get_json_value ${SLUG_JSON} blob.url}
GET_SLUG_ID=${shell node utils/get_json_value ${SLUG_JSON} id}

default: send

slug_json:
	mkdir -p downloads
	curl -X POST \
-H 'Content-Type: application/json' \
-H 'Accept: application/vnd.heroku+json; version=3' \
-d '{"process_types":{"web":"${APP_START}"}}' \
-n https://api.heroku.com/apps/${APP}/slugs -o ${SLUG_JSON}

send: slug_json
	curl -X PUT \
-H "Content-Type:" \
--data-binary @slug.tgz \
"${GET_SLUG_URL}"
	curl -v -X POST \
-H "Accept: application/vnd.heroku+json; version=3" \
-H "Content-Type: application/json" \
-d '{"slug":"${GET_SLUG_ID}"}' \
-n https://api.heroku.com/apps/${APP}/releases -o downloads/published.json

.PHONY: default send slug_json
