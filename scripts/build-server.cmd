set rootpath=%1
set serverpath=%rootpath%\server
set uipath=%serverpath%\ui
set buildpath=%serverpath%\build
set datapath=%serverpath%\data

cd %serverpath%/src
call gox -osarch="linux/amd64" -output "../moneytube"

cd %uipath%
call npm run build
rmdir /s /q %datapath%\ui
move build %datapath%\ui

rmdir /s /q %buildpath%
mkdir %buildpath%

move %serverpath%\moneytube %buildpath%\moneytube
Xcopy /E /I %datapath% %buildpath%\data
rmdir /s /q %buildpath%\data\ffmpeg
del /q %buildpath%\data\keygen\keygen.exe
Xcopy /E /I %serverpath%\serverConfig.json %buildpath%\config.json