package helicon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

func (h *Helicon) GetTweetDetails(request TweetDetailRequest) (*TweetDetailResponse, error) {
	uri, err := request.GetURL()
	if err != nil {
		return nil, err
	}
	body, err := h.hitApi(*uri)
	if err != nil {
		return nil, err
	}
	var response TweetDetailResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &response, nil
}

const QueryId = "1RFzrZSUoVSgHzVK4MHWlg"

type TweetDetailVariables struct {
	// Unique identifier of the main tweet for which you want to fetch details.
	// Every tweet on X has a unique ID. This is the primary input for the TweetDetail query.
	FocalTweetId string `json:"focalTweetId"` // main tweet id
	// Indicates the context or page from which the user navigated to view the tweet details.
	// For example, `home`: The user clicked on the tweet from their main home timeline.
	//
	// X might use this for analytics, for slightly different UI presentations, or to determine what content to load around the tweet.
	//
	// Default: "home"
	Referrer string `json:"referrer,omitempty"`
	// Default: empty
	ControllerData string `json:"controller_data,omitempty"`
	// "RUX" often stands for "Rich User Experience." "Injections" is dynamically adding or modifying content or UI elements.
	//
	//  This flag might control whether certain enhanced UI components, interactive elements, or third-party integrations
	// (like embedded media players with extra features, polls with real-time updates, or special ad formats)
	//	are loaded or initialized as part of the tweet detail view.
	//
	// Setting it to false might request a more basic or "vanilla" version of the tweet detail, potentially for performance
	// reasons or in contexts where these richer features aren't needed or desired.
	//
	// Default: false
	WithRuxInjections bool `json:"with_rux_injections"`
	// Relevance: Replies are sorted by an algorithm that considers factors like engagement (likes, replies, retweets), recency, and possibly user connections.
	// Recency: Latest, from new to old.
	// Likes: Most likes to least.
	//
	// Default: "Relevance"
	RankingMode string `json:"rankingMode"` // reply ranking mode, relevance by default, can be empty
	// Determines whether promoted tweets (advertisements) should be included in the data returned,
	// particularly in the context of replies or related tweets shown alongside the focal tweet.
	//
	// Default: false
	IncludePromotedContent bool `json:"includePromotedContent"`
	// If the tweet is part of a Community, setting this to true might fetch additional
	//	community-related information (e.g., community name, rules, membership status).
	//
	// Default: true
	WithCommunity bool `json:"withCommunity"`
	// Quick Promote allows users to easily boost their tweets.
	// If true, the API response might include additional fields on the tweet data that indicate whether the
	// tweet is eligible for Quick Promote, or perhaps some metadata related to its promotion status.
	//
	// Default: false
	WithQuickPromoteEligibilityTweetFields bool `json:"withQuickPromoteEligibilityTweetFields"`
	//  "Birdwatch" was the original name for X's "Community Notes" feature, which allows contributors to add context or fact-checks to tweets.
	// Setting this to true will fetch any Community Notes associated with the focal tweet or potentially its replies.
	//
	// Default: false
	WithBirdwatchNotes bool `json:"withBirdwatchNotes"`
	// If true, the API response will include the necessary data to play or display information about any voice
	// recording attached to the tweet. This might include URLs to audio files, durations, transcriptions (if available), etc.
	//
	// Default: false
	WithVoice bool `json:"withVoice"` // uhhh
}

func NewTweetDetailVariables(tweetId string) TweetDetailVariables {
	return TweetDetailVariables{
		FocalTweetId:                           tweetId,
		Referrer:                               "",
		ControllerData:                         "",
		WithRuxInjections:                      false,
		RankingMode:                            "Relevance",
		IncludePromotedContent:                 true,
		WithCommunity:                          true,
		WithQuickPromoteEligibilityTweetFields: true,
		WithBirdwatchNotes:                     true,
		WithVoice:                              true,
	}
}

