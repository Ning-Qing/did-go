# did-go

用golang 实现 Decentralized Identifiers 规范

## 具体内容
core：计划可以作为一个模块被其他项目使用，实现自己的did应用
- claim 已完成
- did document 已完成
- did controller 进行中
- did 

provider: 实现持久化
wrapper：包装器

## 快速开始

```go
// 创建存储在内村中的文档集
	provider := core.NewMemoryProvider()

	// 常见did控制器
	controller := core.NewController(provider)

	// 常见密钥对
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	pubKey := privKey.PublicKey

	// 生成did
	did := controller.NewDID(core.NewJwkRSA(&pubKey))
	// 创建claim
	claim := &core.Claim{
		Context: []string{"github.com/Ning-Qing"},
		Doc: &core.Doc{
			Issuer:            did,
			IssuanceDate:      time.Now().String(),
			ExpirationDate:    time.Now().AddDate(1, 0, 0).String(),
			Revocation:        "github.com/Ning-Qing",
			CredentialSubject: nil,
		},
	}

	// 为claim签名并绑定验证方法

	// 获取文档
	doc := controller.GetDocument(did)
	// 获取认证方法,第一个baseauth方法，需要使用创建did文档的私钥签名
	log.Println(doc.Authentication)
	// 签名
	claim.Sign(doc.Authentication[0], privKey)

	// 验证

	log.Println(controller.AuthenticationClaim(claim))
```
