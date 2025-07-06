<div align="center">
  <a href="https://github.com/gkits/pavosql">
    <img src="assets/pavosql-gopher.png" alt="pavosql gopher" width="240" height="240">
  </a>
  <h1 align="center">PavoSQL</h1>
  <p align="center">
  </p>
</div>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build](https://github.com/gkits/pavosql/actions/workflows/build.yaml/badge.svg)](https://github.com/gkits/pavosql/actions/workflows/build.yaml)
[![Test](https://github.com/gkits/pavosql/actions/workflows/test.yaml/badge.svg)](https://github.com/gkits/pavosql/actions/workflows/test.yaml)

**This is project is still work in progress and not supposed to be used in any productive setting.**

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
