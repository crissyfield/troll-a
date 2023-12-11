package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/zricethezav/gitleaks/v8/cmd/generate/config/rules"
	"github.com/zricethezav/gitleaks/v8/config"

	"github.com/crissyfield/troll-a/internal/fetch"
)

// rules ...
var detectionRules = []*config.Rule{
	// rules.AdafruitAPIKey(),
	// rules.AdobeClientID(),
	// rules.AdobeClientSecret(),
	// rules.AgeSecretKey(),
	// rules.Airtable(),
	// rules.AlgoliaApiKey(),
	// rules.AlibabaAccessKey(),
	// rules.AlibabaSecretKey(),
	// rules.AsanaClientID(),
	// rules.AsanaClientSecret(),
	// rules.Atlassian(),
	// rules.Authress(),
	rules.AWS(),
	// rules.BitBucketClientID(),
	// rules.BitBucketClientSecret(),
	// rules.BittrexAccessKey(),
	// rules.BittrexSecretKey(),
	// rules.Beamer(),
	// rules.CodecovAccessToken(),
	// rules.CoinbaseAccessToken(),
	// rules.Clojars(),
	// rules.ConfluentAccessToken(),
	// rules.ConfluentSecretKey(),
	// rules.Contentful(),
	// rules.Databricks(),
	// rules.DatadogtokenAccessToken(),
	// rules.DefinedNetworkingAPIToken(),
	// rules.DigitalOceanPAT(),
	// rules.DigitalOceanOAuthToken(),
	// rules.DigitalOceanRefreshToken(),
	// rules.DiscordAPIToken(),
	// rules.DiscordClientID(),
	// rules.DiscordClientSecret(),
	// rules.Doppler(),
	// rules.DropBoxAPISecret(),
	// rules.DropBoxLongLivedAPIToken(),
	// rules.DropBoxShortLivedAPIToken(),
	// rules.DroneciAccessToken(),
	// rules.Duffel(),
	// rules.Dynatrace(),
	// rules.EasyPost(),
	// rules.EasyPostTestAPI(),
	// rules.EtsyAccessToken(),
	// rules.Facebook(),
	// rules.FastlyAPIToken(),
	// rules.FinicityClientSecret(),
	// rules.FinicityAPIToken(),
	// rules.FlickrAccessToken(),
	// rules.FinnhubAccessToken(),
	// rules.FlutterwavePublicKey(),
	// rules.FlutterwaveSecretKey(),
	// rules.FlutterwaveEncKey(),
	// rules.FrameIO(),
	// rules.FreshbooksAccessToken(),
	// rules.GoCardless(),
	// rules.GCPAPIKey(),
	// rules.GitHubPat(),
	// rules.GitHubFineGrainedPat(),
	// rules.GitHubOauth(),
	// rules.GitHubApp(),
	// rules.GitHubRefresh(),
	// rules.GitlabPat(),
	// rules.GitlabPipelineTriggerToken(),
	// rules.GitlabRunnerRegistrationToken(),
	// rules.GitterAccessToken(),
	// rules.GrafanaApiKey(),
	// rules.GrafanaCloudApiToken(),
	// rules.GrafanaServiceAccountToken(),
	// rules.Hashicorp(),
	// rules.HashicorpField(),
	// rules.Heroku(),
	// rules.HubSpot(),
	// rules.HuggingFaceAccessToken(),
	// rules.HuggingFaceOrganizationApiToken(),
	// rules.Intercom(),
	// rules.JFrogAPIKey(),
	// rules.JFrogIdentityToken(),
	// rules.JWT(),
	// rules.JWTBase64(),
	// rules.KrakenAccessToken(),
	// rules.KucoinAccessToken(),
	// rules.KucoinSecretKey(),
	// rules.LaunchDarklyAccessToken(),
	// rules.LinearAPIToken(),
	// rules.LinearClientSecret(),
	// rules.LinkedinClientID(),
	// rules.LinkedinClientSecret(),
	// rules.LobAPIToken(),
	// rules.LobPubAPIToken(),
	// rules.MailChimp(),
	// rules.MailGunPubAPIToken(),
	// rules.MailGunPrivateAPIToken(),
	// rules.MailGunSigningKey(),
	// rules.MapBox(),
	// rules.MattermostAccessToken(),
	// rules.MessageBirdAPIToken(),
	// rules.MessageBirdClientID(),
	// rules.NetlifyAccessToken(),
	// rules.NewRelicUserID(),
	// rules.NewRelicUserKey(),
	// rules.NewRelicBrowserAPIKey(),
	// rules.NPM(),
	// rules.NytimesAccessToken(),
	// rules.OktaAccessToken(),
	// rules.OpenAI(),
	// rules.PlaidAccessID(),
	// rules.PlaidSecretKey(),
	// rules.PlaidAccessToken(),
	// rules.PlanetScalePassword(),
	// rules.PlanetScaleAPIToken(),
	// rules.PlanetScaleOAuthToken(),
	// rules.PostManAPI(),
	// rules.Prefect(),
	// rules.PrivateKey(),
	// rules.PulumiAPIToken(),
	// rules.PyPiUploadToken(),
	// rules.RapidAPIAccessToken(),
	// rules.ReadMe(),
	// rules.RubyGemsAPIToken(),
	// rules.ScalingoAPIToken(),
	// rules.SendbirdAccessID(),
	// rules.SendbirdAccessToken(),
	// rules.SendGridAPIToken(),
	// rules.SendInBlueAPIToken(),
	// rules.SentryAccessToken(),
	// rules.ShippoAPIToken(),
	// rules.ShopifyAccessToken(),
	// rules.ShopifyCustomAccessToken(),
	// rules.ShopifyPrivateAppAccessToken(),
	// rules.ShopifySharedSecret(),
	// rules.SidekiqSecret(),
	// rules.SidekiqSensitiveUrl(),
	// rules.SlackBotToken(),
	// rules.SlackUserToken(),
	// rules.SlackAppLevelToken(),
	// rules.SlackConfigurationToken(),
	// rules.SlackConfigurationRefreshToken(),
	// rules.SlackLegacyBotToken(),
	// rules.SlackLegacyWorkspaceToken(),
	// rules.SlackLegacyToken(),
	// rules.SlackWebHookUrl(),
	// rules.Snyk(),
	// rules.StripeAccessToken(),
	// rules.SquareAccessToken(),
	// rules.SquareSpaceAccessToken(),
	// rules.SumoLogicAccessID(),
	// rules.SumoLogicAccessToken(),
	// rules.TeamsWebhook(),
	// rules.TelegramBotToken(),
	// rules.TravisCIAccessToken(),
	// rules.Twilio(),
	// rules.TwitchAPIToken(),
	// rules.TwitterAPIKey(),
	// rules.TwitterAPISecret(),
	// rules.TwitterAccessToken(),
	// rules.TwitterAccessSecret(),
	// rules.TwitterBearerToken(),
	// rules.Typeform(),
	// rules.VaultBatchToken(),
	// rules.VaultServiceToken(),
	// rules.YandexAPIKey(),
	// rules.YandexAWSAccessToken(),
	// rules.YandexAccessToken(),
	// rules.ZendeskSecretKey(),
	// rules.GenericCredential(),
	// rules.InfracostAPIToken(),
}

