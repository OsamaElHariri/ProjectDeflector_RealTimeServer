# ProjectDeflector_RealTimeServer

This repo holds the code to manage websocket connections for the `Hit Bounce` mobile game.


## Mobile App

The mobile app, and a high level introduction to this project can be found in the [JsClientGame](https://github.com/OsamaElHariri/ProjectDeflector_JsClientGame) repo, which is intended to be the entry point to understanding this project.


## Overview of This Project

This is a Go server that uses the [Fiber](https://gofiber.io/) web framework.

Note that this project has a `.devcontainer` and is meant to be run inside a dev container.


## Outputting a Binary

To output the binary of this Go code, run the VSCode task using `CTRL+SHIFT+B`. This should be done while inside the dev container.


Once you have the binary, you need to build the docker image _outside_ the dev container. I use this command and just overwrite the image everytime. This keeps the [Infra](https://github.com/OsamaElHariri/ProjectDeflector_Infra) repo simpler.

```
docker build -t project_deflector/realtime_server:1.0 .
```

## Overview of Socket Management

 Each time a request comes in for a websocket connection, the listener for new events is run on one goroutine, and the sender that allows the server to send events to the client keeps the connection open (this is defined in the root `GET` route in `main.go`). The server then keeps track of which connection is for which user by keeping a map of user IDs and connections. There is an internal route `/internal/notify/:id` that other servers can use to send arbitrary data to whatever user needs to be notified.

 Small note: I did not put much effort into this server, while it is working fine, it is probably leaking memory because I did not put much effort into closing connections properly. So yeah, don't use this in production systems.