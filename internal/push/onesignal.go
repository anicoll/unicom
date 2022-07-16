package push

import (
	"context"
	"errors"

	"github.com/OneSignal/onesignal-go-api"
	"github.com/aws/aws-sdk-go-v2/aws"
	"go.uber.org/zap"
)

type Service struct {
	appId     string
	authKey   string
	apiClient *onesignal.DefaultApiService
	logger    *zap.Logger
}

func New(logger *zap.Logger, appId, authKey string) *Service {
	configuration := onesignal.NewConfiguration()
	return &Service{
		appId:     appId,
		authKey:   authKey,
		logger:    logger,
		apiClient: onesignal.NewAPIClient(configuration).DefaultApi,
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

func (s *Service) Send(ctx context.Context, args Notification) (*string, error) {
	notification := *onesignal.NewNotification(s.appId)
	notification.SetIsIos(true)
	notification.SetIsAndroid(true)
	notification.SetIsHuawei(true)
	notification.SetExternalId(args.IdempotencyKey)
	notification.SetAppId(s.appId)
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

	authCtx := context.WithValue(ctx, onesignal.AppAuth, s.authKey)

	resp, _, err := s.apiClient.CreateNotification(authCtx).Notification(notification).Execute()

	if err != nil {
		s.logger.Error("error sending push notification", zap.Error(err))
		return nil, err
	}
	if resp.Errors != nil {
		if resp.Errors.InvalidIdentifierError != nil {
			s.logger.Error("error sending push notification", zap.Strings("invalidExternalUserIds", resp.Errors.InvalidIdentifierError.InvalidExternalUserIds))
			return nil, err
		}
		s.logger.Error("unknown error sending push notification", zap.Any("errors", resp.Errors))
		return nil, errors.New("unknown error occured attempting to send communication")
	}
	return aws.String(resp.GetId()), nil
}
