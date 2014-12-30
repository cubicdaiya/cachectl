# cachectl

`cachectl` is a controller for regular file's page cache. 

## Dependency

`posix_fadvise` is required.

## Install

```bash
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

# cachectld

`cachectld` is a daemon for scheduled purging page cache. Its behavior is described by [TOML](https://github.com/toml-lang/toml).

```
cachectld -c conf/cachectld.toml
```

## Configuration for cachectld

A configuration file format for `cachectld` is [TOML](https://github.com/toml-lang/toml).

A configuration for `cachectld` has one or multiple targets. A example is [here](conf/cachectld.toml).

|name          |type  |description                                  |default|note                                           |
|--------------|------|---------------------------------------------|-------|-----------------------------------------------|
|path          |string|target file path                             |       |directory or file path                         |
|purge_interval|int   |interval for purging page cache for file     |0      |unit is second                                 |
|filter        |string|filtering pattern string for target file path|.*     |regular expression with golang's regexp package|
|rate          |float |rate of puring page cache for file           |1.0    |0.0 < rate <= 1.0                              |
