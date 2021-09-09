#!/bin/bash

docker run -it --rm \
  --publish 8000:3000 \
  -e AUTHN_URL=localhost:8000 \
  -e APP_DOMAINS=localhost \
  -e DATABASE_URL=sqlite3://:memory:?mode=memory\&cache=shared \
  -e SECRET_KEY_BASE='secret' \
  -e HTTP_AUTH_USERNAME=feelguuds \
  -e HTTP_AUTH_PASSWORD=feelguuds \
  --name authentication_service \
  keratin/authn-server:latest \
  sh -c "./authn migrate && ./authn server"
