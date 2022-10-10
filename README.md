# SP1 : GOSCRAP

이 사이드 프로젝트는 사랑스러운 여자친구의 아르바이트 구직 활동을 도와주고 싶은 마음에 구현하게 되었습니다. <br />
원하는 아르바이트의 모집 공고가 언제 올라올지 몰라 항상 직접 들어가서 확인을 해야 했던 번거로움을 해결해주고 싶었습니다. <br />
매일 2번 (오전 10시, 오후 6시) 웹 스크래핑 결과를 문자메세지로 알려주게 된다면 그 번거로움을 해소하고 다른 일에 집중할 수 있게 되지 않을까 생각하게 되어 구현해보았습니다. <br />

---

## 목차

- [아키텍처](#아키텍처)
- [패키지 다운로드 받기](#패키지-다운로드-받기)
- [사용법](#사용법)

## 아키텍처

- 사용한 기술 : aws eventbridge, aws lambda, aws sns, go lang
- 1. EventBridge의 cron 정보(매일 오전 10시, 오후 6시)를 토대로 trigger가 발생합니다.
- 2. 발생한 trigger로 인해 Lambda 함수가 동작합니다.
- 3. 동작한 Lambda 함수는 대상이 되는 SNS에 결과를 보냅니다.
- 4. SNS의 Subscription 정보로 있는 SMS SandBox에 저장되어 있는 전화번호로 메세지를 보냅니다.

## 패키지 다운로드 받기

- 의존성 문제를 해결하기 위해서 아래의 패키지를 다운로드 받습니다.

```bash
# goquery
go get github.com/PuerkitoBio/goquery

# aws
go get github.com/aws/aws-sdk-go-v2/aws
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/sns
go get github.com/aws/aws-lambda-go/lambda
```

- 아마도, go.mod 정보가 있기 때문에 에러가 발생할 가능성은 낮을 것이라 생각이 듭니다.
- 만일, mod에 관한 에러가 출력되면 아래의 명령어를 입력하시면 해결하실 수 있습니다.

```bash
go mod init [YOUR_GOPATH_DIRECTORY]
```

## 사용법

- 아래와 같이 main.go 작성을 합니다.
- **추가로, [] 안에 있는 정보는 각자의 상황에 맞는 term, arn을 기입해주시면 됩니다.**
- **arn의 경우 SMS가 가능한 리전으로 지정해주셔야 합니다. [참고](https://docs.aws.amazon.com/ko_kr/sns/latest/dg/sns-supported-regions-countries.html)**

```go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gyu-beom/sp1/scrapping"
)

type MyEvent struct {
	Time string `json:"time"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, time MyEvent) (string, error) {
	partTime := scrapping.Scrapping([YOUR_PARTTIME_NAME]) // write down part time name

	config, err := config.LoadDefaultConfig(context.TODO())
	scrapping.CheckErr(err)

	client := sns.NewFromConfig(config)

	params := &sns.PublishInput{
		Message:  aws.String(partTime),
		TopicArn: aws.String([YOUR_SNS_TOPIC_ARN]), // write down sns topic arn
	}

	resp, err := client.Publish(context.TODO(), params)
	scrapping.CheckErr(err)

	fmt.Println(resp)

	return fmt.Sprintf("At %s!", time.Time), nil
}
```
