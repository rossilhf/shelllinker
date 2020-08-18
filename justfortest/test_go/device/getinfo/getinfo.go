package getinfo

import (
    "runtime"
)

func Get_curVersion() string {
    version := "V20200818"
	return version
}

func Get_curOs() string {
    sysType := runtime.GOOS
    if (sysType == "linux") || (sysType == "windows") {
        return sysType
    } else {
        sysType = "unknown"
    }
    return sysType
}    