type TweetDetailFeatures struct {
	RwebVideoScreenEnabled                                         bool `json:"rweb_video_screen_enabled"`
	ProfileLabelImprovementsPcfLabelInPostEnabled                  bool `json:"profile_label_improvements_pcf_label_in_post_enabled"`
	RwebTipjarConsumptionEnabled                                   bool `json:"rweb_tipjar_consumption_enabled"`
	VerifiedPhoneLabelEnabled                                      bool `json:"verified_phone_label_enabled"`
	CreatorSubscriptionsTweetPreviewApiEnabled                     bool `json:"creator_subscriptions_tweet_preview_api_enabled"`
	ResponsiveWebGraphqlTimelineNavigationEnabled                  bool `json:"responsive_web_graphql_timeline_navigation_enabled"`
	ResponsiveWebGraphqlSkipUserProfileImageExtensionsEnabled      bool `json:"responsive_web_graphql_skip_user_profile_image_extensions_enabled"`
	PremiumContentApiReadEnabled                                   bool `json:"premium_content_api_read_enabled"`
	CommunitiesWebEnableTweetCommunityResultsFetch                 bool `json:"communities_web_enable_tweet_community_results_fetch"`
	C9STweetAnatomyModeratorBadgeEnabled                           bool `json:"c9s_tweet_anatomy_moderator_badge_enabled"`
	ResponsiveWebGrokAnalyzeButtonFetchTrendsEnabled               bool `json:"responsive_web_grok_analyze_button_fetch_trends_enabled"`
	ResponsiveWebGrokAnalyzePostFollowupsEnabled                   bool `json:"responsive_web_grok_analyze_post_followups_enabled"`
	ResponsiveWebJetfuelFrame                                      bool `json:"responsive_web_jetfuel_frame"`
	ResponsiveWebGrokShareAttachmentEnabled                        bool `json:"responsive_web_grok_share_attachment_enabled"`
	ArticlesPreviewEnabled                                         bool `json:"articles_preview_enabled"`
	ResponsiveWebEditTweetApiEnabled                               bool `json:"responsive_web_edit_tweet_api_enabled"`
	GraphqlIsTranslatableRwebTweetIsTranslatableEnabled            bool `json:"graphql_is_translatable_rweb_tweet_is_translatable_enabled"`
	ViewCountsEverywhereApiEnabled                                 bool `json:"view_counts_everywhere_api_enabled"`
	LongformNotetweetsConsumptionEnabled                           bool `json:"longform_notetweets_consumption_enabled"`
	ResponsiveWebTwitterArticleTweetConsumptionEnabled             bool `json:"responsive_web_twitter_article_tweet_consumption_enabled"`
	TweetAwardsWebTippingEnabled                                   bool `json:"tweet_awards_web_tipping_enabled"`
	ResponsiveWebGrokShowGrokTranslatedPost                        bool `json:"responsive_web_grok_show_grok_translated_post"`
	ResponsiveWebGrokAnalysisButtonFromBackend                     bool `json:"responsive_web_grok_analysis_button_from_backend"`
	CreatorSubscriptionsQuoteTweetPreviewEnabled                   bool `json:"creator_subscriptions_quote_tweet_preview_enabled"`
	FreedomOfSpeechNotReachFetchEnabled                            bool `json:"freedom_of_speech_not_reach_fetch_enabled"`
	StandardizedNudgesMisinfo                                      bool `json:"standardized_nudges_misinfo"`
	TweetWithVisibilityResultsPreferGqlLimitedActionsPolicyEnabled bool `json:"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled"`
	LongformNotetweetsRichTextReadEnabled                          bool `json:"longform_notetweets_rich_text_read_enabled"`
	LongformNotetweetsInlineMediaEnabled                           bool `json:"longform_notetweets_inline_media_enabled"`
	ResponsiveWebGrokImageAnnotationEnabled                        bool `json:"responsive_web_grok_image_annotation_enabled"`
	ResponsiveWebEnhanceCardsEnabled                               bool `json:"responsive_web_enhance_cards_enabled"`
}

