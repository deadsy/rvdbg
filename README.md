[![Go Report Card](https://goreportcard.com/badge/github.com/deadsy/rvdbg)](https://goreportcard.com/report/github.com/deadsy/rvdbg)
[![GoDoc](https://godoc.org/github.com/deadsy/rvdbg?status.svg)](https://godoc.org/github.com/deadsy/rvdbg)

# rvdbg
RISC-V Debugger

```
$ ./cmd/rvdbg/rvdbg --help
Usage of ./cmd/rvdbg/rvdbg:
  -i string
        debug interface name
  -t string
        target name

debug interfaces:
        daplink     ARM DAPLink   
        jlink       Segger J-Link 

targets:
        gd32v       GD32V Board (GigaDevice GD32VF103VBT6 RISC-V RV32)            
        maixgo      SiPeed MaixGo (Kendryte K210, Dual Core RISC-V RV64)          
        redv        SparkFun RED-V RedBoard (SiFive FE310-G002 RISC-V RV32)       
```
