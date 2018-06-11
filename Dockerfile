FROM golang:1.10.2-stretch as build

RUN go get github.com/miekg/dns && \
    go build github.com/miekg/dns

WORKDIR /go/src/facua.org/compartidosd/common
ADD common/src .
RUN go build

WORKDIR /go/src/facua.org/compartidosd/server
ADD server/src .
RUN go build
RUN go install

WORKDIR /go/src/facua.org/compartidosd/client
ADD client/src .
RUN go build
RUN go install

WORKDIR /go/src/facua.org

FROM debian:stretch as package

### APPLICATION RUNTIME ARGUMENTS ###

# The address where the server is meant to be located. The client will connect
# to that address to retreive the share list.
ARG server_address="http://example.com:10000"
# The name of the SMB shared folder in each computer.
ARG shared_folder_name="Shared"
# The tick interval, in milliseconds. The application will perform a tick, then
# wait this amount of time. See client/src/app/tick.go for more information.
ARG tick_interval_ms="30000"
# The folder where the shares will be mounted. This folder will be created on
# startup, and removed on shutdown.
ARG network_folder="/Network"
# The UNIX group which will have ownership of the network folder (the owner
# user will be root). All users who need to access this folders must be members
# of that group.
# After installation of the .deb package, the group will be automatically
# created, and all users that have their home directory under /home will be
# added to it.
ARG network_group="compartidosd"

### DEBIAN PACKAGE ARGUMENTS ###

# The name of the generated .deb package.
ARG debian_package_name="compartidosd"
# The version of the generated .deb package.
# See https://www.debian.org/doc/debian-policy/#version
ARG debian_package_version="0.0.1-1"

WORKDIR /pkg

COPY client/pkg .
COPY --from=build /go/bin /usr/local/bin

RUN mkdir -p usr/local/bin/org.facua && \
    cp /usr/local/bin/client usr/local/bin/org.facua/compartidosd
COPY --from=build /go/bin/client usr/local/bin/org.facua/compartidosd

# Replace the placeholders with their values
RUN f="etc/systemd/system/compartidosd.service" && \
    cat "$f" | \
    sed "s|&PACKAGE_VERSION&|$debian_package_version|g" | \
    sed "s|&SERVER_ADDRESS&|$server_address|g" | \
    sed "s|&SHARED_FOLDER_NAME&|$shared_folder_name|g" | \
    sed "s|&TICK_INTERVAL_MS&|$tick_interval_ms|g" | \
    sed "s|&NETWORK_FOLDER&|$network_folder|g" | \
    sed "s|&NETWORK_GROUP&|$network_group|g" \
        > /tmp/x && mv /tmp/x "$f"

RUN f="DEBIAN/control" && \
    cat "$f" | \
    sed "s|&NAME&|$debian_package_name|g" | \
    sed "s|&VERSION&|$debian_package_version|g" \
        > /tmp/x && mv /tmp/x "$f"

RUN f="DEBIAN/postinst" && \
    cat "$f" | \
    sed "s|&NETWORK_GROUP&|$network_group|g" \
        > /tmp/x && mv /tmp/x "$f" && \
    chmod 0755 "$f"

RUN dpkg-deb --build -Z gzip /pkg

RUN echo "#!/bin/sh" > /usr/local/bin/export-app && \
    echo \
    "cp /pkg.deb /out/${debian_package_name}_${debian_package_version}.deb" \
        >> /usr/local/bin/export-app && \
    echo "cp /usr/local/bin/server /out" >> /usr/local/bin/export-app && \
    echo "cp /usr/local/bin/client /out" >> /usr/local/bin/export-app && \
    chmod +x /usr/local/bin/export-app

CMD [ "export-app" ]
