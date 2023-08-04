# echo-server

## Overview

|| Arguemment || Description || Default Value ||s
| --- | --- | --- |
| port | The port to connect too. | 9001 |
| type | The type of server to run. Options are http, websocket, or all | all |

## Installation

```bash
brew tap jmoney/server-utils
brew install echo-server
```

## Run Locally

```bash
go run cmd/server/main.go -port 9002
```

Starts up on port `9002`.  You can use `ngrok` to expose the server to the internet if so desired with the command `ngrok http 9002`.
