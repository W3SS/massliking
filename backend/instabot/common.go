package instabot

import (
	"net/url"

	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	"github.com/satori/go.uuid"
)

const (
	API_URL = "https://i.instagram.com/api/v1/"

	//MANUFACTURER    = "Xiaomi"
	//MODEL           = "HM 1SW"
	//ANDROID_VERSION = "18"
	//ANDROID_RELEASE = "4.3"
	//RESOLUTION      = "720x1280"
	//DPI             = "320dpi"
	//DEVICE          = "armani"
	//HARDWARE        = "qcom"
	//LOCALE          = "en_US"
	//VERSION         = "9.2.0"

	//TODO: Find out the way to configure this stuff
	ANDROID_VERSION = "19"    //Android.os.Build.VERSION.SDK_INT
	ANDROID_RELEASE = "4.4.2" //Android.os.Build.VERSION.RELEASE
	DPI             = "240dpi"
	RESOLUTION      = "480x800"
	MANUFACTURER    = "FLY"    //Android.os.Build.MANUFACTURER
	MODEL           = "IQ4418" //Android.os.Build.MODEL
	DEVICE          = "IQ4418" //Android.os.Build.DEVICE
	HARDWARE        = "mt6572" //Android.os.Build.HARDWARE
	LOCALE          = "ru_RU"
	VERSION         = "9.2.0" // I'm not sure which value I should pick

	USER_AGENT      = "Instagram " + VERSION + " Android (" + ANDROID_VERSION + "/" + ANDROID_RELEASE + "; " + DPI + "; " + RESOLUTION + "; " + MANUFACTURER + "; " + MODEL + "; " + DEVICE + "; " + HARDWARE + "; " + LOCALE + ")"
	IG_SIG_KEY      = "012a54f51c49aa8c5c322416ab1410909add32c966bbaa0fe3dc58ac43fd7ede"
	EXPERIMENTS     = "ig_android_progressive_jpeg,ig_creation_growth_holdout,ig_android_report_and_hide,ig_android_new_browser,ig_android_enable_share_to_whatsapp,ig_android_direct_drawing_in_quick_cam_universe,ig_android_huawei_app_badging,ig_android_universe_video_production,ig_android_asus_app_badging,ig_android_direct_plus_button,ig_android_ads_heatmap_overlay_universe,ig_android_http_stack_experiment_2016,ig_android_infinite_scrolling,ig_fbns_blocked,ig_android_white_out_universe,ig_android_full_people_card_in_user_list,ig_android_post_auto_retry_v7_21,ig_fbns_push,ig_android_feed_pill,ig_android_profile_link_iab,ig_explore_v3_us_holdout,ig_android_histogram_reporter,ig_android_anrwatchdog,ig_android_search_client_matching,ig_android_high_res_upload_2,ig_android_new_browser_pre_kitkat,ig_android_2fac,ig_android_grid_video_icon,ig_android_white_camera_universe,ig_android_disable_chroma_subsampling,ig_android_share_spinner,ig_android_explore_people_feed_icon,ig_explore_v3_android_universe,ig_android_media_favorites,ig_android_nux_holdout,ig_android_search_null_state,ig_android_react_native_notification_setting,ig_android_ads_indicator_change_universe,ig_android_video_loading_behavior,ig_android_black_camera_tab,liger_instagram_android_univ,ig_explore_v3_internal,ig_android_direct_emoji_picker,ig_android_prefetch_explore_delay_time,ig_android_business_insights_qe,ig_android_direct_media_size,ig_android_enable_client_share,ig_android_promoted_posts,ig_android_app_badging_holdout,ig_android_ads_cta_universe,ig_android_mini_inbox_2,ig_android_feed_reshare_button_nux,ig_android_boomerang_feed_attribution,ig_android_fbinvite_qe,ig_fbns_shared,ig_android_direct_full_width_media,ig_android_hscroll_profile_chaining,ig_android_feed_unit_footer,ig_android_media_tighten_space,ig_android_private_follow_request,ig_android_inline_gallery_backoff_hours_universe,ig_android_direct_thread_ui_rewrite,ig_android_rendering_controls,ig_android_ads_full_width_cta_universe,ig_video_max_duration_qe_preuniverse,ig_android_prefetch_explore_expire_time,ig_timestamp_public_test,ig_android_profile,ig_android_dv2_consistent_http_realtime_response,ig_android_enable_share_to_messenger,ig_explore_v3,ig_ranking_following,ig_android_pending_request_search_bar,ig_android_feed_ufi_redesign,ig_android_video_pause_logging_fix,ig_android_default_folder_to_camera,ig_android_video_stitching_7_23,ig_android_profanity_filter,ig_android_business_profile_qe,ig_android_search,ig_android_boomerang_entry,ig_android_inline_gallery_universe,ig_android_ads_overlay_design_universe,ig_android_options_app_invite,ig_android_view_count_decouple_likes_universe,ig_android_periodic_analytics_upload_v2,ig_android_feed_unit_hscroll_auto_advance,ig_peek_profile_photo_universe,ig_android_ads_holdout_universe,ig_android_prefetch_explore,ig_android_direct_bubble_icon,ig_video_use_sve_universe,ig_android_inline_gallery_no_backoff_on_launch_universe,ig_android_image_cache_multi_queue,ig_android_camera_nux,ig_android_immersive_viewer,ig_android_dense_feed_unit_cards,ig_android_sqlite_dev,ig_android_exoplayer,ig_android_add_to_last_post,ig_android_direct_public_threads,ig_android_prefetch_venue_in_composer,ig_android_bigger_share_button,ig_android_dv2_realtime_private_share,ig_android_non_square_first,ig_android_video_interleaved_v2,ig_android_follow_search_bar,ig_android_last_edits,ig_android_video_download_logging,ig_android_ads_loop_count_universe,ig_android_swipeable_filters_blacklist,ig_android_boomerang_layout_white_out_universe,ig_android_ads_carousel_multi_row_universe,ig_android_mentions_invite_v2,ig_android_direct_mention_qe,ig_android_following_follower_social_context"
	SIG_KEY_VERSION = "4"
)

var JAR_URL = buildJarUrl()
var HEADERS = map[string]string{
	"Connection":      "close",
	"Accept":          "*/*",
	"Content-type":    "application/x-www-form-urlencoded; charset=UTF-8",
	"Cookie2":         "$Version=1",
	"Accept-Language": "en-US",
	"User-Agent":      USER_AGENT,
}

func generateSignature(data []byte) string {
	key := []byte(IG_SIG_KEY)

	sig := hmac.New(sha256.New, key)
	sig.Write(data)

	signedData := hex.EncodeToString(sig.Sum(nil))

	return "ig_sig_key_version=" + SIG_KEY_VERSION + "&signed_body=" + signedData + "." + string(data)
}

func generateSeed(username string, password string) string {
	h := md5.New()
	io.WriteString(h, username+password)
	return hex.EncodeToString(h.Sum(nil))
}

func generateDeviceId(seed string) string {
	volatileSeed := "12345"
	h := md5.New()
	io.WriteString(h, seed+volatileSeed)
	return "android-" + hex.EncodeToString(h.Sum(nil))[0:16]
}

func generateUUID(flag bool) string {
	generatedUUID := uuid.NewV4().String()
	if flag {
		return generatedUUID
	} else {
		return strings.Replace(generatedUUID, "-", "", -1)
	}
}

func buildJarUrl() *url.URL {
	url, _ := url.Parse(API_URL)
	return url
}
