#! /usr/bin/env python
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#
# Copyright Date:2014-05 work <work@cp01-rdqa2014-01.vm.baidu.com>
#
# Distributed under terms of the MIT license.

"""
figure out dependencies
"""

import sh
import os
import json

_goroot = ''

def goroot():
    global _goroot
    if not _goroot:
        _goroot = str(sh.go.env("GOROOT")).strip()
    return _goroot

def main():
    output = sh.go.list('-json')
    s = json.loads(str(output))

    for dep in s.get('Deps'):
        libpath = os.path.join(goroot(), 'src', 'pkg', dep)
        if os.path.isdir(libpath):
            #print 'stand lib %s, skip' %(dep)
            continue
        print 'go get -d -u ', dep

main()