var allowedPayloadTypes = map[string]bool{
	"application/atom+xml":      true, // https://www.rfc-editor.org/rfc/rfc5023.html
	"application/json":          true, // https://www.rfc-editor.org/rfc/rfc8259.html
	"application/mbox":          true, // https://www.rfc-editor.org/rfc/rfc4155.html
	"application/msword":        true, // Microsoft Word Document or Document Template
	"application/pgp-signature": true,
	"application/rdf+xml":       true,
	"application/rss+xml":       true,
	"application/rtf":           true,
	"application/vnd.ms-excel":  true,
	"application/x-sh":          true,
	"application/xhtml+xml":     true,
	"application/xml":           true,
	"image/svg+xml":             true,
	"message/rfc822":            true,
	"text/css":                  true,
	"text/csv":                  true,
	"text/html":                 true,
	"text/plain":                true,
	"text/x-chdr":               true,
	"text/x-diff":               true,
	"text/x-log":                true,
	"text/x-perl":               true,
	"text/x-php":                true,
	"text/x-vcard":              true,
}

// Buffer ...
type Buffer struct {
	TargetURI string
	Content   []byte
}

// CmdTest defines the CLI sub-command 'test'.
var CmdTest = &cobra.Command{
	Use:   "test [flags] [warc url]",
	Short: "...",
	Args:  cobra.ExactArgs(1),
	Run:   runTest,
}

// Initialize CLI options.
func init() {
}

