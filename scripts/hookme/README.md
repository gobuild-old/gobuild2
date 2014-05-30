## hookme
hookme is a program for easily manage hooks from github, gitlab or gogit.

## how to use

first is to install

	go get github.com/gobuild/gobuild2/scripts/hookme

run hookme program to start listening.

	./hookme

**restrict** this program is not designed for windows.

two script is needed.

`authchecker`, this script will be called by hookme with 
	
	authchecker $SECRET

`receiver` will called when authchecker not exit with 0.
	
	receiver $REPOPATH $ZIPBALL_URL

## this is program is part of gobuild2
see <https://github.com/gobuild/gobuild2> for more information.