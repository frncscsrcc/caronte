Caronte
===

*Caronte* allows to "open" webhooks to _target_ services that runs in neworks which do not allow traffic from the outside.

The two main componets of this application are:
- `agent`: it runs as a standalone application inside the network which blocks traffic from the outside. Using long-poll HTTP requests, it connects to the *proxy* and it waits for instructions. When the proxy receives a request, it "forward" it to the agent and the agent "forward" it to the target system. If this feature is enabled, when the target replies, the agent will forward back the response to the proxy, unsing a new connection, and the proxy will forward it back to the external caller.
- `proxy`: it runs as a standalone application outside the network which blocks traffic from the outside. It mantains a connection with one or more *agents* and it exponsed a public REST endpoint. When it receives a request from an external caller, it forward it to the correct agent and it mantains open the connection with the calle until the agent forwards back the reply from the target.

Build:
---

```
go build -o caronte main.go
```

Example:
---

Run the example echo service; this service listens listen on the port 5000 and send back to the caller the body of the request.

```
./caronte echo --delay 1
```

Run the proxy in a network that accept connections from the outside
```
./caronte proxy --secret SECRET
```

In the same machine that runs the echo server run the agent and connect it to the proxy

```
./caronte agent --secret SECRET --code EXAMPLE_TARGET --target-port 5000 --send-reply --proxy-host <proxy_host> --proxy-port <proxy_port>
```

Test the connection to the target, sending any request to the proxy to `<proxy_host>:<proxy_port>/forward/<target_code>/<path_on_the_target>`

```
curl -X POST --data "any data" <proxy_host>:<proxy_port>/forward/EXAMPLE_TARGET/whatever/path
```

The flow of data can be summarised as

```
REQUESTER -> PROXY -> [ -> AGENT -> TARGET -> AGENT -> ] -> PROXY -> REQUESTER
```

In reality all the connections between the proxy and agents are initiated and maintained by the agents (as outbound HTTP requests) and they are formally asynchronous, but from the requester prespective it is juas a synchronous call.


Authentication
---

The proxy request agent authentication. This is done via a shared SECRET. The secret can be passed directly in the `agent` and `proxy` cli, or - by default - it is used the value in the env `CARONTE_SECRET`