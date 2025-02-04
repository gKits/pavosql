# PavoSQL

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build](https://github.com/pavosql/pavosql/actions/workflows/build.yaml/badge.svg)](https://github.com/pavosql/pavosql/actions/workflows/build.yaml)
[![Test](https://github.com/pavosql/pavosql/actions/workflows/test.yaml/badge.svg)](https://github.com/pavosql/pavosql/actions/workflows/test.yaml)

PavoSQL is a SQL database written purely in Go.

## Roadmap

- [ ] Database engine
  - [ ] Single file backend
    - [ ] B+tree structure
    - [ ] Concurrent r/w
    - [ ] Atomic i/o
  - [ ] SQL
    - [ ] Relational model
      - [ ] Tables
      - [ ] Indexes
      - [ ] Metadata
    - [ ] Lexer, parser and AST
    - [ ] Query functionality
- [ ] Network
  - [ ] Server + client
  - [ ] Authentication + Authorization
- [ ] Implement [database/sql](https://pkg.go.dev/database/sql) driver interface
- [ ] Documentation
- [ ] Docker image

## Reference material

> **Build your own Database from Scratch**  
> by James Smith  
> [Book](https://build-your-own.org/database/)

> **GoSQL / Writing a SQL database from Scratch**  
> by Phil Eaton  
> [Blog](https://notes.eatonphil.com/database-basics.html)  
> [Repo](https://github.com/eatonphil/gosql)

> **Lexical Scanning in Go**  
> by Rob Pike  
> [Video](https://www.youtube.com/watch?v=HxaD_trXwRE)  
> [Slides](https://go.dev/talks/2011/lex.slide)

> **SQLite Documentation**  
> by SQLite  
> [Docs](https://sqlite.org/docs.html)
