package sso

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"testing"
	"time"

	dsig "github.com/russellhaering/goxmldsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMetadata = `
<?xml version="1.0" encoding="UTF-8"?>
  <md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://kolide-dev-ed.my.salesforce.com" validUntil="2027-04-29T19:22:40.750Z" xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
      <md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
         <md:KeyDescriptor use="signing">
            <ds:KeyInfo>
               <ds:X509Data>
                  <ds:X509Certificate>MIIErDCCA5SgAwIBAgIOAVuhH3WkAAAAAB5NpvIwDQYJKoZIhvcNAQELBQAwgZAxKDAmBgNVBAMMH1NlbGZTaWduZWRDZXJ0XzI0QXByMjAxN18xODAwNDQxGDAWBgNVBAsMDzAwRDZBMDAwMDAwMTd0ODEXMBUGA1UECgwOU2FsZXNmb3JjZS5jb20xFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xCzAJBgNVBAgMAkNBMQwwCgYDVQQGEwNVU0EwHhcNMTcwNDI0MTgwMDQ1WhcNMTgwNDI0MTIwMDAwWjCBkDEoMCYGA1UEAwwfU2VsZlNpZ25lZENlcnRfMjRBcHIyMDE3XzE4MDA0NDEYMBYGA1UECwwPMDBENkEwMDAwMDAxN3Q4MRcwFQYDVQQKDA5TYWxlc2ZvcmNlLmNvbTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzELMAkGA1UECAwCQ0ExDDAKBgNVBAYTA1VTQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAIOR7h8BF2eFOlQHhV/1S7uOBN22Jv7PDCXMz2fU0uLc+mrv9xDGj6ElfW+9dSdXaCbQzD3+Xq4reS4pYRafJZ/27OtygXl3rpoPjSlhRiW+oYVuDcCURJpu0KuZ4I0fm5q1BDYqxcBxNPSe85OHE3+ucmKqvPozhQgYLPCregMIomC3yyANZnLCoGfCv9TpQl6/+I182tST4WPNhVPxKxijoPU4Rh6xY34Ez8+Jr8KdmzmYSNe4ukkIASplpvG7rKka824Hf8zI1BWnjWLDxb5IAxgUBbdr4x8d8C3kPfTf+3/6yC5wSOm9NSs0BA4OJNowtXZFryMzFfXzDzjl69kCAwEAAaOCAQAwgf0wHQYDVR0OBBYEFO+DkoP6qkysi9ZC74yTPuJVVg2yMA8GA1UdEwEB/wQFMAMBAf8wgcoGA1UdIwSBwjCBv4AU74OSg/qqTKyL1kLvjJM+4lVWDbKhgZakgZMwgZAxKDAmBgNVBAMMH1NlbGZTaWduZWRDZXJ0XzI0QXByMjAxN18xODAwNDQxGDAWBgNVBAsMDzAwRDZBMDAwMDAwMTd0ODEXMBUGA1UECgwOU2FsZXNmb3JjZS5jb20xFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xCzAJBgNVBAgMAkNBMQwwCgYDVQQGEwNVU0GCDgFboR91pAAAAAAeTabyMA0GCSqGSIb3DQEBCwUAA4IBAQAVhYBv5GJvhltks2j7Zc9wdFHW7yB4/hPFo05y0yiOf71tLjOlBucSyxtmXLPjrECJvIJwKhsAIgYXnVp7ditxfauCcxczJgfeL1/dxH/Ge8ePkmH6SdsO71cJL8dXEzOsoF+PAVQzUhqh8zxIipntL0wwNGTD0zIVQeTSozm0KF0SsSHIfbNy279uReGonC61i4Ouk5AMKA7Re9fVeUs6tqM2at22h9Zaj/r/OhXoDcZhzkd8Wq0ER/UKLZA1CyJHgwOC7REEZOuKrqgfWcYt4dGo5q6gqGHHPMv0N7s/MxqCvJCwGA8eJGvOO56I321vhWHQ6ZSJDWUqQFM/Ze7A</ds:X509Certificate>
               </ds:X509Data>
            </ds:KeyInfo>
         </md:KeyDescriptor>
         <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
         <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://kolide-dev-ed.my.salesforce.com/idp/endpoint/HttpPost"/>
         <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://kolide-dev-ed.my.salesforce.com/idp/endpoint/HttpRedirect"/>
      </md:IDPSSODescriptor>
   </md:EntityDescriptor>
