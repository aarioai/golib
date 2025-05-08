#!/bin/bash
set -euo pipefail

handlePrivateRepo(){
    # 可以使用通配符  aarioai/*  或者 多个  xxx/aarioai/golib,xxx/aarioai/xxx
    go env -w GOPRIVATE=github.com/aarioai/golib
    # 通过 https://github.com/settings/tokens  -> Generate new token (classic)
    # 开启 repo 权限即可 --> 其他的通过 ssh 本机权限操作
    git config --global url."https://aarioai:<token>@github.com/".insteadOf "https://github.com/"

    git config --global --list

    # 使用这个检测git能否访问
    git ls-remote git@github.com:aarioai/golib.git

    # 删除某个配置
    # git config --global --unset user.name

    # 临时缓存
    #git config --global credential.helper cache
    # 或使用长期存储（如 keychain）
    #git config --global credential.helper store

    ## goland 配置
    # Setting -> Go -> Go Modules -> Enable Go modules integration -> Environment -> GOPRIVATE=github.com/aarioai/golib
    # 通过power shell 执行 git config --global url."https://aarioai:<token>@github.com/".insteadOf "https://github.com/"

}

main(){
    handlePrivateRepo
}

main