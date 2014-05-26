#! /usr/bin/env python
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#
# Copyright Date:2014-05 work <work@cp01-rdqa2014-01.vm.baidu.com>
#
# Distributed under terms of the MIT license.

"""
build go package
"""

import sh
import json

def main():
    output = sh.go.list('-json')
    s = json.loads(str(output))
    uniq = {}
    for dep in s.get('Deps'):
        if dep.startswith('github.com') or dep.startswith('code.google.com'):
            godep = '/'.join(dep.split('/')[:3])
            if uniq.get(godep):
                continue
            uniq[godep] = True
            print godep

main()
