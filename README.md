### 敏感词验证服务

### 噪声词库
./dicts/noise

#### 启动服务
```code
docker run -d -p 1325:1325 go-toon-sensitive
```

#### 访问
```code
http://localhost:1325/q=敏感词
```
