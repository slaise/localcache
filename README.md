# Golang-localcache

A local cache product using Go. The code implements a simple in-memory cache with limit size, and supports ttl.

## Quick Start

```Go
cache, err := NewCache()

if err != nil {
  log.Fatal("Failed to create the local cache", err)
}

cache.Put("1",1)

i, err := cache.Get("1)

if err != nil {
  log.Fatal("Failed to get key from cache", err)
}

```
