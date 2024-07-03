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

- **临时将`GetL1BlockUpgradeLxLy`的扫描范围调小**
代码位置: `etherman/etherman.go` 第398行

```golang
 func (etherMan *Client) GetL1BlockUpgradeLxLy(ctx context.Context, genesisBlock uint64) (uint64, error) {
        it, err := etherMan.EtrogRollupManager.FilterInitialized(&bind.FilterOpts{
                Start:   6032210,  // 将Start号从1修改为接近Sepolia Rollup的创建区块高度，如(test.genesis.config.json.genesisBlockNumber - 10)
                End:     &genesisBlock,
                Context: ctx,
        })
		// ...
 }
```

- **如果需要连接独立部署的prover，配置test.node.config.toml**

```yaml
[MTClient]
URI = "<prover-IP>:51061" // 修改为独立prover的IP和对应的MT端口

[Executor]
URI = "<prover-IP>:51071" // 修改为独立prover的IP和对应的Executor端口
MaxGRPCMessageSize = 100000000
```

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
	#$(RUNL1NETWORK)
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

- **修改blob交易的DasAddress地址**

```
git clone git@github.com:crustio/zkblob-contracts.git
```
使用sequencer账户部署DAS.sol合约，得到合约地址\<DasAddress\>

将`/test/config/test.node.config.toml`中Blob部分的DasAddress修改成\<DasAddress\>:

```toml
[Blob]
DasAddress = "<DasAddress>"
```

- **重启json-rpc**

```shell
make stop-json-rpc && make run-json-rpc
```

- **修改blob-utils**

```
git clone git@github.com:crustio/blob-utils.git
```

修改blob.go开头的常量为\<ToAddress\>:

```golang
const (
	BlobToAddress = "<ToAddress>"
)
```

给blob-utils仓库重新打tag:

```
git tag v0.1.x 
```

- **cdk-validium-node仓库引用新的blob-utils分支**

```
git clone git@github.com:crustio/cdk-validium-node.git

git checkout ethda
```

更新blob-utils的版本