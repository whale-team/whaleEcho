# WhaleEcho

WhaleEcho is websocket server for handling websocket connection in whale.

## Table of Content
1. **[Installation](#installation)**
2. **[System Design](#system-design)**


## Installation

#### Requirements:

- docker engine
- go 1.14+

#### Setup Steps:

1. Use the go module to get install this project

```bash
git clone https://github.com/whale-team/whaleEcho.git && cd ./whaleEcho
```

2. Setup Environment

```bash
make setup.env
```

3. Run the test

```bash
make test.all
```

4. run websocket server locally

``` bash
make run.ws
```


## System Design
#### System Flow  

![flow chart](./docs/uml/flowchart.drawio.svg)

#### Component Diagram

