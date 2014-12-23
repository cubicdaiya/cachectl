# cachectl

`cachectl` is a controller for regular file's page cache

## Install

```bash
go get -u github.com/cubicdaiya/cachectl
```

## Show Page Cache Stat For File

```
cachectl -f /var/log/access_log
```

## Delete Page Cache For File

```
cachectl -op del -f /var/log/access_log
```

If you want to leave a cache appended recently, assigning a rate for deleting page cache with `-r` is recommended.

```
cachectl -op del -f /var/log/access_log -r 0.9
```
