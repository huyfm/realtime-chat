# Realtime-Chat (RTC)

This web application is written in Go, allows users to login to the server using Github accounts and chat with others. 

### Design choices:
This rtc server initially supports 1000 active users. We choose a monolith architecture which includes an OAuth2 authentication service, a chat service that uses go channels for message queues. All user data and messages are persisted in a PostgreSQL server. Requested data are cached in a Redis server.

This project is in early stage and might have expected changes during development.

### Features:
- [x] Authentication service.
- [ ] Chat service.
- [ ] Caching service.
- [ ] Convert to microservices.

### How to run in local machine?
You can run rtc on localhost with TLS enabled and HTTPS scheme. Use package [mkcert](https://github.com/FiloSottile/mkcert) to install a CA to the local machine's trust system and generate localhost.pem and localhost-key.pem files for the localhost. Then, store these files in tls/ folder. For convinience, you can comment out these settings in .env.dev file which disables TLS and allows the rtc to run with HTTP scheme.
```properties
# TLS Certificate.
TLS_CERT=localhost.pem
TLS_PRIVKEY=localhost-key.pem
```

Follow [this link](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) by Github to register rtc to Github OAuth2 service and put provided CLIENT_ID and CLIENT_SECRET in .env.dev file.
```properties
# Github OAuth2.
CLIENT_ID=00000000000000000000
CLIENT_SECRET=0000000000000000000000000000000000000000
```

Run the whole cluster including the go server and PostgreSQL with Docker.
```shell
# Run containers
docker compose up --build --env-file .env.dev -d
# Drop containers and associated docker resources
docker compose down -v
```

Access to this server at https://localhost:8080 and enjoy.
