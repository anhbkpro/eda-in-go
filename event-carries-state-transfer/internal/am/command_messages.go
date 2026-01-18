package am

import (
	"context"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/registry"
)

type (
	CommandMessage interface {
		Message
		ddd.Command
	}

	IncomingCommandMessage interface {
		IncomingMessage
		ddd.Command
	}

	CommandPublisher  = MessagePublisher[ddd.Command]
	CommandSubscriber interface {
		Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error
	}
	CommandStream interface {
		MessagePublisher[ddd.Command]
		CommandSubscriber
	}

	commandStream struct {
		reg    registry.Registry
		stream RawMessageStream
	}

	commandMessage struct {
		id         string
		name       string
		payload    ddd.CommandPayload
		metadata   ddd.Metadata
		occurredAt time.Time
		msg        IncomingMessage
	}
)

var _ CommandMessage = (*commandMessage)(nil)

var _ CommandStream = (*commandStream)(nil)

func NewCommandStream(reg registry.Registry, stream RawMessageStream) CommandStream {
	return &commandStream{
		reg:    reg,
		stream: stream,
	}
}

func (s commandStream) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	metadata, err := structpb.NewStruct(command.Metadata())
	if err != nil {
		return err
	}

	payload, err := s.reg.Serialize(
		command.CommandName(), command.Payload(),
	)
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(command.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, topicName, rawMessage{
		id:   command.ID(),
		name: command.CommandName(),
		data: data,
	})
}

func (s commandStream) Subscribe(topicName string, handler CommandMessageHandler, options ...SubscriberOption) error {
	cfg := NewSubscriberConfig(options)

	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, key := range cfg.MessageFilters() {
			filters[key] = struct{}{}
		}
	}

	fn := MessageHandlerFunc[IncomingRawMessage](func(ctx context.Context, msg IncomingRawMessage) error {
		var commandData CommandMessageData

		if filters != nil {
			if _, exists := filters[msg.MessageName()]; !exists {
				return nil
			}
		}

		err := proto.Unmarshal(msg.Data(), &commandData)
		if err != nil {
			return err
		}

		commandName := msg.MessageName()

		payload, err := s.reg.Deserialize(commandName, commandData.GetPayload())
		if err != nil {
			return err
		}

		commandMsg := commandMessage{
			id:         msg.ID(),
			name:       commandName,
			payload:    payload,
			metadata:   commandData.GetMetadata().AsMap(),
			occurredAt: commandData.GetOccurredAt().AsTime(),
			msg:        msg,
		}

		// where is the reply message should be sent to?
		destination := commandMsg.Metadata().Get(CommandReplyChannelHandler).(string)

		// create a new reply to store the result of the command
		var reply ddd.Reply
		// handle the command and get the reply
		reply, err = handler.HandleMessage(ctx, commandMsg)
		if err != nil {
			// if there is an error, publish a failure reply
			return s.publishReply(ctx, destination, s.failure(reply, commandMsg))
		}

		// publish the success reply
		return s.publishReply(ctx, destination, s.success(reply, commandMsg))
	})

	return s.stream.Subscribe(topicName, fn, options...)
}

func (s commandStream) publishReply(ctx context.Context, destination string, reply ddd.Reply) error {
	metadata, err := structpb.NewStruct(reply.Metadata())
	if err != nil {
		return err
	}

	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(
			reply.ReplyName(), reply.Payload(),
		)
		if err != nil {
			return err
		}
	}

	data, err := proto.Marshal(&ReplyMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(reply.OccurredAt()),
		Metadata:   metadata,
	})
	if err != nil {
		return err
	}

	return s.stream.Publish(ctx, destination, rawMessage{
		id:   reply.ID(),
		name: reply.ReplyName(),
		data: data,
	})
}

func (s commandStream) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	// if no reply is provided, create a failure reply with no payload
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHandler, OutcomeFailure)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandStream) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	// if no reply is provided, create a success reply with no payload
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHandler, OutcomeSuccess)

	return s.applyCorrelationHeaders(reply, cmd)
}

func (s commandStream) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	// apply the correlation headers to the reply
	// this is used to correlate the reply to the command
	for key, value := range cmd.Metadata() {
		if key == CommandNameHandler {
			continue
		}

		if strings.HasPrefix(key, CommandHandlerPrefix) {
			hdr := ReplyHandlerPrefix + key[len(CommandHandlerPrefix):]
			reply.Metadata().Set(hdr, value)
		}
	}

	return reply
}

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) Metadata() ddd.Metadata      { return c.metadata }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }
func (c commandMessage) MessageName() string         { return c.msg.MessageName() }
func (c commandMessage) Ack() error                  { return c.msg.Ack() }
func (c commandMessage) NAck() error                 { return c.msg.NAck() }
func (c commandMessage) Extend() error               { return c.msg.Extend() }
func (c commandMessage) Kill() error                 { return c.msg.Kill() }
