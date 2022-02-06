# Traefik Plugin: Auth Delay

A [traefik](https://traefik.io/) plugin to add a random delay to failed authentication requests based on status code.

Makes password stuffing and brute force attacks harder for all services using this plugin / middleware by inserting a
random delay before the failed response is returned.

## Example Configuration

TODO

## What is a Traefik Plugin

TL;DR; A Traefik plugin is a custom middleware for Traefik.

[More on Traefik plugins is written here](https://doc.traefik.io/traefik/plugins/).

I also wrote [an init container](https://github.com/colearendt/traefik-plugin-init) that simplifies using "local"
plugins (i.e. plugins without Traefik Pilot) inside of Kubernetes.

## Thanks

Inspired by and much boilerplate
from [traefik-plugin-rewrite-headers](https://github.com/XciD/traefik-plugin-rewrite-headers), which is a fantastically
useful Traefik Plugin.
