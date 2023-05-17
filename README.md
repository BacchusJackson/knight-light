# Knight Light Mono Repo

<!--toc:start-->

- [Knight Light Mono Repo](#knight-light-mono-repo)
  - [Counter API](#counter-api)
  - [Chat API](#chat-api)
  - [Development](#development)
  <!--toc:end-->

The Knight Light Mono Repo demonstrates an application that has multiple micro-services and the front end application in a single repository.
This project can used as a guide to show how other similar applications can be deployed to the batCAVE.

## Counter API

The Counter micro service simply maintains a global, in memory count, which is incremented using the `/add` endpoint.
The front end is setup to receive a Sever Side Event stream to update the counter without needing to refresh the page.

The Counter app has 3 routes:

| Route     | Accepted Methods | Ussage                                        |
| --------- | ---------------- | --------------------------------------------- |
| `/status` | GET              | Returns the JSON response object              |
| `/health` | POST             | Returns a health check JSON Response object   |
| `/add`    | POST             | Receives JSON object to update the count      |
| `/sse`    | GET              | Sever Side Event stream of the global counter |

## Chat API
The Chat micro service is a very simple demonstration of Web Sockets.
The only route for this is `/ws` which starts a web socket connection to the API.

The first prompt will ask for a username which is registered on the server and used as the handle for new messages.
Messages received from a client will be immediately sent to all other active users.


## Development

For local development there are two docker-compose files.
`docker-compose.yml` will build the apps using the Dockerfile in the sub directory for each service.
`local.docker-compose.yml` will expect a binary to exist in the `bin` directory of each sub directory.
This method is faster but requires the user to run `make build-alpine` in the root directory to pull in updates to the servers.
The frontend will be volume mounted so any changes will be reflected when the application is refreshed.

[Local Website](http://localhost:8081)
