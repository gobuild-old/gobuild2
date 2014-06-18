## gobuild2
[![wercker status](https://app.wercker.com/status/33c73c9c4ea5cbc96ca1660d2e1b58a6/m "wercker status")](https://app.wercker.com/project/bykey/33c73c9c4ea5cbc96ca1660d2e1b58a6)

<http://beta.gobuild.io> is a online golang code build service. 
Use this service, you can easily share binary file with your friend.

Also you can got golang lib source package(including dependency) download.

### How to add you repository to `gobuild.io`
for example: if your repo is `github.com/beego/bee`

open your browser, visit <http://beta.gobuild.io/github.com/beego/bee>, then the project will be added automaticly.

### The command line tool(got) for gobuild
how to install

	bash -c "$(curl -s gobuild.io/install_got.sh)"

### Hooks
[click to enter readme](scripts/hookme/README.md)

### gobuild setting .gobuild.yml
sample:

	filesets:
	  includes:
	  - README.md
	  - LICENSE
	  excludes:
	  - .*.go
	settings:
	  targetdir: ""
	  addopts: ""
	  yaml"cgoenable": false

### NOTICES
developing now...

### LICENSE
gobuild2 is under [MIT License](/LICENSE)

### Overview
* See [Trello Board](https://trello.com/b/Ml7fV574/gobuild2) to follow the develop team.

### ChangeLog
    - 2014-06-07 support golang lib source package
    - 2014-06-05 support source code package download
