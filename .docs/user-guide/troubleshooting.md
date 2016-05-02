# Troubleshooting

---

## Log file location

Polly keeps a log file in `/var/log/polly/polly.log`

## Debug output

Polly can be started with the `-l debug` flag to maximize log detail

To capture watch polly live, including error output and traces, you can run it from a command line using

```shell
sudo ~/work/bin/polly service start -f -l debug
```
