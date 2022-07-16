package push

import (
	"context"

	"github.com/OneSignal/onesignal-go-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	appId     string
	authKey   string
	apiClient *onesignal.APIClient
	logger    *zap.Logger
}

func New(logger *zap.Logger, appId, authKey string) *Service {
	configuration := onesignal.NewConfiguration()
	return &Service{
		appId:     appId,
		authKey:   authKey,
		logger:    logger,
		apiClient: onesignal.NewAPIClient(configuration),
	}
}

type LanguageContent struct {
	Arabic  string
	English string
}

type Notification struct {
	IdempotencyKey     string
	ExternalCustomerId string
	Content            LanguageContent
	Heading            LanguageContent
	SubTitle           *LanguageContent
}

func (os *Service) Send(ctx context.Context, args Notification) (*string, error) {
	notification := *onesignal.NewNotification(os.appId)
	notification.SetIsIos(true)
	notification.SetIsAndroid(true)
	notification.SetIsHuawei(true)
	// notification.SetExternalId(args.IdempotencyKey)
	notification.SetExternalId(uuid.NewString())
	notification.SetAppId(os.appId)
	notification.SetIncludeExternalUserIds([]string{args.ExternalCustomerId})
	notification.SetChannelForExternalUserIds("push")

	contents := onesignal.NewStringMap()
	contents.SetAr(args.Content.Arabic)
	contents.SetEn(args.Content.English)
	notification.SetContents(*contents)

	headers := onesignal.NewStringMap()
	headers.SetAr(args.Heading.Arabic)
	headers.SetEn(args.Heading.English)
	notification.SetHeadings(*headers)

	if args.SubTitle != nil {
		subtitles := onesignal.NewStringMap()
		subtitles.SetAr(args.SubTitle.Arabic)
		subtitles.SetEn(args.SubTitle.English)
		notification.SetSubtitle(*subtitles)
	}

	appAuth := context.WithValue(ctx, onesignal.AppAuth, os.authKey)

	resp, _, err := os.apiClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()

	if err != nil {
		os.logger.Error("error sending push notification", zap.Error(err))
		return nil, err
	}
	if resp.Errors != nil {
		if resp.Errors.InvalidIdentifierError != nil {
			os.logger.Error("error sending push notification", zap.Strings("invalidExternalUserIds", resp.Errors.InvalidIdentifierError.InvalidExternalUserIds))
			return nil, err
		}
	}
	return aws.String(resp.GetId()), nil
}
