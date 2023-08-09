set rootpath=%1

set corepath=%rootpath%\client\core
set updatepath=%rootpath%\update
set clipath=%rootpath%\cli

set uipath=%rootpath%\client\ui
set reactpath=%uipath%\react
set electronpath=%uipath%\electron
set fullbuildpath=%electronpath%\MoneyTube-win32-x64

rmdir /s /q %fullbuildpath%
if %errorlevel% neq 0 exit /b %errorlevel%

rmdir /s /q %electronpath%\build
if %errorlevel% neq 0 exit /b %errorlevel%

REM Build client UI
cd %reactpath%
call npm run build
if %errorlevel% neq 0 exit /b %errorlevel%

Xcopy /E /I build %electronpath%\build
if %errorlevel% neq 0 exit /b %errorlevel%

REM Build update exe
cd %updatepath%
call gox -osarch="windows/amd64" -output "update"
if %errorlevel% neq 0 exit /b %errorlevel%

REM Build cli exe
cd %clipath%
call gox -osarch="windows/amd64" -output "cli"
if %errorlevel% neq 0 exit /b %errorlevel%

REM Build and protect client exe
cd %corepath%/src

REM Build client
call gox -osarch="windows/amd64" -output "../core" -cgo -tags="release"
if %errorlevel% neq 0 exit /b %errorlevel%

REM Protect client
"c:\Program Files (x86)\The Enigma Protector\enigma64.exe" %rootpath%\server\backend\data\keygen\protect.enigma
if %errorlevel% neq 0 exit /b %errorlevel%

REM Build electron package
cd %electronpath%
call npm run build
if %errorlevel% neq 0 exit /b %errorlevel%

REM Go to electron build and copy default config
cd %electronpath%\MoneyTube-win32-x64
copy %corepath%\defaultConfig.json config.json
if %errorlevel% neq 0 exit /b %errorlevel%

cd resources

REM Copy or move additional files inside electron package
move %corepath%\core.exe app\core.exe
if %errorlevel% neq 0 exit /b %errorlevel%
copy %corepath%\enigma_ide64.dll app\enigma_ide64.dll
if %errorlevel% neq 0 exit /b %errorlevel%
move %updatepath%\update.exe app\update.exe
if %errorlevel% neq 0 exit /b %errorlevel%
move %clipath%\cli.exe app\cli.exe
if %errorlevel% neq 0 exit /b %errorlevel%
Xcopy /E /I %corepath%\bin app\bin
if %errorlevel% neq 0 exit /b %errorlevel%
copy %corepath%\user-agents.txt app\user-agents.txt
if %errorlevel% neq 0 exit /b %errorlevel%

REM Obfuscate electron code
REM call javascript-obfuscator .\app\app --output .\app\app
REM if %errorlevel% neq 0 exit /b %errorlevel%
REM call javascript-obfuscator .\app\index.js --output .\app\index.js
REM if %errorlevel% neq 0 exit /b %errorlevel%

REM Pack app to asar
call asar p app app.asar
if %errorlevel% neq 0 exit /b %errorlevel%

del /f /s /q app 1>nul
rmdir app /s /q

