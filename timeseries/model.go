package timeseries

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

/**
 * Scraped item
 */

type ScrapedItem struct {
	Source string
	Origin string
	Meta   Data
	Data   Data
}

func NewScrapedItem(source, origin string, data, meta interface{}) (ScrapedItem, error) {
	si := ScrapedItem{
		Source: source,
		Origin: origin,
		Meta:   NewData(meta),
		Data:   NewData(data),
	}
	return si, nil
}

// Data wraps the scraped data (for golang/pq)

type Data struct {
	Data interface{}
}

func (di Data) Value() (driver.Value, error) {
	j, err := json.Marshal(di)
	return j, err
}

func (di *Data) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	v := map[string]interface{}{}
	err := json.Unmarshal(source, &v)
	if err != nil {
		return err
	}

	if di.Data == nil {
		di.Data = v
		return nil
	}

	if err := mapstructure.Decode(v, &di.Data); err != nil {
		return err
	}
	return nil
}

func NewData(data interface{}) Data {
	d := Data{
		Data: data,
	}
	return d
}
