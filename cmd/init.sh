#!/bin/bash
set -euo pipefail

handlePrivateRepo(){
    go env -w GOPRIVATE=github.com/aarioai/golib
    # 通过 https://github.com/settings/tokens 申请 token
    git config --global url."https://aarioai:<token>@github.com".insteadOf "https://github.com"

    git config --global --list


}

main(){
    handlePrivateRepo
}

main