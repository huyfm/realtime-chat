# Realtime-Chat (RTC)

This web application is written in Go, allows users to login to the server using Github accounts and chat with others. 

### Design choices:
This rtc server initially supports 1000 active users. We choose a monolith architecture which includes an OAuth2 authentication service, a chat service that uses go channels for message queues. All user data are persisted in a PostgreSQL server and messages are persisted in a Cassandra server. Requested data are cached in a Redis server.

This project is in early stage and expects abrupted changes during development.

### Features:
- [x] Authentication service.
- [ ] Chat service.
- [ ] Last seen service.
- [ ] Cache service.

### How to run on local machine?
You can run rtc on localhost with TLS enabled and HTTPS scheme. Use package [mkcert](https://github.com/FiloSottile/mkcert) to install a CA to the local machine's trust system and generate localhost.pem and localhost-key.pem files for the localhost. Then, store these files in tls/ folder. For convinience, you can comment out these settings in .env.dev file which disables TLS and allows the rtc to run with HTTP scheme.
```properties
# TLS Certificate.
TLS_CERT=localhost.pem
TLS_PRIVKEY=localhost-key.pem
```

Follow [this link](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app) by Github to register your rtc to Github OAuth2 service and add provided CLIENT_ID and CLIENT_SECRET in .env.dev file.
```properties
# Github OAuth2.
CLIENT_ID=00000000000000000000
CLIENT_SECRET=0000000000000000000000000000000000000000
```

Run the whole cluster where each component runs in an isolated Docker container.
```shell
# Run containers
docker compose up --build --env-file .env.dev -d
# Drop containers and associated docker resources
docker compose down -v
```

Access to your local server at https://localhost:8080 and enjoy.

### How to host on local machine?
We can either use Port Forwarding or a tunnel service such as CloudFlare's.
With Port Forwarding, set the port fowarding rule in your router to forward incoming traffic from router's port 443/80 for https/http to {your-local-ip-address}:8080.