`

func TestNewValidator(t *testing.T) {
	v, err := NewValidator(testMetadata)
	assert.Nil(t, err)
	assert.NotNil(t, v)
}

var testResponse = `PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c2FtbHA6UmVzcG9uc2UgeG1sbnM6c2FtbHA9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCIgRGVzdGluYXRpb249Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvYXBpL3YxL2tvbGlkZS9zc28vY2FsbGJhY2siIElEPSJfYzgyMWM1MmUzZDJkN2NhMjFiMmJlNmJlOTQ4Y2RiYjAxNDkzNTkwMTA2NjM1IiBJblJlc3BvbnNlVG89IjM5MTY5NzllLTNhZTItNGU4My04N2I0LWZlMmViODg0Yjg5MSIgSXNzdWVJbnN0YW50PSIyMDE3LTA0LTMwVDIyOjA4OjI2LjYzNVoiIFZlcnNpb249IjIuMCI+PHNhbWw6SXNzdWVyIHhtbG5zOnNhbWw9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDphc3NlcnRpb24iIEZvcm1hdD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm5hbWVpZC1mb3JtYXQ6ZW50aXR5Ij5odHRwczovL2tvbGlkZS1kZXYtZWQubXkuc2FsZXNmb3JjZS5jb208L3NhbWw6SXNzdWVyPjxkczpTaWduYXR1cmUgeG1sbnM6ZHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPgo8ZHM6U2lnbmVkSW5mbz4KPGRzOkNhbm9uaWNhbGl6YXRpb25NZXRob2QgQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiLz4KPGRzOlNpZ25hdHVyZU1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNyc2Etc2hhMSIvPgo8ZHM6UmVmZXJlbmNlIFVSST0iI19jODIxYzUyZTNkMmQ3Y2EyMWIyYmU2YmU5NDhjZGJiMDE0OTM1OTAxMDY2MzUiPgo8ZHM6VHJhbnNmb3Jtcz4KPGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+CjxkczpUcmFuc2Zvcm0gQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiPjxlYzpJbmNsdXNpdmVOYW1lc3BhY2VzIHhtbG5zOmVjPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiIFByZWZpeExpc3Q9ImRzIHNhbWwgc2FtbHAgeHMgeHNpIi8+PC9kczpUcmFuc2Zvcm0+CjwvZHM6VHJhbnNmb3Jtcz4KPGRzOkRpZ2VzdE1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNzaGExIi8+CjxkczpEaWdlc3RWYWx1ZT5PWnMxdHRQamtKYnF6Mk8yb3UraDJZK3FFYms9PC9kczpEaWdlc3RWYWx1ZT4KPC9kczpSZWZlcmVuY2U+CjwvZHM6U2lnbmVkSW5mbz4KPGRzOlNpZ25hdHVyZVZhbHVlPgpWMmlSWlRIdUpCajBnekYzYzFHVzBGY0JiRlR0QXlWUml4VWdtR3ZrR0xhRzBBeElBZkY4ejdxbUlMTkV4cDlEUjJJU1F5a2lCU3p6CkpBQllRREU2T0V6ak5XTnYyS2NsUTRjQ0VOOUdIbHh6bEo4dkFwQzRsdkV3aGlQL04zd0VKb3RGTlN2MVZRd0YvdWFmZ1Z6b1NIeVIKb3RYaEJ0akFjcktBV25kRWo5L3QvdWV0SkI4dHZ5OXdybzhtS3RIZVNiTmJoZ0dwVEgyVHpuUnFxRnhwS1lRT29adFExaEpOVSsvdApMWGVpbWxjc0QrSHFQejlIT0crN0JuQllZZTkyN1dZRWRqREFLODVmMDY3ekN5T1RmK3pnVFNJNCs5WUVsbUo1OVZ6RFRwR3kyZENVCkwvZHJrSjhDSUJKZzZ4ekc1aVAwbFVGay81TmFrclRxZG0wNkF3PT0KPC9kczpTaWduYXR1cmVWYWx1ZT4KPGRzOktleUluZm8+PGRzOlg1MDlEYXRhPjxkczpYNTA5Q2VydGlmaWNhdGU+TUlJRXJEQ0NBNVNnQXdJQkFnSU9BVnVoSDNXa0FBQUFBQjVOcHZJd0RRWUpLb1pJaHZjTkFRRUxCUUF3Z1pBeEtEQW1CZ05WQkFNTQpIMU5sYkdaVGFXZHVaV1JEWlhKMFh6STBRWEJ5TWpBeE4xOHhPREF3TkRReEdEQVdCZ05WQkFzTUR6QXdSRFpCTURBd01EQXdNVGQwCk9ERVhNQlVHQTFVRUNnd09VMkZzWlhObWIzSmpaUzVqYjIweEZqQVVCZ05WQkFjTURWTmhiaUJHY21GdVkybHpZMjh4Q3pBSkJnTlYKQkFnTUFrTkJNUXd3Q2dZRFZRUUdFd05WVTBFd0hoY05NVGN3TkRJME1UZ3dNRFExV2hjTk1UZ3dOREkwTVRJd01EQXdXakNCa0RFbwpNQ1lHQTFVRUF3d2ZVMlZzWmxOcFoyNWxaRU5sY25SZk1qUkJjSEl5TURFM1h6RTRNREEwTkRFWU1CWUdBMVVFQ3d3UE1EQkVOa0V3Ck1EQXdNREF4TjNRNE1SY3dGUVlEVlFRS0RBNVRZV3hsYzJadmNtTmxMbU52YlRFV01CUUdBMVVFQnd3TlUyRnVJRVp5WVc1amFYTmoKYnpFTE1Ba0dBMVVFQ0F3Q1EwRXhEREFLQmdOVkJBWVRBMVZUUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQwpnZ0VCQUlPUjdoOEJGMmVGT2xRSGhWLzFTN3VPQk4yMkp2N1BEQ1hNejJmVTB1TGMrbXJ2OXhER2o2RWxmVys5ZFNkWGFDYlF6RDMrClhxNHJlUzRwWVJhZkpaLzI3T3R5Z1hsM3Jwb1BqU2xoUmlXK29ZVnVEY0NVUkpwdTBLdVo0STBmbTVxMUJEWXF4Y0J4TlBTZTg1T0gKRTMrdWNtS3F2UG96aFFnWUxQQ3JlZ01Jb21DM3l5QU5abkxDb0dmQ3Y5VHBRbDYvK0kxODJ0U1Q0V1BOaFZQeEt4aWpvUFU0Umg2eApZMzRFejgrSnI4S2Rtem1ZU05lNHVra0lBU3BscHZHN3JLa2E4MjRIZjh6STFCV25qV0xEeGI1SUF4Z1VCYmRyNHg4ZDhDM2tQZlRmCiszLzZ5QzV3U09tOU5TczBCQTRPSk5vd3RYWkZyeU16RmZYekR6amw2OWtDQXdFQUFhT0NBUUF3Z2Ywd0hRWURWUjBPQkJZRUZPK0QKa29QNnFreXNpOVpDNzR5VFB1SlZWZzJ5TUE4R0ExVWRFd0VCL3dRRk1BTUJBZjh3Z2NvR0ExVWRJd1NCd2pDQnY0QVU3NE9TZy9xcQpUS3lMMWtMdmpKTSs0bFZXRGJLaGdaYWtnWk13Z1pBeEtEQW1CZ05WQkFNTUgxTmxiR1pUYVdkdVpXUkRaWEowWHpJMFFYQnlNakF4Ck4xOHhPREF3TkRReEdEQVdCZ05WQkFzTUR6QXdSRFpCTURBd01EQXdNVGQwT0RFWE1CVUdBMVVFQ2d3T1UyRnNaWE5tYjNKalpTNWoKYjIweEZqQVVCZ05WQkFjTURWTmhiaUJHY21GdVkybHpZMjh4Q3pBSkJnTlZCQWdNQWtOQk1Rd3dDZ1lEVlFRR0V3TlZVMEdDRGdGYgpvUjkxcEFBQUFBQWVUYWJ5TUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBVmhZQnY1R0p2aGx0a3MyajdaYzl3ZEZIVzd5QjQvaFBGCm8wNXkweWlPZjcxdExqT2xCdWNTeXh0bVhMUGpyRUNKdklKd0toc0FJZ1lYblZwN2RpdHhmYXVDY3hjekpnZmVMMS9keEgvR2U4ZVAKa21INlNkc083MWNKTDhkWEV6T3NvRitQQVZRelVocWg4enhJaXBudEwwd3dOR1REMHpJVlFlVFNvem0wS0YwU3NTSElmYk55Mjc5dQpSZUdvbkM2MWk0T3VrNUFNS0E3UmU5ZlZlVXM2dHFNMmF0MjJoOVphai9yL09oWG9EY1poemtkOFdxMEVSL1VLTFpBMUN5Skhnd09DCjdSRUVaT3VLcnFnZldjWXQ0ZEdvNXE2Z3FHSEhQTXYwTjdzL014cUN2SkN3R0E4ZUpHdk9PNTZJMzIxdmhXSFE2WlNKRFdVcVFGTS8KWmU3QTwvZHM6WDUwOUNlcnRpZmljYXRlPjwvZHM6WDUwOURhdGE+PC9kczpLZXlJbmZvPjwvZHM6U2lnbmF0dXJlPjxzYW1scDpTdGF0dXM+PHNhbWxwOlN0YXR1c0NvZGUgVmFsdWU9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpzdGF0dXM6U3VjY2VzcyIvPjwvc2FtbHA6U3RhdHVzPjxzYW1sOkFzc2VydGlvbiB4bWxuczpzYW1sPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIiBJRD0iX2YyZTljM2Y3MDAwYzQwZjAyNDEwMjhlYzkwNDMxNDg2MTQ5MzU5MDEwNjYzNiIgSXNzdWVJbnN0YW50PSIyMDE3LTA0LTMwVDIyOjA4OjI2LjYzNloiIFZlcnNpb249IjIuMCI+PHNhbWw6SXNzdWVyIEZvcm1hdD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOm5hbWVpZC1mb3JtYXQ6ZW50aXR5Ij5odHRwczovL2tvbGlkZS1kZXYtZWQubXkuc2FsZXNmb3JjZS5jb208L3NhbWw6SXNzdWVyPjxkczpTaWduYXR1cmUgeG1sbnM6ZHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPgo8ZHM6U2lnbmVkSW5mbz4KPGRzOkNhbm9uaWNhbGl6YXRpb25NZXRob2QgQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiLz4KPGRzOlNpZ25hdHVyZU1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNyc2Etc2hhMSIvPgo8ZHM6UmVmZXJlbmNlIFVSST0iI19mMmU5YzNmNzAwMGM0MGYwMjQxMDI4ZWM5MDQzMTQ4NjE0OTM1OTAxMDY2MzYiPgo8ZHM6VHJhbnNmb3Jtcz4KPGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+CjxkczpUcmFuc2Zvcm0gQWxnb3JpdGhtPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiPjxlYzpJbmNsdXNpdmVOYW1lc3BhY2VzIHhtbG5zOmVjPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxLzEwL3htbC1leGMtYzE0biMiIFByZWZpeExpc3Q9ImRzIHNhbWwgeHMgeHNpIi8+PC9kczpUcmFuc2Zvcm0+CjwvZHM6VHJhbnNmb3Jtcz4KPGRzOkRpZ2VzdE1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNzaGExIi8+CjxkczpEaWdlc3RWYWx1ZT5vbjdQZ1B4TEpiOG1sbS9pK01EUU1yMVZZeVk9PC9kczpEaWdlc3RWYWx1ZT4KPC9kczpSZWZlcmVuY2U+CjwvZHM6U2lnbmVkSW5mbz4KPGRzOlNpZ25hdHVyZVZhbHVlPgpSVEF0MkJHc2ZRRXlJWXU1VUNRZmhYcWwxWE9mWHgraDgrNnNhM1JNMENVNVVYUVBFVjd5RWxFdXZWUWhHZVpNamN0eFZMR1BjZ01WCktYVGowM1Y0RHdmai95aHQwSHVHdktoQXZQUklScnBVYlJOUWEzZWZDMnZDQ05HUjdQR2czSVllZXI4Rk41ZEJiZk5CQW1oQUltdk8KWmdXd2QwSmlpV1FleWNreEdBZmZjZEZ4OC9Ra0dxdVp1S0dGQkMyZ2pJV0ZncXhMcWs1RTJYeDlybW5GamY2TGVJajlIc0duWlUxeApqTE1seTJLRjJyUGRsVmR1UGlTWko3T1JjZVVycmo2YUI4RlNQVGdoa3pnWVV1cU1rbWFjYzlXYVZtc01lTGI3TC9XSTlBOHhIcUhzClZsTkp3ZzZiMWpXMnd5VW1mTzVoVWtxMW9rdDlTZjFST0lNYmFnPT0KPC9kczpTaWduYXR1cmVWYWx1ZT4KPGRzOktleUluZm8+PGRzOlg1MDlEYXRhPjxkczpYNTA5Q2VydGlmaWNhdGU+TUlJRXJEQ0NBNVNnQXdJQkFnSU9BVnVoSDNXa0FBQUFBQjVOcHZJd0RRWUpLb1pJaHZjTkFRRUxCUUF3Z1pBeEtEQW1CZ05WQkFNTQpIMU5sYkdaVGFXZHVaV1JEWlhKMFh6STBRWEJ5TWpBeE4xOHhPREF3TkRReEdEQVdCZ05WQkFzTUR6QXdSRFpCTURBd01EQXdNVGQwCk9ERVhNQlVHQTFVRUNnd09VMkZzWlhObWIzSmpaUzVqYjIweEZqQVVCZ05WQkFjTURWTmhiaUJHY21GdVkybHpZMjh4Q3pBSkJnTlYKQkFnTUFrTkJNUXd3Q2dZRFZRUUdFd05WVTBFd0hoY05NVGN3TkRJME1UZ3dNRFExV2hjTk1UZ3dOREkwTVRJd01EQXdXakNCa0RFbwpNQ1lHQTFVRUF3d2ZVMlZzWmxOcFoyNWxaRU5sY25SZk1qUkJjSEl5TURFM1h6RTRNREEwTkRFWU1CWUdBMVVFQ3d3UE1EQkVOa0V3Ck1EQXdNREF4TjNRNE1SY3dGUVlEVlFRS0RBNVRZV3hsYzJadmNtTmxMbU52YlRFV01CUUdBMVVFQnd3TlUyRnVJRVp5WVc1amFYTmoKYnpFTE1Ba0dBMVVFQ0F3Q1EwRXhEREFLQmdOVkJBWVRBMVZUUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQwpnZ0VCQUlPUjdoOEJGMmVGT2xRSGhWLzFTN3VPQk4yMkp2N1BEQ1hNejJmVTB1TGMrbXJ2OXhER2o2RWxmVys5ZFNkWGFDYlF6RDMrClhxNHJlUzRwWVJhZkpaLzI3T3R5Z1hsM3Jwb1BqU2xoUmlXK29ZVnVEY0NVUkpwdTBLdVo0STBmbTVxMUJEWXF4Y0J4TlBTZTg1T0gKRTMrdWNtS3F2UG96aFFnWUxQQ3JlZ01Jb21DM3l5QU5abkxDb0dmQ3Y5VHBRbDYvK0kxODJ0U1Q0V1BOaFZQeEt4aWpvUFU0Umg2eApZMzRFejgrSnI4S2Rtem1ZU05lNHVra0lBU3BscHZHN3JLa2E4MjRIZjh6STFCV25qV0xEeGI1SUF4Z1VCYmRyNHg4ZDhDM2tQZlRmCiszLzZ5QzV3U09tOU5TczBCQTRPSk5vd3RYWkZyeU16RmZYekR6amw2OWtDQXdFQUFhT0NBUUF3Z2Ywd0hRWURWUjBPQkJZRUZPK0QKa29QNnFreXNpOVpDNzR5VFB1SlZWZzJ5TUE4R0ExVWRFd0VCL3dRRk1BTUJBZjh3Z2NvR0ExVWRJd1NCd2pDQnY0QVU3NE9TZy9xcQpUS3lMMWtMdmpKTSs0bFZXRGJLaGdaYWtnWk13Z1pBeEtEQW1CZ05WQkFNTUgxTmxiR1pUYVdkdVpXUkRaWEowWHpJMFFYQnlNakF4Ck4xOHhPREF3TkRReEdEQVdCZ05WQkFzTUR6QXdSRFpCTURBd01EQXdNVGQwT0RFWE1CVUdBMVVFQ2d3T1UyRnNaWE5tYjNKalpTNWoKYjIweEZqQVVCZ05WQkFjTURWTmhiaUJHY21GdVkybHpZMjh4Q3pBSkJnTlZCQWdNQWtOQk1Rd3dDZ1lEVlFRR0V3TlZVMEdDRGdGYgpvUjkxcEFBQUFBQWVUYWJ5TUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCQVFBVmhZQnY1R0p2aGx0a3MyajdaYzl3ZEZIVzd5QjQvaFBGCm8wNXkweWlPZjcxdExqT2xCdWNTeXh0bVhMUGpyRUNKdklKd0toc0FJZ1lYblZwN2RpdHhmYXVDY3hjekpnZmVMMS9keEgvR2U4ZVAKa21INlNkc083MWNKTDhkWEV6T3NvRitQQVZRelVocWg4enhJaXBudEwwd3dOR1REMHpJVlFlVFNvem0wS0YwU3NTSElmYk55Mjc5dQpSZUdvbkM2MWk0T3VrNUFNS0E3UmU5ZlZlVXM2dHFNMmF0MjJoOVphai9yL09oWG9EY1poemtkOFdxMEVSL1VLTFpBMUN5Skhnd09DCjdSRUVaT3VLcnFnZldjWXQ0ZEdvNXE2Z3FHSEhQTXYwTjdzL014cUN2SkN3R0E4ZUpHdk9PNTZJMzIxdmhXSFE2WlNKRFdVcVFGTS8KWmU3QTwvZHM6WDUwOUNlcnRpZmljYXRlPjwvZHM6WDUwOURhdGE+PC9kczpLZXlJbmZvPjwvZHM6U2lnbmF0dXJlPjxzYW1sOlN1YmplY3Q+PHNhbWw6TmFtZUlEIEZvcm1hdD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6dW5zcGVjaWZpZWQiPmpvaG5Aa29saWRlLmNvPC9zYW1sOk5hbWVJRD48c2FtbDpTdWJqZWN0Q29uZmlybWF0aW9uIE1ldGhvZD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmNtOmJlYXJlciI+PHNhbWw6U3ViamVjdENvbmZpcm1hdGlvbkRhdGEgSW5SZXNwb25zZVRvPSIzOTE2OTc5ZS0zYWUyLTRlODMtODdiNC1mZTJlYjg4NGI4OTEiIE5vdE9uT3JBZnRlcj0iMjAxNy0wNC0zMFQyMjoxMzoyNi42NDNaIiBSZWNpcGllbnQ9Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvYXBpL3YxL2tvbGlkZS9zc28vY2FsbGJhY2siLz48L3NhbWw6U3ViamVjdENvbmZpcm1hdGlvbj48L3NhbWw6U3ViamVjdD48c2FtbDpDb25kaXRpb25zIE5vdEJlZm9yZT0iMjAxNy0wNC0zMFQyMjowNzo1Ni42NDNaIiBOb3RPbk9yQWZ0ZXI9IjIwMTctMDQtMzBUMjI6MTM6MjYuNjQzWiI+PHNhbWw6QXVkaWVuY2VSZXN0cmljdGlvbj48c2FtbDpBdWRpZW5jZT5rb2xpZGU8L3NhbWw6QXVkaWVuY2U+PC9zYW1sOkF1ZGllbmNlUmVzdHJpY3Rpb24+PC9zYW1sOkNvbmRpdGlvbnM+PHNhbWw6QXV0aG5TdGF0ZW1lbnQgQXV0aG5JbnN0YW50PSIyMDE3LTA0LTMwVDIyOjA4OjI2LjYzN1oiPjxzYW1sOkF1dGhuQ29udGV4dD48c2FtbDpBdXRobkNvbnRleHRDbGFzc1JlZj51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YWM6Y2xhc3Nlczp1bnNwZWNpZmllZDwvc2FtbDpBdXRobkNvbnRleHRDbGFzc1JlZj48L3NhbWw6QXV0aG5Db250ZXh0Pjwvc2FtbDpBdXRoblN0YXRlbWVudD48c2FtbDpBdHRyaWJ1dGVTdGF0ZW1lbnQ+PHNhbWw6QXR0cmlidXRlIE5hbWU9InVzZXJJZCIgTmFtZUZvcm1hdD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmF0dHJuYW1lLWZvcm1hdDp1bnNwZWNpZmllZCI+PHNhbWw6QXR0cmlidXRlVmFsdWUgeG1sbnM6eHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvWE1MU2NoZW1hIiB4bWxuczp4c2k9Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvWE1MU2NoZW1hLWluc3RhbmNlIiB4c2k6dHlwZT0ieHM6YW55VHlwZSI+MDA1NkEwMDAwMDBRNlJsPC9zYW1sOkF0dHJpYnV0ZVZhbHVlPjwvc2FtbDpBdHRyaWJ1dGU+PHNhbWw6QXR0cmlidXRlIE5hbWU9InVzZXJuYW1lIiBOYW1lRm9ybWF0PSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXR0cm5hbWUtZm9ybWF0OnVuc3BlY2lmaWVkIj48c2FtbDpBdHRyaWJ1dGVWYWx1ZSB4bWxuczp4cz0iaHR0cDovL3d3dy53My5vcmcvMjAwMS9YTUxTY2hlbWEiIHhtbG5zOnhzaT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS9YTUxTY2hlbWEtaW5zdGFuY2UiIHhzaTp0eXBlPSJ4czphbnlUeXBlIj5qb2huQGtvbGlkZS5jbzwvc2FtbDpBdHRyaWJ1dGVWYWx1ZT48L3NhbWw6QXR0cmlidXRlPjxzYW1sOkF0dHJpYnV0ZSBOYW1lPSJlbWFpbCIgTmFtZUZvcm1hdD0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmF0dHJuYW1lLWZvcm1hdDp1bnNwZWNpZmllZCI+PHNhbWw6QXR0cmlidXRlVmFsdWUgeG1sbnM6eHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvWE1MU2NoZW1hIiB4bWxuczp4c2k9Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvWE1MU2NoZW1hLWluc3RhbmNlIiB4c2k6dHlwZT0ieHM6YW55VHlwZSI+am9obkBrb2xpZGUuY288L3NhbWw6QXR0cmlidXRlVmFsdWU+PC9zYW1sOkF0dHJpYnV0ZT48c2FtbDpBdHRyaWJ1dGUgTmFtZT0iaXNfcG9ydGFsX3VzZXIiIE5hbWVGb3JtYXQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDphdHRybmFtZS1mb3JtYXQ6dW5zcGVjaWZpZWQiPjxzYW1sOkF0dHJpYnV0ZVZhbHVlIHhtbG5zOnhzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYSIgeG1sbnM6eHNpPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSIgeHNpOnR5cGU9InhzOmFueVR5cGUiPmZhbHNlPC9zYW1sOkF0dHJpYnV0ZVZhbHVlPjwvc2FtbDpBdHRyaWJ1dGU+PC9zYW1sOkF0dHJpYnV0ZVN0YXRlbWVudD48L3NhbWw6QXNzZXJ0aW9uPjwvc2FtbHA6UmVzcG9uc2U+`

func TestValidate(t *testing.T) {
	tm, err := time.Parse(time.UnixDate, "Sun Apr 30 22:10:00 UTC 2017")
	require.Nil(t, err)
	clock := dsig.NewFakeClockAt(tm)
	validator, err := NewValidator(testMetadata, Clock(clock))
	require.Nil(t, err)
	require.NotNil(t, validator)
	auth, err := DecodeAuthResponse(testResponse)
	signed, err := validator.ValidateSignature(auth)
	require.Nil(t, err)
	require.NotNil(t, signed)

	err = validator.ValidateResponse(auth)
	assert.Nil(t, err)
}

func tamperedResponse(original string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(original)
	if err != nil {
		return "", err
	}
	var resp Response
	rdr := bytes.NewBuffer(decoded)
	err = xml.NewDecoder(rdr).Decode(&resp)
	if err != nil {
		return "", err
	}
	// change name
	resp.Assertion.Subject.NameID.Value = "bob@kolide.co"
	var wrtr bytes.Buffer
	err = xml.NewEncoder(&wrtr).Encode(resp)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(wrtr.Bytes()), nil
}

func TestVerfiyValidTamperedWithDocFails(t *testing.T) {
	tampered, err := tamperedResponse(testResponse)
	require.Nil(t, err)
	tm, err := time.Parse(time.UnixDate, "Sun Apr 30 22:10:00 UTC 2017")
	require.Nil(t, err)
	clock := dsig.NewFakeClockAt(tm)
	validator, err := NewValidator(testMetadata, Clock(clock))
	require.Nil(t, err)
	require.NotNil(t, validator)
	auth, err := DecodeAuthResponse(tampered)
	_, err = validator.ValidateSignature(auth)
	require.NotNil(t, err)
}

// Message hasn't been tampered with but is stale
func TestVerfiyStaleMessageFails(t *testing.T) {
	tm, err := time.Parse(time.UnixDate, "Sun Apr 30 22:14:00 UTC 2017")
	require.Nil(t, err)
	clock := dsig.NewFakeClockAt(tm)
	validator, err := NewValidator(testMetadata, Clock(clock))
	require.Nil(t, err)
	require.NotNil(t, validator)

	auth, err := DecodeAuthResponse(testResponse)
	require.Nil(t, err)

	signed, err := validator.ValidateSignature(auth)
	require.Nil(t, err)
	require.NotNil(t, signed)

	err = validator.ValidateResponse(auth)
	assert.NotNil(t, err)
}

var testGoogleMetadata = `
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://accounts.google.com/o/saml2?idpid=C0171bstf" validUntil="2022-07-16T20:07:43.000Z">
  <md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:KeyDescriptor use="signing">
      <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <ds:X509Data>
          <ds:X509Certificate>MIIDdDCCAlygAwIBAgIGAV1SKeijMA0GCSqGSIb3DQEBCwUAMHsxFDASBgNVBAoTC0dvb2dsZSBJ
bmMuMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MQ8wDQYDVQQDEwZHb29nbGUxGDAWBgNVBAsTD0dv
b2dsZSBGb3IgV29yazELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWEwHhcNMTcwNzE3
MjAwNzQzWhcNMjIwNzE2MjAwNzQzWjB7MRQwEgYDVQQKEwtHb29nbGUgSW5jLjEWMBQGA1UEBxMN
TW91bnRhaW4gVmlldzEPMA0GA1UEAxMGR29vZ2xlMRgwFgYDVQQLEw9Hb29nbGUgRm9yIFdvcmsx
CzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAzLXNn7VmJBkvVNYHffTzDoow/8eSklauVeYjhELY6dtFv56wAQsFNeMovFUPxPeG
7Fci50/KStvoNZOdKqZFCwYkfI2ssXuMpBP37x2iprV7moVwGdGJb52elMNe0DesgTPbJ/IWIvzF
3GYxqYCHUlHuzJEzBYsdtvM8T/PClBxiLXRNbnjotzleFqb25w3XRfayOZg5GdQPeEmceWXDBhCa
eQyEPOrUTZ+//pZXSuKnOyaFfESNFNgvQJlYQQukjnhPtf674eWT6OdgZHyq8EBbZKfEhs5+KiAN
U43bDh9rpTJCB7rAKk1BFAW3r72pggwN9Z/sfp/C5B7uKAM5hwIDAQABMA0GCSqGSIb3DQEBCwUA
A4IBAQAZXypikbbRzichNXLdK96M/do9nGS5Q3xVgA2uxTzm/6qNkAfOSGSk8OcLrppPonbohkeZ
WVnNB5VZZava4DoSZ6OZsvKc1FM0wKvPJd83KUb7Syk1bV7TkT8DPEclfsLnn5s5g0oHlhsqkNly
0WPFTAoGHXYyOKGEARPoC/o+ZfgfvoMNyZkSQHiRboVVP2cT1ckJt4iCA65hNGXte29hSGmnX7QG
QyrBRp8n4UR9PjoeIy0tTCmG0tqu/NackFH4PkamY84Etxe9uH0StmkhID46QTT4Cv2+jqCaklg+
7VYqXbY64Wc/k0sK7WI1o3IVLWAPNb8ajV6Eo0Y8u+1N</ds:X509Certificate>
        </ds:X509Data>
      </ds:KeyInfo>
    </md:KeyDescriptor>
    <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat>
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://accounts.google.com/o/saml2/idp?idpid=C0171bstf"/>
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://accounts.google.com/o/saml2/idp?idpid=C0171bstf"/>
  </md:IDPSSODescriptor>
