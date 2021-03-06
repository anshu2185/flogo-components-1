package slack

import (
	"context"
	"errors"
	"net/url"
)

type listPinsResponseFull struct {
	Items  []Item
	Paging `json:"paging"`
	SlackResponse
}

// AddPin pins an item in a channel
func (api *Client) AddPin(channel string, item ItemRef) error {
	return api.AddPinContext(context.Background(), channel, item)
}

// AddPinContext pins an item in a channel with a custom context
func (api *Client) AddPinContext(ctx context.Context, channel string, item ItemRef) error {
	values := url.Values{
		"channel": {channel},
		"token":   {api.token},
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.File != "" {
		values.Set("file", string(item.File))
	}
	if item.Comment != "" {
		values.Set("file_comment", string(item.Comment))
	}

	response := &SlackResponse{}
	if err := post(ctx, api.httpclient, "pins.add", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// RemovePin un-pins an item from a channel
func (api *Client) RemovePin(channel string, item ItemRef) error {
	return api.RemovePinContext(context.Background(), channel, item)
}

// RemovePinContext un-pins an item from a channel with a custom context
func (api *Client) RemovePinContext(ctx context.Context, channel string, item ItemRef) error {
	values := url.Values{
		"channel": {channel},
		"token":   {api.token},
	}
	if item.Timestamp != "" {
		values.Set("timestamp", string(item.Timestamp))
	}
	if item.File != "" {
		values.Set("file", string(item.File))
	}
	if item.Comment != "" {
		values.Set("file_comment", string(item.Comment))
	}

	response := &SlackResponse{}
	if err := post(ctx, api.httpclient, "pins.remove", values, response, api.debug); err != nil {
		return err
	}
	if !response.Ok {
		return errors.New(response.Error)
	}
	return nil
}

// ListPins returns information about the items a user reacted to.
func (api *Client) ListPins(channel string) ([]Item, *Paging, error) {
	return api.ListPinsContext(context.Background(), channel)
}

// ListPinsContext returns information about the items a user reacted to with a custom context.
func (api *Client) ListPinsContext(ctx context.Context, channel string) ([]Item, *Paging, error) {
	values := url.Values{
		"channel": {channel},
		"token":   {api.token},
	}

	response := &listPinsResponseFull{}
	err := post(ctx, api.httpclient, "pins.list", values, response, api.debug)
	if err != nil {
		return nil, nil, err
	}
	if !response.Ok {
		return nil, nil, errors.New(response.Error)
	}
	return response.Items, &response.Paging, nil
}
