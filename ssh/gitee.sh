pacman -S openssh

# golang安装gitee插件并登录，位于设置-->版本控制
 ssh-keygen -t rsa

# /home/ddk/.ssh/id_rsa.pub 粘贴到下面的这个网址添加公钥
# https://gitee.com/profile/sshkeys

ssh git@gitee.com
# 首次使用需要确认并添加主机到本机SSH可信列表 /home/ddk/.gitconfig

#验证信任列表
ssh -T git@gitee.com

#强制使用ssh get
git config --global url."git@gitee.com:".insteadOf "https://gitee.com/"

#pacman -S go
# fn 1.17

# sudo pacman -S mysql mysql-workbench sudo mysqld --initialize --user=mysql --basedir=/usr --datadir=/var/lib/mysql # 设置开机启动MySQL服务 systemctl enable mysqld.service systemctl daemon-reload systemctl start mysqld.service # 登录数据库 mysql -u root -p

# ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'root';

go env -w GOPROXY=goproxy.cn
go env -w GOPRIVATE=gitee.com

#go clean --modcache

#[[ $(command -v curl) ]] || sudo pacman -Syu curl
 #   bash -c "$(curl -L https://gitee.com/mo2/linux/raw/2/2)"
