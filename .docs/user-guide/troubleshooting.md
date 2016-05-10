# Troubleshooting

---

## Log file location

Polly keeps a log file in `/var/log/polly/polly.log`.

## Debug output

Polly can be started with the `-l debug` flag to maximize log detail

To watch polly live, including error output and traces, you can run it from a
command line using the following command.

```shell
$ sudo polly start -f -l debug
```

Following this you can either run another `polly` process to instantiate
commands or make requests from a libStorage compatible client.
