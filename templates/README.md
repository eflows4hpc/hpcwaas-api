# Configuration template files

This folder contains two template config files, both in YAML format:

- `hpcwaas-api.yml` is the config file for the HPCWaaS server. It is mandatory in order to allow SSO authentication. It is managed by the sysadmin and the server needs to be restarted after each modification. Its default location is `/etc/hpcwaas-api/hpcwaas-api.yml`.  
- `.waas` is an optional config file for the `waas` command-line utility. It is managed by the users. Its default locatizon is `$HOME/.waas`.