</md:EntityDescriptor>

`

func TestVerifyValidGoogleResponse(t *testing.T) {
	samlResponse :=
		`PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+PHNhbWwycDpSZXNwb25zZSB4bWxuczpzYW1sMnA9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCIgRGVzdGluYXRpb249Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvYXBpL3YxL2tvbGlkZS9zc28vY2FsbGJhY2siIElEPSJfODM1NzlhOTAwOGVmNzI2Zjg3YzUyYWFkNGI2ZGNjMDQiIEluUmVzcG9uc2VUbz0iU0dKaGkxZzVENC9ucE93WGF3OHQ2QT09IiBJc3N1ZUluc3RhbnQ9IjIwMTctMDctMThUMTQ6NDc6MDguMDM1WiIgVmVyc2lvbj0iMi4wIj48c2FtbDI6SXNzdWVyIHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIj5odHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9zYW1sMj9pZHBpZD1DMDE3MWJzdGY8L3NhbWwyOklzc3Vlcj48c2FtbDJwOlN0YXR1cz48c2FtbDJwOlN0YXR1c0NvZGUgVmFsdWU9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpzdGF0dXM6U3VjY2VzcyIvPjwvc2FtbDJwOlN0YXR1cz48c2FtbDI6QXNzZXJ0aW9uIHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIiBJRD0iXzUwMDA2MTk5MGFjYzAwNzIzMjg4ODMzYTMyN2NjOTg2IiBJc3N1ZUluc3RhbnQ9IjIwMTctMDctMThUMTQ6NDc6MDguMDM1WiIgVmVyc2lvbj0iMi4wIj48c2FtbDI6SXNzdWVyPmh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL3NhbWwyP2lkcGlkPUMwMTcxYnN0Zjwvc2FtbDI6SXNzdWVyPjxkczpTaWduYXR1cmUgeG1sbnM6ZHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPjxkczpTaWduZWRJbmZvPjxkczpDYW5vbmljYWxpemF0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8xMC94bWwtZXhjLWMxNG4jIi8+PGRzOlNpZ25hdHVyZU1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMDQveG1sZHNpZy1tb3JlI3JzYS1zaGEyNTYiLz48ZHM6UmVmZXJlbmNlIFVSST0iI181MDAwNjE5OTBhY2MwMDcyMzI4ODgzM2EzMjdjYzk4NiI+PGRzOlRyYW5zZm9ybXM+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMTAveG1sLWV4Yy1jMTRuIyIvPjwvZHM6VHJhbnNmb3Jtcz48ZHM6RGlnZXN0TWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjc2hhMjU2Ii8+PGRzOkRpZ2VzdFZhbHVlPm5abWdLOVh0anlUN3NCQXBVMHR5WmJVRTRXV013Q3NEejhqNklaRTVJeHc9PC9kczpEaWdlc3RWYWx1ZT48L2RzOlJlZmVyZW5jZT48L2RzOlNpZ25lZEluZm8+PGRzOlNpZ25hdHVyZVZhbHVlPkRIZFUrTG5PWC91OEh1angrSXBEbW96dDl1MlJPRDlVVTJPYjVFbDBaakVwQUVTcXlZMlBqOVk0S2QwMUlzRFRmL2dGS0pXT3lWTXoKUFAzaW81UDRlaUE5NnArMGcwWU51TzZpY2tWRjlCSEFKeWpFVDM4QzNwQjk1cmdxVWI3ckxhRDZYZGZBWEZRN2wyZGFsSFM5eUxhLwpLQnRUM2YzeWtZUGI3NE5yQWhpaFY4WjBndlBweVdxQkRnMjNCNzZ0SWVyV24yNkxvb1prUE5YUFRHdi9zeThvY1k1b3o1NnBsS3ZaCk9tVmR3cHp3SDcvN2kvVUVuTnY2c2lzMy9lczBPbW01Z3hlS0xQNDB2V2I5bFRtMUhtdkxUVjNzWmlIWlFRbVV3bWZjc1pMNmd5VkUKZWFKTkRRUDR5T3crdlhLZGV5QWxWQzZqdHQwNk1nWTlWMHpqNWc9PTwvZHM6U2lnbmF0dXJlVmFsdWU+PGRzOktleUluZm8+PGRzOlg1MDlEYXRhPjxkczpYNTA5U3ViamVjdE5hbWU+U1Q9Q2FsaWZvcm5pYSxDPVVTLE9VPUdvb2dsZSBGb3IgV29yayxDTj1Hb29nbGUsTD1Nb3VudGFpbiBWaWV3LE89R29vZ2xlIEluYy48L2RzOlg1MDlTdWJqZWN0TmFtZT48ZHM6WDUwOUNlcnRpZmljYXRlPk1JSURkRENDQWx5Z0F3SUJBZ0lHQVYxU0tlaWpNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1Ic3hGREFTQmdOVkJBb1RDMGR2YjJkc1pTQkoKYm1NdU1SWXdGQVlEVlFRSEV3MU5iM1Z1ZEdGcGJpQldhV1YzTVE4d0RRWURWUVFERXdaSGIyOW5iR1V4R0RBV0JnTlZCQXNURDBkdgpiMmRzWlNCR2IzSWdWMjl5YXpFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V3SGhjTk1UY3dOekUzCk1qQXdOelF6V2hjTk1qSXdOekUyTWpBd056UXpXakI3TVJRd0VnWURWUVFLRXd0SGIyOW5iR1VnU1c1akxqRVdNQlFHQTFVRUJ4TU4KVFc5MWJuUmhhVzRnVm1sbGR6RVBNQTBHQTFVRUF4TUdSMjl2WjJ4bE1SZ3dGZ1lEVlFRTEV3OUhiMjluYkdVZ1JtOXlJRmR2Y21zeApDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpMWE5uN1ZtSkJrdlZOWUhmZlR6RG9vdy84ZVNrbGF1VmVZamhFTFk2ZHRGdjU2d0FRc0ZOZU1vdkZVUHhQZUcKN0ZjaTUwL0tTdHZvTlpPZEtxWkZDd1lrZkkyc3NYdU1wQlAzN3gyaXByVjdtb1Z3R2RHSmI1MmVsTU5lMERlc2dUUGJKL0lXSXZ6RgozR1l4cVlDSFVsSHV6SkV6QllzZHR2TThUL1BDbEJ4aUxYUk5ibmpvdHpsZUZxYjI1dzNYUmZheU9aZzVHZFFQZUVtY2VXWERCaENhCmVReUVQT3JVVForLy9wWlhTdUtuT3lhRmZFU05GTmd2UUpsWVFRdWtqbmhQdGY2NzRlV1Q2T2RnWkh5cThFQmJaS2ZFaHM1K0tpQU4KVTQzYkRoOXJwVEpDQjdyQUtrMUJGQVczcjcycGdnd045Wi9zZnAvQzVCN3VLQU01aHdJREFRQUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQQpBNElCQVFBWlh5cGlrYmJSemljaE5YTGRLOTZNL2RvOW5HUzVRM3hWZ0EydXhUem0vNnFOa0FmT1NHU2s4T2NMcnBwUG9uYm9oa2VaCldWbk5CNVZaWmF2YTREb1NaNk9ac3ZLYzFGTTB3S3ZQSmQ4M0tVYjdTeWsxYlY3VGtUOERQRWNsZnNMbm41czVnMG9IbGhzcWtObHkKMFdQRlRBb0dIWFl5T0tHRUFSUG9DL28rWmZnZnZvTU55WmtTUUhpUmJvVlZQMmNUMWNrSnQ0aUNBNjVoTkdYdGUyOWhTR21uWDdRRwpReXJCUnA4bjRVUjlQam9lSXkwdFRDbUcwdHF1L05hY2tGSDRQa2FtWTg0RXR4ZTl1SDBTdG1raElENDZRVFQ0Q3YyK2pxQ2FrbGcrCjdWWXFYYlk2NFdjL2swc0s3V0kxbzNJVkxXQVBOYjhhalY2RW8wWTh1KzFOPC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9kczpTaWduYXR1cmU+PHNhbWwyOlN1YmplY3Q+PHNhbWwyOk5hbWVJRCBGb3JtYXQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjEuMTpuYW1laWQtZm9ybWF0OmVtYWlsQWRkcmVzcyI+am9obkBlZGlsb2submV0PC9zYW1sMjpOYW1lSUQ+PHNhbWwyOlN1YmplY3RDb25maXJtYXRpb24gTWV0aG9kPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6Y206YmVhcmVyIj48c2FtbDI6U3ViamVjdENvbmZpcm1hdGlvbkRhdGEgSW5SZXNwb25zZVRvPSJTR0poaTFnNUQ0L25wT3dYYXc4dDZBPT0iIE5vdE9uT3JBZnRlcj0iMjAxNy0wNy0xOFQxNDo1MjowOC4wMzVaIiBSZWNpcGllbnQ9Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvYXBpL3YxL2tvbGlkZS9zc28vY2FsbGJhY2siLz48L3NhbWwyOlN1YmplY3RDb25maXJtYXRpb24+PC9zYW1sMjpTdWJqZWN0PjxzYW1sMjpDb25kaXRpb25zIE5vdEJlZm9yZT0iMjAxNy0wNy0xOFQxNDo0MjowOC4wMzVaIiBOb3RPbk9yQWZ0ZXI9IjIwMTctMDctMThUMTQ6NTI6MDguMDM1WiI+PHNhbWwyOkF1ZGllbmNlUmVzdHJpY3Rpb24+PHNhbWwyOkF1ZGllbmNlPmtvbGlkZS5lZGlsb2submV0PC9zYW1sMjpBdWRpZW5jZT48L3NhbWwyOkF1ZGllbmNlUmVzdHJpY3Rpb24+PC9zYW1sMjpDb25kaXRpb25zPjxzYW1sMjpBdHRyaWJ1dGVTdGF0ZW1lbnQ+PHNhbWwyOkF0dHJpYnV0ZSBOYW1lPSJteWF0dHJpYnV0ZSI+PHNhbWwyOkF0dHJpYnV0ZVZhbHVlIHhtbG5zOnhzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYSIgeG1sbnM6eHNpPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSIgeHNpOnR5cGU9InhzOmFueVR5cGUiPmpvaG5AZWRpbG9rLm5ldDwvc2FtbDI6QXR0cmlidXRlVmFsdWU+PC9zYW1sMjpBdHRyaWJ1dGU+PC9zYW1sMjpBdHRyaWJ1dGVTdGF0ZW1lbnQ+PHNhbWwyOkF1dGhuU3RhdGVtZW50IEF1dGhuSW5zdGFudD0iMjAxNy0wNy0xOFQxNDozMzo0MS4wMDBaIiBTZXNzaW9uSW5kZXg9Il81MDAwNjE5OTBhY2MwMDcyMzI4ODgzM2EzMjdjYzk4NiI+PHNhbWwyOkF1dGhuQ29udGV4dD48c2FtbDI6QXV0aG5Db250ZXh0Q2xhc3NSZWY+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmFjOmNsYXNzZXM6dW5zcGVjaWZpZWQ8L3NhbWwyOkF1dGhuQ29udGV4dENsYXNzUmVmPjwvc2FtbDI6QXV0aG5Db250ZXh0Pjwvc2FtbDI6QXV0aG5TdGF0ZW1lbnQ+PC9zYW1sMjpBc3NlcnRpb24+PC9zYW1sMnA6UmVzcG9uc2U+`
	tm, err := time.Parse(time.RFC3339, "2017-07-18T14:47:08.035Z")
	require.Nil(t, err)
	clock := dsig.NewFakeClockAt(tm)
	validator, err := NewValidator(testGoogleMetadata, Clock(clock))
	require.Nil(t, err)
	require.NotNil(t, validator)
	auth, err := DecodeAuthResponse(samlResponse)
	require.Nil(t, err)
	require.NotNil(t, auth)
	signed, err := validator.ValidateSignature(auth)
	require.Nil(t, err)
	require.NotNil(t, signed)
	err = validator.ValidateResponse(auth)
	assert.Nil(t, err)
}

func TestVerifyInvalidSignatureGoogleResponse(t *testing.T) {
	samlResponse :=
		`PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+PHNhbWwycDpSZXNwb25zZSB4bWxuczpzYW1sMnA9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpwcm90b2NvbCIgRGVzdGluYXRpb249Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvYXBpL3YxL2tvbGlkZS9zc28vY2FsbGJhY2siIElEPSJfODM1NzlhOTAwOGVmNzI2Zjg3YzUyYWFkNGI2ZGNjMDQiIEluUmVzcG9uc2VUbz0iU0dKaGkxZzVENC9ucE93WGF3OHQ2QT09IiBJc3N1ZUluc3RhbnQ9IjIwMTctMDctMThUMTQ6NDc6MDguMDM1WiIgVmVyc2lvbj0iMi4wIj48c2FtbDI6SXNzdWVyIHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIj5odHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9zYW1sMj9pZHBpZD1DMDE3MWJzdGY8L3NhbWwyOklzc3Vlcj48c2FtbDJwOlN0YXR1cz48c2FtbDJwOlN0YXR1c0NvZGUgVmFsdWU9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpzdGF0dXM6U3VjY2VzcyIvPjwvc2FtbDJwOlN0YXR1cz48c2FtbDI6QXNzZXJ0aW9uIHhtbG5zOnNhbWwyPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6YXNzZXJ0aW9uIiBJRD0iXzUwMDA2MTk5MGFjYzAwNzIzMjg4ODMzYTMyN2NjOTg2IiBJc3N1ZUluc3RhbnQ9IjIwMTctMDctMThUMTQ6NDc6MDguMDM1WiIgVmVyc2lvbj0iMi4wIj48c2FtbDI6SXNzdWVyPmh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL3NhbWwyP2lkcGlkPUMwMTcxYnN0Zjwvc2FtbDI6SXNzdWVyPjxkczpTaWduYXR1cmUgeG1sbnM6ZHM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyMiPjxkczpTaWduZWRJbmZvPjxkczpDYW5vbmljYWxpemF0aW9uTWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8xMC94bWwtZXhjLWMxNG4jIi8+PGRzOlNpZ25hdHVyZU1ldGhvZCBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMDQveG1sZHNpZy1tb3JlI3JzYS1zaGEyNTYiLz48ZHM6UmVmZXJlbmNlIFVSST0iI181MDAwNjE5OTBhY2MwMDcyMzI4ODgzM2EzMjdjYzk4NiI+PGRzOlRyYW5zZm9ybXM+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvMDkveG1sZHNpZyNlbnZlbG9wZWQtc2lnbmF0dXJlIi8+PGRzOlRyYW5zZm9ybSBBbGdvcml0aG09Imh0dHA6Ly93d3cudzMub3JnLzIwMDEvMTAveG1sLWV4Yy1jMTRuIyIvPjwvZHM6VHJhbnNmb3Jtcz48ZHM6RGlnZXN0TWV0aG9kIEFsZ29yaXRobT0iaHR0cDovL3d3dy53My5vcmcvMjAwMS8wNC94bWxlbmMjc2hhMjU2Ii8+PGRzOkRpZ2VzdFZhbHVlPm5abWdLOVh0anlUN3NCQXBVMHR5WmJVRTRXV013Q3NEejhqNklaRTVJeHc9PC9kczpEaWdlc3RWYWx1ZT48L2RzOlJlZmVyZW5jZT48L2RzOlNpZ25lZEluZm8+PGRzOlNpZ25hdHVyZVZhbHVlPkRIZFUrTG5PWC91OEh1angrSXBEbW96dDl1MlJPRDlVVTJPYjVFbDBaakVwQUVTcXlZMlBqOVk0S2QwMUlzRFRmL2dGS0pXT3lWTXoKUFAzaW81UDRlaUE5NnArMGcwWU51TzZpY2tWRjlCSEFKeWpFVDM4QzNwQjk1cmdxVWI3ckxhRDZYZGZBWEZRN2wyZGFsSFM5eUxhLwpLQnRUM2YzeWtZUGI3NE5yQWhpaFY4WjBndlBweVdxQkRnMjNCNzZ0SWVyV24yNkxvb1prUE5YUFRHdi9zeThvY1k1b3o1NnBsS3ZaCk9tVmR3cHp3SDcvN2kvVUVuTnY2c2lzMy9lczBPbW01Z3hlS0xQNDB2V2I5bFRtMUhtdkxUVjNzWmlIWlFRbVV3bWZjc1pMNmd5VkUKZWFKTkRRUDR5T3crdlhLZGV5QWxWQzZqdHQwNk1nWTlWMHpqNWc9PTwvZHM6U2lnbmF0dXJlVmFsdWU+PGRzOktleUluZm8+PGRzOlg1MDlEYXRhPjxkczpYNTA5U3ViamVjdE5hbWU+U1Q9Q2FsaWZvcm5pYSxDPVVTLE9VPUdvb2dsZSBGb3IgV29yayxDTj1Hb29nbGUsTD1Nb3VudGFpbiBWaWV3LE89R29vZ2xlIEluYy48L2RzOlg1MDlTdWJqZWN0TmFtZT48ZHM6WDUwOUNlcnRpZmljYXRlPk1JSURkRENDQWx5Z0F3SUJBZ0lHQVYxU0tlaWpNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1Ic3hGREFTQmdOVkJBb1RDMGR2YjJkc1pTQkoKYm1NdU1SWXdGQVlEVlFRSEV3MU5iM1Z1ZEdGcGJpQldhV1YzTVE4d0RRWURWUVFERXdaSGIyOW5iR1V4R0RBV0JnTlZCQXNURDBkdgpiMmRzWlNCR2IzSWdWMjl5YXpFTE1Ba0dBMVVFQmhNQ1ZWTXhFekFSQmdOVkJBZ1RDa05oYkdsbWIzSnVhV0V3SGhjTk1UY3dOekUzCk1qQXdOelF6V2hjTk1qSXdOekUyTWpBd056UXpXakI3TVJRd0VnWURWUVFLRXA0SGIyOW5iR1VnU1c1akxqRVdNQlFHQTFVRUJ4TU4KVFc5MWJuUmhhVzRnVm1sbGR6RVBNQTBHQTFVRUF4TUdSMjl2WjJ4bE1SZ3dGZ1lEVlFRTEV3OUhiMjluYkdVZ1JtOXlJRmR2Y21zeApDekFKQmdOVkJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXpMWE5uN1ZtSkJrdlZOWUhmZlR6RG9vdy84ZVNrbGF1VmVZamhFTFk2ZHRGdjU2d0FRc0ZOZU1vdkZVUHhQZUcKN0ZjaTUwL0tTdHZvTlpPZEtxWkZDd1lrZkkyc3NYdU1wQlAzN3gyaXByVjdtb1Z3R2RHSmI1MmVsTU5lMERlc2dUUGJKL0lXSXZ6RgozR1l4cVlDSFVsSHV6SkV6QllzZHR2TThUL1BDbEJ4aUxYUk5ibmpvdHpsZUZxYjI1dzNYUmZheU9aZzVHZFFQZUVtY2VXWERCaENhCmVReUVQT3JVVForLy9wWlhTdUtuT3lhRmZFU05GTmd2UUpsWVFRdWtqbmhQdGY2NzRlV1Q2T2RnWkh5cThFQmJaS2ZFaHM1K0tpQU4KVTQzYkRoOXJwVEpDQjdyQUtrMUJGQVczcjcycGdnd045Wi9zZnAvQzVCN3VLQU01aHdJREFRQUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQQpBNElCQVFBWlh5cGlrYmJSemljaE5YTGRLOTZNL2RvOW5HUzVRM3hWZ0EydXhUem0vNnFOa0FmT1NHU2s4T2NMcnBwUG9uYm9oa2VaCldWbk5CNVZaWmF2YTREb1NaNk9ac3ZLYzFGTTB3S3ZQSmQ4M0tVYjdTeWsxYlY3VGtUOERQRWNsZnNMbm41czVnMG9IbGhzcWtObHkKMFdQRlRBb0dIWFl5T0tHRUFSUG9DL28rWmZnZnZvTU55WmtTUUhpUmJvVlZQMmNUMWNrSnQ0aUNBNjVoTkdYdGUyOWhTR21uWDdRRwpReXJCUnA4bjRVUjlQam9lSXkwdFRDbUcwdHF1L05hY2tGSDRQa2FtWTg0RXR4ZTl1SDBTdG1raElENDZRVFQ0Q3YyK2pxQ2FrbGcrCjdWWXFYYlk2NFdjL2swc0s3V0kxbzNJVkxXQVBOYjhhalY2RW8wWTh1KzFOPC9kczpYNTA5Q2VydGlmaWNhdGU+PC9kczpYNTA5RGF0YT48L2RzOktleUluZm8+PC9kczpTaWduYXR1cmU+PHNhbWwyOlN1YmplY3Q+PHNhbWwyOk5hbWVJRCBGb3JtYXQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjEuMTpuYW1laWQtZm9ybWF0OmVtYWlsQWRkcmVzcyI+am9obkBlZGlsb2submV0PC9zYW1sMjpOYW1lSUQ+PHNhbWwyOlN1YmplY3RDb25maXJtYXRpb24gTWV0aG9kPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6Y206YmVhcmVyIj48c2FtbDI6U3ViamVjdENvbmZpcm1hdGlvbkRhdGEgSW5SZXNwb25zZVRvPSJTR0poaTFnNUQ0L25wT3dYYXc4dDZBPT0iIE5vdE9uT3JBZnRlcj0iMjAxNy0wNy0xOFQxNDo1MjowOC4wMzVaIiBSZWNpcGllbnQ9Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvYXBpL3YxL2tvbGlkZS9zc28vY2FsbGJhY2siLz48L3NhbWwyOlN1YmplY3RDb25maXJtYXRpb24+PC9zYW1sMjpTdWJqZWN0PjxzYW1sMjpDb25kaXRpb25zIE5vdEJlZm9yZT0iMjAxNy0wNy0xOFQxNDo0MjowOC4wMzVaIiBOb3RPbk9yQWZ0ZXI9IjIwMTctMDctMThUMTQ6NTI6MDguMDM1WiI+PHNhbWwyOkF1ZGllbmNlUmVzdHJpY3Rpb24+PHNhbWwyOkF1ZGllbmNlPmtvbGlkZS5lZGlsb2submV0PC9zYW1sMjpBdWRpZW5jZT48L3NhbWwyOkF1ZGllbmNlUmVzdHJpY3Rpb24+PC9zYW1sMjpDb25kaXRpb25zPjxzYW1sMjpBdHRyaWJ1dGVTdGF0ZW1lbnQ+PHNhbWwyOkF0dHJpYnV0ZSBOYW1lPSJteWF0dHJpYnV0ZSI+PHNhbWwyOkF0dHJpYnV0ZVZhbHVlIHhtbG5zOnhzPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYSIgeG1sbnM6eHNpPSJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSIgeHNpOnR5cGU9InhzOmFueVR5cGUiPmpvaG5AZWRpbG9rLm5ldDwvc2FtbDI6QXR0cmlidXRlVmFsdWU+PC9zYW1sMjpBdHRyaWJ1dGU+PC9zYW1sMjpBdHRyaWJ1dGVTdGF0ZW1lbnQ+PHNhbWwyOkF1dGhuU3RhdGVtZW50IEF1dGhuSW5zdGFudD0iMjAxNy0wNy0xOFQxNDozMzo0MS4wMDBaIiBTZXNzaW9uSW5kZXg9Il81MDAwNjE5OTBhY2MwMDcyMzI4ODgzM2EzMjdjYzk4NiI+PHNhbWwyOkF1dGhuQ29udGV4dD48c2FtbDI6QXV0aG5Db250ZXh0Q2xhc3NSZWY+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmFjOmNsYXNzZXM6dW5zcGVjaWZpZWQ8L3NhbWwyOkF1dGhuQ29udGV4dENsYXNzUmVmPjwvc2FtbDI6QXV0aG5Db250ZXh0Pjwvc2FtbDI6QXV0aG5TdGF0ZW1lbnQ+PC9zYW1sMjpBc3NlcnRpb24+PC9zYW1sMnA6UmVzcG9uc2U+`
	tm, err := time.Parse(time.RFC3339, "2017-07-18T14:47:08.035Z")
	require.Nil(t, err)
	clock := dsig.NewFakeClockAt(tm)
	validator, err := NewValidator(testGoogleMetadata, Clock(clock))
	require.Nil(t, err)
	require.NotNil(t, validator)
	auth, err := DecodeAuthResponse(samlResponse)
	require.Nil(t, err)
	require.NotNil(t, auth)
	_, err = validator.ValidateSignature(auth)
	require.NotNil(t, err)
}

// validate id's are unique and that I didn't screw up my maths
func TestIDGenerator(t *testing.T) {
	idTable := make(map[string]struct{})
	for i := 0; i < 100; i++ {
		id, err := generateSAMLValidID()
		require.Nil(t, err)
		assert.Subset(t, []byte(idAlphabet), []byte(id))
		_, ok := idTable[id]
		assert.False(t, ok)
		idTable[id] = struct{}{}
	}
}

func TestIDPrefix(t *testing.T) {
	// Ensure ID comes with the appropriate prefix
	id, err := generateSAMLValidID()
	require.Nil(t, err)
	assert.Equal(t, idPrefix, id[:2])
}
