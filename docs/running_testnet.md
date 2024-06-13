# 部署ethda-node流程
## 生成2个账户

```
	git clone git@github.com:crustio/evmkey.git
	go install -ldflags="-w -s" github.com/crustio/evmkey@v1.0.2
	执行2次，输入keystore的密码后续使用，分别生成sequencer和aggregater账户
	evmkey account new
	
```
2个账户的助记词和keystore文件都保存在keystore目录下

## 在l1部署合约
- **部署zkevm-contract**
  与部署[cdk-validum-node](https://github.com/crustio/cdk-validium-node/blob/ethda/docs/running_ethda_sepolia.md)相同流程,**dataAvailabilityProtocol无需填写**
```
npm install && npm run deploy:testnet:v2:sepolia
```

- **部署zkblob**

```
git clone git@github.com:crustio/zkblob-contracts.git
```
将script/ZkBlob.s.sol中 new ZkBlob(0x20574f8eb8B7Bd3f6E3C0Aa749681290BB8308e9) 地址修改为和部署zkevm-contract相同的sequencer

## 部署ethda-node

- **使用最新代码**

```
git clone git@github.com:crustio/ethda-node.git
```

- **修改test.genesis.config.json(与cdk-validum-node相同)**
  增加一项l1Config.polygonZkBlobAddress 填写zkblob合约地址

- **修改test.node.config.toml(与cdk-validum-node相同)**

- **将test目录下sequencer.keystore和aggregator.keystore替换为部署合约的2个账户**

- **构建docker镜像，部署node**
```
cd test/ && make ship

.PHONY: ship
ship: ## Builds docker images and run them
	cd .. && make build-docker && cd ./test && make run

.PHONY: run
run: ## Runs a full node
	$(RUNSTATEDB)
	$(RUNPOOLDB)
	$(RUNEVENTDB)
	$(RUNL1NETWORK)
	sleep 1
	$(RUNZKPROVER)
	$(RUNAPPROVE)
	sleep 3
	$(RUNSYNC)
	sleep 4
	$(RUNETHTXMANAGER)
	$(RUNSEQUENCER)
	$(RUNSEQUENCESENDER)
	$(RUNL2GASPRICER)
	$(RUNAGGREGATOR)
	$(RUNJSONRPC)
```