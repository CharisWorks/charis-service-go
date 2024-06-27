package mail

import (
	"fmt"
	"time"
)

func PurchasedAdminEmailFactory(
	name string,
	email string,
	address string,
	itemName string,
	price int,
	quantity int,
	amount int,
	income int,
	purchasedAt time.Time,
) string {
	return fmt.Sprintf(`
購入が完了しました。
購入者情報：
名前： %v  
メールアドレス： %v 
住所： %v 

--------------------------

商品情報：
商品名： %v 
値段： %d
数量： %d
合計金額： %d

--------------------------
合計売上： %d
購入日時： %v 
	`,
		name,
		email,
		address,
		itemName,
		price,
		quantity,
		amount,
		income,
		convertToJST(purchasedAt),
	)
}

func PurchasedCustomerEmailFactory(
	name string,
	transactionId string,
	itemName string,
	price int,
	quantity int,
	shippingFee int,
	amount int,
	purchasedAt time.Time,
	to string,
	address string,
) string {
	return fmt.Sprintf(`
%v 様 


この度はお買い上げいただき、誠にありがとうございます。
お客様のご注文を確認いたしましたので、ご連絡いたします。
以下に、ご注文の詳細情報を記載いたします。

注文ID： %v

--------------------------

【ご注文情報】
商品名		値段		数量
%v		%v円		%v個
送料： %v円
合計金額： %v 円
購入日時： %v 

--------------------------

【お届け先】
お名前： %v様
住所： %v

--------------------------

商品の発送準備が整いましたら、別途メールにてご連絡いたします。通常、商品の発送には、2,3日程度かかりますので、ご了承ください。

ご質問やご不明な点がございましたら、いつでもお気軽にお客様相談室からお問い合わせください。お手続きや配送に関する詳細情報は、ご注文IDを教えていただくとスムーズに対応できます。

また、商品の受け取り後に何かお気づきの点やご意見がございましたら、お知らせいただけると幸いです。お客様のご意見は、弊社のサービス向上につながりますので、ぜひお聞かせください。

改めまして、ご購入いただきありがとうございます。今後とも、より良い商品とサービスをご提供できるよう努めてまいりますので、どうぞよろしくお願い申し上げます。


--------------------------
CharisWorks

お客様相談室:contact@charis.works

`,
		name,
		transactionId,
		itemName,
		price,
		quantity,
		shippingFee,
		amount,
		convertToJST(purchasedAt),
		to,
		address)
}

func PurchasedWorkerEmailFactory(
	name string,
	itemName string,
	price int,
	quantity int,
	amount int,
	purchasedAt time.Time,
) string {
	return fmt.Sprintf(`
%v 様


出品された商品が購入されました。
以下に、ご注文の詳細情報を記載いたします。

--------------------------

【商品情報】
商品名		値段		数量
%v		%v円		%v個
売上： %v 円
購入日時： %v

--------------------------
CharisWorks

お客様相談室:contact@charis.works
`,
		name,
		itemName,
		price,
		quantity,
		amount,
		convertToJST(purchasedAt),
	)
}

func SendAdminShippedEmail(
	name string,
	mail string,
	address string,
	itemName string,
	price int,
	quantity int,
	amount int,
	income int,
	purchasedAt time.Time,
) string {
	return fmt.Sprintf(`
発送が完了しました。

購入者情報：
名前： %v  
メールアドレス： %v 
住所： %v 

--------------------------

商品情報：
商品名： %v 
値段： %v 
数量： %v 
合計金額： %v 

--------------------------
合計売上： %v 
購入日時： %v `,
		name,
		mail,
		address,
		itemName,
		price,
		quantity,
		amount,
		income,
		convertToJST(purchasedAt),
	)
}
func SendCustomerShippedEmail(
	transactionId string,
	trackingId string,
	itemName string,
	price int,
	quantity int,
	shippingFee int,
	totalAmount int,
	purchasedAt time.Time,
	name string,
	address string,
) string {
	return fmt.Sprintf(`
%v 様 


この度はお買い上げいただき、誠にありがとうございます。
お客様から注文のあった商品を発送しましたので、ご連絡いたします。
以下に、ご注文の詳細情報を記載いたします。

注文ID： %v
追跡番号： %v

--------------------------

【ご注文情報】
商品名		値段		数量
%v		%v円		%v個
			
送料： %v円
合計金額： %v 円
購入日時： %v 
		
--------------------------

【お届け先】
お名前： %v様
住所： %v

--------------------------

商品の返品・返金に致しましては、商品到着後7日以内にお問い合わせフォームよりご連絡ください。商品の状態を確認の上、返品・返金の手続きをさせていただきます。

また、お客様自身の都合による返品・返金については、致しかねる場合がございますので、予めご了承ください。

--------------------------
CharisWorks

お客様相談室:contact@charis.works
お問い合わせフォーム:[link]
`,
		name,
		transactionId,
		trackingId,
		itemName,
		price,
		quantity,
		shippingFee,
		totalAmount,
		convertToJST(purchasedAt),
		name,
		address)
}

func convertToJST(utcTime time.Time) string {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// タイムゾーンの読み込みに失敗した場合のエラーハンドリング
		fmt.Println("タイムゾーンの読み込みに失敗しました:", err)
		return ""
	}
	jstTime := utcTime.In(loc)
	return jstTime.Format("2006-01-02")
}
