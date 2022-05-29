package example

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"time"

	"github.com/Ning-Qing/did-go/core"
)

func Authentication() {
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
}
