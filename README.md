# Relase notes parser

This tool was created to cover my personal pain points

1) Plugin dependency management

its hard to track plugin dependencies. you install one plugin and it requires 2nd plugin and 2nd plugin requires 5 more and so on
each plugin has required version of each dep.
during upgrade of 1 plugin you can miss to upgrade dependent plugin

2) Release notes

each time when i upgrade 1 plugin jenkins i might have to upgrade all deps of this plugin to satisfied versions
and sometimes those versions have braking changes or just important features and i have to manually go and scroll 50 repos release notes in github

I DON'T WANT TO SCROLL 50+ PROJECTS/ 10+ RELEASE NOTES in each project. WANT EVERYTHING ON ONE PAGE :)


# Instructions


Plugin Manager - main page. you can add plugins which you have to it

supports add one plugin or list of plugins

Plugin Dependencies - get all plugins dependencies for plugin manager list,
will update plugin if needed and will show release notes for it

Updated Plugins - plugins from "Plugin Dependencies" tab, but you can edit it for use caseses like

* you want to check release notes between plugin version A and B, so you will add version A to Plugin Manager and version B to Updated Plugins

as result you will see diff in Updated Plugins

`Get Txt file` - option button will give you file with list of plugin manager plugins/updated plugins


## Log levels

set via `RN_DEBUG` environment variable. Can be

* panic
* fatal
* error
* warn or warning
* info
* debug
* trace


## Deploy

```
docker-compose up -d
```

TODO: should be deployment to docker/k8s

## Github personal token

by default github has limit to 60 request per hour

you can increase this limit with personal token

1. go to page https://github.com/settings/tokens (you should be autorized)

2. generate personal classic token with permission:
```
public_repo
```
3. set env var GITHUB_PERSONAL_TOKEN=<YOUR_TOKEN> or .env file

will take env var and if not set .env file else will use public api without token (60 requests per hour)

## Local development


### docker-compose
1) build images
```
make build-images
```
2) execute:
```
docker-compose up
```

ps. if you using linux you can go to docker-compose.yml file and
uncomment build sections for "controller" and "web"
comment "image" sections

and do
```
docker-compose up
```
without images prebuild commands

i use windows and it has issues with "build" section in docker-compose which buils images too long.

#### main services without docker-compose

run compose
```
docker-compose -f docker-compose-limited.yml up
```

1st terminal

```
cd frontend && npm install && npm start
```

2nd terminal

```
go run .
```

TODO: use skaffold+k8s in future?

https://skaffold.dev/

https://kubernetes.io/



