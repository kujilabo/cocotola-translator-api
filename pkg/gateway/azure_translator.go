package gateway

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/translatortext"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/kujilabo/cocotola-translator-api/pkg/domain"
	"github.com/kujilabo/cocotola-translator-api/pkg/service"
	"github.com/kujilabo/cocotola-translator-api/pkg_lib/log"
)

type azureTranslationClient struct {
	client translatortext.TranslatorClient
}

type AzureDisplayTranslation struct {
	Pos        int
	Target     string
	Confidence float64
}

func NewAzureTranslationClient(subscriptionKey string) service.AzureTranslationClient {
	client := translatortext.NewTranslatorClient("https://api.cognitive.microsofttranslator.com")
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscriptionKey)
	return &azureTranslationClient{
		client: client,
	}

}

func (c *azureTranslationClient) DictionaryLookup(ctx context.Context, text string, fromLang, toLang domain.Lang2) ([]service.AzureTranslation, error) {
	logger := log.FromContext(ctx)
	result, err := c.client.DictionaryLookup(context.Background(), fromLang.String(), toLang.String(), []translatortext.DictionaryLookupTextInput{{Text: to.StringPtr(text)}}, "")
	if err != nil {
		return nil, err
	}
	if result.Value == nil {
		logger.Info("a")
		return nil, nil
	}

	translations := make([]service.AzureTranslation, 0)
	for _, v := range *result.Value {
		if v.Translations == nil {
			continue
		}

		for _, t := range *v.Translations {
			pos, err := domain.ParsePos(c.pointerToString(t.PosTag))
			if err != nil {
				return nil, err
			}
			if pos == domain.PosOther {
				logger.Warnf("PosOther. text: %s, pos: %s", text, *t.PosTag)
			}
			translations = append(translations, service.AzureTranslation{
				Pos:        pos,
				Target:     c.pointerToString(t.DisplayTarget),
				Confidence: c.pointerToFloat64(t.Confidence),
			})
		}
	}
	logger.Info("b")
	return translations, nil
}

func (c *azureTranslationClient) pointerToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (c *azureTranslationClient) pointerToFloat64(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

// func (c *azureTranslationClient) stringToPos(value string) int {
// 	switch value {
// 	case "ADJ":
// 		return 1
// 	case "ADV":
// 		return 2
// 	case "CONJ":
// 		return 3
// 	case "DET":
// 		return 4
// 	case "MODAL":
// 		return 5
// 	case "NOUN":
// 		return 6
// 	case "PREP":
// 		return 7
// 	case "PRON":
// 		return 8
// 	case "VERB":
// 		return 9
// 	default:
// 		return 99
// 	}
// }
