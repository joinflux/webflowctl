# Webflow Command Line Utility

Webflowctl is a command line utility to manage Webflow site configurations. It is being used specifically to manage webhooks.

## Installation

Download a binary from the [Release](https://github.com/joinflux/webflowctl/releases) page.

## Usage

```
❯ ./webflowctl
A tool to help manage webhooks in the Webflow API

Usage:
  webflowctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  webhooks    Manage webhooks

Flags:
  -a, --api-token string   Webflow API Token
  -h, --help               help for webflowctl

Use "webflowctl [command] --help" for more information about a command.
```

## Available Commands

- [x] List webhooks
- [x] Create webhook
- [x] Delete webhook
- [x] Get webhook

## Development

This tool was created with [Cobra CLI](https://cobra.dev/) so please make sure to leverage it when adding new functionality.
