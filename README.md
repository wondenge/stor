# Stor

`stor` provides a Go native library to store metadata.

The goal of `stor` is to abstract common store operations for multiple distributed and/or local Key/Value store backends.

As of now, `stor` offers support for `Consul`, `Etcd`, `Zookeeper` (**Distributed** store) and `BoltDB` (**Local** store).

## Usage

`stor` is meant to be used as an abstraction layer over existing distributed Key/Value stores. It is especially useful if you plan to support `consul`, `etcd` and `zookeeper` using the same codebase.

It is ideal if you plan for something written in Go that should support:

- A simple metadata storage, distributed or local
- A lightweight discovery service for your nodes
- A distributed lock mechanism

You can find examples of usage for `stor` under in `docs/examples.go`.

## Supported versions

`stor` supports:

- Consul versions >= `0.5.1` because it uses Sessions with `Delete` behavior for the use of `TTLs` (mimics zookeeper's Ephemeral node support), If you don't plan to use `TTLs`: you can use Consul version `0.4.0+`.
- Etcd versions >= `2.0` because it uses the new `coreos/etcd/client`, this might change in the future as the support for `APIv3` comes along and adds more capabilities.
- Zookeeper versions >= `3.4.5`. Although this might work with previous version but this remains untested as of now.
- Boltdb, which shouldn't be subject to any version dependencies.

## Interface

A **storage backend** in `stor` should implement (fully or partially) this interface:

```go
type Store interface {
	Put(key string, value []byte, options *WriteOptions) error
	Get(key string) (*KVPair, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Watch(key string, stopCh <-chan struct{}) (<-chan *KVPair, error)
	WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*KVPair, error)
	NewLock(key string, options *LockOptions) (Locker, error)
	List(directory string) ([]*KVPair, error)
	DeleteTree(directory string) error
	AtomicPut(key string, value []byte, previous *KVPair, options *WriteOptions) (bool, *KVPair, error)
	AtomicDelete(key string, previous *KVPair) (bool, error)
	Close()
}
```

## Compatibility matrix

Backend drivers in `stor` are generally divided between **local drivers** and **distributed drivers**. Distributed backends offer enhanced capabilities like `Watches` and/or distributed `Locks`.

Local drivers are usually used in complement to the distributed drivers to store informations that only needs to be available locally.

| Calls                 | Consul | Etcd | Zookeeper | BoltDB |
| --------------------- | :----: | :--: | :-------: | :----: |
| Put                   |   X    |  X   |     X     |   X    |
| Get                   |   X    |  X   |     X     |   X    |
| Delete                |   X    |  X   |     X     |   X    |
| Exists                |   X    |  X   |     X     |   X    |
| Watch                 |   X    |  X   |     X     |        |
| WatchTree             |   X    |  X   |     X     |        |
| NewLock (Lock/Unlock) |   X    |  X   |     X     |        |
| List                  |   X    |  X   |     X     |   X    |
| DeleteTree            |   X    |  X   |     X     |   X    |
| AtomicPut             |   X    |  X   |     X     |   X    |
| Close                 |   X    |  X   |     X     |   X    |

## Limitations

Distributed Key/Value stores often have different concepts for managing and formatting keys and their associated values. Even though `stor` tries to abstract those stores aiming for some consistency, in some cases it can't be applied easily.

Please refer to the `docs/compatibility.md` to see what are the special cases for cross-backend compatibility.

Other than those special cases, you should expect the same experience for basic operations like `Get`/`Put`, etc.

Calls like `WatchTree` may return different events (or number of events) depending on the backend (for now, `Etcd` and `Consul` will likely return more events than `Zookeeper` that you should triage properly). Although you should be able to use it successfully to watch on events in an interchangeable way (see the **docker/leadership** repository or the **pkg/discovery/kv** package in **docker/docker**).

## TLS

Only `Consul` and `etcd` have support for TLS and you should build and provide your own `config.TLS` object to feed the client. Support is planned for `zookeeper`.

## Roadmap

- Make the API nicer to use (using `options`)
- Provide more options (`consistency` for example)
- Improve performance (remove extras `Get`/`List` operations)
- Better key formatting
- New backends?
