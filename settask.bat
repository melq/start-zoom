@echo off
set NAME=%~1
set ID=%~2
set PASS=%~3
set TIME=%~4
set DATE=%~5

schtasks -create -sc once -tn startzoom_%NAME% -tr "D:/myzoom.bat %ID% %PASS%" -st %TIME% -sd %DATE%