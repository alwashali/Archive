# Requires -Version 3.0
# 
# This Script is used to create a event at MISP using LOgrhythm ALaram Metadata
#
# The following steps are performed:
#
# 
#==========================================#
# LogRhythm SmartResponse Plugin           #
# Hive Create Alert - SmartResponse        #
# Author Ali Alwashali                     #
# v1.0  -- OCT 2022                        #
#==========================================#

[CmdletBinding()]
param(
[Parameter(Mandatory=$true)]
[ValidateNotNullOrEmpty()]
[string]$AlertSeverity,
[Parameter(Mandatory=$true)]
[ValidateNotNullOrEmpty()]
[string]$TrafficLightProtocol,
[string]$AlarmID,
[string]$AlarmDate,
[string]$AlarmRuleName,
[string]$OriginHost,
[string]$ImpactedHost,
[string]$KnownApplication,
[string]$OriginHostIP,
[string]$ImpactedHostIP,
[string]$Url,
[string]$Hash,
[string]$Domain,
[string]$SourcePort,
[string]$DestinationPort
)


# Trap for an exception during the Script
trap [Exception]
{
    if ($PSItem.ToString() -eq "ExecutionFailure")
	{
		exit 1
	}
	else
	{
		write-error $("Trapped: $_")
		write-host "Aborting Operation."
		exit
	}
}


# Function to Disable SSL Certificate Error and Enable Tls12

function Disable-SSLError
{
	# Disabling SSL certificate error
    add-type @"
        using System.Net;
        using System.Security.Cryptography.X509Certificates;
        public class TrustAllCertsPolicy : ICertificatePolicy {
            public bool CheckValidationResult(
                ServicePoint srvPoint, X509Certificate certificate,
                WebRequest request, int certificateProblem) {
                return true;
            }
        }
"@
    [System.Net.ServicePointManager]::CertificatePolicy = New-Object TrustAllCertsPolicy

    # Forcing to use TLS1.2
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
}


# Function to fetch saved parameter data

function Get-ConfigFileData
{
	try{
		if (!(Test-Path -Path $global:ConfigurationFilePath)){
			write-host "Config File Not Found. Please run 'Create The Hive Configuration File' action."
			write-error "Error: Config File Not Found. Please run 'Create The Hive Configuration File' action."
			throw "ExecutionFailure"
		}
		else{
			$ConfigFileContent = Import-Clixml -Path $global:ConfigurationFilePath
			$EncryptedHiveServerIP = $ConfigFileContent.hiveIP
			$EncryptedHiveApiKey = $ConfigFileContent.hiveApiKey
			$EncryptedHivePort = $ConfigFileContent.port
			$EncryptedCaseIP = $ConfigFileContent.CaseIP
			$global:HiveServerIP = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR((($EncryptedHiveServerIP))))
			$global:HiveApiKey = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR((($EncryptedHiveApiKey))))
			$global:HivePort = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR((($EncryptedHivePort))))
			$global:CaseIP = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR((($EncryptedCaseIP))))
		}
	}
	catch{
		$message = $_.Exception.message
		if($message -eq "ExecutionFailure"){
			throw "ExecutionFailure"
		}
		else{
			write-host "Error: User does not have access to Config File."
			write-error $message
			throw "ExecutionFailure"
		}
	}
}


# Function to get Alert Level

function Get-AlertSeverity{
	if (!(($AlertSeverity -eq "low") -or ($AlertSeverity -eq "medium") -or ($AlertSeverity -eq "high"))){
		write-host "Error: Not a valid value for Alert Level."
		write-error "Invalid Alert Level. Please Specify from 'High, Medium, Low'."
		throw "ExecutionFailure"
	}
	else{
		if ($AlertSeverity -eq "low"){
			$global:AlertSeverityid = 1
		}
		elseif ($AlertSeverity -eq "medium"){
			$global:AlertSeverityid = 2
		}
		elseif ($AlertSeverity -eq "high"){
			$global:AlertSeverityid = 3
		}
	}
}


# Function to get tlp id

function Get-Tlp{
	if (!(($TrafficLightProtocol -eq "white") -or ($TrafficLightProtocol -eq "green") -or ($TrafficLightProtocol -eq "amber") -or ($TrafficLightProtocol -eq "red"))){
		write-host "Error: Not a valid value for Traffic Light Protocol."
		write-error "Invalid Traffic Light Protocol. Please Specify from 'White, Green, Amber and Red'."
		throw "ExecutionFailure"
	}
	else {
		if ($TrafficLightProtocol -eq "white"){
			$global:TrafficLightProtocolid = 0
		}
		elseif ($TrafficLightProtocol -eq "green"){
			$global:TrafficLightProtocolid = 1
		}
		elseif ($TrafficLightProtocol -eq "amber"){
			$global:TrafficLightProtocolid = 2
		}
		elseif ($TrafficLightProtocol -eq "red"){
			$global:TrafficLightProtocolid = 3
		}
	}
}


