---
title: ${COMMAND}(1) ${RELEASE_VERSION} | Manual
version: ${RELEASE_VERSION}
date: ${RELEASE_DATE}
adjusting: both
---
NAME
====

**gocp(1)** -- posix cp in go

SYNOPSIS
========

| **gocp** `[options]` source_file target_file
| **gocp** `[options]` source_file ... target

DESCRIPTION
===========

The **cp** utility in go, per the [POSIX specification](http://pubs.opengroup.org/onlinepubs/000095399/utilities/cp.html).

OPTIONS
=======

These options control how `gocp` operates.

`-R`
  ~ Copy recursive.

`-r`
  ~ Copy recursive (OB).

WWW
===

${PACKAGE_URL}
