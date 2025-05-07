$ErrorActionPreference = 'Stop' # stop on all errors
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  file64        = "$toolsDir\bunyan-logviewer-0.3.1-windows-amd64.7z"
  softwareName  = 'lv*'
  checksum64    = '7f6ad4a411b75087157a2d05e23d816afccd48a5b726683077c7244c652d3ee1'
  checksumType64= 'sha256'
}

Get-ChocolateyUnzip @packageArgs
Remove-Item -Path $packageArgs.file64 -Force

Write-Output "To load tab completion in your current PowerShell session, please run:"
Write-Output "  lv completion powershell | Out-String | Invoke-Expression"
Write-Output " "
Write-Output "To load completions for every new session, add the output of the above command to your powershell profile."
