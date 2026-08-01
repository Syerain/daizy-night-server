# Project Daizy Night 

Designed for ATOM Reforge Community's daily activities.

The project name Daisy Night is a homophone of The Hazy Night.

It's also my first golang project and I wished to learn something during the work.

## Repo

server & devCliTool
https://github.com/Syerain/daizy-night-server

scheduled client
https://github.com/Syerain/daizy-night-app

cli test client
https://github.com/Syerain/daizy-night-appcli

## Deployment

### Configuration

during alpha period there are some files for test placed under the root dir

before first launch, you must :

- edit file `./Config.yaml`

- delete dir `./test`

## Tech stack & packages

**golang 1.26.5**

Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.

**viper**

Viper is a complete configuration solution for Go applications including 12-Factor apps. 
It is designed to work within any application, and can handle all types of configuration needs anclsd formats.

**tinter**

Package tint implements a zero-dependency slog.Handler that writes tinted (colorized) logs. 
Its output format is inspired by the zerolog.ConsoleWriter and slog.TextHandler.

**echo**

high performance, minimalist Go web framework

**sqlite**

SQLite is a C-language library that implements a small, fast, self-contained, high-reliability, full-featured, SQL database engine.

**gorm**

GORM, a popular ORM library for Golang, provides a flexible plugin system that allows developers to extend its functionality. 

**argon2id**

This package provides a convenience wrapper around Go's argon2 implementation, making it simpler to securely hash and verify passwords using Argon2.