// NewTweetDetailFeatures returns new [TweetDetailFeatures] All true by default, please adjust the features you dont want.
func NewTweetDetailFeatures() TweetDetailFeatures {
	return TweetDetailFeatures{
		RwebVideoScreenEnabled:                                         true,
		ProfileLabelImprovementsPcfLabelInPostEnabled:                  true,
		RwebTipjarConsumptionEnabled:                                   true,
		VerifiedPhoneLabelEnabled:                                      true,
		CreatorSubscriptionsTweetPreviewApiEnabled:                     true,
		ResponsiveWebGraphqlTimelineNavigationEnabled:                  true,
		ResponsiveWebGraphqlSkipUserProfileImageExtensionsEnabled:      true,
		PremiumContentApiReadEnabled:                                   true,
		CommunitiesWebEnableTweetCommunityResultsFetch:                 true,
		C9STweetAnatomyModeratorBadgeEnabled:                           true,
		ResponsiveWebGrokAnalyzeButtonFetchTrendsEnabled:               true,
		ResponsiveWebGrokAnalyzePostFollowupsEnabled:                   true,
		ResponsiveWebJetfuelFrame:                                      true,
		ResponsiveWebGrokShareAttachmentEnabled:                        true,
		ArticlesPreviewEnabled:                                         true,
		ResponsiveWebEditTweetApiEnabled:                               true,
		GraphqlIsTranslatableRwebTweetIsTranslatableEnabled:            true,
		ViewCountsEverywhereApiEnabled:                                 true,
		LongformNotetweetsConsumptionEnabled:                           true,
		ResponsiveWebTwitterArticleTweetConsumptionEnabled:             true,
		TweetAwardsWebTippingEnabled:                                   true,
		ResponsiveWebGrokShowGrokTranslatedPost:                        true,
		ResponsiveWebGrokAnalysisButtonFromBackend:                     true,
		CreatorSubscriptionsQuoteTweetPreviewEnabled:                   true,
		FreedomOfSpeechNotReachFetchEnabled:                            true,
		StandardizedNudgesMisinfo:                                      true,
		TweetWithVisibilityResultsPreferGqlLimitedActionsPolicyEnabled: true,
		LongformNotetweetsRichTextReadEnabled:                          true,
		LongformNotetweetsInlineMediaEnabled:                           true,
		ResponsiveWebGrokImageAnnotationEnabled:                        true,
		ResponsiveWebEnhanceCardsEnabled:                               true,
	}
}

type TweetDetailFieldToggles struct {
	WithArticleRichContentState bool `json:"withArticleRichContentState"`
	WithArticlePlainText        bool `json:"withArticlePlainText"`
	WithGrokAnalyze             bool `json:"withGrokAnalyze"`
	WithDisallowedReplyControls bool `json:"withDisallowedReplyControls"`
}

// NewTweetDetailFieldToggles returns [TweetDetailFieldToggles], by default, everything is false.
func NewTweetDetailFieldToggles() TweetDetailFieldToggles {
	return TweetDetailFieldToggles{
		WithArticleRichContentState: true,
		WithArticlePlainText:        false,
		WithGrokAnalyze:             false,
		WithDisallowedReplyControls: false}
}

type TweetDetailRequest struct {
	Variables    TweetDetailVariables
	Features     TweetDetailFeatures
	FieldToggles TweetDetailFieldToggles
}

func NewTweetDetailRequest(variables TweetDetailVariables, features TweetDetailFeatures, fieldToggles TweetDetailFieldToggles) *TweetDetailRequest {
	return &TweetDetailRequest{Variables: variables, Features: features, FieldToggles: fieldToggles}
}

