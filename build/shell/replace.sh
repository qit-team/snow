#/bin/bash
# 将项目的包命名空间统一替换成目标目录空间
target=$1
default="github.com/qit-team/snow"
if [ "$target" == "" ]; then
    target="$default"
fi

#回到根目录
rootPath=$(cd `dirname $0`/../../; pwd)

#替换
find . -type f -name "*.*" ! -path "./vendor/*" ! -path "./logs/*" ! -path "./docs/*"  ! -path "./.git/*" ! -path "./build/shell/replace.sh" | xargs sed -i "s|${default}|${target}|g"
