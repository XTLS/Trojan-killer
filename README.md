# Trojan-Killer

这个 POC 是为了 **狠 狠 打 脸** 某些认为 TLS in TLS 检测不存在或成本很高的人。

该程序在 `127.0.0.1:12345` 接收 TLS 流量，并用 **非 常 廉 价** 的方式检测出其中的 Trojan 代理。

1. 设置浏览器的 HTTP 代理至 `127.0.0.1:12345`，观察该程序的输出。
2. 设置 Trojan 链式 HTTP 代理至 `127.0.0.1:12345`，观察该程序的输出。

我们的测试结果如下：

1. 对于浏览器的 HTTPS 流量，**几乎没有阳性结果**。
2. 对于 Trojan 的 TLS in TLS 流量，**Trojan 字样直接刷屏**。

这与我们多次收到的 Trojan 被封、XTLS Vision 存活的反馈相符（它们均可选 Golang 指纹）。

值得一提的是，根据我们的观察，目前 REALITY 的“白名单域名”会被豁免于这样的检测。

## License

[GNU AFFERO GENERAL PUBLIC LICENSE](https://github.com/XTLS/Trojan-killer/blob/main/LICENSE)

## Compilation

```bash
go build -trimpath -ldflags "-s -w -buildid=" .
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/XTLS/Trojan-killer.svg)](https://starchart.cc/XTLS/Trojan-killer)
