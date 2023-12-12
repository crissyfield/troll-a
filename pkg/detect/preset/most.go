package preset

import (
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/rules"

	"github.com/crissyfield/troll-a/pkg/detect"
)

// Most is a list of most detection rules.
var Most = []detect.RuleFunction{
	rules.AWS,                             // AWS
	rules.AdafruitAPIKey,                  // Adafruit API Key
	rules.AdobeClientID,                   // Adobe Client ID (OAuth Web)
	rules.AdobeClientSecret,               // Adobe Client Secret
	rules.AgeSecretKey,                    // Age secret key
	rules.Airtable,                        // Airtable API Key
	rules.AlgoliaApiKey,                   // Algolia API Key
	rules.AlibabaAccessKey,                // Alibaba AccessKey ID
	rules.AlibabaSecretKey,                // Alibaba Secret Key
	rules.AsanaClientID,                   // Asana Client ID
	rules.AsanaClientSecret,               // Asana Client Secret
	rules.Atlassian,                       // Atlassian API token
	rules.Authress,                        // Authress Service Client Access Key
	rules.Beamer,                          // Beamer API token
	rules.BitBucketClientID,               // Bitbucket Client ID
	rules.BitBucketClientSecret,           // Bitbucket Client Secret
	rules.BittrexAccessKey,                // Bittrex Access Key
	rules.BittrexSecretKey,                // Bittrex Secret Key
	rules.Clojars,                         // Clojars API token
	rules.CodecovAccessToken,              // Codecov Access Token
	rules.CoinbaseAccessToken,             // Coinbase Access Token
	rules.ConfluentAccessToken,            // Confluent Access Token
	rules.ConfluentSecretKey,              // Confluent Secret Key
	rules.Contentful,                      // Contentful delivery API token
	rules.Databricks,                      // Databricks API token
	rules.DatadogtokenAccessToken,         // Datadog Access Token
	rules.DefinedNetworkingAPIToken,       // Defined Networking API token
	rules.DigitalOceanOAuthToken,          // DigitalOcean OAuth Access Token
	rules.DigitalOceanPAT,                 // DigitalOcean Personal Access Token
	rules.DigitalOceanRefreshToken,        // DigitalOcean OAuth Refresh Token
	rules.DiscordAPIToken,                 // Discord API key
	rules.DiscordClientID,                 // Discord client ID
	rules.DiscordClientSecret,             // Discord client secret
	rules.Doppler,                         // Doppler API token
	rules.DroneciAccessToken,              // Droneci Access Token
	rules.DropBoxAPISecret,                // Dropbox API secret
	rules.DropBoxLongLivedAPIToken,        // Dropbox long lived API token
	rules.DropBoxShortLivedAPIToken,       // Dropbox short lived API token
	rules.Duffel,                          // Duffel API token
	rules.Dynatrace,                       // Dynatrace API token
	rules.EasyPost,                        // EasyPost API token
	rules.EasyPostTestAPI,                 // EasyPost test API token
	rules.EtsyAccessToken,                 // Etsy Access Token
	rules.Facebook,                        // Facebook Access Token
	rules.FastlyAPIToken,                  // Fastly API key
	rules.FinicityAPIToken,                // Finicity API token
	rules.FinicityClientSecret,            // Finicity Client Secret
	rules.FinnhubAccessToken,              // Finnhub Access Token
	rules.FlickrAccessToken,               // Flickr Access Token
	rules.FlutterwaveEncKey,               // Flutterwave Encryption Key
	rules.FlutterwavePublicKey,            // Finicity Public Key
	rules.FlutterwaveSecretKey,            // Flutterwave Secret Key
	rules.FrameIO,                         // Frame.io API token
	rules.FreshbooksAccessToken,           // Freshbooks Access Token
	rules.GitHubApp,                       // GitHub App Token
	rules.GitHubFineGrainedPat,            // GitHub Fine-Grained Personal Access Token
	rules.GitHubOauth,                     // GitHub OAuth Access Token
	rules.GitHubPat,                       // GitHub Personal Access Token
	rules.GitHubRefresh,                   // GitHub Refresh Token
	rules.GitlabPat,                       // GitLab Personal Access Token
	rules.GitlabPipelineTriggerToken,      // GitLab Pipeline Trigger Token
	rules.GitlabRunnerRegistrationToken,   // GitLab Runner Registration Token
	rules.GitterAccessToken,               // Gitter Access Token
	rules.GoCardless,                      // GoCardless API token
	rules.GrafanaApiKey,                   // Grafana api key (or Grafana cloud api key)
	rules.GrafanaCloudApiToken,            // Grafana cloud api token
	rules.GrafanaServiceAccountToken,      // Grafana service account token
	rules.Hashicorp,                       // HashiCorp Terraform user/org API token
	rules.HashicorpField,                  // HashiCorp Terraform password field
	rules.Heroku,                          // Heroku API Key
	rules.HubSpot,                         // HubSpot API Token
	rules.HuggingFaceAccessToken,          // Hugging Face Access token
	rules.HuggingFaceOrganizationApiToken, // Hugging Face Organization API token
	rules.InfracostAPIToken,               // Infracost API Token
	rules.Intercom,                        // Intercom API Token
	rules.JFrogAPIKey,                     // JFrog API Key
	rules.JFrogIdentityToken,              // JFrog Identity Token
	rules.KrakenAccessToken,               // Kraken Access Token
	rules.KucoinAccessToken,               // Kucoin Access Token
	rules.KucoinSecretKey,                 // Kucoin Secret Key
	rules.LaunchDarklyAccessToken,         // Launchdarkly Access Token
	rules.LinearAPIToken,                  // Linear API Token
	rules.LinearClientSecret,              // Linear Client Secret
	rules.LinkedinClientID,                // LinkedIn Client ID
	rules.LinkedinClientSecret,            // LinkedIn Client secret
	rules.LobAPIToken,                     // Lob API Key
	rules.LobPubAPIToken,                  // Lob Publishable API Key
	rules.MailChimp,                       // Mailchimp API key
	rules.MailGunPrivateAPIToken,          // Mailgun private API token
	rules.MailGunPubAPIToken,              // Mailgun public validation key
	rules.MailGunSigningKey,               // Mailgun webhook signing key
	rules.MapBox,                          // MapBox API token
	rules.MattermostAccessToken,           // Mattermost Access Token
	rules.MessageBirdAPIToken,             // MessageBird API token
	rules.MessageBirdClientID,             // MessageBird client ID
	rules.NPM,                             // npm access token
	rules.NetlifyAccessToken,              // Netlify Access Token
	rules.NewRelicBrowserAPIKey,           // New Relic ingest browser API token
	rules.NewRelicUserID,                  // New Relic user API Key
	rules.NewRelicUserKey,                 // New Relic user API ID
	rules.NytimesAccessToken,              // Nytimes Access Token
	rules.OktaAccessToken,                 // Okta Access Token
	rules.OpenAI,                          // OpenAI API Key
	rules.PlaidAccessID,                   // Plaid Client ID
	rules.PlaidAccessToken,                // Plaid API Token
	rules.PlaidSecretKey,                  // Plaid Secret key
	rules.PlanetScaleAPIToken,             // PlanetScale API token
	rules.PlanetScaleOAuthToken,           // PlanetScale OAuth token
	rules.PlanetScalePassword,             // PlanetScale password
	rules.PostManAPI,                      // Postman API token
	rules.Prefect,                         // Prefect API token
	rules.PrivateKey,                      // Private Key
	rules.PulumiAPIToken,                  // Pulumi API token
	rules.PyPiUploadToken,                 // PyPI upload token
	rules.RapidAPIAccessToken,             // RapidAPI Access Token
	rules.ReadMe,                          // Readme API token
	rules.RubyGemsAPIToken,                // Rubygem API token
	rules.ScalingoAPIToken,                // Scalingo API token
	rules.SendGridAPIToken,                // SendGrid API token
	rules.SendInBlueAPIToken,              // Sendinblue API token
	rules.SendbirdAccessID,                // Sendbird Access ID
	rules.SendbirdAccessToken,             // Sendbird Access Token
	rules.SentryAccessToken,               // Sentry Access Token
	rules.ShippoAPIToken,                  // Shippo API token
	rules.ShopifyAccessToken,              // Shopify access token
	rules.ShopifyCustomAccessToken,        // Shopify custom access token
	rules.ShopifyPrivateAppAccessToken,    // Shopify private app access token
	rules.ShopifySharedSecret,             // Shopify shared secret
	rules.SidekiqSecret,                   // Sidekiq Secret
	rules.SidekiqSensitiveUrl,             // Sidekiq Sensitive URL
	rules.SlackAppLevelToken,              // Slack App-level token
	rules.SlackBotToken,                   // Slack Bot token
	rules.SlackConfigurationRefreshToken,  // Slack Configuration refresh token
	rules.SlackConfigurationToken,         // Slack Configuration access token
	rules.SlackLegacyBotToken,             // Slack Legacy bot token
	rules.SlackLegacyToken,                // Slack Legacy token
	rules.SlackLegacyWorkspaceToken,       // Slack Legacy Workspace token
	rules.SlackUserToken,                  // Slack User
	rules.SlackWebHookUrl,                 // Slack Webhook
	rules.Snyk,                            // Snyk API token
	rules.SquareAccessToken,               // Square Access Token
	rules.SquareSpaceAccessToken,          // Squarespace Access Token
	rules.SumoLogicAccessID,               // SumoLogic Access ID
	rules.SumoLogicAccessToken,            // SumoLogic Access Token
	rules.TeamsWebhook,                    // Microsoft Teams Webhook
	rules.TelegramBotToken,                // Telegram Bot API Token
	rules.TravisCIAccessToken,             // Travis CI Access Token
	rules.Twilio,                          // Twilio API Key
	rules.TwitchAPIToken,                  // Twitch API token
	rules.TwitterAPIKey,                   // Twitter API Key
	rules.TwitterAPISecret,                // Twitter API Secret
	rules.TwitterAccessSecret,             // Twitter Access Secret
	rules.TwitterAccessToken,              // Twitter Access Token
	rules.TwitterBearerToken,              // Twitter Bearer Token
	rules.Typeform,                        // Typeform API token
	rules.VaultBatchToken,                 // Vault Batch Token
	rules.VaultServiceToken,               // Vault Service Token
	rules.YandexAPIKey,                    // Yandex API Key
	rules.YandexAWSAccessToken,            // Yandex AWS Access Token
	rules.YandexAccessToken,               // Yandex Access Token
	rules.ZendeskSecretKey,                // Zendesk Secret Key
}
