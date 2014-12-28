# cachectl

`cachectl` is a controller for regular file's page cache. 

## Dependency

`posix_fadvise` is required.

## Install

```bash
git clone https://github.com/cubicdaiya/cachectl.git
cd cachectl
make gom
make bundle
make
```

## Show Page Cache Stat For File

```
cachectl -f /var/log/access_log
```

## Purge Page Cache For File

```
cachectl -op purge -f /var/log/access_log
```

If you want to leave a cache appended recently, assigning a rate for purging page cache with `-r` is recommended.

```
cachectl -op del -f /var/log/access_log -r 0.9
```
