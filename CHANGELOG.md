# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned for 0.2.0

- `Changed` Command structure is now owl verb object, ex: `owl create user` instead of `owl user create` (#6)
- `Changed` Usage of `realm` command is simplified, ex: `owl realm ls` becomes `owl realms` and `owl realm set` becomes `owl realm` (#6)
- `Changed` Usage of `unit` command is simplified, ex: `owl unit use` becomes `owl unit` (#6)
- `Added` New verbs (`update`, `upsert`, `append`, `remove`, `get`) and aliases (#6)
- `Added` Unit selection is now optional, use `owl unit -` or `owl unit default` to unselect current unit (#6)
- `Added` Configuration for LDAP backend via environment variables (`OWL_BACKEND_SUBUNIT=false` to disable subunits, `OWL_BACKEND_SUBUNIT_USER` and `OWL_BACKEND_SUBUNIT_GROUP` to customize subunits names)

## 0.1.0 "MVP"

### Owl CLI

- `Added` LDAP backend connector
- `Added` Create, update, delete and list realms
- `Added` Create, update, delete and list units
- `Added` Create, update, delete and list users
- `Added` Create, update, delete and list groups
- `Added` Add and remove users from groups
- `Added` Export one or all units
- `Added` Import units
- `Added` Set user password available hash formats : SHA224 SSHA224 SHA256 SSHA256 SHA384 SSHA384 SHA512 SSHA512 SHA3-224 SSHA3-224 SHA3-256 SSHA3-256 SHA3-384 SSHA3-384 SHA3-512 SSHA3-512 SHAKE128 SSHAKE128 SHAKE256 SSHAKE256
- `Added` Password sent by mail to user when changed
- `Added` JSON, YAML, Table outputs for all commands
- `Added` JSON input for all commands
