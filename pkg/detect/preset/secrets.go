package preset

import (
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/rules"

	"github.com/crissyfield/troll-a/pkg/detect"
)

// Secret is a list of secrets, access and refresh token detection rules.
var Secret = []detect.RuleFunction{
	rules.AdobeClientSecret,              // Adobe Client Secret
	rules.AgeSecretKey,                   // Age secret key
	rules.AlibabaSecretKey,               // Alibaba Secret Key
	rules.AsanaClientSecret,              // Asana Client Secret
	rules.BitBucketClientSecret,          // Bitbucket Client Secret
	rules.BittrexSecretKey,               // Bittrex Secret Key
	rules.CodecovAccessToken,             // Codecov Access Token
	rules.CoinbaseAccessToken,            // Coinbase Access Token
	rules.ConfluentAccessToken,           // Confluent Access Token
	rules.ConfluentSecretKey,             // Confluent Secret Key
	rules.DatadogtokenAccessToken,        // Datadog Access Token
	rules.DigitalOceanOAuthToken,         // DigitalOcean OAuth Access Token
	rules.DigitalOceanPAT,                // DigitalOcean Personal Access Token
	rules.DigitalOceanRefreshToken,       // DigitalOcean OAuth Refresh Token
	rules.DiscordClientSecret,            // Discord client secret
	rules.DroneciAccessToken,             // Droneci Access Token
	rules.DropBoxAPISecret,               // Dropbox API secret
	rules.DropBoxLongLivedAPIToken,       // Dropbox long lived API token
	rules.DropBoxShortLivedAPIToken,      // Dropbox short lived API token
	rules.EtsyAccessToken,                // Etsy Access Token
	rules.Facebook,                       // Facebook Access Token
	rules.FinicityClientSecret,           // Finicity Client Secret
	rules.FinnhubAccessToken,             // Finnhub Access Token
	rules.FlickrAccessToken,              // Flickr Access Token
	rules.FlutterwaveSecretKey,           // Flutterwave Secret Key
	rules.FreshbooksAccessToken,          // Freshbooks Access Token
	rules.GitHubFineGrainedPat,           // GitHub Fine-Grained Personal Access Token
	rules.GitHubOauth,                    // GitHub OAuth Access Token
	rules.GitHubPat,                      // GitHub Personal Access Token
	rules.GitHubRefresh,                  // GitHub Refresh Token
	rules.GitlabPat,                      // GitLab Personal Access Token
	rules.GitterAccessToken,              // Gitter Access Token
	rules.HuggingFaceAccessToken,         // Hugging Face Access token
	rules.KrakenAccessToken,              // Kraken Access Token
	rules.KucoinAccessToken,              // Kucoin Access Token
	rules.KucoinSecretKey,                // Kucoin Secret Key
	rules.LaunchDarklyAccessToken,        // Launchdarkly Access Token
	rules.LinearClientSecret,             // Linear Client Secret
	rules.LinkedinClientSecret,           // LinkedIn Client secret
	rules.MailGunPrivateAPIToken,         // Mailgun private API token
	rules.MattermostAccessToken,          // Mattermost Access Token
	rules.NPM,                            // npm access token
	rules.NetlifyAccessToken,             // Netlify Access Token
	rules.NytimesAccessToken,             // Nytimes Access Token
	rules.OktaAccessToken,                // Okta Access Token
	rules.PlaidAccessToken,               // Plaid API Token
	rules.PlaidSecretKey,                 // Plaid Secret key
	rules.PlanetScalePassword,            // PlanetScale password
	rules.PrivateKey,                     // Private Key
	rules.RapidAPIAccessToken,            // RapidAPI Access Token
	rules.SendbirdAccessToken,            // Sendbird Access Token
	rules.SentryAccessToken,              // Sentry Access Token
	rules.ShopifyAccessToken,             // Shopify access token
	rules.ShopifyCustomAccessToken,       // Shopify custom access token
	rules.ShopifyPrivateAppAccessToken,   // Shopify private app access token
	rules.ShopifySharedSecret,            // Shopify shared secret
	rules.SidekiqSecret,                  // Sidekiq Secret
	rules.SlackConfigurationRefreshToken, // Slack Configuration refresh token
	rules.SlackConfigurationToken,        // Slack Configuration access token
	rules.SquareAccessToken,              // Square Access Token
	rules.SquareSpaceAccessToken,         // Squarespace Access Token
	rules.SumoLogicAccessToken,           // SumoLogic Access Token
	rules.TravisCIAccessToken,            // Travis CI Access Token
	rules.TwitterAPISecret,               // Twitter API Secret
	rules.TwitterAccessSecret,            // Twitter Access Secret
	rules.TwitterAccessToken,             // Twitter Access Token
	rules.TwitterBearerToken,             // Twitter Bearer Token
	rules.YandexAWSAccessToken,           // Yandex AWS Access Token
	rules.YandexAccessToken,              // Yandex Access Token
	rules.ZendeskSecretKey,               // Zendesk Secret Key
}
