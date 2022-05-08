module notionboy

go 1.16

require (
	github.com/argoproj/pkg v0.9.0
	github.com/aws/aws-lambda-go v1.26.0
	github.com/gin-gonic/gin v1.7.2
	github.com/jomei/notionapi v1.8.2
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	github.com/silenceper/wechat/v2 v2.0.6
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/text v0.3.6 // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.22.4
)

require gorm.io/driver/mysql v1.2.1

require golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
