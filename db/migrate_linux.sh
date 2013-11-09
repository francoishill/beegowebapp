#!/bin/sh
# Current directory without trailing slash
curDir=$(pwd)
migrateDir="${curDir}/migrations/linux"
mkdir -p ${migrateDir}
NowString=$(date +"%Y_%m_%d__%H_%M_%S")
FileName=migrate_${NowString}.log
cd ..
rm -f temp_migrate
go build -o temp_migrate main.go > "${migrateDir}/${FileName}"
./temp_migrate orm syncdb -v=true >> "${migrateDir}/${FileName}"
grep "." ${migrateDir}/${FileName}
echo ""
echo Contents of migration file printed, file path: ${migrateDir}/${FileName}
cd ${curDir}