// runTest is called when the "test" command is used.
func runTest(_ *cobra.Command, args []string) {
	// Read URL
	ur, err := fetch.URL(args[0], fetch.WithTimeout(4*time.Hour))
	if err != nil {
		slog.Error("Unable to fetch WARC file", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	defer ur.Close()

	// Decompress
	gr, err := gzip.NewReader(ur)
	if err != nil {
		slog.Error("Unable to decompress WARC body", slog.Any("error", err))
		os.Exit(1) //nolint
	}

	// Spawn go routines to check buffers for secrets
	bufferCh := make(chan *Buffer)
	wg := &sync.WaitGroup{}

	for j := 0; j < 8; j++ {
		wg.Add(1)
		go findSecret(wg, bufferCh)
	}

	// Buffered IO
	br := bufio.NewReaderSize(gr, 4*1024*1024)

	for {
		// Read version
		version, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			slog.Error("Error while reading WARC record version", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		if !strings.HasPrefix(string(version), "WARC/") {
			slog.Error("Unknown WARC record version", slog.String("version", string(version)))
			os.Exit(1) //nolint
		}

		// Read headers
		headers := make(map[string]string)

		for {
			// Read header
			header, isPrefix, err := br.ReadLine()
			if err != nil {
				slog.Error("Error while processing WARC record header", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			// Exit if the buffer is not big enough (32KiB)
			if isPrefix {
				slog.Error("WARC record header seems too big", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			// Stop reading headers on empty line
			if len(header) == 0 {
				break
			}

			// Split header into key and value
			parts := strings.SplitN(string(header), ":", 2)
			if len(parts) == 2 {
				key := strings.ToLower(parts[0])
				value := strings.TrimSpace(parts[1])

				headers[key] = value
			}
		}

		// Extract length of record content
		length, err := strconv.Atoi(headers["content-length"])
		if err != nil {
			slog.Error("Unable to read record content length", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		// ...
		cr := io.LimitReader(br, int64(length))

		if (headers["warc-type"] == "response") && allowedPayloadTypes[headers["warc-identified-payload-type"]] {
			// ...
			content, err := io.ReadAll(cr)
			if err != nil {
				slog.Error("Unable to read content block", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			bufferCh <- &Buffer{
				TargetURI: headers["warc-target-uri"],
				Content:   content,
			}
		}

		// Discard remaining content block
		_, err = io.Copy(io.Discard, cr)
		if err != nil {
			slog.Error("Unable to discard remaining content block", slog.Any("error", err))
			os.Exit(1) //nolint
		}

		// Skip two empty lines
		for i := 0; i < 2; i++ {
			boundary, _, err := br.ReadLine()
			if (err != nil) && (err != io.EOF) {
				slog.Error("Unable to read WARC record boundary", slog.Any("error", err))
				os.Exit(1) //nolint
			}

			if len(boundary) != 0 {
				slog.Error("WARC record boundary not empty", slog.Any("error", err))
				os.Exit(1) //nolint
			}
		}
	}

	slog.Info("Done", slog.String("url", args[0]))
}

// ...
func findSecret(wg *sync.WaitGroup, bufferCh chan *Buffer) {
	// ...
	defer wg.Done()

	for buffer := range bufferCh {
		// ...
		for _, r := range detectionRules {
			// ...
			idxs := r.Regex.FindAllIndex(buffer.Content, -1)
			for _, idx := range idxs {
				// ...
				if (idx[0] > 0) && isAlphaNum(buffer.Content[idx[0]-1]) {
					continue
				}

				if (idx[1] < len(buffer.Content)-1) && isAlphaNum(buffer.Content[idx[1]+1]) {
					continue
				}

				// ...
				fmt.Printf(
					"\033[96m%s:%d\033[0m: \033[91m%s\033[0m: \033[37m%s\033[93m%s\033[37m%s\033[0m\n",
					buffer.TargetURI,
					idx[0],
					r.RuleID,
					cleanUpStrings(string(buffer.Content[max(0, idx[0]-20):idx[0]])),
					string(buffer.Content[idx[0]:idx[1]]),
					cleanUpStrings(string(buffer.Content[idx[1]:min(len(buffer.Content), idx[1]+20)])),
				)
			}
		}
	}
}

// ...
func cleanUpStrings(in string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		},
		in,
	)
}

// ...
func isAlphaNum(in byte) bool {
	return (in >= 48 && in <= 57) || (in >= 65 && in <= 90) || (in >= 97 && in <= 122)
}
