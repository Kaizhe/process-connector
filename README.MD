# Process Connector

A lightweight process connector agent to emit process execution information enriched with container information.

## High Level Design

- Listen on process execution event from netlink

- Enrich with container information (container ID, image name etc.)

- Emit enriched information

## TODO List

- Add MD5 value

- Add K8s related information from kubelet
