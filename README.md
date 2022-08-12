# Thought Machine SRE Challenge

`cpxctl` is a CLI tool for querying the health of services running in Cloud Provider X and is my solution to the Thought Machine SRE take home challenge (see [brief](brief/SRE_take_home_challenge.txt)).

## Installation

`cpxctl` is written in Go and you will thus need the Go toolchain installed in order to build the binary. The given [Makefile](Makefile) contains some useful commands for working with `cpxctl`.

1. To install dependencies:
   
```
make deps
```

2. To build the `cpxctl` binary into a local `bin` directory:

```
make bin/cpxctl
```

3. To run tests:

```
make test
```

## Usage

Once the `cpxctl` binary is installed, start the CPX server by running:

```
make cpx-server
```

### `cpxctl` Commands

#### `list-instances`

```
$ bin/cpxctl list-instances --help
usage: cpxctl list-instances [<flags>]

Show the status of all running instances.

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -s, --service=SERVICE  Only show instances for a given service.
  -t, --tail             Tail periodic queries to list instances.
```

To view the status of all running instances:
  
```
$ bin/cpxctl list-instances
      IP      |      SERVICE       |  STATUS   |  % CPU  | % MEMORY
--------------+--------------------+-----------+---------+-----------
   10.58.1.1  |   StorageService   |  Healthy  | 15.00%  |  62.00%
   10.58.1.2  |    TimeService     |  Healthy  | 12.00%  |  53.00%
   10.58.1.3  |     MLService      | Unhealthy | 82.00%  |  94.00%
   10.58.1.4  |   StorageService   |  Healthy  | 21.00%  |  28.00%
   10.58.1.5  |   StorageService   | Unhealthy | 96.00%  |  56.00%
...
```

- Instances that are unhealthy will be highlighted in red

#### `list-services`

```
$ bin/cpxctl list-services --help
usage: cpxctl list-services [<flags>]

Show the status of all running services.

Flags:
      --help  Show context-sensitive help (also try --help-long and --help-man).
  -t, --tail  Tail periodic queries to list services.
```

To view the status of all services:

```
$ bin/cpxctl list-services
       SERVICE       | HEALTHY | % CPU AVG | % MEMORY AVG
---------------------+---------+-----------+---------------
     AuthService     |  14/17  |  55.65%   |    55.12%
      GeoService     |  16/17  |  53.65%   |    40.94%
      IdService      |  16/21  |  50.14%   |    54.48%
      MLService      |  14/16  |  44.50%   |    36.81%
  PermissionsService |  8/10   |  50.00%   |    63.30%
...
```

- Services that have less than 2 healthy instances will be highlighted in red

*Note: The `--tail` options work most of the time but are a little glitchy and might sometimes show output that's not formatted correctly.*

## Decisions & Trade-offs

Project structure:
- [`main.go`](main.go) contains the top-level CLI code that glues the various packages together
- The [`internal/commands`](internal/commands) package contains the business logic for each command, including displaying results to stdout
- The [`internal/domain`](internal/domain) package contains the core interface and data model for our business logic that the other packages rely on
- The [`internal/httpapi`](internal/httpapi) package contains a HTTP client implementation of the CPX API that satisfies our core domain interface

Using our domain [`MonitoringService`](internal/domain/cpx.go) interface, we decouple our CPX client implementation from our core business logic, allowing us to easily move to another cloud provider, like Cloud Provider Y (CPY), if needed - all we need to do is implement a client of the CPY API that satisfies the `MonitoringService` interface. This also allows us to easily test our code in [`internal/commands`](internal/commands) through the use of [fakes](internal/domain/domainfakes), where we can stub responses from the downstream API rather than needing to issue actual HTTP requests.

One trade-off of the current solution is that we are making downstream API requests sequentially rather than concurrently. For a large number of servers, this would not be ideal in terms of performance. As the number of servers in this example is in the hundreds, I think this is a reasonable choice to make for a simpler solution where we don't have to deal with concurrency.

Another downside of the above approach is that one failed `GetServer` call to the downstream API causes our entire command to return nothing and error out entirely. A more preferable solution would be to handle downstream errors more gracefully and return partial results to the user, perhaps even logging the error for debugging. I have chosen to leave it as is for now due to lack of time but it is a worthy future improvement we can make.

## Future Improvements

- [ ] Replace `github.com/gosuri/uilive` with a library that's less glitchy
- [ ] Make concurrent `GetServer` calls to the CPX server using goroutines rather than doing them sequentially
- [ ] Allow the user to configure the threshold (currently 2) for the minimum number of healthy instances in order for a service to be considered healthy
- [ ] Return partial results to the user when one or more `GetServer` calls to the CPX server fail to respond, instead of erroring out entirely
  - One option is to show an `Unknown` value to the user under each column, ie.
    ```
    10.58.1.5 | Unknown | Unknown | Unknown | Unknown 
    ```