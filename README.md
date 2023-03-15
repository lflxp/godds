# 资料

```bash
gcc -E HelloWorldData.c -o h.i -I /opt/gopath/src/github.com/eclipse-cyclonedds/cyclonedds/install/include
gcc -S h.i -o h.s
gcc -c h.s -o HelloWorldData.o 
```

# LD_LIBRARY_PATH配置

> export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:./library/lib

# golang cross complie

1. 可选为：windows, unix, posix, plan9, darwin, bsd, linux, freebsd, nacl, netbsd, openbsd, solaris, dragonfly, bsd, notbsd， android，stubs

## 构建约束(https://learnku.com/go/t/34696)

Go 语言中注释的第二个特殊用途是构建约束。

Go 作为一种编程语言的关键特性是它支持各种操作系统和体系结构。通常情况下，相同的代码可以用于多个平台。但是在某些情况下，特定于操作系统或体系结构的代码应该只用于特定的目标。标准的 go build 工具可以通过理解以操作系统名称和 / 或体系结构结尾的程序应该只用于匹配这些标记的目标来处理某些情况。例如，一个名为 foo_linux.go 的文件，将只会被 Linux 系统编译。 foo_amd64.go 用于 AMD64 架构，foo_windows_amd64.go 用于运行在 AMD64 架构上的 64 位 Windows 系统。

然而，这些命名约束在更复杂的情况下会失败，例如当同一代码可以用于多个 (但不是所有) 操作系统时。在这些情况下，Go 有构建约束的概念 —— 在编译 Go 程序时，go build 将读取经过特殊设计的注释，以确定要引用哪些文件。

构建约束的注释遵循以下规则:

以前缀 +build 开始，后跟一个或多个空格
位于文件顶部在包声明之前
在它和包声明之间至少有一个空行，以防止它被视为包文档
而不是将文件命名为 foo_linux.go。我们可以把下面的注释放在文件 foo.go 的开头：

```
// +build linux
```

然而，当引用多个体系结构和 / 或操作系统时，构建约束的威力就显现出来了。Go 为组合构建约束制定了以下规则：

以 ！开头的构建标记是无效的
用空格分隔的构建标记在逻辑上是或
用逗号分隔的构建标记在逻辑上是和
在多行上构建约束是逻辑和
根据上面的规则，下面的约束将把文件限制为 Linux 或 Darwin (MacOS)：

```
// +build linux darwin
```

而这个约束同时需要 Windows 和 i386：
```
// +build windows,386
```
上面的约束也可以写在下面的两行中：
```
// +build windows
// +build 386
```
除了指定操作系统和体系结构之外，构建约束可以通过 ignore 标记的常见用法来完全忽略文件（任何不匹配有效架构或操作系统的文本是可以工作的）：

应该注意的是，这些构建约束（以及前面提到的命名约定）也适用于测试文件，因此可以以类似的方式执行特定于体系结构 / 操作系统的测试。

构建约束的全部功能在这里的 go build 文档中有详细说明。