# 部署 production prover

- **在目标机器创建文件目录**

```
> tree .

.
├── config
│   └── prover.config.json
├── ethda-zkevm-prover.yml
```

- **复制`test/config/test.prover.config.json`到`config/prover.config.json`并修改配置**

```json
{
    // ...
    "runAggregatorServer": false,
    "runAggregatorClient": true, // 开启AggregatorClient
    "runAggregatorClientMock": false, // 关闭mock prover

    // ...
    "outputPath": "output",
    "configPath": "/config",  // 新增configPath指定config路径

    // ...

    "aggregatorServerPort": 50081,
    "aggregatorClientPort": 50081,
    "aggregatorClientHost": "<zkevm-aggregator-ip>", // 修改为ethda-node.zkevm-aggregator对应的IP

    // ...

    "databaseURL": "postgresql://prover_user:prover_pass@<zkevm-state-db-ip>:5432/prover_db", // 修改为ethda-node.zkevm-state-db的数据库访问url
    
    // ...
}
```

- **`ethda-zkevm-prover.yml`示例**

```yaml
version: "3.5"
networks:
  default:
    name: ethda-zkevm

services:
  zkevm-prover:
    container_name: ethda-zkevm-prover
    image: hermeznetwork/zkevm-prover:v6.0.0
    environment:
      - EXPERIMENTAL_DOCKER_DESKTOP_FORCE_QEMU=1
    ports:
      - 51061:50061 # MT
      - 51071:50071 # Executor
    volumes:
      - ./v6.0.0-rc.1-fork.9/config:/config
      - ./config/prover.config.json:/usr/src/app/config.json
    command: >
      zkProver -c /usr/src/app/config.json
```

- **启动prover**

```
docker-compose -f ethda-zkevm-prover.yml up -d
```