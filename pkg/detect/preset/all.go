package preset

import (
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/rules"

	"github.com/crissyfield/troll-a/pkg/detect"
)

// All is a list of all detection rules.
var All = []detect.GitleaksRuleFunction{
	rules.AWS,                              // AWS
	rules.AdobeClientID,                    // Adobe Client ID (OAuth Web)
	rules.AdobeClientSecret,                // Adobe Client Secret
	rules.AgeSecretKey,                     // Age Secret Key
	rules.Airtable,                         // Airtable API Key
	rules.AlgoliaApiKey,                    // Algolia API Key
	rules.AlibabaAccessKey,                 // Alibaba AccessKey ID
	rules.AlibabaSecretKey,                 // Alibaba Secret Key
	rules.AsanaClientID,                    // Asana Client ID
	rules.AsanaClientSecret,                // Asana Client Secret
	rules.Atlassian,                        // Atlassian API Token
	rules.Authress,                         // Authress Service Client Access Key
	rules.AzureActiveDirectoryClientSecret, // Azure Active Directory
	rules.Beamer,                           // Beamer API Token
	rules.BitBucketClientID,                // Bitbucket Client ID
	rules.BitBucketClientSecret,            // Bitbucket Client Secret
	rules.BittrexAccessKey,                 // Bittrex Access Key
	rules.BittrexSecretKey,                 // Bittrex Secret Key
	rules.Clojars,                          // Clojars API Token
	rules.CloudflareAPIKey,                 // Cloudflare API Key
	rules.CloudflareGlobalAPIKey,           // Cloudflare Global API Key
	rules.CloudflareOriginCAKey,            // Cloudflare Origin CA Key
	rules.CodecovAccessToken,               // Codecov Access Token
	rules.CohereAPIToken,                   // Cohere API Token
	rules.CoinbaseAccessToken,              // Coinbase Access Token
	rules.ConfluentAccessToken,             // Confluent Access Token
	rules.ConfluentSecretKey,               // Confluent Secret Key
	rules.Contentful,                       // Contentful delivery API Token
	rules.Databricks,                       // Databricks API Token
	rules.DatadogtokenAccessToken,          // Datadog Access Token
	rules.DefinedNetworkingAPIToken,        // Defined Networking API Token
	rules.DigitalOceanOAuthToken,           // DigitalOcean OAuth Access Token
	rules.DigitalOceanPAT,                  // DigitalOcean Personal Access Token
	rules.DigitalOceanRefreshToken,         // DigitalOcean OAuth Refresh Token
	rules.DiscordAPIToken,                  // Discord API Key
	rules.DiscordClientID,                  // Discord Client ID
	rules.DiscordClientSecret,              // Discord Client Secret
	rules.Doppler,                          // Doppler API Token
	rules.DroneciAccessToken,               // Droneci Access Token
	rules.DropBoxAPISecret,                 // Dropbox API Secret
	rules.DropBoxLongLivedAPIToken,         // Dropbox long lived API Token
	rules.DropBoxShortLivedAPIToken,        // Dropbox short lived API Token
	rules.Duffel,                           // Duffel API Token
	rules.Dynatrace,                        // Dynatrace API Token
	rules.EasyPost,                         // EasyPost API Token
	rules.EasyPostTestAPI,                  // EasyPost test API Token
	rules.EtsyAccessToken,                  // Etsy Access Token
	rules.FacebookAccessToken,              // Facebook Access Token
	rules.FacebookPageAccessToken,          // Facebook Page Access Token
	rules.FacebookSecret,                   // Facebook Secret
	rules.FastlyAPIToken,                   // Fastly API Key
	rules.FinicityAPIToken,                 // Finicity API Token
	rules.FinicityClientSecret,             // Finicity Client Secret
	rules.FinnhubAccessToken,               // Finnhub Access Token
	rules.FlickrAccessToken,                // Flickr Access Token
	rules.FlutterwaveEncKey,                // Flutterwave Encryption Key
	rules.FlutterwavePublicKey,             // Finicity Public Key
	rules.FlutterwaveSecretKey,             // Flutterwave Secret Key
	rules.FlyIOAccessToken,                 // Fly.io Access Token
	rules.FrameIO,                          // Frame.io API Token
	rules.FreshbooksAccessToken,            // Freshbooks Access Token
	rules.GCPAPIKey,                        // GCP API Key
	rules.GenericCredential,                // Generic API Key
	rules.GitHubApp,                        // GitHub App Token
	rules.GitHubFineGrainedPat,             // GitHub Fine-Grained Personal Access Token
	rules.GitHubOauth,                      // GitHub OAuth Access Token
	rules.GitHubPat,                        // GitHub Personal Access Token
	rules.GitHubRefresh,                    // GitHub Refresh Token
	rules.GitlabPat,                        // GitLab Personal Access Token
	rules.GitlabPipelineTriggerToken,       // GitLab Pipeline Trigger Token
	rules.GitlabRunnerRegistrationToken,    // GitLab Runner Registration Token
	rules.GitterAccessToken,                // Gitter Access Token
	rules.GoCardless,                       // GoCardless API Token
	rules.GrafanaApiKey,                    // Grafana api Key (or Grafana cloud api Key)
	rules.GrafanaCloudApiToken,             // Grafana cloud api Token
	rules.GrafanaServiceAccountToken,       // Grafana service account Token
	rules.HarnessApiKey,                    // Harness API Key
	rules.HashiCorpTerraform,               // HashiCorp Terraform User/org API Token
	rules.HashicorpField,                   // HashiCorp Terraform password field
	rules.Heroku,                           // Heroku API Key
	rules.HubSpot,                          // HubSpot API Token
	rules.HuggingFaceAccessToken,           // Hugging Face Access Token
	rules.HuggingFaceOrganizationApiToken,  // Hugging Face Organization API Token
	rules.InfracostAPIToken,                // Infracost API Token
	rules.Intercom,                         // Intercom API Token
	rules.Intra42ClientSecret,              // Intra 42 Client Secret
	rules.JFrogAPIKey,                      // JFrog API Key
	rules.JFrogIdentityToken,               // JFrog Identity Token
	rules.JWT,                              // JSON Web Token
	rules.JWTBase64,                        // Base64-encoded JSON Web Token
	rules.KrakenAccessToken,                // Kraken Access Token
	rules.KubernetesSecret,                 // Kubernetes Secret
	rules.KucoinAccessToken,                // Kucoin Access Token
	rules.KucoinSecretKey,                  // Kucoin Secret Key
	rules.LaunchDarklyAccessToken,          // Launchdarkly Access Token
	rules.LinearAPIToken,                   // Linear API Token
	rules.LinearClientSecret,               // Linear Client Secret
	rules.LinkedinClientID,                 // LinkedIn Client ID
	rules.LinkedinClientSecret,             // LinkedIn Client Secret
	rules.LobAPIToken,                      // Lob API Key
	rules.LobPubAPIToken,                   // Lob Publishable API Key
	rules.MailChimp,                        // Mailchimp API Key
	rules.MailGunPrivateAPIToken,           // Mailgun private API Token
	rules.MailGunPubAPIToken,               // Mailgun public validation Key
	rules.MailGunSigningKey,                // Mailgun webhook signing Key
	rules.MapBox,                           // MapBox API Token
	rules.MattermostAccessToken,            // Mattermost Access Token
	rules.MessageBirdAPIToken,              // MessageBird API Token
	rules.MessageBirdClientID,              // MessageBird Client ID
	rules.NPM,                              // npm Access Token
	rules.NetlifyAccessToken,               // Netlify Access Token
	rules.NewRelicBrowserAPIKey,            // New Relic ingest browser API Token
	rules.NewRelicInsertKey,                // New Relic insert Key
	rules.NewRelicUserID,                   // New Relic User API Key
	rules.NewRelicUserKey,                  // New Relic User API ID
	rules.NugetConfigPassword,              // Nuget config password
	rules.NytimesAccessToken,               // Nytimes Access Token
	rules.OktaAccessToken,                  // Okta Access Token
	rules.OpenAI,                           // OpenAI API Key
	rules.OpenshiftUserToken,               // Openshift User Token
	rules.PlaidAccessID,                    // Plaid Client ID
	rules.PlaidAccessToken,                 // Plaid API Token
	rules.PlaidSecretKey,                   // Plaid Secret Key
	rules.PlanetScaleAPIToken,              // PlanetScale API Token
	rules.PlanetScaleOAuthToken,            // PlanetScale OAuth Token
	rules.PlanetScalePassword,              // PlanetScale password
	rules.PostManAPI,                       // Postman API Token
	rules.Prefect,                          // Prefect API Token
	rules.PrivateAIToken,                   // Private AI Token
	rules.PrivateKey,                       // Private Key
	rules.PulumiAPIToken,                   // Pulumi API Token
	rules.PyPiUploadToken,                  // PyPI upload Token
	rules.RapidAPIAccessToken,              // RapidAPI Access Token
	rules.ReadMe,                           // Readme API Token
	rules.RubyGemsAPIToken,                 // Rubygem API Token
	rules.ScalingoAPIToken,                 // Scalingo API Token
	rules.SendGridAPIToken,                 // SendGrid API Token
	rules.SendInBlueAPIToken,               // Sendinblue API Token
	rules.SendbirdAccessID,                 // Sendbird Access ID
	rules.SendbirdAccessToken,              // Sendbird Access Token
	rules.SentryAccessToken,                // Sentry Access Token
	rules.ShippoAPIToken,                   // Shippo API Token
	rules.ShopifyAccessToken,               // Shopify Access Token
	rules.ShopifyCustomAccessToken,         // Shopify custom Access Token
	rules.ShopifyPrivateAppAccessToken,     // Shopify private app Access Token
	rules.ShopifySharedSecret,              // Shopify shared Secret
	rules.SidekiqSecret,                    // Sidekiq Secret
	rules.SidekiqSensitiveUrl,              // Sidekiq Sensitive URL
	rules.SlackAppLevelToken,               // Slack App-level Token
	rules.SlackBotToken,                    // Slack Bot Token
	rules.SlackConfigurationRefreshToken,   // Slack Configuration refresh Token
	rules.SlackConfigurationToken,          // Slack Configuration Access Token
	rules.SlackLegacyBotToken,              // Slack Legacy bot Token
	rules.SlackLegacyToken,                 // Slack Legacy Token
	rules.SlackLegacyWorkspaceToken,        // Slack Legacy Workspace Token
	rules.SlackUserToken,                   // Slack User
	rules.SlackWebHookUrl,                  // Slack Webhook
	rules.Snyk,                             // Snyk API Token
	rules.SquareAccessToken,                // Square Access Token
	rules.SquareSpaceAccessToken,           // Squarespace Access Token
	rules.StripeAccessToken,                // Stripe Access Token
	rules.SumoLogicAccessID,                // SumoLogic Access ID
	rules.SumoLogicAccessToken,             // SumoLogic Access Token
	rules.TeamsWebhook,                     // Microsoft Teams Webhook
	rules.TelegramBotToken,                 // Telegram Bot API Token
	rules.TravisCIAccessToken,              // Travis CI Access Token
	rules.Twilio,                           // Twilio API Key
	rules.TwitchAPIToken,                   // Twitch API Token
	rules.TwitterAPIKey,                    // Twitter API Key
	rules.TwitterAPISecret,                 // Twitter API Secret
	rules.TwitterAccessSecret,              // Twitter Access Secret
	rules.TwitterAccessToken,               // Twitter Access Token
	rules.TwitterBearerToken,               // Twitter Bearer Token
	rules.Typeform,                         // Typeform API Token
	rules.VaultBatchToken,                  // Vault Batch Token
	rules.VaultServiceToken,                // Vault Service Token
	rules.YandexAPIKey,                     // Yandex API Key
	rules.YandexAWSAccessToken,             // Yandex AWS Access Token
	rules.YandexAccessToken,                // Yandex Access Token
	rules.ZendeskSecretKey,                 // Zendesk Secret Key
}
