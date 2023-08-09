Dim Arg, script
Set Arg = WScript.Arguments
script = Arg(0)
Set ps = CreateObject("Photoshop.Application")
dim fso, fullPath
set fso = CreateObject("Scripting.FileSystemObject")
fullPath = fso.GetAbsolutePathName(script)
ps.DoJavaScriptFile fullPath