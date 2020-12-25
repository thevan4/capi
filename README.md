# Consul Api worker for lbos

Simple consul cli kv worker.

Created for [read and write configs for lbos](https://github.com/khannz/crispy-palm-tree)

## Help

```bash
consul api worker for lbos 😉

Usage:
  capi [command]

Available Commands:
  get         get services at consul root path
  help        Help about any command
  put         put service in consul

Flags:
  -s, --app-servers-folder string    folder and consul key-folder name for app servers. Example: app-servers (default "app-servers")
  -c, --config string                Path to config file. Example value: './capi.yaml' (default "./capi.yaml")
  -a, --consul-address string        consul address (default "127.0.0.1:8500")
  -r, --consul-root-path string      consul root path for lbos cluster. Example: lbos/t1-cluster-0/ (default "lbos/t1-cluster-0/")
  -t, --consul-timeout duration      consul timeout (default 2s)
  -d, --data-dir-path-names string   folder for json files for send to consul. Example: ./json-services (default "./json-services")
  -f, --force-update-keys            force update keys (bool)
  -h, --help                         help for capi
      --log-event-location           Log event location (like python) (default true)
      --log-format string            Log format. Example values: 'text', 'json' (default "text")
      --log-level string             Log level. Example values: 'info', 'debug', 'trace' (default "trace")
      --log-output string            Log output. Example values: 'stdout', 'syslog' (default "stdout")
  -m, --manifest-name string         manifest key name for service. Example: manifest (default "manifest")
      --syslog-tag string            Syslog tag. Example: 'trac-dgen'
```

## Read consul keys example

```bash
./capi get -a 10.10.10.10:1234 -r lbos/t1-cluster-22/
```

## Put consul keys examples

Don't update keys:

```bash
./capi put -a 10.10.10.10:1234 -r lbos/t1-cluster-22/ -d ./new-servers
```

Force update keys:

```bash
./capi put -a 10.10.10.10:1234 -r lbos/t1-cluster-22/ -d ./new-servers -f true
```