func (r TweetDetailRequest) GetURL() (*string, error) {
	variablesJSON, err := json.Marshal(r.Variables)
	if err != nil {
		return nil, fmt.Errorf("error marshalling variables: %w", err)
	}
	encodedVariables := url.QueryEscape(string(variablesJSON))

	featuresJSON, err := json.Marshal(r.Features)
	if err != nil {
		return nil, fmt.Errorf("error marshalling features: %w", err)
	}
	encodedFeatures := url.QueryEscape(string(featuresJSON))

	fieldTogglesJSON, err := json.Marshal(r.FieldToggles)
	if err != nil {
		return nil, fmt.Errorf("error marshalling fieldToggles: %w", err)
	}
	encodedFieldToggles := url.QueryEscape(string(fieldTogglesJSON))

	baseURL := fmt.Sprintf("https://x.com/i/api/graphql/%s/TweetDetail", QueryId)
	fullURL := fmt.Sprintf("%s?variables=%s&features=%s&fieldToggles=%s",
		baseURL, encodedVariables, encodedFeatures, encodedFieldToggles)
	return &fullURL, nil
}

// TweetDetailResponse example => https://paste.rs/GSkCG.json
type TweetDetailResponse struct {
	Data struct {
		ThreadedConversationWithInjectionsV2 struct {
			Instructions []struct {
				Type    string `json:"type"`
				Entries []struct {
					EntryId   string `json:"entryId"`
					SortIndex string `json:"sortIndex"`
					Content   struct {
						EntryType   string `json:"entryType"`
						Typename    string `json:"__typename"`
						ItemContent struct {
							ItemType     string `json:"itemType"`
							Typename     string `json:"__typename"`
							TweetResults struct {
								Result struct {
									Typename string `json:"__typename"`
									Tweet    struct {
										RestId            string `json:"rest_id"`
										HasBirdwatchNotes bool   `json:"has_birdwatch_notes"`
										Core              struct {
											UserResults struct {
												Result struct {
													Typename                   string `json:"__typename"`
													Id                         string `json:"id"`
													RestId                     string `json:"rest_id"`
													AffiliatesHighlightedLabel struct {
													} `json:"affiliates_highlighted_label"`
													HasGraduatedAccess bool `json:"has_graduated_access"`
													IsBlueVerified     bool `json:"is_blue_verified"`
													Legacy             struct {
														CanDm               bool   `json:"can_dm"`
														CanMediaTag         bool   `json:"can_media_tag"`
														CreatedAt           string `json:"created_at"`
														DefaultProfile      bool   `json:"default_profile"`
														DefaultProfileImage bool   `json:"default_profile_image"`
														Description         string `json:"description"`
														Entities            struct {
															Description struct {
																Urls []interface{} `json:"urls"`
															} `json:"description"`
															Url struct {
																Urls []struct {
																	DisplayUrl  string `json:"display_url"`
																	ExpandedUrl string `json:"expanded_url"`
																	Url         string `json:"url"`
																	Indices     []int  `json:"indices"`
																} `json:"urls"`
															} `json:"url"`
														} `json:"entities"`
														FastFollowersCount      int           `json:"fast_followers_count"`
														FavouritesCount         int           `json:"favourites_count"`
														FollowersCount          int           `json:"followers_count"`
														Following               bool          `json:"following"`
														FriendsCount            int           `json:"friends_count"`
														HasCustomTimelines      bool          `json:"has_custom_timelines"`
														IsTranslator            bool          `json:"is_translator"`
														ListedCount             int           `json:"listed_count"`
														MediaCount              int           `json:"media_count"`
														Name                    string        `json:"name"`
														NormalFollowersCount    int           `json:"normal_followers_count"`
														PinnedTweetIdsStr       []string      `json:"pinned_tweet_ids_str"`
														PossiblySensitive       bool          `json:"possibly_sensitive"`
														ProfileBannerUrl        string        `json:"profile_banner_url"`
														ProfileImageUrlHttps    string        `json:"profile_image_url_https"`
														ProfileInterstitialType string        `json:"profile_interstitial_type"`
														ScreenName              string        `json:"screen_name"`
														StatusesCount           int           `json:"statuses_count"`
														TranslatorType          string        `json:"translator_type"`
														Url                     string        `json:"url"`
														WantRetweets            bool          `json:"want_retweets"`
														WithheldInCountries     []interface{} `json:"withheld_in_countries"`
													} `json:"legacy"`
													Location struct {
														Location string `json:"location"`
													} `json:"location"`
													ParodyCommentaryFanLabel string `json:"parody_commentary_fan_label"`
													ProfileImageShape        string `json:"profile_image_shape"`
													Privacy                  struct {
														Protected bool `json:"protected"`
													} `json:"privacy"`
													TipjarSettings struct {
													} `json:"tipjar_settings"`
													Verification struct {
														Verified bool `json:"verified"`
													} `json:"verification"`
													VerifiedPhoneStatus bool `json:"verified_phone_status"`
												} `json:"result"`
											} `json:"user_results"`
										} `json:"core"`
										UnmentionData struct {
										} `json:"unmention_data"`
										EditControl struct {
											EditTweetIds       []string `json:"edit_tweet_ids"`
											EditableUntilMsecs string   `json:"editable_until_msecs"`
											IsEditEligible     bool     `json:"is_edit_eligible"`
											EditsRemaining     string   `json:"edits_remaining"`
										} `json:"edit_control"`
										IsTranslatable bool `json:"is_translatable"`
										Views          struct {
											Count string `json:"count"`
											State string `json:"state"`
										} `json:"views"`
										Source        string `json:"source"`
										AwardEligible bool   `json:"award_eligible"`
										GrantedAwards struct {
										} `json:"granted_awards"`
										GrokAnalysisButton bool `json:"grok_analysis_button"`
										Legacy             struct {
											BookmarkCount       int    `json:"bookmark_count"`
											Bookmarked          bool   `json:"bookmarked"`
											CreatedAt           string `json:"created_at"`
											ConversationControl struct {
												Policy                   string `json:"policy"`
												ConversationOwnerResults struct {
													Result struct {
														Typename string `json:"__typename"`
														Legacy   struct {
															ScreenName string `json:"screen_name"`
														} `json:"legacy"`
													} `json:"result"`
												} `json:"conversation_owner_results"`
											} `json:"conversation_control"`
											ConversationIdStr string `json:"conversation_id_str"`
											DisplayTextRange  []int  `json:"display_text_range"`
											Entities          struct {
												Hashtags     []interface{} `json:"hashtags"`
												Symbols      []interface{} `json:"symbols"`
												Timestamps   []interface{} `json:"timestamps"`
												Urls         []interface{} `json:"urls"`
												UserMentions []interface{} `json:"user_mentions"`
											} `json:"entities"`
											FavoriteCount int    `json:"favorite_count"`
											Favorited     bool   `json:"favorited"`
											FullText      string `json:"full_text"`
											IsQuoteStatus bool   `json:"is_quote_status"`
											Lang          string `json:"lang"`
											QuoteCount    int    `json:"quote_count"`
											ReplyCount    int    `json:"reply_count"`
											RetweetCount  int    `json:"retweet_count"`
											Retweeted     bool   `json:"retweeted"`
											UserIdStr     string `json:"user_id_str"`
											IdStr         string `json:"id_str"`
										} `json:"legacy"`
										QuickPromoteEligibility struct {
											Eligibility string `json:"eligibility"`
										} `json:"quick_promote_eligibility"`
									} `json:"tweet"`
									LimitedActionResults struct {
										LimitedActions []struct {
											Action string `json:"action"`
											Prompt struct {
												Typename string `json:"__typename"`
												CtaType  string `json:"cta_type"`
												Headline struct {
													Text     string        `json:"text"`
													Entities []interface{} `json:"entities"`
												} `json:"headline"`
												Subtext struct {
													Text     string        `json:"text"`
													Entities []interface{} `json:"entities"`
												} `json:"subtext"`
											} `json:"prompt"`
										} `json:"limited_actions"`
									} `json:"limitedActionResults"`
								} `json:"result"`
							} `json:"tweet_results"`
							TweetDisplayType    string `json:"tweetDisplayType"`
							HasModeratedReplies bool   `json:"hasModeratedReplies"`
						} `json:"itemContent,omitempty"`
						ClientEventInfo struct {
							Component string `json:"component"`
							Element   string `json:"element"`
							Details   struct {
								ConversationDetails struct {
									ConversationSection string `json:"conversationSection"`
								} `json:"conversationDetails"`
							} `json:"details,omitempty"`
						} `json:"clientEventInfo"`
						Items []struct {
							EntryId string `json:"entryId"`
							Item    struct {
								ItemContent struct {
									ItemType     string `json:"itemType"`
									Typename     string `json:"__typename"`
									TweetResults struct {
										Result struct {
											Typename          string `json:"__typename"`
											RestId            string `json:"rest_id"`
											HasBirdwatchNotes bool   `json:"has_birdwatch_notes"`
											Core              struct {
												UserResults struct {
													Result struct {
														Typename                   string `json:"__typename"`
														Id                         string `json:"id"`
														RestId                     string `json:"rest_id"`
														AffiliatesHighlightedLabel struct {
														} `json:"affiliates_highlighted_label"`
														HasGraduatedAccess bool `json:"has_graduated_access"`
														IsBlueVerified     bool `json:"is_blue_verified"`
														Legacy             struct {
															CanDm               bool   `json:"can_dm"`
															CanMediaTag         bool   `json:"can_media_tag"`
															CreatedAt           string `json:"created_at"`
															DefaultProfile      bool   `json:"default_profile"`
															DefaultProfileImage bool   `json:"default_profile_image"`
															Description         string `json:"description"`
															Entities            struct {
																Description struct {
																	Urls []interface{} `json:"urls"`
																} `json:"description"`
																Url struct {
																	Urls []struct {
																		DisplayUrl  string `json:"display_url"`
																		ExpandedUrl string `json:"expanded_url"`
																		Url         string `json:"url"`
																		Indices     []int  `json:"indices"`
																	} `json:"urls"`
																} `json:"url"`
															} `json:"entities"`
															FastFollowersCount      int           `json:"fast_followers_count"`
															FavouritesCount         int           `json:"favourites_count"`
															FollowersCount          int           `json:"followers_count"`
															Following               bool          `json:"following"`
															FriendsCount            int           `json:"friends_count"`
															HasCustomTimelines      bool          `json:"has_custom_timelines"`
															IsTranslator            bool          `json:"is_translator"`
															ListedCount             int           `json:"listed_count"`
															MediaCount              int           `json:"media_count"`
															Name                    string        `json:"name"`
															NormalFollowersCount    int           `json:"normal_followers_count"`
															PinnedTweetIdsStr       []interface{} `json:"pinned_tweet_ids_str"`
															PossiblySensitive       bool          `json:"possibly_sensitive"`
															ProfileBannerUrl        string        `json:"profile_banner_url"`
															ProfileImageUrlHttps    string        `json:"profile_image_url_https"`
															ProfileInterstitialType string        `json:"profile_interstitial_type"`
															ScreenName              string        `json:"screen_name"`
															StatusesCount           int           `json:"statuses_count"`
															TranslatorType          string        `json:"translator_type"`
															Url                     string        `json:"url"`
															WantRetweets            bool          `json:"want_retweets"`
															WithheldInCountries     []interface{} `json:"withheld_in_countries"`
														} `json:"legacy"`
														Location struct {
															Location string `json:"location"`
														} `json:"location"`
														ParodyCommentaryFanLabel string `json:"parody_commentary_fan_label"`
														ProfileImageShape        string `json:"profile_image_shape"`
														Privacy                  struct {
															Protected bool `json:"protected"`
														} `json:"privacy"`
														TipjarSettings struct {
															IsEnabled     bool   `json:"is_enabled"`
															PatreonHandle string `json:"patreon_handle,omitempty"`
															CashAppHandle string `json:"cash_app_handle,omitempty"`
															VenmoHandle   string `json:"venmo_handle,omitempty"`
														} `json:"tipjar_settings"`
														Verification struct {
															Verified bool `json:"verified"`
														} `json:"verification"`
														VerifiedPhoneStatus bool `json:"verified_phone_status"`
													} `json:"result"`
												} `json:"user_results"`
											} `json:"core"`
											UnmentionData struct {
											} `json:"unmention_data"`
											EditControl struct {
												EditTweetIds       []string `json:"edit_tweet_ids"`
												EditableUntilMsecs string   `json:"editable_until_msecs"`
												IsEditEligible     bool     `json:"is_edit_eligible"`
												EditsRemaining     string   `json:"edits_remaining"`
											} `json:"edit_control"`
											IsTranslatable bool `json:"is_translatable"`
											Views          struct {
												Count string `json:"count"`
												State string `json:"state"`
											} `json:"views"`
											Source        string `json:"source"`
											AwardEligible bool   `json:"award_eligible"`
											GrantedAwards struct {
											} `json:"granted_awards"`
											GrokAnalysisButton bool `json:"grok_analysis_button"`
											Legacy             struct {
												BookmarkCount       int    `json:"bookmark_count"`
												Bookmarked          bool   `json:"bookmarked"`
												CreatedAt           string `json:"created_at"`
												ConversationControl struct {
													Policy                   string `json:"policy"`
													ConversationOwnerResults struct {
														Result struct {
															Typename string `json:"__typename"`
															Legacy   struct {
																ScreenName string `json:"screen_name"`
															} `json:"legacy"`
														} `json:"result"`
													} `json:"conversation_owner_results"`
												} `json:"conversation_control"`
												ConversationIdStr string `json:"conversation_id_str"`
												DisplayTextRange  []int  `json:"display_text_range"`
												Entities          struct {
													Hashtags     []interface{} `json:"hashtags"`
													Symbols      []interface{} `json:"symbols"`
													Timestamps   []interface{} `json:"timestamps"`
													Urls         []interface{} `json:"urls"`
													UserMentions []struct {
														IdStr      string `json:"id_str"`
														Name       string `json:"name"`
														ScreenName string `json:"screen_name"`
														Indices    []int  `json:"indices"`
													} `json:"user_mentions"`
												} `json:"entities"`
												FavoriteCount        int    `json:"favorite_count"`
												Favorited            bool   `json:"favorited"`
												FullText             string `json:"full_text"`
												InReplyToScreenName  string `json:"in_reply_to_screen_name"`
												InReplyToStatusIdStr string `json:"in_reply_to_status_id_str"`
												InReplyToUserIdStr   string `json:"in_reply_to_user_id_str"`
												IsQuoteStatus        bool   `json:"is_quote_status"`
												Lang                 string `json:"lang"`
												QuoteCount           int    `json:"quote_count"`
												ReplyCount           int    `json:"reply_count"`
												RetweetCount         int    `json:"retweet_count"`
												Retweeted            bool   `json:"retweeted"`
												UserIdStr            string `json:"user_id_str"`
												IdStr                string `json:"id_str"`
											} `json:"legacy"`
											QuickPromoteEligibility struct {
												Eligibility string `json:"eligibility"`
											} `json:"quick_promote_eligibility"`
										} `json:"result"`
									} `json:"tweet_results"`
									TweetDisplayType string `json:"tweetDisplayType"`
									SocialContext    struct {
										Type        string `json:"type"`
										ContextType string `json:"contextType"`
										Text        string `json:"text"`
									} `json:"socialContext"`
								} `json:"itemContent"`
								ClientEventInfo struct {
									Component string `json:"component"`
									Element   string `json:"element"`
									Details   struct {
										ConversationDetails struct {
											ConversationSection string `json:"conversationSection"`
										} `json:"conversationDetails"`
										TimelinesDetails struct {
											ControllerData string `json:"controllerData"`
										} `json:"timelinesDetails"`
									} `json:"details"`
								} `json:"clientEventInfo"`
							} `json:"item"`
						} `json:"items,omitempty"`
						DisplayType string `json:"displayType,omitempty"`
					} `json:"content"`
				} `json:"entries,omitempty"`
				Direction string `json:"direction,omitempty"`
			} `json:"instructions"`
		} `json:"threaded_conversation_with_injections_v2"`
	} `json:"data"`
}
