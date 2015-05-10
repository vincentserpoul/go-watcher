watcher
=======

dev tools that watches file changes on a folder and recompile live

1. Install gometalinter

```
https://github.com/alecthomas/gometalinter
```

```
go build;go install
```

Then from your app folder:

```
watcher -pkg="catalog-api" -env="dev"
``` 
