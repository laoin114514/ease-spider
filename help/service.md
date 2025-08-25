
### 1. 服务文件的位置

Systemd 会在多个路径中查找服务文件，优先级从高到低：

1.  **系统目录（最高优先级）**: `/etc/systemd/system/` - **这是你放置自定义服务文件的最佳位置**，因为它不会因为系统更新而被覆盖。
2.  **系统运行时目录**: `/run/systemd/system/`
3.  **安装目录（最低优先级）**: `/usr/lib/systemd/system/` - 软件包安装的默认服务文件通常在这里。不要直接修改这里的文件。

### 2. 服务文件的基本结构和语法

-   **文件名**: 必须以 `.service` 结尾，例如 `myapp.service`。
-   **编码**: 必须是 UTF-8。
-   **结构**: 文件由 `[Unit]`, `[Service]`, `[Install]` 三个主要区块（Section）组成，每个区块包含一系列的指令（Directives）。

### 3. 常用指令详解

#### [Unit] 区块
用于定义服务的元数据以及与其他单元的关系。

-   `Description`: 服务的描述信息，用于标识服务是什么。
-   `After`: 定义在哪些目标（target）或服务之后启动。这并不构成强依赖关系，只是调整启动顺序。例如 `After=network.target` 表示在网络就绪后启动。
-   `Requires`: 定义强依赖关系。如果指定的服务失败，本服务也会停止。
-   `Wants`: 定义弱依赖关系。希望某个服务启动，但即使那个服务启动失败，本服务仍然继续启动。这是更常用的选项。

#### [Service] 区块
这是定义服务行为的核心区块。

-   `Type`: 定义进程的启动类型，非常重要。
    -   `simple`（默认值）： systemd 认为你启动的命令就是主进程，它会一直运行。
    -   `forking`: 程序会通过 fork 系统调用在后台运行（守护进程）。**systemd 会期望父进程退出，而子进程继续运行**。你必须同时指定 `PIDFile=` 来帮助 systemd 跟踪主进程。
    -   `oneshot`: 命令执行后就会退出，不会长期运行。通常需要与 `RemainAfterExit=yes` 配合使用，让 systemd 认为服务即使进程退出也仍处于“活跃状态”。
    -   `idle`: 类似于 simple，但 systemd 会等待所有活动任务处理完毕后再启动该服务，以避免输出与 shell 混合。
-   `ExecStart`: **最重要的指令**。指定启动服务所要执行的完整命令和参数。必须是绝对路径。
-   `ExecStop`: （可选）指定停止服务所要执行的命令。
-   `ExecReload`: （可选）指定重载服务配置所要执行的命令（例如发送 SIGHUP）。
-   `Restart`: 定义何时自动重启服务。常用值：
    -   `no`（默认）： 不自动重启。
    -   `on-failure`: 仅在进程以非零状态退出时重启。
    -   `on-abnormal`: 因信号终止或超时退出时重启。
    -   `always`: 总是重启。
    -   `unless-stopped`: 总是重启，除非服务被手动停止。
-   `RestartSec`: 在重启服务前，等待的秒数（例如 `RestartSec=5`）。
-   `User` 和 `Group`: 指定以哪个用户和组的身份运行进程。为了安全，不应使用 root。
-   `WorkingDirectory`: 指定进程的工作目录。
-   `Environment`: 设置环境变量（例如 `Environment="NODE_ENV=production"`）。
-   `StandardOutput` / `StandardError`: 指定标准输出和错误输出的重定向。
    -   `journal`（默认）： 输出到 systemd journal。
    -   `syslog`: 输出到 syslog。
    -   `file:/path/to/file`: 输出到文件。
    -   `inherit`: 继承 systemd 本身的文件描述符（通常意味着输出到控制台，但取决于启动 systemd 的方式）。

#### [Install] 区块
定义服务的安装信息，即如何使用 `systemctl enable` 命令启用它。

-   `WantedBy`: 最常用的指令。指定一个“目标”（target），当使用 `systemctl enable` 时，会在这个目标目录下创建一个符号链接。最常见的是 `multi-user.target`，表示多用户命令行模式。
-   `RequiredBy`: 指定服务被某个目标强依赖。
-   `Alias`: 为服务启用时提供一个别名。

### 4. 示例

假设我们有一个简单的 Python HTTP 服务器应用，位于 `/opt/myapp/app.py`，我们希望它以 `myuser` 用户身份运行。

#### 简单示例

`/etc/systemd/system/myapp.service`:

```ini
[Unit]
Description=My Simple Python Web Application
After=network.target

[Service]
Type=simple
User=myuser
Group=myuser
WorkingDirectory=/opt/myapp
ExecStart=/usr/bin/python3 app.py
Restart=on-failure
RestartSec=5

# 可选：输出到 journal 的同时也写入一个文件
Environment=PYTHONUNBUFFERED=1
StandardOutput=file:/var/log/myapp.log
StandardError=inherit

[Install]
WantedBy=multi-user.target
```

#### 复杂示例（仿照 nginx）

对于像 Nginx 这样会 fork 到后台的程序，需要使用 `Type=forking` 并指定 PID 文件。

```ini
[Unit]
Description=The Awesome Nginx Web Server
Documentation=man:nginx(8)
After=network.target nss-lookup.target

[Service]
Type=forking
PIDFile=/run/nginx.pid
ExecStartPre=/usr/sbin/nginx -t
ExecStart=/usr/sbin/nginx
ExecReload=/usr/sbin/nginx -s reload
ExecStop=/bin/kill -s QUIT $MAINPID
TimeoutStopSec=5
KillMode=mixed
Restart=on-failure
RestartSec=1

[Install]
WantedBy=multi-user.target
```

### 5. 如何使用写好的服务文件

1.  **保存文件**: 将文件保存到 `/etc/systemd/system/`，例如 `sudo vim /etc/systemd/system/myapp.service`。
2.  **重新加载 systemd 配置**: 每次修改服务文件后都需要执行，让 systemd 知道变化。
    ```bash
    sudo systemctl daemon-reload
    ```
3.  **启用服务（开机自启）**:
    ```bash
    sudo systemctl enable myapp.service
    ```
4.  **启动服务**:
    ```bash
    sudo systemctl start myapp.service
    ```
5.  **检查状态**:
    ```bash
    sudo systemctl status myapp.service
    ```
6.  **查看日志**:
    ```bash
    sudo journalctl -u myapp.service -f  # -f 表示持续输出
    ```
7.  **停止/重启/重载服务**:
    ```bash
    sudo systemctl stop myapp.service
    sudo systemctl restart myapp.service
    sudo systemctl reload myapp.service  # 如果支持的话
    ```

### 总结与最佳实践

-   **放置位置**: 始终使用 `/etc/systemd/system/`。
-   **安全第一**: 使用 `User` 和 `Group` 以非 root 用户运行你的服务。
-   **自动重启**: 对生产环境服务，合理使用 `Restart=on-failure` 或 `always`/`unless-stopped`。
-   **日志是朋友**: 利用 `journalctl` 来调试问题。如果程序本身不处理日志，可以通过 `StandardOutput` 和 `StandardError` 重定向。
-   **测试**: 写完配置文件后，务必先 `daemon-reload`，然后启动并检查状态和日志，确保一切按预期运行。
-   **`Type` 是关键**: 正确设置 `Type` 指令，尤其是对于会 fork 的守护进程，一定要用 `forking`。

通过这个指南，你应该能够为大多数应用程序编写出合适的 systemd 服务文件了。