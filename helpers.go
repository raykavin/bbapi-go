package bbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AccountLookupParams holds optional account lookup fields used by several
// detail and solicitation endpoints.
type AccountLookupParams struct {
	Agency     *int64
	Account    *int64
	CheckDigit *string
}

func decode[T any](data []byte) (T, error) {
	var value T
	if len(data) == 0 {
		return value, nil
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("bbapi: decode response: %w", err)
	}
	return value, nil
}

func get[T any](c *Client, ctx context.Context, path string) (T, error) {
	raw, err := c.doJSON(ctx, http.MethodGet, c.apiURL(path), nil)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](raw)
}

func post[T any](c *Client, ctx context.Context, path string, body any) (T, error) {
	raw, err := c.doJSON(ctx, http.MethodPost, c.apiURL(path), body)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](raw)
}

func put[T any](c *Client, ctx context.Context, path string, body any) (T, error) {
	raw, err := c.doJSON(ctx, http.MethodPut, c.apiURL(path), body)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](raw)
}

func buildPath(path string, params url.Values) string {
	if len(params) == 0 {
		return path
	}
	return path + "?" + params.Encode()
}

func setInt64(values url.Values, key string, value *int64) {
	if value != nil {
		values.Set(key, strconv.FormatInt(*value, 10))
	}
}

func setString(values url.Values, key string, value *string) {
	if value != nil {
		values.Set(key, *value)
	}
}

func setAccountLookupQuery(
	values url.Values,
	params *AccountLookupParams,
	agencyKey string,
	accountKey string,
	digitKey string,
) {
	if params == nil {
		return
	}
	setInt64(values, agencyKey, params.Agency)
	setInt64(values, accountKey, params.Account)
	setString(values, digitKey, params.CheckDigit)
}
