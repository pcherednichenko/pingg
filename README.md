# Pingg

Using this golang project you can simply debug ping to Stadia compare to you router
to determine where is the problem.

This application sends ping requests to the stage and to your router, 
and then builds graphs.

## How to run

```bash
go run main.go -routerIP=192.168.0.1
```

in routerIP you should specify your router IP
You need to have [Golang](https://golang.org/doc/install) installed to run application
To stop press ctrl + c

## Result

As result graph with results will open in your browser:

![graph](/.github/graph.png)