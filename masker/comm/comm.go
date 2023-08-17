package masker

import (
	"context"
	"git.multiverse.io/eventkit/kit/common/msg"
	"git.multiverse.io/eventkit/kit/log"
	masker "git.multiverse.io/eventkit/kit/masker/json"
)

func MaskMessageIfNecessary(ctx context.Context, commonHeaderMaskRules []string, headerMaskRules []string,
	commonBodyMaskRules []string, bodyMaskRules []string, message *msg.Message, bodyType... string) (*msg.Message, error) {
	var finalAppProps map[string]string
	var finalBody []byte
	var gerr error
	var deduplicationMap = make(map[string]bool, 0)

	if nil == message {
		return message, nil
	}

	// combine the header mask rules
	for _, v := range commonHeaderMaskRules {
		deduplicationMap[v] = true
	}

	for _, v := range headerMaskRules {
		deduplicationMap[v] = true
	}

	finalHeaderMaskRules := make([]string, 0)
	for k, _ := range deduplicationMap {
		finalHeaderMaskRules = append(finalHeaderMaskRules, k)
	}

	finalAppProps = message.GetAppProps()
	if len(finalAppProps) > 0 {
		if len(finalHeaderMaskRules) > 0 {
			finalAppProps, _, gerr = masker.MapValueMask(finalAppProps, finalHeaderMaskRules)
			if nil != gerr {
				log.Errorf(ctx, "Failed to do the header mask for request, error:%++v", gerr)
			}
		}
	}

	// combine the body mask rules
	for _, v := range commonBodyMaskRules {
		deduplicationMap[v] = true
	}

	for _, v := range bodyMaskRules {
		deduplicationMap[v] = true
	}

	finalBodyMaskRules := make([]string, 0)
	for k, _ := range deduplicationMap {
		finalBodyMaskRules = append(finalBodyMaskRules, k)
	}
	finalBody = message.Body
	if len(finalBody) > 0 {
		if len(finalBodyMaskRules) > 0 {
			finalBody, _, gerr = masker.JsonBodyMask(finalBody, finalBodyMaskRules)
			if nil != gerr {
				return nil, gerr
			}
		}
	}

	finalMessage := &msg.Message{
		ID:             message.ID,
		TopicAttribute: message.TopicAttribute,
		RequestURL:     message.RequestURL,
		NeedReply:      message.NeedReply,
		NeedAck:        message.NeedAck,
		SessionName:    message.SessionName,
		Body:           finalBody,
	}
	finalMessage.SetAppProps(finalAppProps)

	return finalMessage, nil
}
