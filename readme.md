# socketcluster-client-go
Editted from `sacOO7/socketcluster-client-go`:
* Add a marshaller compatible with
[JS sc-codec-min-bin](https://github.com/SocketCluster/sc-codec-min-bin)
* Fix concurrent io in `scclient/event_listener.go`
* Rewrite `gowebsocket` to add feature: read write deadline,
notify when disconnected suddenly 
