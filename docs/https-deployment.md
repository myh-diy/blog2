# 阿里云 HTTPS 部署

证书签发、Nginx 切换和请求转发的原理见 [HTTPS 工作流程说明](https-flow.md)。

## 前置条件

1. 准备域名 `myh97.xyz`。
2. 在 DNS 中添加 `A` 记录，指向 ECS 公网 IP。
3. 在阿里云 ECS 安全组开放入方向 TCP `80` 和 `443`。
4. 确认服务器没有其他程序占用这两个端口。

如果 ECS 位于中国大陆，域名对公网提供网站服务前还需要完成 ICP 备案；香港或海外地域通常不需要工信部 ICP 备案。

HTTP-01 验证必须能从公网访问 TCP 80。部署完成后，80 端口只提供 ACME 验证和到 HTTPS 的跳转。

## 首次部署

将项目部署到服务器后执行：

```bash
cd /opt/blog
cp .env.production.example .env.production
chmod 600 .env.production
vi .env.production
```

至少填写：

```dotenv
DOMAIN=myh97.xyz
CERTBOT_EMAIL=you@example.com
JWT_SECRET=替换为长随机字符串
ADMIN_USERNAME=admin
ADMIN_PASSWORD=替换为强密码
```

初始化 HTTPS：

```bash
chmod +x deploy/init-https.sh deploy/renew-cert.sh
./deploy/init-https.sh
```

脚本会先启动 HTTP Nginx，使用 Certbot webroot 申请证书，再自动切换到 HTTPS 配置。

## 日常更新

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml pull
docker compose --env-file .env.production -f docker-compose.prod.yml up -d --remove-orphans
docker compose --env-file .env.production -f docker-compose.prod.yml ps
```

`PUBLIC_URL` 会自动设置为 `https://DOMAIN`。应用容器的 8080 和 9101 端口不会暴露到公网。

## 自动续期

先测试一次：

```bash
./deploy/renew-cert.sh
```

然后执行 `crontab -e`，每天凌晨检查一次：

```cron
17 3 * * * cd /opt/blog && ./deploy/renew-cert.sh >> /var/log/blog-certbot.log 2>&1
```

Certbot 只会在证书接近过期时续期；脚本结束后会重新加载 Nginx。

## 检查与排错

```bash
docker compose --env-file .env.production -f docker-compose.prod.yml ps
docker compose --env-file .env.production -f docker-compose.prod.yml logs --tail=100 nginx
docker compose --env-file .env.production -f docker-compose.prod.yml logs --tail=100 certbot
curl -I http://myh97.xyz
curl -I https://myh97.xyz
```

正常情况下，HTTP 返回 `301`，HTTPS 返回 `200`。证书申请失败时，优先检查域名解析、阿里云安全组的 80/443 规则、本机防火墙和端口占用。
