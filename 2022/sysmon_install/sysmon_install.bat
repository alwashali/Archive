if not exist "C:\windows\sysmonconfig.xml" (
echo copying configuration ... 
copy /z /y "\\ad01.hdlabs.local\sysmon\sysmonconfig.xml" "C:\windows\"
)

sc query Sysmon64 | Find "RUNNING"
If %ERRORLEVEL% == 0 (
goto eexit
)
echo service is not running
goto startsysmon

:startsysmon
net start Sysmon64

If %ERRORLEVEL% == 0 (
goto eexit
)

echo service not found, going to install sysmon 
"\\ad01.hdlabs.local\sysmon\Sysmon64.exe" /accepteula -i c:\windows\sysmonconfig.xml

:eexit
