package upload

import (
	"encoding/json"
	"errors"
	"github.com/meandrewdev/logger"
	"strings"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

func (vu *VideoUploader) updateMetadata(metadataParams VideoMetadataParams) (*req.Resp, error) {
	request := req.New()
	request.EnableCookie(false)
	if vu.ytBrowser.Account.Proxy != "" {
		request.SetProxyUrl(vu.ytBrowser.Account.Proxy)
	}

	resp, err := request.Post("https://studio.youtube.com/youtubei/v1/video_manager/metadata_update?"+metadataParams.Query, metadataParams.Headers, req.BodyJSON(metadataParams.Metadata))
	if err != nil {
		return nil, err
	}

	body := resp.String()
	jbody := gjson.Parse(body)

	if jbody.Get("error").Exists() {
		return nil, errors.New("Ошибка при изменении видео: " + jbody.Get("error.message").String())
	}

	result := jbody.Get("overallResult.resultCode").String()

	if result == "" {
		d, _ := json.Marshal(metadataParams.Metadata)
		logger.Warning("unknown answer from metadata updater\nbody:\r\n" + string(d) + "answer:\r\n" + body)
		return nil, errors.New("unknown error")
	}

	if result == "SOME_ERRORS" {
		fields := []string{}
		jbody.ForEach(func(key, value gjson.Result) bool {
			success := value.Get("success")
			if success.Exists() && !success.Bool() {
				fields = append(fields, key.String())
			}
			return true
		})
		d, _ := json.Marshal(metadataParams.Metadata)
		return nil, errors.New("error when changing fields: " + strings.Join(fields, ", ") + "body: " + string(d))
	}

	return resp, nil
}

func GetVideoReadMask() VideoReadMask {
	all := AllMask{
		All: true,
	}
	return VideoReadMask{
		VideoID:                    true,
		Title:                      true,
		TitleFormattedString:       all,
		Description:                true,
		DescriptionFormattedString: all,
		ThumbnailDetails:           all,
		Status:                     true,
		StatusDetails:              all,
		DraftStatus:                true,
		DateRecorded:               all,
		Location:                   all,
		Category:                   true,
		GameTitle:                  all,
		AudioLanguage:              all,
		CrowdsourcingEnabled:       true,
		UncaptionedReason:          true,
		OwnedClaimDetails:          all,
		Tags:                       all,
		AllowEmbed:                 true,
		Publishing:                 all,
		PaidProductPlacement:       true,
		SerializedShareEntity:      true,
		License:                    true,
		AllowComments:              true,
		CommentFilter:              true,
		DefaultCommentSortOrder:    true,
		AllowRatings:               true,
		ClaimDetails:               all,
		CommentsDisabledInternally: true,
		Music:                      all,
		Features:                   all,
		AudienceRestriction:        all,
		Livestream:                 all,
		Origin:                     true,
		Premiere:                   all,
		ThumbnailEditorState:       all,
		Permissions:                all,
		ChannelID:                  true,
		OriginalFilename:           true,
		VideoStreamURL:             true,
		VideoResolutions:           all,
		InlineEditProcessingStatus: true,
		SelfCertification:          all,
		Monetization:               all,
		ResponseStatus:             all,
		AdSettings:                 all,
		MonetizedStatus:            true,
		VideoDurationMs:            true,
		VideoEditorProject:         all,
		Privacy:                    true,
		ScheduledPublishingDetails: all,
		TimePublishedSeconds:       true,
		Visibility:                 all,
		PrivateShare:               all,
		SponsorsOnly:               all,
		UnlistedExpired:            true,
		Metrics:                    all,
		TimeCreatedSeconds:         true,
		LengthSeconds:              true,
		MetadataLanguage:           all,
	}
}

type VideoMetadataParams struct {
	Metadata VideoMetadata
	Query    string
	Headers  req.Header
}

type VideoMetadata struct {
	Context           interface{}       `json:"context,omitempty"`
	DelegationContext DelegationContext `json:"delegationContext,omitempty"`
	EncryptedVideoID  string            `json:"encryptedVideoId,omitempty"`

	VideoReadMask VideoReadMask `json:"videoReadMask,omitempty"`

	Title               *Title               `json:"title,omitempty"`
	Description         *Description         `json:"description,omitempty"`
	Tags                *Tags                `json:"tags,omitempty"`
	CommentOptions      *CommentOptions      `json:"commentOptions,omitempty"`
	TargetedAudience    *TargetedAudience    `json:"targetedAudience,omitempty"`
	AudioLanguage       *AudioLanguage       `json:"audioLanguage,omitempty"`
	PublishingOptions   *PublishingOptions   `json:"publishingOptions,omitempty"`
	DraftState          *DraftState          `json:"draftState,omitempty"`
	PrivacyState        *PrivacyState        `json:"privacyState,omitempty"`
	Category            *Category            `json:"category,omitempty"`
	ScheduledPublishing *ScheduledPublishing `json:"scheduledPublishing,omitempty"`
	VideoStill          *VideoStill          `json:"videoStill,omitempty"`
}
type AllMask struct {
	All bool `json:"all,omitempty"`
}
type VideoReadMask struct {
	VideoID                    bool    `json:"videoId,omitempty"`
	Title                      bool    `json:"title,omitempty"`
	TitleFormattedString       AllMask `json:"titleFormattedString,omitempty"`
	Description                bool    `json:"description,omitempty"`
	DescriptionFormattedString AllMask `json:"descriptionFormattedString,omitempty"`
	ThumbnailDetails           AllMask `json:"thumbnailDetails,omitempty"`
	Status                     bool    `json:"status,omitempty"`
	StatusDetails              AllMask `json:"statusDetails,omitempty"`
	DraftStatus                bool    `json:"draftStatus,omitempty"`
	DateRecorded               AllMask `json:"dateRecorded,omitempty"`
	Location                   AllMask `json:"location,omitempty"`
	Category                   bool    `json:"category,omitempty"`
	GameTitle                  AllMask `json:"gameTitle,omitempty"`
	AudioLanguage              AllMask `json:"audioLanguage,omitempty"`
	CrowdsourcingEnabled       bool    `json:"crowdsourcingEnabled,omitempty"`
	UncaptionedReason          bool    `json:"uncaptionedReason,omitempty"`
	OwnedClaimDetails          AllMask `json:"ownedClaimDetails,omitempty"`
	Tags                       AllMask `json:"tags,omitempty"`
	AllowEmbed                 bool    `json:"allowEmbed,omitempty"`
	Publishing                 AllMask `json:"publishing,omitempty"`
	PaidProductPlacement       bool    `json:"paidProductPlacement,omitempty"`
	SerializedShareEntity      bool    `json:"serializedShareEntity,omitempty"`
	License                    bool    `json:"license,omitempty"`
	AllowComments              bool    `json:"allowComments,omitempty"`
	CommentFilter              bool    `json:"commentFilter,omitempty"`
	DefaultCommentSortOrder    bool    `json:"defaultCommentSortOrder,omitempty"`
	AllowRatings               bool    `json:"allowRatings,omitempty"`
	ClaimDetails               AllMask `json:"claimDetails,omitempty"`
	CommentsDisabledInternally bool    `json:"commentsDisabledInternally,omitempty"`
	Music                      AllMask `json:"music,omitempty"`
	Features                   AllMask `json:"features,omitempty"`
	AudienceRestriction        AllMask `json:"audienceRestriction,omitempty"`
	Livestream                 AllMask `json:"livestream,omitempty"`
	Origin                     bool    `json:"origin,omitempty"`
	Premiere                   AllMask `json:"premiere,omitempty"`
	ThumbnailEditorState       AllMask `json:"thumbnailEditorState,omitempty"`
	Permissions                AllMask `json:"permissions,omitempty"`
	ChannelID                  bool    `json:"channelId,omitempty"`
	OriginalFilename           bool    `json:"originalFilename,omitempty"`
	VideoStreamURL             bool    `json:"videoStreamUrl,omitempty"`
	VideoResolutions           AllMask `json:"videoResolutions,omitempty"`
	InlineEditProcessingStatus bool    `json:"inlineEditProcessingStatus,omitempty"`
	SelfCertification          AllMask `json:"selfCertification,omitempty"`
	Monetization               AllMask `json:"monetization,omitempty"`
	ResponseStatus             AllMask `json:"responseStatus,omitempty"`
	AdSettings                 AllMask `json:"adSettings,omitempty"`
	MonetizedStatus            bool    `json:"monetizedStatus,omitempty"`
	VideoDurationMs            bool    `json:"videoDurationMs,omitempty"`
	VideoEditorProject         AllMask `json:"videoEditorProject,omitempty"`
	Privacy                    bool    `json:"privacy,omitempty"`
	ScheduledPublishingDetails AllMask `json:"scheduledPublishingDetails,omitempty"`
	TimePublishedSeconds       bool    `json:"timePublishedSeconds,omitempty"`
	Visibility                 AllMask `json:"visibility,omitempty"`
	PrivateShare               AllMask `json:"privateShare,omitempty"`
	SponsorsOnly               AllMask `json:"sponsorsOnly,omitempty"`
	UnlistedExpired            bool    `json:"unlistedExpired,omitempty"`
	Metrics                    AllMask `json:"metrics,omitempty"`
	TimeCreatedSeconds         bool    `json:"timeCreatedSeconds,omitempty"`
	LengthSeconds              bool    `json:"lengthSeconds,omitempty"`
	MetadataLanguage           AllMask `json:"metadataLanguage,omitempty"`
}
type Title struct {
	NewTitle      string `json:"newTitle,omitempty"`
	ShouldSegment bool   `json:"shouldSegment,omitempty"`
}
type Description struct {
	NewDescription string `json:"newDescription,omitempty"`
	ShouldSegment  bool   `json:"shouldSegment,omitempty"`
}
type Tags struct {
	NewTags []string `json:"newTags,omitempty"`
}
type CommentOptions struct {
	NewAllowComments     bool   `json:"newAllowComments"`
	NewAllowCommentsMode string `json:"newAllowCommentsMode,omitempty"`
	NewCanViewRatings    bool   `json:"newCanViewRatings"`
	NewDefaultSortOrder  string `json:"newDefaultSortOrder,omitempty"`
}
type TargetedAudience struct {
	Operation           string `json:"operation,omitempty"`
	NewTargetedAudience string `json:"newTargetedAudience,omitempty"`
}
type AudioLanguage struct {
	NewAudioLanguage string `json:"newAudioLanguage,omitempty"`
}
type PublishingOptions struct {
	NewPostToFeed bool `json:"newPostToFeed"`
}

type RoleType struct {
	ChannelRoleType string `json:"channelRoleType,omitempty"`
}
type DelegationContext struct {
	ExternalChannelID string   `json:"externalChannelId,omitempty"`
	RoleType          RoleType `json:"roleType,omitempty"`
}
type DraftState struct {
	Operation string `json:"operation,omitempty"`
}
type PrivacyState struct {
	NewPrivacy string `json:"newPrivacy,omitempty"`
}
type Category struct {
	NewCategoryID int `json:"newCategoryId,omitempty"`
}
type ScheduledPublishing struct {
	Set Set `json:"set,omitempty"`
}
type Set struct {
	TimeSec string `json:"timeSec,omitempty"`
	Privacy string `json:"privacy,omitempty"`
}
type Image struct {
	DataURI string `json:"dataUri,omitempty"`
}
type VideoStill struct {
	NewStillID int    `json:"newStillId,omitempty"`
	Operation  string `json:"operation,omitempty"`
	Image      *Image `json:"image,omitempty"`
}
