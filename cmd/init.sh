#!/bin/bash
set -euo pipefail

handlePrivateRepo(){
    # 可以使用通配符  aarioai/*  或者 多个  xxx/aarioai/golib,xxx/aarioai/xxx
    go env -w GOPRIVATE=github.com/aarioai/golib
}

main(){
    handlePrivateRepo
}

main