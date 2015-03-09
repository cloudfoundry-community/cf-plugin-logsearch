Overview
========

Fetches logs from logstash services which uses elastic search

Installation
------------

```
$ go get github.com/cloudfoundry-community/cf-plugin-logsearch
$ cf install-plugin $GOPATH/bin/logsearch
```

Usage
-----

```
$ cf search-logs <appname>
```

Development
-----------

```
cf uninstall-plugin logsearch; go get ./...; cf install-plugin $GOPATH/bin/cf-plugin-logsearch
```
