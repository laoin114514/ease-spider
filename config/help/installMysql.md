在Ubuntu服务器上手动安装MySQL（不使用APT包管理器），可通过以下步骤完成。这里以**MySQL 8.0**为例，使用官方二进制压缩包安装：

---

### **步骤 1：安装依赖**
```bash
sudo apt update
sudo apt install libaio1 libnuma1 wget
```

---

### **步骤 2：下载MySQL二进制包**
```bash
# 进入临时目录
cd /tmp

# 下载MySQL社区版压缩包（替换版本号为你需要的版本）
wget https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.xx-linux-glibc2.17-x86_64.tar.xz

# 解压文件
tar -xvf mysql-8.0.xx-linux-glibc2.17-x86_64.tar.xz

# 移动到安装目录（通常为/usr/local）
sudo mv mysql-8.0.xx-linux-glibc2.17-x86_64 /usr/local/mysql
```

---

### **步骤 3：创建MySQL用户和组**
```bash
sudo groupadd mysql
sudo useradd -r -g mysql -s /bin/false mysql
```

---

### **步骤 4：配置文件和目录**
```bash
# 创建数据目录
sudo mkdir -p /var/lib/mysql
sudo chown mysql:mysql /var/lib/mysql

# 创建配置文件
sudo cp /usr/local/mysql/support-files/mysql.server /etc/init.d/mysqld
sudo chmod +x /etc/init.d/mysqld
```

---

### **步骤 5：初始化MySQL**
```bash
cd /usr/local/mysql
sudo bin/mysqld --initialize --user=mysql --basedir=/usr/local/mysql --datadir=/var/lib/mysql
```
**注意**：初始化后终端会显示**临时root密码**（如 `[Note] [MY-010454] A temporary password is generated for root@localhost: Abc123!`），务必保存！

---

### **步骤 6：配置SSL（可选但推荐）**
```bash
sudo bin/mysql_ssl_rsa_setup --datadir=/var/lib/mysql
```

---

### **步骤 7：启动MySQL服务**
```bash
sudo /etc/init.d/mysqld start
```

---

### **步骤 8：修改root密码**
```bash
# 使用临时密码登录
sudo bin/mysql -u root -p

# 在MySQL命令行中修改密码
ALTER USER 'root'@'localhost' IDENTIFIED BY '你的新密码';
FLUSH PRIVILEGES;
EXIT;
```

---

### **步骤 9：设置环境变量（可选）**
```bash
echo 'export PATH=/usr/local/mysql/bin:$PATH' | sudo tee -a /etc/profile
source /etc/profile
```

---

### **步骤 10：配置开机自启**
```bash
# 创建systemd服务文件
sudo nano /etc/systemd/system/mysql.service
```
粘贴以下内容：
```ini
[Unit]
Description=MySQL Server
After=network.target

[Service]
User=mysql
Group=mysql
WorkingDirectory=/usr/local/mysql
ExecStart=/usr/local/mysql/bin/mysqld --basedir=/usr/local/mysql --datadir=/var/lib/mysql
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
启用服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable mysql
sudo systemctl start mysql
```

---

### **验证安装**
```bash
mysql -u root -p -e "SELECT VERSION();"
```
应输出类似：`8.0.xx`

---

### **常见问题解决**
1. **启动失败**：
   - 检查日志：`sudo tail -f /var/lib/mysql/主机名.err`
   - 确保目录权限：`sudo chown -R mysql:mysql /var/lib/mysql`

2. **忘记临时密码**：
   ```bash
   sudo rm -rf /var/lib/mysql/*  # 删除数据目录
   sudo bin/mysqld --initialize ... # 重新初始化
   ```

---

> **强烈建议**：生产环境推荐使用APT安装（`sudo apt install mysql-server`），手动安装仅适用于特殊需求。如需卸载，删除相关目录并移除服务即可。 