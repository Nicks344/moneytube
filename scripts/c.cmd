set rootpath=E:\Projects\Go\src\moneytube
set clientpath=%rootpath%\client
set updatepath=%rootpath%\update
set electronpath=%clientpath%\ui\electron

REM Build electron package
cd %electronpath%
call npm run build

REM Go to electron build and copy default config
cd %electronpath%\MoneyTube-win32-x64
copy %clientpath%\defaultConfig.json config.json

REM Copy or move additional files inside electron package
cd resources
move %clientpath%\client.exe app\client.exe
copy %clientpath%\enigma_ide64.dll app\enigma_ide64.dll
move %updatepath%\update.exe app\update.exe
Xcopy /E /I %clientpath%\bin app\bin

REM Obfuscate electron code
call javascript-obfuscator .\app\app --output .\app
call javascript-obfuscator .\app\index.js --output .\app\index.js

REM Pack app to asar
call asar p app app.asar
rmdir app /s /q