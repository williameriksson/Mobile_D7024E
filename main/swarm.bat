@ECHO OFF
SETLOCAL ENABLEDELAYEDEXPANSION

SET /A OWNADDRPORT=8101
SET OWNADDRBASE=127.0.0.1:
SET BOOTSTRAPIP=127.0.0.1:
SET /A BOOTSTRAPPORT=8100
FOR /L %%a in (1,1,50) DO (
Start main_http.exe %OWNADDRBASE%!OWNADDRPORT! %BOOTSTRAPIP%!BOOTSTRAPPORT!
SET /A TEMPPORT=!OWNADDRPORT!-8100
SET /A BOOTSTRAPPORT=!RANDOM!*!TEMPPORT!/32768+8100
SET /A OWNADDRPORT+=1
timeout 1
)