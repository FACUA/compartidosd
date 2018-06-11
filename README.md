# compartidosd

`compartidosd` (from Spanish: "shared files daemon"), is an UNIX daemon that
mounts SMB shared folders over the network in real time.

## How it works

`compartidosd` consists of a server and a client.

The server is a lighweight binary that sends the clients the list of computers
offering shared folders over the network, and the client is a daemon that
mounts them to a "Network" folder whenever they go online, and unmounts them
after they disconnect.

The clients will periodically query the server to keep their share list up to
date, but the only time this connection is mandatory is when it boots for the
first time. Afterwards, the clients can perform just fine even if the server is
offline.

The server acts as an authority. If a share is added, it is added across all the
clients, and if a share is removed, it is removed across all as well. Therefore
the clients will not mount any available shared folder on the network, if it
hasn't been approved by the server first.

## Tested systems

The server should work on any UNIX system.

The client has been tested on:

* Ubuntu 18.04
* Xubuntu 18.04

It may not work on other systems.

## Why compartidosd

`compartidosd` might be useful in offices with a medium to large number of
employees, where the workflow for file sharing is using a distributed network
with a protocol like SMB.

At FACUA, we found that, while SMB usually worked fine, the network discovery
was incredibly unreliable on both Windows and Linux systems. So we wrote
`compartidosd` in order to solve this problem.

On distributions using the [Nautilus](https://github.com/GNOME/nautilus)
file manager such as Ubuntu, we also found that it wasn't able to mount 
network locations on the "Save file" dialog, which also was a minor workflow
disruption, because you had to open another Nautilus instance to mount the
location first, then go back to the original dialog to save the file.

Becuase `compartidosd` keeps the online shares mounted in real time, this issue
is also solved.

## How to build

This guide contains instructions on how to build the project, and pakcage it
for Debian system

In order to build the project, you must have [Docker](https://www.docker.com)
installed in your system.

Start by cloning the repo:

```bash
$ git clone https://github.com/FACUA/compartidosd
```

And then build the Docker image:

```bash
$ cd compartidosd
$ docker build \
    --build-arg server_address="http://example.com:10000" \
    --build-arg shared_folder_name="Shared" \
    --build-arg tick_interval_ms="30000" \
    --build-arg network_folder="/Network" \
    --build-arg network_group="compartidosd" \
    --build-arg debian_package_name="compartidosd" \
    --build-arg debian_package_version="1.0.0-1" \
    -t compartidosd .
```

The Dockerfile accepts multiple build arguments. You can find their
documentation on it. These arguments allow to customize some runtime variables,
and some variables of the Debian package of the client.

After the image is built, you may extract the build artifacts to your
filesystem:

```bash
$ docker run -v $(pwd)/build:/out compartidosd
```

This will create a `build` folder with the following artifacts:

* `client`: the naked client executable, in case you want to package it
manually (for distributions other than Debian, for example).
* `compartidosd_<version>.deb`: the Debian package file containing the client.
It will install the client executable, along with a
[Systemd](https://www.freedesktop.org/wiki/Software/systemd/) service that
will start the deamon automatically on system boot, and will stop it on
shutdown.
* `server`: the server executable.

## How to use

### Using the server

The `compartidosd` server takes an index file to watch for changes, a port to
listen for HTTP connections, and optionally a DNS server to translate host's
domain names to IP addresses.

The index file is formatted like so:

```json
[
	{
		"Name": "Alice",
		"Host": "alice.company"
	},
	{
		"Name": "Bob",
		"Host": "bob.company"
	},
]
```

Or if you're not using DNS to resolve domain names:

```json
[
	{
		"Name": "Alice",
		"Host": "192.168.0.100"
	},
	{
		"Name": "Bob",
		"Host": "192.168.0.101"
	},
]
```

This, however, will require that your clients are not using DHCP, or that your
DHCP is configured to reserve static IP addresses for them.

To run the server, upload it to your server and simply execute it:

```bash
$ scp build/server ubuntu@my.server.com:~
$ ssh ubuntu@my.server.com
$ chmod +x server
$ sudo mv server /usr/local/bin/compartidosd-server
$ compartidosd-server \
        --index /path/to/your/index.json \
        # You may omit this one
        --dns 192.168.0.1 \
        --port 10000
```

The server is provided as-is, with no packaging whasoever. The implementation
of monitoring, supervising and logging and HTTP security is up to you. We recommend to create a Systemd unit file for the server too. You may also want
to add a reverse proxy like [Nginx](https://nginx.org) to add TLS support.

### Using the client

On Debian based distributions, the client is much more easy to use. Just
install the package, and the daemon will start immediatly:

```bash
$ sudo apt install ./build/compartidosd-<version>.deb
```

Or better yet, create your own APT repo with something like
[Aptly](https://www.aptly.info) and install it from there.
