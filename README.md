# PavoSQL

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build](https://github.com/gKits/PavoSQL/actions/workflows/gobuild.yaml/badge.svg)](https://github.com/gKits/PavoSQL/actions/workflows/gobuild.yaml)
[![Test](https://github.com/gKits/PavoSQL/actions/workflows/gotest.yaml/badge.svg)](https://github.com/gKits/PavoSQL/actions/workflows/gotest.yaml)
[![Build Hugo docs and deploy to pages](https://github.com/gKits/PavoSQL/actions/workflows/hugo.yaml/badge.svg)](https://gkits.github.io/PavoSQL)

A simple SQL database written in pure Go.

**This is a learning project and is not meant to be run in production environments.**  
**This project is stil w.i.p.**

## Roadmap

- [x] Atomic backend store on single file
- [ ] Relational model build on KV Store
- [ ] Lexer and Parser for SQL queries
- [ ] Database server and client to use PavoSQL over the network
- [ ] Implement [database/sql](https://pkg.go.dev/database/sql) driver interface
- [ ] Database Management System in single directory
- [ ] Windows compatibilty of backend store (remain atomic)
- [ ] Documentation
- [ ] Installable as service/daemon (e.g. systemd)
- [ ] Create and release Docker image
- [ ] 80% Test coverage (not needed but nice to have)
