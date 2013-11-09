:: Current directory without trailing slash
set curDir=%CD%
set migrateDir=%curDir%\migrations\windows
mkdir %migrateDir%

:: cut off fractional seconds
set t=%time:~0,8%
:: replace colons with underscores
set t=%t::=_%
set d=%date%
:: replace dashes with underscores
set d=%d:-=_%
set FileName=migrate_%d%__%t%.log
cd ..
go run main.go orm syncdb -v=true > "%migrateDir%\%FileName%"
cd %curDir%