# Function to create Metadata fields mapping with Hive attribute type

function Type-Mapping{
	$global:MappingObject = @{
		"Origin Host" = "Origin Host"
		"Impacted Host" = "Impacted Host"
		"Known Application" = "Known Application"
		"Origin Host IP" = "Origin Host IP"
		"Impacted Host IP" = "Impacted Host IP"
		"Url" = "url"
		"Hash" = "hash"
		"Domain" = "domain"
		"Source Port" = "other"
		"Destination Port" = "Destination Port"
	}
}


# Function to Create Metadata object

function Create-MetadataObject{
	if (!(($OriginHost -eq "") -or ($OriginHost -eq $null))){
		$OriginHost = $OriginHost.TrimEnd(" *")
	}
	if (!(($ImpactedHost -eq "") -or ($ImpactedHost -eq $null))){
		$ImpactedHost = $ImpactedHost.TrimEnd(" *")
	}
	$global:MetadataObject = [ordered]@{
		"Destination Port" = $DestinationPort
		"Source Port" = $SourcePort
		"Domain" = $Domain
		"Hash" = $Hash
		"Url" = $Url
		"Origin Host IP" = $OriginHostIP
		"Impacted Host IP" = $ImpactedHostIP
		"Known Application" = $KnownApplication
		"Impacted Host" = $ImpactedHost
		"Origin Host" = $OriginHost
		}
		#$global:MetadataObject
}





# Function to Create Json Body 

function Create-JsonBody{
	$Artifact = @()
	$Metadatkeys = $global:MetadataObject.keys
	foreach ($key in $Metadatkeys){
		if (($global:MetadataObject[$key] -eq "") -or ($global:MetadataObject[$key] -eq $null)){
			continue
		}
		$ArtifactBody = [ordered]@{
				dataType = $global:MappingObject[$key];
				data = $global:MetadataObject[$key];
				message = $key
				}
		$Artifact += $ArtifactBody
	}
	$AlertBody = [ordered]@{
			title = $AlarmID + ": " + $AlarmRuleName ;
			description = "Alarm Id - " + $AlarmID + "`n" + " Alarm Date - " + $AlarmDate + " Observables Fetched by SRP" + "`n" + " LR- Alarm Link - " + "http://" + $global:CaseIP + ":8443/alarms/" + $AlarmID;
			type = "external";
			source = "LR-SIEM";
			sourceRef = "LR-SIEM: ID- " + $AlarmID;
			severity = $global:AlertSeverityid;
			tlp = $global:TrafficLightProtocolid;
			artifacts = $Artifact
			caseTemplate = "external-alert"
			}
	$global:JsonBody = $AlertBody | ConvertTo-Json -Depth 5
}


# Function to call Hive API to create an alert

function Create-HiveAlert{
	try{
		$Header = @{
			"Accept" = "application/json"
			"Authorization" = "Bearer " + $global:HiveApiKey
		}
		$CreateAlertUrl = $HiveBaseUrl + "api/alert"
		$CreateAlertResponse = Invoke-RestMethod -uri $CreateAlertUrl -Headers $Header -Method Post -Body $global:JsonBody -ContentType "application/json"
		write-host "Hive Alert Created with Title : Observable from Alarm ID -  $AlarmID"
	}
	catch {
		$ErrorMessage = $_.Exception.Message
		if ($ErrorMessage -eq "The remote server returned an error: (400) Bad Request."){
			write-host "Error: Invalid Json body for Alert."
			write-error "Error: (400) Bad Request."
			throw "ExecutionFailure"
		}
		else {
			$ErrorMessage
			throw "ExecutionFailure"
		}
		
	}
}


$global:ConfigurationFilePath = "C:\Program Files\LogRhythm\SmartResponse Plugins\TheHiveConfigFile.xml"
$AlertSeverity = $AlertSeverity.Trim()
$TrafficLightProtocol = $TrafficLightProtocol.Trim()


Disable-SSLError
Get-ConfigFileData

$HiveBaseUrl = "http://" + $global:HiveServerIP + ":" + $global:HivePort + "/"

Get-AlertSeverity
Get-Tlp
Type-Mapping
Create-MetadataObject
Create-JsonBody
Create-HiveAlert

